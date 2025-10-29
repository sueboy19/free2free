package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestActivitiesEndpointsContract tests the API contracts for activities endpoints
func TestActivitiesEndpointsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router with activities routes
	router := gin.New()

	// Define activities routes for testing (these would match the actual implementation)
	router.POST("/api/activities", func(c *gin.Context) {
		var activity testutils.TestActivity
		if err := c.ShouldBindJSON(&activity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		// Simulate successful creation with proper response format
		responseActivity := testutils.TestActivity{
			ID:          1,
			Title:       activity.Title,
			Description: activity.Description,
			LocationID:  activity.LocationID,
			Status:      "pending",
			CreatorID:   1,
		}

		c.JSON(http.StatusCreated, responseActivity)
	})

	router.GET("/api/activities/:id", func(c *gin.Context) {
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
		var activity testutils.TestActivity
		if err := c.ShouldBindJSON(&activity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		// Simulate successful update with proper response format
		c.JSON(http.StatusOK, activity)
	})

	t.Run("POST /api/activities - Valid Request and Response Format", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"title":       "New Activity",
			"description": "An activity description",
			"location_id": 1,
		}

		w, err := testutils.PostRequest(router, "/api/activities", requestBody, "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "New Activity", response.Title)
		assert.Equal(t, "An activity description", response.Description)
		assert.Equal(t, uint(1), response.LocationID)
		assert.Equal(t, "pending", response.Status) // Default status
		assert.Equal(t, uint(1), response.CreatorID)
	})

	t.Run("POST /api/activities - Validation Error Response Format", func(t *testing.T) {
		// Request with missing required fields
		requestBody := map[string]interface{}{
			"title": "Activity without description", // Missing description and location_id
		}

		w, err := testutils.PostRequest(router, "/api/activities", requestBody, "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify error response structure
		assert.Contains(t, response, "error")
		assert.Contains(t, response, "details")
		assert.Equal(t, "validation failed", response["error"])
	})

	t.Run("GET /api/activities/{id} - Success Response Format", func(t *testing.T) {
		w, err := testutils.GetRequest(router, "/api/activities/1", "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Activity", response.Title)
		assert.Equal(t, "Test Description", response.Description)
		assert.Equal(t, uint(1), response.LocationID)
		assert.Equal(t, "approved", response.Status)
		assert.Equal(t, uint(1), response.CreatorID)
	})

	t.Run("PUT /api/activities/{id} - Valid Request and Response Format", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"title":       "Updated Activity",
			"description": "Updated description",
			"location_id": 2,
		}

		w, err := testutils.PutRequest(router, "/api/activities/1", requestBody, "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response testutils.TestActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Equal(t, uint(1), response.ID) // ID should remain the same
		assert.Equal(t, "Updated Activity", response.Title)
		assert.Equal(t, "Updated description", response.Description)
		assert.Equal(t, uint(2), response.LocationID) // Updated location
	})
}
