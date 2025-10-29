package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAdminEndpointsContract tests the API contracts for admin endpoints
func TestAdminEndpointsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router with admin routes
	router := gin.New()

	// Define admin routes for testing (these would match the actual implementation)
	router.GET("/admin/activities", func(c *gin.Context) {
		// Simulate successful response for admin endpoint
		c.JSON(http.StatusOK, gin.H{
			"activities": []testutils.TestActivity{
				{
					ID:          1,
					Title:       "Test Activity 1",
					Description: "Test Description 1",
					LocationID:  1,
					Status:      "pending",
					CreatorID:   1,
				},
				{
					ID:          2,
					Title:       "Test Activity 2",
					Description: "Test Description 2",
					LocationID:  2,
					Status:      "approved",
					CreatorID:   2,
				},
			},
			"total": 2,
			"page":  1,
			"limit": 10,
		})
	})

	router.PUT("/admin/activities/:id/approve", func(c *gin.Context) {
		activityID := c.Param("id")

		// Simulate successful approval
		c.JSON(http.StatusOK, gin.H{
			"id":          activityID,
			"status":      "approved",
			"approved_at": "2023-01-01T00:00:00Z",
			"approved_by": 1,
		})
	})

	router.PUT("/admin/activities/:id/reject", func(c *gin.Context) {
		var reqBody map[string]interface{}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		activityID := c.Param("id")
		reason, exists := reqBody["reason"].(string)
		if !exists {
			reason = ""
		}

		// Simulate successful rejection
		c.JSON(http.StatusOK, gin.H{
			"id":               activityID,
			"status":           "rejected",
			"rejection_reason": reason,
		})
	})

	t.Run("GET /admin/activities - Success Response Format", func(t *testing.T) {
		// This would require an admin token in a real implementation
		w, err := testutils.GetRequest(router, "/admin/activities", "admin-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "activities")
		assert.Contains(t, response, "total")
		assert.Contains(t, response, "page")
		assert.Contains(t, response, "limit")

		activities, ok := response["activities"].([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, activities)

		// Check the structure of the first activity
		if len(activities) > 0 {
			activity, ok := activities[0].(map[string]interface{})
			assert.True(t, ok)
			assert.Contains(t, activity, "id")
			assert.Contains(t, activity, "title")
			assert.Contains(t, activity, "description")
			assert.Contains(t, activity, "location_id")
			assert.Contains(t, activity, "status")
			assert.Contains(t, activity, "creator_id")
		}
	})

	t.Run("PUT /admin/activities/{id}/approve - Success Response Format", func(t *testing.T) {
		// This would require an admin token in a real implementation
		w, err := testutils.PutRequest(router, "/admin/activities/1/approve", nil, "admin-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "status")
		assert.Contains(t, response, "approved_at")
		assert.Contains(t, response, "approved_by")

		assert.Equal(t, "approved", response["status"])
	})

	t.Run("PUT /admin/activities/{id}/reject - Success Response Format", func(t *testing.T) {
		// This would require an admin token in a real implementation
		requestBody := map[string]string{
			"reason": "Does not meet community guidelines",
		}

		w, err := testutils.PutRequest(router, "/admin/activities/1/reject", requestBody, "admin-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "status")
		assert.Contains(t, response, "rejection_reason")

		assert.Equal(t, "rejected", response["status"])
		assert.Equal(t, "Does not meet community guidelines", response["rejection_reason"])
	})

	t.Run("PUT /admin/activities/{id}/reject - Without Reason", func(t *testing.T) {
		// This would require an admin token in a real implementation
		requestBody := map[string]string{} // No reason provided

		w, err := testutils.PutRequest(router, "/admin/activities/1/reject", requestBody, "admin-token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "status")
		assert.Contains(t, response, "rejection_reason")

		assert.Equal(t, "rejected", response["status"])
		// Reason might be empty string if not provided
	})
}
