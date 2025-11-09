package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProfileEndpointAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("Profile endpoint returns correct status codes", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Profile endpoint should not return 404 anymore
		resp, err := ts.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)

		// Should return either 200 (if we provide auth) or 401 (if no auth), but not 404
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "Profile endpoint should be registered and not return 404")
		
		// It's acceptable for profile to return 401 Unauthorized when no auth is provided
		assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized}, resp.StatusCode)
	})

	t.Run("Profile endpoint with proper authentication", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Create a request with mock authentication header
		headers := map[string]string{
			"Authorization": "Bearer mock-token",
		}

		resp, err := ts.DoRequest("GET", "/profile", headers, nil)
		assert.NoError(t, err)

		// With auth header, should return 200 or properly structured 401/403
		// (exact behavior depends on how the token validation works)
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "Profile endpoint should not return 404")
		
		// If it's a 200, check the response structure
		if resp.StatusCode == http.StatusOK {
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			
			// Response should contain user profile fields
			assert.Contains(t, response, "id")
			assert.Contains(t, response, "email")
			assert.Contains(t, response, "name")
		}
	})

	t.Run("Profile endpoint integration with session middleware", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test that profile endpoint integrates properly with our middleware
		// and doesn't cause runtime panics as it did before our fixes
		resp, err := ts.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)

		// Should not crash or panic
		assert.Condition(t, func() bool {
			return resp.StatusCode != http.StatusNotFound && 
				   resp.StatusCode != 500 // Internal server error indicating a panic
		}, "Profile endpoint should not return 404 or 500 error")
	})
}