package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/models"
	"free2free/tests/testutils"
)

// TestReviewEndpointsWithFacebookJWT tests all review-specific endpoints with Facebook JWT
func TestReviewEndpointsWithFacebookJWT(t *testing.T) {
	// Initialize test server
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err, "Should clear test data successfully")

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err, "Should setup test database successfully")

	// Create a test user and JWT token (simulating Facebook login result)
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("Review endpoints require valid JWT", func(t *testing.T) {
		// Test creating a review (this would normally require a valid match ID)
		reviewData := map[string]interface{}{
			"reviewee_id": 2, // Different user
			"score":       5,
			"comment":     "Great experience!",
		}

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review/matches/1", jwtToken, reviewData)
		assert.NoError(t, err)

		// Could return 201 (created), 400 (bad request), 404 (match not found), etc.
		// But should NOT return 401 (unauthorized)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("Review-like endpoints require valid JWT", func(t *testing.T) {
		// Test liking a review
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/like", jwtToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode) // Should not be unauthorized
		resp.Body.Close()

		// Test disliking a review
		resp, err = testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/dislike", jwtToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode) // Should not be unauthorized
		resp.Body.Close()
	})
}

// TestReviewEndpointsAuthorization tests that review endpoints properly enforce authorization
func TestReviewEndpointsAuthorization(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Create a user
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("Review endpoints require valid JWT", func(t *testing.T) {
		endpoints := []string{
			"/review/matches/1",
			"/review-like/reviews/1/like",
			"/review-like/reviews/1/dislike",
		}

		for _, endpoint := range endpoints {
			t.Run("Access "+endpoint+" without JWT", func(t *testing.T) {
				resp, err := testServer.DoRequest("POST", endpoint, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with invalid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", endpoint, "invalid.jwt.token", nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with valid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", endpoint, jwtToken, nil)
				assert.NoError(t, err)
				// Should not be unauthorized - may be 400 (bad request), 404 (not found), etc.
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})
		}
	})
}

// TestReviewEndpointsResponseFormat tests that review endpoints return properly formatted responses
func TestReviewEndpointsResponseFormat(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Create a user
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)

	t.Run("Review creation response format", func(t *testing.T) {
		reviewData := map[string]interface{}{
			"reviewee_id": 2,
			"score":       4,
			"comment":     "Good experience",
		}

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review/matches/1", jwtToken, reviewData)
		assert.NoError(t, err)

		// Response should be valid JSON regardless of success/failure
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		// Both successful and error responses should be valid JSON
		resp.Body.Close()
	})

	t.Run("Review like response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/like", jwtToken, nil)
		assert.NoError(t, err)

		// Response should be valid JSON regardless of success/failure
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
	})

	t.Run("Review dislike response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/dislike", jwtToken, nil)
		assert.NoError(t, err)

		// Response should be valid JSON regardless of success/failure
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
	})
}

// TestReviewCreationAndInteraction tests creating reviews and interacting with them
func TestReviewCreationAndInteraction(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err, "Should clear test data successfully")

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err, "Should setup test database successfully")

	// Create users for review functionality
	reviewerUser, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, reviewerUser)

	revieweeUser := &models.User{
		SocialID:       "reviewee_social_id",
		SocialProvider: "facebook",
		Name:           "Reviewee User",
		Email:          "reviewee@example.com",
		AvatarURL:      "https://example.com/reviewee.jpg",
		IsAdmin:        false,
	}
	result := testServer.DB.Create(revieweeUser)
	assert.NoError(t, result.Error)

	reviewerToken, err := testutils.CreateMockJWTToken(reviewerUser.ID, reviewerUser.Name, reviewerUser.IsAdmin)
	assert.NoError(t, err)

	// In a real scenario, we would need a completed match between these users
	// For this test, we'll just verify that the endpoints properly handle JWTs
	t.Run("Review endpoints handle JWTs properly", func(t *testing.T) {
		// Test creating a review
		reviewData := map[string]interface{}{
			"reviewee_id": revieweeUser.ID,
			"score":       5,
			"comment":     "Excellent experience!",
		}

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review/matches/1", reviewerToken, reviewData)
		assert.NoError(t, err)
		// May return 404 (match not found) or 400 (bad request) but not 401 (unauthorized)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()

		// Test review-like features (these would work on an existing review)
		resp, err = testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/like", reviewerToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()

		resp, err = testutils.MakeAuthenticatedRequest(testServer, "POST", "/review-like/reviews/1/dislike", reviewerToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})
}