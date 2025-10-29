package testutils

import (
	"context"
	"fmt"
	"time"
)

// TestTimeoutConfig holds timeout configurations for tests
type TestTimeoutConfig struct {
	OverallTimeout time.Duration // Overall test suite timeout
	RequestTimeout time.Duration // Timeout for individual API requests
	DBTimeout      time.Duration // Timeout for database operations
	JWTTimeout     time.Duration // Timeout for JWT operations
}

// GetDefaultTestTimeoutConfig returns default timeout configuration
func GetDefaultTestTimeoutConfig() TestTimeoutConfig {
	return TestTimeoutConfig{
		OverallTimeout: 300 * time.Second, // 5 minutes for overall test suite
		RequestTimeout: 30 * time.Second,  // 30 seconds for API requests
		DBTimeout:      10 * time.Second,  // 10 seconds for DB operations
		JWTTimeout:     5 * time.Second,   // 5 seconds for JWT operations
	}
}

// TestContext provides a context with timeout for tests
type TestContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Config TestTimeoutConfig
}

// NewTestContext creates a new test context with timeout
func NewTestContext() *TestContext {
	config := GetDefaultTestTimeoutConfig()
	ctx, cancel := context.WithTimeout(context.Background(), config.OverallTimeout)

	return &TestContext{
		Ctx:    ctx,
		Cancel: cancel,
		Config: config,
	}
}

// WithRequestTimeout returns a context with request-specific timeout
func (tc *TestContext) WithRequestTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(tc.Ctx, tc.Config.RequestTimeout)
}

// WithDBTimeout returns a context with database-specific timeout
func (tc *TestContext) WithDBTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(tc.Ctx, tc.Config.DBTimeout)
}

// WithJWTTimeout returns a context with JWT-specific timeout
func (tc *TestContext) WithJWTTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(tc.Ctx, tc.Config.JWTTimeout)
}

// TestCleanup handles cleanup operations after tests
type TestCleanup struct {
	cleanupFuncs []func() error // Cleanup functions to run
}

// NewTestCleanup creates a new test cleanup handler
func NewTestCleanup() *TestCleanup {
	return &TestCleanup{
		cleanupFuncs: []func() error{},
	}
}

// AddCleanup adds a cleanup function to be executed later
func (tc *TestCleanup) AddCleanup(cleanupFunc func() error) {
	tc.cleanupFuncs = append(tc.cleanupFuncs, cleanupFunc)
}

// AddCleanupFunc adds a cleanup function with no return value
func (tc *TestCleanup) AddCleanupFunc(cleanupFunc func()) {
	tc.cleanupFuncs = append(tc.cleanupFuncs, func() error {
		cleanupFunc()
		return nil
	})
}

// Execute executes all cleanup functions
func (tc *TestCleanup) Execute() []error {
	var errors []error

	for i := len(tc.cleanupFuncs) - 1; i >= 0; i-- {
		if err := tc.cleanupFuncs[i](); err != nil {
			errors = append(errors, fmt.Errorf("cleanup function %d failed: %w", len(tc.cleanupFuncs)-i, err))
		}
	}

	return errors
}

// TestWithTimeout runs a test function with a timeout
func TestWithTimeout(timeout time.Duration, testFunc func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan error, 1)

	go func() {
		resultChan <- testFunc()
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("test timed out after %v", timeout)
	}
}

// TestOperationWithTimeout runs an operation with a specific timeout
func TestOperationWithTimeout(operationName string, timeout time.Duration, operation func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan error, 1)

	go func() {
		resultChan <- operation()
	}()

	select {
	case err := <-resultChan:
		if err != nil {
			return fmt.Errorf("%s failed: %w", operationName, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s timed out after %v", operationName, timeout)
	}
}

// CleanupTestDatabase clears test data from the database
func CleanupTestDatabase(db interface{}) func() error {
	return func() error {
		// In a real implementation, this would be more specific to the database type
		// For now, we'll return a function that just reports what it would do
		fmt.Println("Cleaning up test database...")
		// In actual implementation, this would clear test data
		return nil
	}
}

// CleanupTestServer stops the test server
func CleanupTestServer(server interface{}) func() error {
	return func() error {
		fmt.Println("Cleaning up test server...")
		// In actual implementation, this would properly close the test server
		return nil
	}
}

// CleanupMockOAuthProvider stops the mock OAuth provider
func CleanupMockOAuthProvider(provider interface{}) func() error {
	return func() error {
		fmt.Println("Cleaning up mock OAuth provider...")
		// In actual implementation, this would properly close the mock provider
		return nil
	}
}

// TestDeadline checks if the test deadline has been exceeded
func TestDeadline(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// RunTestsWithCleanup runs tests with proper cleanup
func RunTestsWithCleanup(cleanup *TestCleanup, testFunc func() error) error {
	// Run the test function
	err := testFunc()

	// Execute cleanup regardless of test result
	cleanupErrors := cleanup.Execute()

	// Report any cleanup errors
	for _, cleanupErr := range cleanupErrors {
		fmt.Printf("Cleanup error: %v\n", cleanupErr)
	}

	// Return the original test error if there was one
	if err != nil {
		return err
	}

	// If cleanup had errors but test passed, return cleanup errors
	if len(cleanupErrors) > 0 {
		return fmt.Errorf("test completed with %d cleanup error(s)", len(cleanupErrors))
	}

	return nil
}
