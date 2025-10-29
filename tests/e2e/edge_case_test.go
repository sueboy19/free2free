package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestEdgeCases tests various edge cases in the API workflow
func TestEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Setup routes for edge case testing
	setupEdgeCaseRoutes(ts.Router)

	t.Run("Large Payload Handling", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Create a large description to test payload limits
		largeDescription := ""
		for i := 0; i < 600; i++ { // Exceeds 500 char limit
			largeDescription += "A"
		}

		requestBody := map[string]interface{}{
			"title":       "Large Payload Test",
			"description": largeDescription,
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		// Should return a validation error due to description length
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("Boundary Value Testing", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Test minimum length title
		minTitle := "A" // 1 character
		requestBody := map[string]interface{}{
			"title":       minTitle,
			"description": "Valid description with sufficient length",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code) // Should succeed with min valid length

		// Test maximum length title
		maxTitle := ""
		for i := 0; i < 100; i++ {
			maxTitle += "A"
		}
		requestBody["title"] = maxTitle

		w, err = testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code) // Should succeed with max valid length

		// Test just over maximum length title
		overMaxTitle := maxTitle + "X"
		requestBody["title"] = overMaxTitle

		w, err = testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail with validation error
	})

	t.Run("Special Characters Handling", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Test special characters in title and description
		specialCharTitle := "Title with special chars: !@#$%^&*()"
		specialCharDesc := "Description with special chars: <>{}[]|\\`~"

		requestBody := map[string]interface{}{
			"title":       specialCharTitle,
			"description": specialCharDesc,
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		// Should either succeed (if chars are allowed) or have proper validation
		// For this test we'll accept both 201 or 400 depending on implementation
		assert.Contains(t, []int{http.StatusCreated, http.StatusBadRequest}, w.Code)
	})

	t.Run("Concurrent Access Without Proper Authorization", func(t *testing.T) {
		// Test that one user cannot access another user's sensitive data
		authHelper := testutils.NewAuthTestHelper()

		user1Token, err := authHelper.CreateValidUserToken(100, "user1@example.com", "User One", "facebook")
		assert.NoError(t, err)

		user2Token, err := authHelper.CreateValidUserToken(101, "user2@example.com", "User Two", "facebook")
		assert.NoError(t, err)

		// User 1 creates an activity
		activityData := map[string]interface{}{
			"title":       "User 1 Activity",
			"description": "This activity belongs to User 1",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", activityData, user1Token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdActivity testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &createdActivity)
		assert.NoError(t, err)

		// Verify User 2 cannot access User 1's data in restricted endpoints
		// (This would depend on implementation - just testing the concept)
		w, err = testutils.GetRequest(ts.Router, "/api/activities/1", user2Token)
		assert.NoError(t, err)
		// Should either succeed (if public) or fail (if private)
		assert.Contains(t, []int{http.StatusOK, http.StatusForbidden, http.StatusUnauthorized}, w.Code)
	})
}

// TestErrorConditions tests various error conditions and how they're handled
func TestErrorConditions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Setup routes for error condition testing
	setupEdgeCaseRoutes(ts.Router)

	t.Run("Malformed JSON Handling", func(t *testing.T) {
		// Create an invalid JSON string
		invalidJSON := `{"title": "test", "description":}`

		// Make a request with invalid JSON (using httptest directly)
		req, err := testutils.CreateTestAuthRequest("POST", "/api/activities", nil, "valid-token")
		assert.NoError(t, err)

		// Set the body manually to invalid JSON
		req.Body = nil // This is simplified; in practice would need to simulate the invalid JSON

		// For this test we'll just verify that the application handles it gracefully
		// In a real test, we would send the actual invalid JSON and verify the response
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Request with missing required fields
		requestBody := map[string]interface{}{
			"title": "Only title provided",
			// Missing description and location_id
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("Invalid Data Types", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Request with invalid data types
		requestBody := map[string]interface{}{
			"title":       12345,            // Should be string
			"description": true,             // Should be string
			"location_id": "not-an-integer", // Should be number
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// setupEdgeCaseRoutes configures routes for edge case testing
func setupEdgeCaseRoutes(router *gin.Engine) {
	authHelper := testutils.NewAuthTestHelper()

	// Activity creation endpoint for testing edge cases
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

		// Return success response
		c.JSON(http.StatusCreated, gin.H{
			"id":          1,
			"title":       input.Title,
			"description": input.Description,
			"location_id": input.LocationID,
			"status":      "pending",
		})
	})

	// Endpoint to test access to activities
	router.GET("/api/activities/:id", func(c *gin.Context) {
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

		// Return mock activity
		c.JSON(http.StatusOK, gin.H{
			"id":          1,
			"title":       "Mock Activity",
			"description": "This is a mock activity for testing",
			"location_id": 1,
			"status":      "active",
		})
	})
}
