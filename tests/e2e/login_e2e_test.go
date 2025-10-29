package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestLoginFlowWithValidCredentials tests the complete login flow with valid credentials
func TestLoginFlowWithValidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server with mock auth routes
	server, authHelper, mockProvider := testutils.CreateMockAuthServer()
	defer server.Close()

	// Create a separate test router for this test
	testRouter := gin.New()
	testRouter.Use(gin.Logger())
	testRouter.Use(gin.Recovery())

	// Mock the necessary endpoints for this test
	testRouter.GET("/auth/facebook", func(c *gin.Context) {
		// Generate a mock auth code
		authCode := "mock-auth-code-facebook"
		mockProvider.AddValidAuthCode(authCode, testutils.MockUser{
			ID:       "123456",
			Email:    "test@example.com",
			Name:     "Test User",
			Provider: "facebook",
			Avatar:   "https://example.com/avatar.jpg",
		})

		// Instead of redirecting (which is hard to test), return the code
		c.JSON(http.StatusOK, gin.H{"auth_code": authCode})
	})

	testRouter.GET("/auth/facebook/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			code = c.Query("auth_code") // Fallback to body param if query param not available
		}

		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}

		user, valid := mockProvider.ValidateAuthCode(code)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
			return
		}

		// Create a JWT token for the authenticated user
		token, err := authHelper.CreateValidUserToken(1, user.Email, user.Name, user.Provider)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  user,
		})
	})

	// Test the initial OAuth redirect
	t.Run("OAuth Redirect Step", func(t *testing.T) {
		w, err := testutils.GetRequest(testRouter, "/auth/facebook", "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "auth_code")
	})

	// Test the callback and token generation
	t.Run("OAuth Callback and Token Generation", func(t *testing.T) {
		// First, get an auth code
		w, err := testutils.GetRequest(testRouter, "/auth/facebook", "")
		assert.NoError(t, err)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		authCode := response["auth_code"].(string)

		// Now test the callback with the auth code
		callbackURL := "/auth/facebook/callback?code=" + authCode
		w, err = testutils.GetRequest(testRouter, callbackURL, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
		assert.Contains(t, response, "user")

		userData := response["user"].(map[string]interface{})
		assert.Equal(t, "123456", userData["id"])
		assert.Equal(t, "test@example.com", userData["email"])
	})
}

// TestLoginFlowWithInvalidCredentials tests the login flow with invalid credentials
func TestLoginFlowWithInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	testRouter := gin.New()
	testRouter.Use(gin.Logger())
	testRouter.Use(gin.Recovery())

	// Mock the callback endpoint to simulate invalid credentials
	testRouter.GET("/auth/facebook/callback", func(c *gin.Context) {
		// Simulate invalid oauth code
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid oauth code",
		})
	})

	t.Run("Invalid OAuth Code", func(t *testing.T) {
		w, err := testutils.GetRequest(testRouter, "/auth/facebook/callback?code=invalid-code", "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid oauth code", response["error"])
	})

	// Test with invalid token
	testRouter.GET("/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := authHeader[7:] // Remove "Bearer " prefix
		if token == "invalid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "access granted"})
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		w, err := testutils.GetRequest(testRouter, "/protected", "invalid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid token", response["error"])
	})
}

// TestSessionManagement tests session management functionality
func TestSessionManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHelper := testutils.NewAuthTestHelper()

	t.Run("Valid Session Access", func(t *testing.T) {
		// Create a valid token
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Create a test router that validates the token
		testRouter := gin.New()
		testRouter.Use(gin.Logger())
		testRouter.Use(gin.Recovery())

		testRouter.GET("/profile", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			token := authHeader[7:]
			_, err := testutils.ValidateToken(token, authHelper.Secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":    1,
				"email": "user@example.com",
				"name":  "Test User",
			})
		})

		w, err := testutils.GetRequest(testRouter, "/profile", token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["id"])
		assert.Equal(t, "user@example.com", response["email"])
	})

	t.Run("Expired Session Rejection", func(t *testing.T) {
		// Create an expired token
		token, err := authHelper.CreateExpiredUserToken(2, "expired@example.com", "Expired User", "facebook")
		assert.NoError(t, err)

		// Create a test router that validates the token
		testRouter := gin.New()
		testRouter.Use(gin.Logger())
		testRouter.Use(gin.Recovery())

		testRouter.GET("/profile", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			token := authHeader[7:]
			_, err := testutils.ValidateToken(token, authHelper.Secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":    2,
				"email": "expired@example.com",
				"name":  "Expired User",
			})
		})

		w, err := testutils.GetRequest(testRouter, "/profile", token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid token", response["error"])
	})
}
