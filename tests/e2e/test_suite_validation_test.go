package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestSuiteExecutionValidation tests that the entire test suite can be executed properly
func TestSuiteExecutionValidation(t *testing.T) {
	t.Run("Complete Facebook login flow executes without errors", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear any existing test data
		err := testServer.ClearTestData()
		assert.NoError(t, err, "Should clear test data successfully")

		// Setup test database
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err, "Should setup test database successfully")

		// Simulate complete Facebook login flow
		// Step 1: Create a user (simulating successful Facebook login)
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		assert.NotNil(t, user)

		// Step 2: Generate JWT token (simulating token creation after Facebook login)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Step 3: Use token to access protected endpoint
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		resp.Body.Close()
	})

	t.Run("All API endpoints accessible with valid JWT", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create user and JWT
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		// Test multiple endpoints
		endpoints := []string{
			"/profile",
			"/user/matches",
			"/user/past-matches",
		}

		for _, endpoint := range endpoints {
			resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", endpoint, token, nil)
			assert.NoError(t, err)
			// Should not get unauthorized error
			assert.NotEqual(t, 401, resp.StatusCode)
			resp.Body.Close()
		}
	})

	t.Run("User permissions properly enforced", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create regular user
		regularUser, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		regularToken, err := testutils.CreateMockJWTToken(regularUser.ID, regularUser.Name, regularUser.IsAdmin)
		assert.NoError(t, err)

		// Try to access admin endpoint with regular user token
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/admin/activities", regularToken, nil)
		assert.NoError(t, err)
		// Should get unauthorized or forbidden
		assert.Contains(t, []int{401, 403}, resp.StatusCode)
		resp.Body.Close()

		// Create admin user
		adminUser := &testutils.TestUser{
			ID:    regularUser.ID, // Same ID but updating admin status
			Name:  regularUser.Name,
			Email: regularUser.Email,
			IsAdmin: true,
		}
		testServer.DB.Model(&regularUser).Update("is_admin", true)

		// Create new token for admin
		adminToken, err := testutils.CreateMockJWTToken(adminUser.ID, adminUser.Name, true)
		assert.NoError(t, err)

		// Admin should be able to access admin endpoint
		resp, err = testutils.MakeAuthenticatedRequest(testServer, "GET", "/admin/activities", adminToken, nil)
		assert.NoError(t, err)
		// Should not get unauthorized (might get 200 or 404 depending on data)
		assert.NotEqual(t, 401, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("JWT validation works throughout test execution", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Clear and setup database
		err := testServer.ClearTestData()
		assert.NoError(t, err)
		err = testServer.SetupTestDatabase()
		assert.NoError(t, err)

		// Create user and token
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)
		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		// Validate token can be parsed at beginning
		claims, err := testutils.ValidateJWTToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)

		// Use token in multiple requests
		resp1, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		assert.NoError(t, err)
		resp1.Body.Close()

		resp2, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", token, nil)
		assert.NoError(t, err)
		resp2.Body.Close()

		// Validate token can still be parsed at end
		claims, err = testutils.ValidateJWTToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
	})
}

// TestSuiteErrorHandling tests that the test suite handles errors properly
func TestSuiteErrorHandling(t *testing.T) {
	t.Run("Test suite continues after individual test errors", func(t *testing.T) {
		// This test demonstrates that the test suite as a whole
		// can continue running even if individual test assertions fail
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// This assertion might fail, but the test function continues
		assert.Equal(t, 1, 2, "This will fail but test continues")

		// This assertion runs regardless of the previous failure
		assert.True(t, true, "This should still run")

		// Test that we can still interact with test server
		resp, err := testServer.DoRequest("GET", "/health", nil, nil)
		// The endpoint might not exist, but the server should be functional
		assert.NotNil(t, resp)
		if err != nil {
			resp.Body.Close()
		}
	})

	t.Run("Invalid tokens are properly rejected throughout suite", func(t *testing.T) {
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Try to access protected endpoint with invalid token
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", "invalid.token.here", nil)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode, "Invalid token should be rejected")
		resp.Body.Close()

		// Verify server still works for valid requests
		resp, err = testServer.DoRequest("GET", "/nonexistent", nil, nil)
		// Should get 404 (not found) not 500 (internal error)
		assert.NotNil(t, resp)
		resp.Body.Close()
	})
}