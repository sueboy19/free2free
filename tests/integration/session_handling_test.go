package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSessionHandlingInAuthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("Authentication endpoint with proper session", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test OAuth begin endpoint
		resp, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth begin", response["message"])
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("OAuth callback endpoint", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test OAuth callback endpoint
		resp, err := ts.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OAuth callback", response["message"])
		assert.Equal(t, "facebook", response["provider"])
	})

	t.Run("Logout endpoint", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test logout endpoint
		resp, err := ts.DoRequest("GET", "/logout", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "logged out", response["message"])
	})

	t.Run("Token exchange endpoint", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test token exchange endpoint
		resp, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")
	})

	t.Run("Profile endpoint should return 200 not 404", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Test profile endpoint (should now return 200 instead of 404)
		resp, err := ts.DoRequest("GET", "/profile", nil, nil)
		assert.NoError(t, err)
		// This should now return 200 instead of 404
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "email")
		assert.Contains(t, response, "name")
	})
}