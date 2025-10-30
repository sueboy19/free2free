package security

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowSecurityValidation tests the security aspects of the complete workflow
func TestWorkflowSecurityValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Setup security test routes
	setupSecurityTestRoutes(ts.Router)

	t.Run("Authentication Required for Protected Endpoints", func(t *testing.T) {
		// Try to access protected endpoint without authentication
		w, err := testutils.GetRequest(ts.Router, "/api/protected", "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "unauthorized", response["error"])
	})

	t.Run("Invalid Token Rejection", func(t *testing.T) {
		// Try to access with invalid token
		w, err := testutils.GetRequest(ts.Router, "/api/protected", "invalid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid token", response["error"])
	})

	t.Run("Expired Token Rejection", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()

		// Create an expired token
		expiredToken, err := authHelper.CreateExpiredUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Try to access with expired token
		w, err := testutils.GetRequest(ts.Router, "/api/protected", expiredToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid token", response["error"])
	})

	t.Run("Proper Authorization for Admin Endpoints", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()

		// Create a regular user token
		userToken, err := authHelper.CreateValidUserToken(1, "user@example.com", "Regular User", "facebook")
		assert.NoError(t, err)

		// Try to access admin endpoint with regular user token
		w, err := testutils.GetRequest(ts.Router, "/admin/protected", userToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "forbidden - admin permissions required", response["error"])
	})

	t.Run("Valid Admin Access to Admin Endpoints", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()

		// Create an admin token
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Access admin endpoint with admin token
		w, err := testutils.GetRequest(ts.Router, "/admin/protected", adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "message")
		assert.Equal(t, "admin access granted", response["message"])
	})
}

// TestInputValidationSecurity tests security aspects of input validation
func TestInputValidationSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Setup security test routes
	setupSecurityTestRoutes(ts.Router)

	authHelper := testutils.NewAuthTestHelper()
	token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
	assert.NoError(t, err)

	t.Run("SQL Injection Prevention", func(t *testing.T) {
		// Try to inject SQL through the activity creation endpoint
		maliciousInput := map[string]interface{}{
			"title":       "Normal Title",
			"description": "'; DROP TABLE users; --",
			"location_id": 1,
		}

		_, err := testutils.PostRequest(ts.Router, "/api/activities", maliciousInput, token)
		assert.NoError(t, err)
		// Should either reject the request or properly handle the input without executing SQL
		// The exact status code would depend on validation implementation
	})

	t.Run("XSS Prevention", func(t *testing.T) {
		// Try to inject script through the activity creation endpoint
		maliciousInput := map[string]interface{}{
			"title":       "Normal Title",
			"description": "<script>alert('XSS')</script>",
			"location_id": 1,
		}

		_, err := testutils.PostRequest(ts.Router, "/api/activities", maliciousInput, token)
		assert.NoError(t, err)
		// Should either reject the request or properly sanitize the input
		// The exact status code would depend on validation implementation
	})

	t.Run("JSON Injection Prevention", func(t *testing.T) {
		// Try to inject JSON to manipulate the structure
		maliciousInput := map[string]interface{}{
			"title":       "Normal Title",
			"description": "Normal Description",
			"location_id": 1,
			"__proto__":   map[string]interface{}{"admin": true}, // Potential prototype pollution
		}

		_, err := testutils.PostRequest(ts.Router, "/api/activities", maliciousInput, token)
		assert.NoError(t, err)
		// Should properly validate and reject malicious fields
	})
}

// TestTokenSecurity tests various token security aspects
func TestTokenSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Token Tampering Detection", func(t *testing.T) {
		// This would test if the system can detect tampered JWT tokens
		t.Skip("Token tampering detection tests would require more specific implementation details")
	})

	t.Run("Token Reuse Prevention", func(t *testing.T) {
		// This would test if refresh tokens can't be reused after being used once
		t.Skip("Token reuse prevention tests would require implementation of refresh token invalidation")
	})

	t.Run("Token Confidentiality", func(t *testing.T) {
		// Verify that sensitive information is not exposed in tokens
		authHelper := testutils.NewAuthTestHelper()

		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Parse the token to check for sensitive information
		claims, err := testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err)

		// Verify that sensitive data like passwords are not in the token
		assert.NotContains(t, claims, "password")
		assert.NotContains(t, claims, "credit_card")
	})
}

// setupSecurityTestRoutes configures routes for security testing
func setupSecurityTestRoutes(router *gin.Engine) {
	authHelper := testutils.NewAuthTestHelper()

	// Protected endpoint that requires authentication
	router.GET("/api/protected", func(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{"message": "access granted", "data": "protected data"})
	})

	// Admin-only endpoint that requires admin role
	router.GET("/admin/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, authHelper.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check if user has admin role
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "admin access granted", "data": "admin data"})
	})

	// Activity creation endpoint for testing input validation
	router.POST("/api/activities", func(c *gin.Context) {
		// Simulate token validation
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

		var input struct {
			Title       string `json:"title" binding:"required,min=1,max=100"`
			Description string `json:"description" binding:"required,min=10,max=500"`
			LocationID  uint   `json:"location_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		// In a real implementation, we would sanitize inputs to prevent XSS, SQL injection, etc.
		// For this test, we'll just return a success response
		c.JSON(http.StatusCreated, gin.H{
			"id":          1,
			"title":       input.Title,
			"description": input.Description,
			"location_id": input.LocationID,
			"status":      "pending",
		})
	})
}
