package security

import (
	"testing"

	"free2free/tests/testutils"
	"github.com/stretchr/testify/assert"
)

// TestJWTSecurity validates security aspects of JWT implementation
func TestJWTSecurity(t *testing.T) {
	authHelper := testutils.NewAuthTestHelper()

	t.Run("Token Expiration Validation", func(t *testing.T) {
		// Create an expired token
		expiredToken, err := authHelper.CreateExpiredUserToken(1, "user@example.com", "Test User", "facebook")
		assert.NoError(t, err)

		// Verify that token is properly rejected
		_, err = testutils.ValidateToken(expiredToken, authHelper.Secret)
		assert.Error(t, err)
	})

	t.Run("Token Signature Verification", func(t *testing.T) {
		// Create a valid token
		token, err := authHelper.CreateValidUserToken(2, "user2@example.com", "Test User 2", "facebook")
		assert.NoError(t, err)

		// Try to validate with wrong secret
		_, err = testutils.ValidateToken(token, "wrong-secret")
		assert.Error(t, err)
	})

	t.Run("Token Tampering Detection", func(t *testing.T) {
		// Create a valid token
		token, err := authHelper.CreateValidUserToken(5, "user5@example.com", "Test User 5", "facebook")
		assert.NoError(t, err)

		// Step 1: Verify that original token is valid
		_, err = testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err, "Original token should be valid")

		// Step 2: Verify original token signature is valid
		isValid, err := testutils.ValidateTokenSignature(token, authHelper.Secret)
		assert.NoError(t, err)
		assert.True(t, isValid, "Original token signature should be valid")

		// Step 3: Tamper with token (modify user_id in payload)
		tamperedToken, err := testutils.TamperWithJWTToken(token, map[string]interface{}{
			"user_id": 999, // Change to a different user ID
		})
		assert.NoError(t, err, "Should successfully tamper with token")

		// Step 4: Verify that original token is still valid (we didn't modify it)
		_, err = testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err, "Original token should still be valid")

		// Step 5: Verify that tampered token signature is invalid (signature mismatch)
		isValid, err = testutils.ValidateTokenSignature(tamperedToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.False(t, isValid, "Tampered token signature should be invalid")
	})

	t.Run("Role-based Access Validation", func(t *testing.T) {
		// Create tokens with different roles
		userToken, err := authHelper.CreateValidUserToken(3, "user3@example.com", "Regular User", "facebook")
		assert.NoError(t, err)

		adminToken, err := authHelper.CreateValidAdminToken(4, "admin@example.com", "Admin User", "facebook")
		assert.NoError(t, err)

		// Verify role extraction
		userClaims, err := testutils.ValidateToken(userToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.Equal(t, "user", userClaims["role"])

		adminClaims, err := testutils.ValidateToken(adminToken, authHelper.Secret)
		assert.NoError(t, err)
		assert.Equal(t, "admin", adminClaims["role"])
	})
}

// TestOAuthSecurity validates security aspects of OAuth implementation
func TestOAuthSecurity(t *testing.T) {
	// Set up mock OAuth provider
	mockProvider := testutils.NewMockAuthProvider()

	t.Run("Invalid OAuth Code Rejection", func(t *testing.T) {
		// Try to validate an invalid OAuth code
		_, valid := mockProvider.ValidateAuthCode("invalid-code")
		assert.False(t, valid)
	})

	t.Run("OAuth Code Reuse Prevention", func(t *testing.T) {
		// Create a mock user
		mockUser := testutils.MockUser{
			ID:       "123456",
			Email:    "test@example.com",
			Name:     "Test User",
			Provider: "facebook",
		}

		// Add a valid auth code
		authCode := "test-auth-code"
		mockProvider.AddValidAuthCode(authCode, mockUser)

		// First validation should succeed
		returnedUser, valid := mockProvider.ValidateAuthCode(authCode)
		assert.True(t, valid)
		assert.Equal(t, mockUser, returnedUser)

		// MockAuthProvider should already delete code on first validation
		// Verify it's no longer in ValidAuthCodes
		_, exists := mockProvider.ValidAuthCodes[authCode]
		assert.False(t, exists, "Auth code should be removed after first validation")
	})

	t.Run("State Parameter Validation", func(t *testing.T) {
		mockProvider := testutils.NewMockAuthProvider()

		// Step 1: Generate and save a state parameter
		state := mockProvider.StateManager.GenerateState()
		assert.NotEmpty(t, state, "State should be generated")
		assert.True(t, mockProvider.StateManager.ValidStates[state], "State should be saved in valid states")

		// Step 2: Validate correct state (should succeed)
		isValid := mockProvider.ValidateState(state)
		assert.True(t, isValid, "Valid state should be accepted")
		_, exists := mockProvider.StateManager.ValidStates[state]
		assert.False(t, exists, "State should be consumed after validation (one-time use)")

		// Step 3: Generate another state for CSRF protection test
		state2 := mockProvider.StateManager.GenerateState()
		assert.NotEqual(t, state, state2, "Each state should be unique")

		// Step 4: Validate second state (should succeed on first use)
		isValid = mockProvider.ValidateState(state2)
		assert.True(t, isValid, "Second state should be accepted on first use")

		// Step 5: Try to reuse second state (should fail - CSRF protection)
		isValid = mockProvider.ValidateState(state2)
		assert.False(t, isValid, "State should not be reusable (prevents CSRF attacks)")

		// Step 6: Try to validate an invalid state (should fail)
		invalidState := "invalid-state-12345"
		isValid = mockProvider.ValidateState(invalidState)
		assert.False(t, isValid, "Invalid state should be rejected")
	})

	t.Run("PKCE Validation", func(t *testing.T) {
		// Set up mock OAuth provider with PKCE
		mockProvider := testutils.NewMockAuthProvider()

		// Step 1: Generate code verifier and challenge
		verifier, challenge, err := mockProvider.PKCEManager.GenerateCodeChallenge()
		assert.NoError(t, err, "Should generate code challenge successfully")
		assert.NotEmpty(t, verifier, "Code verifier should not be empty")
		assert.NotEmpty(t, challenge, "Code challenge should not be empty")
		assert.NotEqual(t, verifier, challenge, "Verifier and challenge should be different")

		// Step 2: Add a valid auth code with PKCE verifier
		authCode := "pkce-auth-code-123"
		mockProvider.PKCEManager.AddValidAuthCodeWithPKCE(authCode, verifier)
		assert.NotEmpty(t, mockProvider.PKCEManager.Codes[authCode], "Auth code with verifier should be stored")

		// Step 3: Validate correct verifier (should succeed)
		isValid := mockProvider.PKCEManager.ValidateCodeVerifier(verifier)
		assert.True(t, isValid, "Valid code verifier should be accepted")

		// Step 4: Verify that verifier was stored before validation
		_, exists := mockProvider.PKCEManager.Codes[authCode]
		assert.True(t, exists, "Code verifier should exist in storage before validation")

		// Step 5: Verify that verifier is now consumed (one-time use)
		_, exists = mockProvider.PKCEManager.Codes[authCode]
		assert.False(t, exists, "Code verifier should be consumed after validation")

		// Step 6: Generate another verifier for replay attack test
		verifier2, _, err := mockProvider.PKCEManager.GenerateCodeChallenge()
		assert.NoError(t, err)

		authCode2 := "pkce-auth-code-456"
		mockProvider.PKCEManager.AddValidAuthCodeWithPKCE(authCode2, verifier2)

		// Step 7: Try to validate an incorrect verifier (should fail)
		invalidVerifier := "invalid-verifier-xyz"
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(invalidVerifier)
		assert.False(t, isValid, "Invalid code verifier should be rejected")

		// Step 8: Verify correct verifier2 is still valid (not yet used)
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(verifier2)
		assert.True(t, isValid, "Second valid verifier should be accepted")

		// Step 9: Consume verifier2 and verify it can't be reused
		_ = mockProvider.PKCEManager.ValidateCodeVerifier(verifier2)
		_, exists = mockProvider.PKCEManager.Codes[authCode2]
		assert.False(t, exists, "Code verifier2 should be consumed after use")

		// Step 10: Try to validate verifier2 again (should fail - one-time use)
		isValid = mockProvider.PKCEManager.ValidateCodeVerifier(verifier2)
		assert.False(t, isValid, "Code verifier2 should not be reusable")
	})
}

// TestTokenLeakagePrevention ensures sensitive information is not leaked
func TestTokenLeakagePrevention(t *testing.T) {
	authHelper := testutils.NewAuthTestHelper()

	t.Run("Sensitive Data Not in Tokens", func(t *testing.T) {
		token, err := authHelper.CreateValidUserToken(5, "user5@example.com", "Test User 5", "facebook")
		assert.NoError(t, err)

		claims, err := testutils.ValidateToken(token, authHelper.Secret)
		assert.NoError(t, err)

		// Verify sensitive data is not in token
		assert.NotContains(t, claims, "password")
		assert.NotContains(t, claims, "credit_card")
		assert.NotContains(t, claims, "secret_key")

		// Verify expected data is present
		assert.Contains(t, claims, "user_id")
		assert.Contains(t, claims, "email")
		assert.Contains(t, claims, "role")
		assert.Contains(t, claims, "exp")
		assert.Contains(t, claims, "iat")
	})

	t.Run("Error Messages Don't Reveal Sensitive Info", func(t *testing.T) {
		// Create a token with an invalid secret
		_, err := testutils.ValidateToken("invalid.token.format", "some-secret")
		assert.Error(t, err)

		// Verify error message doesn't reveal internal details
		errMsg := err.Error()
		assert.NotContains(t, errMsg, "secret")
		assert.NotContains(t, errMsg, "internal")
		assert.NotContains(t, errMsg, "database")
	})
}
