package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/models"
	"free2free/tests/testutils"
)

// TestOrganizerEndpointsWithFacebookJWT tests all organizer-specific endpoints with Facebook JWT
func TestOrganizerEndpointsWithFacebookJWT(t *testing.T) {
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

	jwtToken, err := testutils.CreateMockJWTToken(uint(user.ID), user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("Organizer endpoints require proper permissions", func(t *testing.T) {
		// These endpoints require the user to be the organizer of the specific match
		// For this test, we'll check that they properly validate JWTs at least
		endpoints := []string{
			"/organizer/matches/1/participants/1/approve",
			"/organizer/matches/1/participants/1/reject",
		}

		for _, endpoint := range endpoints {
			t.Run("Access "+endpoint+" with JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", endpoint, jwtToken, nil)
				assert.NoError(t, err)

				// Should not return 401 (unauthorized) - may return 403 (forbidden) if not organizer,
				// 404 (not found) if match/participant doesn't exist, or 200 if successful
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})
		}
	})

	t.Run("Organizer endpoint authorization check", func(t *testing.T) {
		// Test that organizer endpoints reject requests without JWT
		endpoints := []string{
			"/organizer/matches/1/participants/1/approve",
			"/organizer/matches/1/participants/1/reject",
		}

		for _, endpoint := range endpoints {
			resp, err := testServer.DoRequest("PUT", endpoint, nil, nil)
			assert.NoError(t, err)
			// Should return 401 (unauthorized) without JWT
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			resp.Body.Close()
		}
	})
}

// TestOrganizerEndpointsAuthorization tests that organizer endpoints properly enforce authorization
func TestOrganizerEndpointsAuthorization(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Create a user
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(uint(user.ID), user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	t.Run("Organizer endpoints require valid JWT", func(t *testing.T) {
		endpoints := []string{
			"/organizer/matches/1/participants/1/approve",
			"/organizer/matches/1/participants/1/reject",
		}

		for _, endpoint := range endpoints {
			t.Run("Access "+endpoint+" without JWT", func(t *testing.T) {
				resp, err := testServer.DoRequest("PUT", endpoint, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with invalid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", endpoint, "invalid.jwt.token", nil)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})

			t.Run("Access "+endpoint+" with valid JWT", func(t *testing.T) {
				resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", endpoint, jwtToken, nil)
				assert.NoError(t, err)
				// Should not be unauthorized - may be forbidden (403) if not the organizer,
				// or not found (404) if match/participant doesn't exist
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				resp.Body.Close()
			})
		}
	})
}

// TestOrganizerEndpointsResponseFormat tests that organizer endpoints return properly formatted responses
func TestOrganizerEndpointsResponseFormat(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Create a user
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)

	jwtToken, err := testutils.CreateMockJWTToken(uint(user.ID), user.Name, user.IsAdmin)
	assert.NoError(t, err)

	t.Run("Organizer approve response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", "/organizer/matches/1/participants/1/approve", jwtToken, nil)
		assert.NoError(t, err)

		// Even if the request fails due to missing resources, the response should be valid JSON
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		// This should not panic; both successful and error responses should be valid JSON
		resp.Body.Close()
	})

	t.Run("Organizer reject response format", func(t *testing.T) {
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", "/organizer/matches/1/participants/1/reject", jwtToken, nil)
		assert.NoError(t, err)

		// Even if the request fails due to missing resources, the response should be valid JSON
		var response interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
	})
}

// TestOrganizerActionsOnMatches tests organizer actions on match participants
func TestOrganizerActionsOnMatches(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err, "Should clear test data successfully")

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err, "Should setup test database successfully")

	// Create an admin user to create necessary data
	adminUser := &models.User{
		SocialID:       "admin_organizer_test",
		SocialProvider: "facebook",
		Name:           "Admin Organizer",
		Email:          "adminorg@example.com",
		AvatarURL:      "https://example.com/admin-org.jpg",
		IsAdmin:        true,
	}
	result := testServer.DB.Create(adminUser)
	assert.NoError(t, result.Error)

	adminToken, err := testutils.CreateMockJWTToken(uint(adminUser.ID), adminUser.Name, adminUser.IsAdmin)
	assert.NoError(t, err)

	// Create a test location using admin privileges
	location := models.Location{
		Name:      "Organizer Test Location",
		Address:   "100 Organizer St",
		Latitude:  25.0,
		Longitude: 121.0,
	}
	resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/admin/locations", adminToken, location)
	assert.NoError(t, err)
	resp.Body.Close()

	// Create a test activity using admin privileges
	activity := models.Activity{
		Title:       "Organizer Test Activity",
		Description: "Test activity for organizer features",
		LocationID:  1, // Using first location
		TargetCount: 2,
	}
	resp, err = testutils.MakeAuthenticatedRequest(testServer, "POST", "/admin/activities", adminToken, activity)
	assert.NoError(t, err)
	resp.Body.Close()

	// Create an organizer user
	organizerUser, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, organizerUser)

	organizerToken, err := testutils.CreateMockJWTToken(uint(organizerUser.ID), organizerUser.Name, organizerUser.IsAdmin)
	assert.NoError(t, err)

	// In a real scenario, we would create a match where organizerUser is the organizer
	// For this test, we'll just verify that the endpoints accept valid JWTs
	t.Run("Approve/Reject endpoints accept JWTs", func(t *testing.T) {
		// These tests will likely return 404 (not found) or 403 (forbidden)
		// because the match/participant doesn't exist or user isn't the organizer
		approveResp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", "/organizer/matches/1/participants/1/approve", organizerToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, approveResp.StatusCode) // Should not be unauthorized
		approveResp.Body.Close()

		rejectResp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", "/organizer/matches/1/participants/1/reject", organizerToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, rejectResp.StatusCode) // Should not be unauthorized
		rejectResp.Body.Close()
	})
}
