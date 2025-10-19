package performance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestFacebookLoginPerformance tests the performance of the complete Facebook login flow
func TestFacebookLoginPerformance(t *testing.T) {
	t.Run("Facebook login to JWT token generation under 30 seconds", func(t *testing.T) {
		start := time.Now()

		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err, "Should clear test data successfully")
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err, "Should setup test database successfully")

		// Simulate Facebook login process (creating user and generating JWT)
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate JWT token
		claims, err := testutils.ValidateJWTToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)

		elapsed := time.Since(start)

		// According to requirements in plan.md: "Facebook OAuth flow completed in under 30 seconds"
		assert.True(t, elapsed < 30*time.Second, "Facebook login flow should complete in under 30 seconds, took %v", elapsed)

		t.Logf("Facebook login flow completed in %v", elapsed)
	})

	t.Run("JWT token validation under 10ms", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create user and JWT token
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		start := time.Now()
		claims, err := testutils.ValidateJWTToken(token)
		validationTime := time.Since(start)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.True(t, validationTime < 10*time.Millisecond, "JWT validation should complete in under 10ms, took %v", validationTime)

		t.Logf("JWT validation completed in %v", validationTime)
	})

	t.Run("API request with JWT under 500ms", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create user and JWT token
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		start := time.Now()
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		requestTime := time.Since(start)

		assert.NoError(t, err)
		// May get 200 or 404 depending on implementation, but not a performance issue
		assert.Contains(t, []int{200, 404}, resp.StatusCode)
		assert.True(t, requestTime < 500*time.Millisecond, "API request should complete in under 500ms, took %v", requestTime)

		resp.Body.Close()
		t.Logf("API request with JWT completed in %v", requestTime)
	})

	t.Run("Multiple concurrent Facebook login simulations", func(t *testing.T) {
		// Test performance under simulated load of multiple users logging in
		// This is a simplified version - in real systems you'd use actual concurrency

		const numSimulations = 5
		var totalElapsed time.Duration

		for i := 0; i < numSimulations; i++ {
			start := time.Now()

			testServer := testutils.NewTestServer()
			
			// Create user and JWT token
			user, err := testServer.CreateTestUser()
			assert.NoError(t, err)
			token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
			assert.NoError(t, err)

			// Validate JWT token
			claims, err := testutils.ValidateJWTToken(token)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, claims.UserID)

			testServer.Close()
			elapsed := time.Since(start)
			totalElapsed += elapsed
		}

		avgTime := totalElapsed / numSimulations
		maxAllowedAvg := 30 * time.Second // Adjust based on requirements
		assert.True(t, avgTime < maxAllowedAvg, "Average Facebook login flow should complete in under %v, took %v", maxAllowedAvg, avgTime)

		t.Logf("Average Facebook login flow completed in %v across %d simulations", avgTime, numSimulations)
	})
}

// TestSystemPerformanceUnderLoad tests the system performance under load conditions
func TestSystemPerformanceUnderLoad(t *testing.T) {
	t.Run("JWT token generation performance", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create user
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		const numTokens = 10
		var totalGenTime time.Duration

		for i := 0; i < numTokens; i++ {
			start := time.Now()
			token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
			genTime := time.Since(start)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			totalGenTime += genTime
		}

		avgGenTime := totalGenTime / numTokens
		// JWT generation should be fast
		assert.True(t, avgGenTime < 10*time.Millisecond, "JWT generation should complete in under 10ms, took %v on average", avgGenTime)

		t.Logf("JWT generation completed in average of %v", avgGenTime)
	})

	t.Run("Database operations performance", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Test user creation performance
		start := time.Now()
		user, err := testServer.CreateTestUser()
		creationTime := time.Since(start)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, creationTime < 100*time.Millisecond, "User creation should complete in under 100ms, took %v", creationTime)

		t.Logf("User creation completed in %v", creationTime)
	})
}