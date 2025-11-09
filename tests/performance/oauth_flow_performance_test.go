package performance

import (
	"net/http"
	"testing"
	"time"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuthFlowPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("OAuth begin endpoint response time", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Measure response time
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Check that response time is under 500ms as per requirements
		assert.Less(t, duration.Milliseconds(), int64(500), 
			"OAuth begin endpoint should respond in under 500ms, took %dms", duration.Milliseconds())
	})

	t.Run("OAuth callback endpoint response time", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Measure response time
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Check that response time is under 500ms as per requirements
		assert.Less(t, duration.Milliseconds(), int64(500), 
			"OAuth callback endpoint should respond in under 500ms, took %dms", duration.Milliseconds())
	})

	t.Run("Token exchange endpoint response time", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Measure response time
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Check that response time is under 500ms as per requirements
		assert.Less(t, duration.Milliseconds(), int64(500), 
			"Token exchange endpoint should respond in under 500ms, took %dms", duration.Milliseconds())
	})

	t.Run("Logout endpoint response time", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Measure response time
		start := time.Now()
		resp, err := ts.DoRequest("GET", "/logout", nil, nil)
		duration := time.Since(start)

		assert.NoError(t, err)
		// Logout returns redirect, so check for redirect status
		assert.Condition(t, func() bool {
			return resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound
		}, "Logout should return redirect status")
		
		// Check that response time is under 500ms as per requirements
		assert.Less(t, duration.Milliseconds(), int64(500), 
			"Logout endpoint should respond in under 500ms, took %dms", duration.Milliseconds())
	})

	t.Run("OAuth flow completion time under 10 seconds", func(t *testing.T) {
		// Create test server
		ts := testutils.NewTestServer()
		defer ts.Close()

		// Measure the time for a complete OAuth flow simulation
		start := time.Now()
		
		// Call the OAuth begin endpoint
		resp1, err := ts.DoRequest("GET", "/auth/facebook", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)
		
		// Call the OAuth callback endpoint
		resp2, err := ts.DoRequest("GET", "/auth/facebook/callback", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)
		
		// Call the token exchange endpoint
		resp3, err := ts.DoRequest("GET", "/auth/token", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp3.StatusCode)
		
		duration := time.Since(start)

		// The complete OAuth flow should complete in under 10 seconds as per requirements
		assert.Less(t, duration.Seconds(), float64(10), 
			"Complete OAuth flow should complete in under 10 seconds, took %f seconds", duration.Seconds())
	})
}