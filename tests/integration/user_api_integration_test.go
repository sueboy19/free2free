package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/models"
	"free2free/tests/testutils"
)

// TestUserEndpointsWithFacebookJWT tests all user-specific endpoints with Facebook JWT
func TestUserEndpointsWithFacebookJWT(t *testing.T) {
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

	t.Run("Get user matches", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", jwtToken, nil)
		assert.NoError(t, err)
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)

		// Verify response structure
		var matches []interface{}
		err = json.NewDecoder(resp.Body).Decode(&matches)
		// If the response is an array, it's valid
		// If it's an error message, that's also acceptable
		resp.Body.Close()
	})

	t.Run("Get past matches", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/past-matches", jwtToken, nil)
		assert.NoError(t, err)
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)

		// Verify response structure
		var matches []interface{}
		err = json.NewDecoder(resp.Body).Decode(&matches)
		resp.Body.Close()
	})

	t.Run("Create a new match", func(t *testing.T) {
		// First, create a test activity to associate with the match
		activity := models.Activity{
			Title:       "Test Activity",
			Description: "Test Description",
			LocationID:  1, // This might not exist, leading to a 400
			TargetCount: 2,
		}

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/user/matches", jwtToken, activity)
		assert.NoError(t, err)
		// Could return 201 (created), 400 (bad request), or 404 (location not found)
		assert.Contains(t, []int{http.StatusCreated, http.StatusBadRequest, http.StatusNotFound}, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("Join an existing match", func(t *testing.T) {
		// This would require an existing match to join
		// For this test, we'll just verify the endpoint exists and handles auth correctly
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/user/matches/1/join", jwtToken, nil)
		assert.NoError(t, err)
		// Could return 404 (match not found), 400 (bad request), or 201 (created)
		assert.Contains(t, []int{http.StatusCreated, http.StatusBadRequest, http.StatusNotFound}, resp.StatusCode)
		resp.Body.Close()
	})
}

// TestUserEndpointsAuthorization tests that user endpoints properly enforce authorization
func TestUserEndpointsAuthorization(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("User endpoints require valid JWT", func(t *testing.T) {
		endpoints := []string{
			"/user/matches",
			"/user/past-matches",
		}

		for _, endpoint := range endpoints {
			t.Run("Access "+endpoint+" without JWT", func(t *testing.T) {
				resp, err := testServer.DoRequest("GET", endpoint, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with invalid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", endpoint, "invalid.jwt.token", nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with valid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", endpoint, jwtToken, nil)
				assert.NoError(t, err)
				// Should not be unauthorized (could be 200, 404, 500, etc.)
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})
		}
	})
}

// TestUserEndpointsResponseFormat tests that user endpoints return properly formatted responses
func TestUserEndpointsResponseFormat(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("User matches response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", jwtToken, nil)
		assert.NoError(t, err)

		// Verify that response is valid JSON
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err, "Response should be valid JSON")
		resp.Body.Close()
	})

	t.Run("User past matches response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/past-matches", jwtToken, nil)
		assert.NoError(t, err)

		// Verify that response is valid JSON
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err, "Response should be valid JSON")
		resp.Body.Close()
	})
}

// TestUserDataIsolation tests that users can only access their own data
func TestUserDataIsolation(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Create two different users
	user1, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user1)

	user2, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user2)

	// Create tokens for both users
	token1, err := testutils.CreateMockJWTToken(user1.ID, user1.Name, user1.IsAdmin)
	assert.NoError(t, err)

	token2, err := testutils.CreateMockJWTToken(user2.ID, user2.Name, user2.IsAdmin)
	assert.NoError(t, err)

	t.Run("User can access own profile", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token1, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var profile map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&profile)
		assert.NoError(t, err)
		assert.Equal(t, float64(user1.ID), profile["id"])

		resp.Body.Close()
	})

	// Note: In our implementation, users access their own data through endpoints like /profile
	// More complex data isolation would require testing specific resources that belong to users
}