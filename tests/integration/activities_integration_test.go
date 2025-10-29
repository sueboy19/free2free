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

// TestActivitiesIntegration tests the integration of activities functionality
func TestActivitiesIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup routes for the test
	setupActivitiesRoutes(ts.Router, db)

	t.Run("Create Activity Integration", func(t *testing.T) {
		// Create a valid JWT token for the test
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		requestBody := map[string]interface{}{
			"title":       "Integration Test Activity",
			"description": "This is a test activity for integration testing",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Integration Test Activity", response.Title)
		assert.Equal(t, "This is a test activity for integration testing", response.Description)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, uint(1), response.CreatorID)
	})

	t.Run("Get Activity Integration", func(t *testing.T) {
		// Test getting an activity after creating it
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// First, create an activity
		requestBody := map[string]interface{}{
			"title":       "Get Test Activity",
			"description": "Activity for get test",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdActivity testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &createdActivity)
		assert.NoError(t, err)

		// Now get the activity by ID
		getURL := "/api/activities/" + string(rune('0'+int(createdActivity.ID)))
		w, err = testutils.GetRequest(ts.Router, getURL, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedActivity testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &retrievedActivity)
		assert.NoError(t, err)

		assert.Equal(t, createdActivity.ID, retrievedActivity.ID)
		assert.Equal(t, createdActivity.Title, retrievedActivity.Title)
		assert.Equal(t, createdActivity.Description, retrievedActivity.Description)
	})

	t.Run("Update Activity Integration", func(t *testing.T) {
		// Test updating an activity
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// First, create an activity to update
		requestBody := map[string]interface{}{
			"title":       "Original Title",
			"description": "Original Description",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var originalActivity testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &originalActivity)
		assert.NoError(t, err)

		// Now update the activity
		updateBody := map[string]interface{}{
			"title":       "Updated Title",
			"description": "Updated Description",
			"location_id": 2,
		}

		updateURL := "/api/activities/" + string(rune('0'+int(originalActivity.ID)))
		w, err = testutils.PutRequest(ts.Router, updateURL, updateBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var updatedActivity testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &updatedActivity)
		assert.NoError(t, err)

		assert.Equal(t, originalActivity.ID, updatedActivity.ID)
		assert.Equal(t, "Updated Title", updatedActivity.Title)
		assert.Equal(t, "Updated Description", updatedActivity.Description)
		assert.Equal(t, uint(2), updatedActivity.LocationID)
	})
}

// TestActivitiesValidationIntegration tests validation in the activities integration
func TestActivitiesValidationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Create test database
	db, err := testutils.CreateTestDB()
	assert.NoError(t, err)

	// Setup routes for the test
	setupActivitiesRoutes(ts.Router, db)

	t.Run("Validation with Invalid Data", func(t *testing.T) {
		// Create a valid JWT token for the test
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Try to create an activity with invalid data (e.g., empty title)
		requestBody := map[string]interface{}{
			"title":       "", // Empty title should fail validation
			"description": "Valid description",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		// We expect a validation error (typically 400)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("Validation with Missing Required Fields", func(t *testing.T) {
		// Create a valid JWT token for the test
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Try to create an activity with missing required fields
		requestBody := map[string]interface{}{
			"title": "Valid Title", // Title is present
			// Missing description and location_id
		}

		w, err := testutils.PostRequest(ts.Router, "/api/activities", requestBody, token)
		assert.NoError(t, err)
		// We expect a validation error (typically 400)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// setupActivitiesRoutes configures the routes for activities testing
func setupActivitiesRoutes(router *gin.Engine, db *gorm.DB) {
	// For integration testing, we're simulating the actual routes
	// In a real implementation, this would connect to the database and models

	router.POST("/api/activities", func(c *gin.Context) {
		// Simulate token validation (in real app, this would be middleware)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

		// Simulate creating an activity in the database
		// In real implementation, we would use the db connection
		activity := testutils.TestActivity{
			ID:          1, // In real app, this would be generated
			Title:       input.Title,
			Description: input.Description,
			LocationID:  input.LocationID,
			Status:      "pending", // Default status
			CreatorID:   1,         // Extracted from token in real app
		}

		c.JSON(http.StatusCreated, activity)
	})

	router.GET("/api/activities/:id", func(c *gin.Context) {
		// Simulate token validation
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Simulate retrieving an activity from the database
		activity := testutils.TestActivity{
			ID:          1,
			Title:       "Test Activity",
			Description: "Test Description",
			LocationID:  1,
			Status:      "approved",
			CreatorID:   1,
		}

		c.JSON(http.StatusOK, activity)
	})

	router.PUT("/api/activities/:id", func(c *gin.Context) {
		// Simulate token validation
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

		// Simulate updating an activity in the database
		activity := testutils.TestActivity{
			ID:          1, // From URL parameter in real app
			Title:       input.Title,
			Description: input.Description,
			LocationID:  input.LocationID,
			Status:      "pending", // Status might not change
			CreatorID:   1,         // From token in real app
		}

		c.JSON(http.StatusOK, activity)
	})
}
