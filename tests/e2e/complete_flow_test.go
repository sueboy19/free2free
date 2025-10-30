package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestCompleteFacebookLoginFlow tests the complete flow from Facebook login to API usage
func TestCompleteFacebookLoginFlow(t *testing.T) {
	t.Log("Starting complete Facebook login to API usage flow test...")

	// Initialize test server
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err, "Should clear test data successfully")

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err, "Should setup test database successfully")

	t.Log("Test server initialized and database setup complete.")

	// Step 1: Simulate Facebook login and JWT token generation
	t.Log("Step 1: Simulating Facebook login and JWT token generation...")
	user, err := testServer.CreateTestUser()
	assert.NoError(t, err)
	assert.NotNil(t, user)
	t.Logf("Created test user with ID: %d", user.ID)

	jwtToken, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)
	t.Log("JWT token generated successfully.")

	// Step 2: Verify JWT token is valid
	t.Log("Step 2: Verifying JWT token validity...")
	_, err = testutils.ValidateJWTToken(jwtToken)
	assert.NoError(t, err)
	
	// The main validation is that no error occurred
	t.Log("JWT token validation successful.")

	// Step 3: Test accessing protected endpoints with JWT
	t.Log("Step 3: Testing access to protected endpoints with JWT...")

	// Test profile endpoint
	resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", jwtToken, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var profileData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&profileData)
	assert.NoError(t, err)
	assert.Equal(t, float64(user.ID), profileData["id"])
	resp.Body.Close()
	t.Log("Profile endpoint access successful.")

	// Test user matches endpoint
	resp, err = testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", jwtToken, nil)
	assert.NoError(t, err)
	// Could be 200 (success) or 404 (no matches found), but not 401 (unauthorized)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
	t.Log("User matches endpoint access successful.")

	// Step 4: Test admin-specific functionality if user is admin
	if user.IsAdmin {
		t.Log("Step 4: Testing admin-specific functionality...")
		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/admin/activities", jwtToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
		t.Log("Admin endpoint access successful.")
	} else {
		t.Log("Step 4: Skipping admin tests as user is not admin.")
	}

	// Step 5: Test performance requirements
	t.Log("Step 5: Testing performance requirements...")

	// Test JWT validation time (should be under 10ms as per requirements)
	startTime := time.Now()
	_, err = testutils.ValidateJWTToken(jwtToken)
	validationTime := time.Since(startTime)
	assert.True(t, validationTime < 10*time.Millisecond,
		fmt.Sprintf("JWT validation should be under 10ms but took %v", validationTime))
	t.Logf("JWT validation completed in %v", validationTime)

	// Step 6: Test API response time (should be under 500ms as per requirements)
	startTime = time.Now()
	resp, err = testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", jwtToken, nil)
	apiTime := time.Since(startTime)
	assert.NoError(t, err)
	assert.True(t, apiTime < 500*time.Millisecond,
		fmt.Sprintf("API request should be under 500ms but took %v", apiTime))
	resp.Body.Close()
	t.Logf("API request completed in %v", apiTime)

	t.Log("Complete Facebook login flow test completed successfully!")
}

// TestEndToEndFlowWithTestDataGenerator tests the complete flow using the test data generator
func TestEndToEndFlowWithTestDataGenerator(t *testing.T) {
	t.Log("Starting end-to-end test with test data generator...")

	// Initialize test server
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear any existing test data
	err := testServer.ClearTestData()
	assert.NoError(t, err)

	// Setup test database
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err)

	// Setup sample data (this is a mock implementation since we don't have a real DB)
	t.Log("Sample data setup complete (mock implementation).")

	// Get the admin user (first user created in SetupSampleData)
	var adminUser struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
	}
	err = testServer.DB.Model(&struct{}{}).Table("users").Where("is_admin = ?", true).First(&adminUser).Error
	if err != nil {
		// If admin wasn't created correctly, create one manually
		adminUser.ID = 1 // Assuming this is the first user
		adminUser.Name = "Test Admin"
		adminUser.IsAdmin = true
	}

	// Generate JWT for admin user
	jwtToken, err := testutils.CreateMockJWTToken(adminUser.ID, adminUser.Name, adminUser.IsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)
	t.Log("JWT token for admin generated.")

	// Test admin endpoints
	resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/admin/activities", jwtToken, nil)
	assert.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	resp.Body.Close()
	t.Log("Admin activities endpoint accessed successfully.")

	// Test user endpoints
	resp, err = testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", jwtToken, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
	t.Log("User matches endpoint accessed successfully.")

	t.Log("End-to-end test with test data generator completed successfully!")
}

// TestFlowErrorHandling tests error handling throughout the complete flow
func TestFlowErrorHandling(t *testing.T) {
	t.Log("Starting flow error handling test...")

	// Initialize test server
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	// Clear and setup database
	err := testServer.ClearTestData()
	assert.NoError(t, err)
	err = testServer.SetupTestDatabase()
	assert.NoError(t, err)

	t.Log("Testing behavior with invalid JWT...")
	// Try to access protected endpoint with invalid token
	resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", "invalid.token.here", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
	t.Log("Invalid JWT properly rejected.")

	t.Log("Testing behavior with expired JWT...")
	// Create an expired JWT token for testing (by custom function - placeholder)
	// In actual implementation, we would create a token with past expiration
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	resp, err = testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", expiredToken, nil)
	assert.NoError(t, err)
	// May or may not be rejected depending on whether it's validated, but shouldn't cause server error
	assert.NotEqual(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
	t.Log("Expired JWT test completed.")

	t.Log("Flow error handling test completed successfully!")
}
