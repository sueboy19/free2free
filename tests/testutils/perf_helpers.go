package testutils

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// PerformanceTestConfig holds configuration for performance tests
type PerformanceTestConfig struct {
	Timeout              time.Duration
	MaxConcurrentReqs    int
	TokenValidationLimit time.Duration
}

// DefaultPerfConfig returns default performance test configuration
func DefaultPerfConfig() PerformanceTestConfig {
	return PerformanceTestConfig{
		Timeout:              500 * time.Millisecond, // 500ms API response time limit
		MaxConcurrentReqs:    10,                     // Max concurrent requests for load testing
		TokenValidationLimit: 10 * time.Millisecond,  // 10ms token validation limit
	}
}

// MeasureResponseTime measures the response time of a function
func MeasureResponseTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

// RunPerformanceTest executes a performance test with timeout validation
func RunPerformanceTest(t *testing.T, testName string, f func(), timeout time.Duration) {
	duration := MeasureResponseTime(f)
	assert.Less(t, duration, timeout, "%s took %v, expected less than %v", testName, duration, timeout)
}

// ConcurrentResult represents the result of a concurrent request
type ConcurrentResult struct {
	ReqID    int
	Duration time.Duration
}

// RunConcurrentPerformanceTest executes a performance test with concurrent requests
func RunConcurrentPerformanceTest(t *testing.T, testName string, numRequests int, requestFunc func(int), timeout time.Duration) {
	var wg sync.WaitGroup
	results := make(chan ConcurrentResult, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()
			duration := MeasureResponseTime(func() { requestFunc(reqID) })
			results <- ConcurrentResult{ReqID: reqID, Duration: duration}
		}(i)
	}

	wg.Wait()
	close(results)

	// Check that all requests completed within the timeout
	maxDuration := time.Duration(0)
	for result := range results {
		assert.Less(t, result.Duration, timeout, "%s (req %d) took %v, expected less than %v", testName, result.ReqID, result.Duration, timeout)
		if result.Duration > maxDuration {
			maxDuration = result.Duration
		}
	}

	t.Logf("%s: Max response time under load: %v", testName, maxDuration)
}

// ValidateTokenPerformance validates JWT token validation performance
func ValidateTokenPerformance(t *testing.T, token, secret string, limit time.Duration) {
	// Measure token validation time for multiple iterations
	start := time.Now()
	iterations := 100
	for i := 0; i < iterations; i++ {
		_, err := ValidateToken(token, secret)
		assert.NoError(t, err)
	}
	duration := time.Since(start)

	avgDuration := duration / time.Duration(iterations)
	assert.Less(t, avgDuration, limit,
		"Average JWT validation took %v, expected less than %v", avgDuration, limit)

	t.Logf("Average JWT validation time: %v", avgDuration)
}
