package performance

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"free2free/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowPerformance tests the performance of the complete API workflow
func TestWorkflowPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	ts := testutils.NewTestServer()
	defer ts.Close()

	// Use default performance config
	config := testutils.DefaultPerfConfig()

	t.Run("Single Request Performance", func(t *testing.T) {
		// Create a valid JWT token for the test
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Test API call with performance validation
		testutils.RunPerformanceTest(t, "API Request", func() {
			w, err := testutils.GetRequest(ts.Router, "/health", token)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		}, config.Timeout)
	})

	t.Run("Concurrent Request Performance", func(t *testing.T) {
		// Create a valid JWT token for the test
		authHelper := testutils.NewAuthTestHelper()
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Execute concurrent requests with performance validation
		testutils.RunConcurrentPerformanceTest(t, "Concurrent API Requests", config.MaxConcurrentReqs,
			func(reqID int) {
				_, err := testutils.GetRequest(ts.Router, "/health", token)
				assert.NoError(t, err)
			},
			config.Timeout)
	})

	t.Run("JWT Token Validation Performance", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()

		// Create a valid token
		token, err := authHelper.CreateValidUserToken(1, "test@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Validate token validation performance using helper
		testutils.ValidateTokenPerformance(t, token, authHelper.Secret, config.TokenValidationLimit)
	})

	t.Run("Token Creation Performance", func(t *testing.T) {
		authHelper := testutils.NewAuthTestHelper()

		// Measure token creation time
		start := time.Now()
		iterations := 100
		for i := 0; i < iterations; i++ {
			_, err := testutils.CreateValidToken(1, "perf@example.com", "user", authHelper.Secret)
			assert.NoError(t, err)
		}
		duration := time.Since(start)

		avgDuration := duration / time.Duration(iterations)
		// Requirement: JWT token creation <20ms average
		tokenCreationLimit := 20 * time.Millisecond
		assert.Less(t, avgDuration, tokenCreationLimit,
			"Average JWT creation took %v, expected less than %v", avgDuration, tokenCreationLimit)

		t.Logf("Average JWT creation time: %v", avgDuration)
	})
}

// BenchmarkTokenValidation benchmarks the performance of token validation
func BenchmarkTokenValidation(b *testing.B) {
	authHelper := testutils.NewAuthTestHelper()

	// Create a valid token
	token, err := authHelper.CreateValidUserToken(1, "benchmark@example.com", "Benchmark User", "facebook")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testutils.ValidateToken(token, authHelper.Secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTokenCreation benchmarks the performance of token creation
func BenchmarkTokenCreation(b *testing.B) {
	authHelper := testutils.NewAuthTestHelper()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testutils.CreateValidToken(1, "benchmark@example.com", "user", authHelper.Secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// setupPerformanceTestRoutes configures routes for performance testing
func setupPerformanceTestRoutes(router *gin.Engine) {
	// Add a health check endpoint for performance testing
	router.GET("/health", func(c *gin.Context) {
		// Simulate token validation
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := authHeader[7:]
		_, err := testutils.ValidateToken(token, "test-secret-change-in-production")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().Unix()})
	})
}
