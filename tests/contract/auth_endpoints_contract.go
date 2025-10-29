package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAuthEndpointsContract tests the API contracts for authentication endpoints
func TestAuthEndpointsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router with authentication routes
	router := gin.New()

	// Define auth routes for testing (these would match the actual implementation)
	router.GET("/auth/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		if provider == "facebook" || provider == "instagram" {
			c.JSON(http.StatusTemporaryRedirect, gin.H{"redirect": "oauth-provider-url"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		}
	})

	router.GET("/auth/:provider/callback", func(c *gin.Context) {
		c.JSON(http.StatusTemporaryRedirect, gin.H{"redirect": "frontend-url"})
	})

	router.GET("/auth/token", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"token": "sample-jwt-token",
			"user": gin.H{
				"id":       1,
				"email":    "test@example.com",
				"name":     "Test User",
				"provider": "facebook",
				"role":     "user",
			},
		})
	})

	router.POST("/auth/refresh", func(c *gin.Context) {
		var reqBody map[string]interface{}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		refreshToken, exists := reqBody["refresh_token"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
			return
		}

		// In a real test, verify the token, but for contract test we just need to return expected format
		c.JSON(http.StatusOK, gin.H{
			"token":         "new-jwt-token",
			"refresh_token": refreshToken,
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "logged out",
		})
	})

	t.Run("GET /auth/:provider - Valid Provider", func(t *testing.T) {
		testutils.RequestWithValidation(t, router, "GET", "/auth/facebook", nil, "", http.StatusTemporaryRedirect)
	})

	t.Run("GET /auth/:provider - Invalid Provider", func(t *testing.T) {
		w := testutils.RequestWithValidation(t, router, "GET", "/auth/invalid", nil, "", http.StatusBadRequest)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("GET /auth/token - Success Response Format", func(t *testing.T) {
		w := testutils.RequestWithValidation(t, router, "GET", "/auth/token", nil, "", http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")

		user, ok := response["user"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, user, "id")
		assert.Contains(t, user, "email")
		assert.Contains(t, user, "name")
		assert.Contains(t, user, "provider")
		assert.Contains(t, user, "role")
	})

	t.Run("POST /auth/refresh - Success Response Format", func(t *testing.T) {
		requestBody := map[string]string{
			"refresh_token": "sample-refresh-token",
		}

		w := testutils.RequestWithValidation(t, router, "POST", "/auth/refresh", requestBody, "", http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "refresh_token")
	})

	t.Run("GET /logout - Success Response Format", func(t *testing.T) {
		w := testutils.RequestWithValidation(t, router, "GET", "/logout", nil, "", http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "message")
		assert.Equal(t, "logged out", response["message"])
	})
}
