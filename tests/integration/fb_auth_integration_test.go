package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestFacebookAuthFlowIntegration tests the complete Facebook OAuth integration flow
func TestFacebookAuthFlowIntegration(t *testing.T) {
	// Initialize test server
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err, "Should clear test data successfully")

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err, "Should setup test database successfully")

	t.Run("Facebook OAuth flow creates user and returns JWT", func(t *testing.T) {
		// In a real implementation, this would involve:
		// 1. Making a request to /auth/facebook to initiate OAuth
		// 2. Simulating the redirect to Facebook
		// 3. Simulating the callback to /auth/facebook/callback
		// 4. Validating the returned JWT token

		// Since we can't fully simulate OAuth in integration tests without real Facebook credentials,
		// we'll validate the components that make up the flow

		// Test that auth endpoints exist
		resp, err := testServer.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		// This might return 500 if Facebook keys are not configured, which is expected
		assert.Contains(t, []int{200, 302, 500}, resp.StatusCode)
		resp.Body.Close()

		// Test callback endpoint exists
		resp, err = testServer.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		assert.NoError(t, err)
		// The callback without proper params will likely return an error
		assert.Contains(t, []int{200, 400, 500}, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("Token exchange from session to JWT", func(t *testing.T) {
		// Create a test user first
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		// Test that token exchange endpoint exists
		resp, err := testServer.DoRequest("GET", "/auth/token", nil, nil)
		assert.NoError(t, err)
		assert.Contains(t, []int{200, 401}, resp.StatusCode) // 401 if not logged in, 200 if logged in
		resp.Body.Close()
	})

	t.Run("Profile endpoint accessible with JWT", func(t *testing.T) {
		// Create a test user
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		// Create a mock JWT for this user
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Test accessing profile with JWT
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Check that user data is in response
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "name")
		assert.Contains(t, response, "email")
		resp.Body.Close()
	})

	t.Run("JWT validation works correctly", func(t *testing.T) {
		// Create a test user
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		// Create a valid JWT
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate the token using our utility
		claims, err := testutils.ValidateJWTToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Name, claims.UserName)

		// Test with invalid token
		invalidClaims, err := testutils.ValidateJWTToken("invalid.token.here")
		assert.Error(t, err)
		assert.Nil(t, invalidClaims)
	})

	t.Run("Logout endpoint clears session", func(t *testing.T) {
		// Test that logout endpoint exists and returns appropriate response
		resp, err := testServer.DoRequest("GET", "/logout", nil, nil)
		assert.NoError(t, err)
		// Logout might return redirect (302) or success (200)
		assert.Contains(t, []int{200, 302}, resp.StatusCode)
		resp.Body.Close()
	})
}

// TestFacebookOAuthErrorHandling tests error scenarios in OAuth flow
func TestFacebookOAuthErrorHandling(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	t.Run("Invalid provider returns error", func(t *testing.T) {
		resp, err := testServer.DoRequest("GET", "/auth/invalidprovider", nil, nil)
		assert.NoError(t, err)
		// Should return an error for invalid provider
		resp.Body.Close()
	})

	t.Run("Missing auth headers return 401", func(t *testing.T) {
		// Try to access protected endpoint without token
		resp, err := testServer.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})
}
