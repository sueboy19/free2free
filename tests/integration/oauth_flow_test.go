package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuthFlowCompletion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("OAuth begin endpoint works", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test OAuth begin for Facebook
		resp, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth begin", response["message"])
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("OAuth begin endpoint works for Instagram", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test OAuth begin for Instagram
		resp, err := ts.DoRequest("GET", "/auth/instagram", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth begin", response["message"])
		assert.Equal(t, "instagram", response["provider"])
	})

	t.Run("OAuth callback endpoint works", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test OAuth callback for Facebook
		resp, err := ts.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response structure
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth callback", response["message"])
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("Token exchange endpoint works", func(t *testing.T) {
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
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")
	})

	t.Run("Logout endpoint works", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test logout endpoint
		resp, err := ts.DoRequest("GET", "/logout", nil, nil)
		assert.NoError(t, err)

		// Should be a redirect (302 or 307)
		assert.Condition(t, func() bool {
			return resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound
		}, "Logout should return redirect status, got %d", resp.StatusCode)
	})

	t.Run("OAuth flow endpoints integration", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test that all OAuth-related endpoints are properly registered and don't return 404
		endpoints := []string{
			"/auth/facebook",
			"/auth/instagram", 
			"/auth/facebook/callback",
			"/auth/instagram/callback",
			"/auth/token",
			"/logout",
		}

		for _, endpoint := range endpoints {
			resp, err := ts.DoRequest("GET", endpoint, nil, nil)
			assert.NoError(t, err, "Endpoint %s should not return error", endpoint)
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "Endpoint %s should not return 404", endpoint)
		}
	})
}