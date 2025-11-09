package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuthEndpointsContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("OAuth begin endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test Facebook OAuth begin endpoint
		resp, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Should have expected fields in response
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "provider")
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("Instagram OAuth begin endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test Instagram OAuth begin endpoint
		resp, err := ts.DoRequest("GET", "/auth/instagram", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Should have expected fields in response
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "provider")
		assert.Equal(t, "instagram", response["provider"])
	})

	t.Run("OAuth callback endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test Facebook OAuth callback endpoint
		resp, err := ts.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Should have expected fields in response
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "provider")
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("Token exchange endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test token exchange endpoint
		resp, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Should have expected fields in response
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user") // token exchange should return user info
	})

	t.Run("Logout endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test logout endpoint
		resp, err := ts.DoRequest("GET", "/logout", nil, nil)
		assert.NoError(t, err)

		// Should be a redirect response
		assert.Condition(t, func() bool {
			return resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound
		}, "Logout should return redirect status")
	})

	t.Run("OAuth endpoints return consistent error format", func(t *testing.T) {
		// For endpoints that might have errors, check they return consistent format
		// This test is more about contract compliance

		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test an invalid provider (this might return an error depending on implementation)
		resp, err := ts.DoRequest("GET", "/auth/unknown_provider", nil, nil)
		assert.NoError(t, err)

		// Should not return 404 since the route exists with parameter
		// Could return 200, 400, or another status but not 404
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "OAuth routes should be parameterized and not return 404 for invalid providers")
	})

	t.Run("All OAuth endpoints are properly registered", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Verify that all OAuth-related endpoints are registered (not returning 404)
		oauthEndpoints := []struct {
			method   string
			endpoint string
		}{
			{"GET", "/auth/facebook"},
			{"GET", "/auth/instagram"},
			{"GET", "/auth/facebook/callback"},
			{"GET", "/auth/instagram/callback"},
			{"GET", "/auth/token"},
			{"GET", "/logout"},
		}

		for _, endpoint := range oauthEndpoints {
			resp, err := ts.DoRequest(endpoint.method, endpoint.endpoint, nil, nil)
			assert.NoError(t, err, "Endpoint %s %s should not error", endpoint.method, endpoint.endpoint)
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, 
				"OAuth endpoint %s should be registered, got %d", endpoint.endpoint, resp.StatusCode)
		}
	})
}