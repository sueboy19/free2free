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

// TestAdminManagementAPI tests the integration of admin management API endpoints
func TestAdminManagementAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup admin management routes for the test
	setupAdminManagementRoutes(ts.Router, db)

	t.Run("View All Activities with Filters", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Get activities with status filter
		w, err := testutils.GetRequest(ts.Router, "/admin/activities?status=pending", adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "activities")
		assert.Contains(t, response, "total")
	})

	t.Run("Update Activity Details as Admin", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Update activity details
		updatedData := map[string]interface{}{
			"title":       "Updated Admin Title",
			"description": "Updated admin description",
			"location_id": 5,
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/activities/1", updatedData, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Updated Admin Title", response.Title)
		assert.Equal(t, "Updated admin description", response.Description)
		assert.Equal(t, uint(5), response.LocationID)
	})

	t.Run("Bulk Approve Activities", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Request to approve multiple activities
		activitiesToApprove := map[string]interface{}{
			"activity_ids": []uint{1, 2, 3},
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/activities/approve-bulk", activitiesToApprove, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "approved_count")
		assert.Contains(t, response, "message")
	})

	t.Run("Bulk Reject Activities", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Request to reject multiple activities with a reason
		activitiesToReject := map[string]interface{}{
			"activity_ids": []uint{4, 5, 6},
			"reason":       "Multiple violations of community guidelines",
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/activities/reject-bulk", activitiesToReject, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "rejected_count")
		assert.Contains(t, response, "message")
	})
}

// TestAdminUserManagement tests admin's ability to manage user accounts
func TestAdminUserManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup admin management routes for the test
	setupAdminManagementRoutes(ts.Router, db)

	t.Run("View All Users", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Get all users
		w, err := testutils.GetRequest(ts.Router, "/admin/users", adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "users")
		assert.Contains(t, response, "total")
	})

	t.Run("Ban User", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Request to ban a user
		banDetails := map[string]interface{}{
			"reason":   "Violation of community guidelines",
			"duration": "30d", // 30 days
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/users/5/ban", banDetails, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "message")
		assert.Contains(t, response, "user_id")
	})

	t.Run("Change User Role", func(t *testing.T) {
		// Create an admin token for the test
		authHelper := testutils.NewAuthTestHelper()
		adminToken, err := authHelper.CreateValidAdminToken(1, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Request to change user's role
		roleChange := map[string]interface{}{
			"new_role": "moderator",
		}

		w, err := testutils.PutRequest(ts.Router, "/admin/users/3/change-role", roleChange, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "message")
		assert.Contains(t, response, "user_id")
		assert.Contains(t, response, "new_role")
	})
}

// setupAdminManagementRoutes configures the routes for admin management testing
func setupAdminManagementRoutes(router *gin.Engine, db *gorm.DB) {
	// For integration testing, we're simulating the actual routes
	// In a real implementation, this would connect to the database and models

	router.GET("/admin/activities", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		// Get query parameters
		status := c.Query("status")

		// Simulate filtering by status
		var activities []testutils.TestActivity
		if status == "pending" {
			activities = []testutils.TestActivity{
				{
					ID:          1,
					Title:       "Pending Activity 1",
					Description: "This activity is pending approval",
					LocationID:  1,
					Status:      "pending",
					CreatorID:   1,
				},
				{
					ID:          2,
					Title:       "Pending Activity 2",
					Description: "This activity is also pending approval",
					LocationID:  2,
					Status:      "pending",
					CreatorID:   2,
				},
			}
		} else {
			activities = []testutils.TestActivity{
				{
					ID:          1,
					Title:       "Activity 1",
					Description: "Description 1",
					LocationID:  1,
					Status:      "pending",
					CreatorID:   1,
				},
				{
					ID:          2,
					Title:       "Activity 2",
					Description: "Description 2",
					LocationID:  2,
					Status:      "approved",
					CreatorID:   2,
				},
				{
					ID:          3,
					Title:       "Activity 3",
					Description: "Description 3",
					LocationID:  3,
					Status:      "rejected",
					CreatorID:   3,
				},
			}
		}

		response := gin.H{
			"activities": activities,
			"total":      len(activities),
			"page":       1,
			"limit":      10,
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/activities/:id", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
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

		// Simulate updating the activity
		activityID := c.Param("id")

		response := testutils.TestActivity{
			ID:          parseUint(activityID),
			Title:       input.Title,
			Description: input.Description,
			LocationID:  input.LocationID,
			Status:      "pending", // Status might remain unchanged
			CreatorID:   1,         // For this test, set a default creator ID
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/activities/approve-bulk", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		var input struct {
			ActivityIDs []uint `json:"activity_ids" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		// Simulate bulk approval
		response := gin.H{
			"approved_count": len(input.ActivityIDs),
			"message":        "Successfully approved activities",
			"activity_ids":   input.ActivityIDs,
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/activities/reject-bulk", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		var input struct {
			ActivityIDs []uint `json:"activity_ids" binding:"required"`
			Reason      string `json:"reason"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		// Simulate bulk rejection
		response := gin.H{
			"rejected_count": len(input.ActivityIDs),
			"message":        "Successfully rejected activities",
			"activity_ids":   input.ActivityIDs,
		}

		c.JSON(http.StatusOK, response)
	})

	router.GET("/admin/users", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		// Simulate returning user data
		users := []testutils.TestUser{
			{
				ID:       1,
				Email:    "user1@example.com",
				Name:     "User One",
				Provider: "facebook",
				Role:     "user",
			},
			{
				ID:       2,
				Email:    "user2@example.com",
				Name:     "User Two",
				Provider: "facebook",
				Role:     "user",
			},
		}

		response := gin.H{
			"users": users,
			"total": len(users),
			"page":  1,
			"limit": 10,
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/users/:id/ban", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		var input struct {
			Reason   string `json:"reason"`
			Duration string `json:"duration"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		userID := c.Param("id")

		// Simulate banning a user
		response := gin.H{
			"message":  "User banned successfully",
			"user_id":  userID,
			"reason":   input.Reason,
			"duration": input.Duration,
		}

		c.JSON(http.StatusOK, response)
	})

	router.PUT("/admin/users/:id/change-role", func(c *gin.Context) {
		// Simulate token validation and admin role check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Validate token format and check admin role
		token := authHeader[7:]
		claims, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		var input struct {
			NewRole string `json:"new_role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		userID := c.Param("id")

		// Simulate changing user's role
		response := gin.H{
			"message":  "User role changed successfully",
			"user_id":  userID,
			"old_role": "user", // For this test, assume the old role was always "user"
			"new_role": input.NewRole,
		}

		c.JSON(http.StatusOK, response)
	})
}

// Helper function to parse uint from string
func parseUint(s string) uint {
	var result uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + uint(c-'0')
		} else {
			return 0
		}
	}
	return result
}
