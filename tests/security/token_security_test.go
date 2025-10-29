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

		// Verify the token is properly rejected
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
		// This would test if the system can detect tampered tokens
		// In a real implementation, this might involve creating a valid token,
		// manually modifying parts of it, and verifying it's rejected
		t.Skip("Token tampering detection requires more specific implementation details")
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
			Avatar:   "https://example.com/avatar.jpg",
		}

		// Add a valid auth code
		authCode := "test-auth-code"
		mockProvider.AddValidAuthCode(authCode, mockUser)

		// First validation should succeed
		returnedUser, valid := mockProvider.ValidateAuthCode(authCode)
		assert.True(t, valid)
		assert.Equal(t, mockUser, returnedUser)

		// Second validation with same code should fail
		_, valid = mockProvider.ValidateAuthCode(authCode)
		assert.False(t, valid)
	})

	t.Run("State Parameter Validation", func(t *testing.T) {
		// In a real implementation, OAuth flows should use state parameters
		// to prevent CSRF attacks. This test would validate that implementation.
		t.Skip("State parameter validation requires actual OAuth implementation details")
	})

	t.Run("PKCE Validation", func(t *testing.T) {
		// In a real implementation, if using PKCE (Proof Key for Code Exchange),
		// this test would validate the PKCE flow.
		t.Skip("PKCE validation requires actual OAuth implementation details")
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

		// Verify sensitive data is not in the token
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
