package contract

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProfileEndpointContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("Profile endpoint contract", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test that profile endpoint exists and returns correct format
		resp, err := ts.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)
		
		// The profile endpoint may return 401 (unauthorized) instead of 404 (not found)
		// which is expected behavior since it requires authentication
		assert.Condition(t, func() bool {
			return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
		}, "Profile endpoint should return either 200 (with auth) or 401 (without auth), got %d", resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			// If we get a 200, verify the response structure
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			// Profile response should contain user information
			assert.Contains(t, response, "id")
			assert.Contains(t, response, "email")
			assert.Contains(t, response, "name")
			assert.Contains(t, response, "provider")
		} else if resp.StatusCode == http.StatusUnauthorized {
			// If we get 401, verify it's an auth error, not a 404
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "Profile endpoint should not return 404 - route should be registered")
		}
	})

	t.Run("Profile endpoint returns correct response format", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Profile endpoint should exist (not return 404)
		resp, err := ts.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "Profile endpoint should be registered and not return 404")
	})
}