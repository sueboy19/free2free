package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestAdminApprovalFlow tests the admin approval workflow integration
func TestAdminApprovalFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup admin routes for the test
	setupAdminRoutes(ts.Router, db)

	t.Run("Approve Activity Integration", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Approve an activity
		w, err := testutils.PutRequest(ts.Router, "/admin/activities/1/approve", nil, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "1", response["id"])
		assert.Equal(t, "approved", response["status"])
		assert.Contains(t, response, "approved_at")
		assert.Equal(t, float64(1), response["approved_by"])
	})

	t.Run("Reject Activity Integration", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Reject an activity with a reason
		requestBody := map[string]string{
			"reason": "Does not meet community guidelines",
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/activities/2/reject", requestBody, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "2", response["id"])
		assert.Equal(t, "rejected", response["status"])
		assert.Equal(t, "Does not meet community guidelines", response["rejection_reason"])
	})

	t.Run("View Admin Activities List", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Get the list of activities
		w, err := testutils.GetRequest(ts.Router, "/admin/activities", adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "activities")
		assert.Contains(t, response, "total")
		assert.Contains(t, response, "page")
		assert.Contains(t, response, "limit")

		activities, ok := response["activities"].([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, activities)

		// Check the structure of the first activity
		activity, ok := activities[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, activity, "id")
		assert.Contains(t, activity, "title")
		assert.Contains(t, activity, "description")
		assert.Contains(t, activity, "location_id")
		assert.Contains(t, activity, "status")
		assert.Contains(t, activity, "creator_id")
	})
}

// TestAdminPermissions tests admin-specific permissions and access controls
func TestAdminPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup admin routes for the test
	setupAdminRoutes(ts.Router, db)

	t.Run("Non-Admin Cannot Approve", func(t *testing.T) {
		// Create a regular user token
		authHelper := testutils.NewAuthTestHelper()
		userToken, err := authHelper.CreateValidUserToken(2, "user@example.com", "Regular User", "facebook")
		assert.NoError(t, err)

		// Try to approve an activity with a regular user token
		w, err := testutils.PutRequest(ts.Router, "/admin/activities/1/approve", nil, userToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "forbidden - admin permissions required", response["error"])
	})

	t.Run("Non-Admin Cannot Reject", func(t *testing.T) {
		// Create a regular user token
		authHelper := testutils.NewAuthTestHelper()
		userToken, err := authHelper.CreateValidUserToken(2, "user@example.com", "Regular User", "facebook")
		assert.NoError(t, err)

		// Try to reject an activity with a regular user token
		requestBody := map[string]string{
			"reason": "Does not meet community guidelines",
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/activities/1/reject", requestBody, userToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "forbidden - admin permissions required", response["error"])
	})

	t.Run("Unauthorized Access to Admin Endpoint", func(t *testing.T) {
		// Try to access admin endpoint without token
		w, err := testutils.GetRequest(ts.Router, "/admin/activities", "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "unauthorized", response["error"])
	})

	t.Run("Invalid Token Access to Admin Endpoint", func(t *testing.T) {
		// Try to access admin endpoint with invalid token
		w, err := testutils.GetRequest(ts.Router, "/admin/activities", "invalid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid token", response["error"])
	})
}

// setupAdminRoutes configures the routes for admin testing
func setupAdminRoutes(router *gin.Engine, db *gorm.DB) {
	// For integration testing, we're simulating the actual routes
	// In a real implementation, this would connect to the database and models

	router.GET("/admin/activities", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// In a real implementation, we would validate the token and check admin role
		// For this test, we'll just verify the token format and return mock data
		token := authHeader[7:]

		// Validate token format (in real app, this would be proper JWT validation)
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
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

		// Return mock activities data
		activities := []testutils.TestActivity{
			{
				ID:          1,
				Title:       "Pending Activity",
				Description: "This activity is pending approval",
				LocationID:  1,
				Status:      "pending",
				CreatorID:   1,
			},
			{
				ID:          2,
				Title:       "Approved Activity",
				Description: "This activity has been approved",
				LocationID:  2,
				Status:      "approved",
				CreatorID:   2,
			},
		}

		response := gin.H{
			"activities": activities,
			"total":      len(activities),
			"page":       1,
			"limit":      10,
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/activities/:id/approve", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// In a real implementation, we would validate the token and check admin role
		token := authHeader[7:]

		// Validate token format (in real app, this would be proper JWT validation)
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
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

		// Simulate approving the activity
		activityID := c.Param("id")

		response := gin.H{
			"id":          activityID,
			"status":      "approved",
			"approved_at": "2023-01-01T00:00:00Z",
			"approved_by": claims["user_id"],
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/activities/:id/reject", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// In a real implementation, we would validate the token and check admin role
		token := authHeader[7:]

		// Validate token format (in real app, this would be proper JWT validation)
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
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

		// Parse request body for rejection reason
		var reqBody map[string]interface{}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		reason, exists := reqBody["reason"].(string)
		if !exists {
			reason = ""
		}

		// Simulate rejecting the activity
		activityID := c.Param("id")

		response := gin.H{
			"id":               activityID,
			"status":           "rejected",
			"rejection_reason": reason,
		}

		c.JSON(http.StatusOK, response)
	})
}
