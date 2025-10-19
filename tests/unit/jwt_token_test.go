package unit

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestJWTTokenGeneration tests the JWT token generation functionality
func TestJWTTokenGeneration(t *testing.T) {
	t.Run("Create valid JWT token", func(t *testing.T) {
		userID := int64(12345)
		userName := "Test User"
		isAdmin := false

		tokenString, err := testutils.CreateMockJWTToken(userID, userName, isAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Verify the token has the expected structure (3 parts separated by '.')
		parts := splitToken(tokenString)
		assert.Equal(t, 3, len(parts), "JWT should have 3 parts: header.payload.signature")
	})

	t.Run("Generated token contains correct claims", func(t *testing.T) {
		userID := int64(67890)
		userName := "Another Test User"
		isAdmin := true

		tokenString, err := testutils.CreateMockJWTToken(userID, userName, isAdmin)
		assert.NoError(t, err)

		// Validate the token and check claims
		claims, err := testutils.ValidateJWTToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userName, claims.UserName)
		assert.Equal(t, isAdmin, claims.IsAdmin)
		assert.NotEmpty(t, claims.Issuer)
		assert.NotEmpty(t, claims.Subject)
	})

	t.Run("Generated token expires after set time", func(t *testing.T) {
		userID := int64(11111)
		userName := "Expiring User"
		isAdmin := false

		tokenString, err := testutils.CreateMockJWTToken(userID, userName, isAdmin)
		assert.NoError(t, err)

		// Validate the token is initially valid
		claims, err := testutils.ValidateJWTToken(tokenString)
		assert.NoError(t, err)
		assert.False(t, time.Now().After(claims.ExpiresAt.Time))

		// The token should expire in 15 minutes as set in the function
		assert.WithinDuration(t, time.Now().Add(15*time.Minute), claims.ExpiresAt.Time, 1*time.Minute)
	})

	t.Run("Invalid token fails validation", func(t *testing.T) {
		// Test with completely invalid token
		invalidToken := "this.is.not.a.valid.token"
		claims, err := testutils.ValidateJWTToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Expired token fails validation", func(t *testing.T) {
		// Create a custom token with past expiration time
		jwtSecret := "test-jwt-secret-key-32-chars-long-enough!!"

		claims := &testutils.JWTClaims{
			UserID:   99999,
			UserName: "Expired User",
			IsAdmin:  false,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Issuer:    "free2free-test",
				Subject:   "user:99999",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, err := token.SignedString([]byte(jwtSecret))
		assert.NoError(t, err)

		// This token should be expired
		isExpired, err := testutils.IsTokenExpired(expiredToken)
		assert.NoError(t, err)
		assert.True(t, isExpired)

		// Validation should fail
		_, err = testutils.ValidateJWTToken(expiredToken)
		assert.Error(t, err)
	})
}

// TestJWTTokenExtraction tests extracting information from JWT tokens
func TestJWTTokenExtraction(t *testing.T) {
	t.Run("Extract user ID from valid token", func(t *testing.T) {
		expectedUserID := int64(54321)
		userName := "Extract User"
		isAdmin := false

		tokenString, err := testutils.CreateMockJWTToken(expectedUserID, userName, isAdmin)
		assert.NoError(t, err)

		userID, err := testutils.GetUserIDFromToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("Extract user name from valid token", func(t *testing.T) {
		userID := int64(98765)
		expectedUserName := "Extract Name User"
		isAdmin := true

		tokenString, err := testutils.CreateMockJWTToken(userID, expectedUserName, isAdmin)
		assert.NoError(t, err)

		userName, err := testutils.GetUserNameFromToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, expectedUserName, userName)
	})

	t.Run("Check admin status from token", func(t *testing.T) {
		userID := int64(12345)
		userName := "Admin Check User"
		expectedIsAdmin := true

		tokenString, err := testutils.CreateMockJWTToken(userID, userName, expectedIsAdmin)
		assert.NoError(t, err)

		isAdmin, err := testutils.IsUserAdminFromToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, expectedIsAdmin, isAdmin)

		// Test with non-admin user
		nonAdminToken, err := testutils.CreateMockJWTToken(userID, userName, false)
		assert.NoError(t, err)

		isAdmin, err = testutils.IsUserAdminFromToken(nonAdminToken)
		assert.NoError(t, err)
		assert.False(t, isAdmin)
	})

	t.Run("Extraction fails with invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		// All extraction functions should fail with invalid token
		_, err := testutils.GetUserIDFromToken(invalidToken)
		assert.Error(t, err)

		_, err = testutils.GetUserNameFromToken(invalidToken)
		assert.Error(t, err)

		_, err = testutils.IsUserAdminFromToken(invalidToken)
		assert.Error(t, err)
	})
}

// TestJWTTokenSecurity tests security aspects of JWT tokens
func TestJWTTokenSecurity(t *testing.T) {
	t.Run("Token signature verification", func(t *testing.T) {
		userID := int64(111222)
		userName := "Secure User"
		isAdmin := false

		tokenString, err := testutils.CreateMockJWTToken(userID, userName, isAdmin)
		assert.NoError(t, err)

		// The token should validate with the correct secret
		claims, err := testutils.ValidateJWTToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)

		// Modifying any part of the token should cause validation to fail
		parts := splitToken(tokenString)
		modifiedToken := parts[0] + "." + parts[1] + ".totallydifferentandsignature"
		_, err = testutils.ValidateJWTToken(modifiedToken)
		assert.Error(t, err)
	})
}

// Helper function to split a JWT token
func splitToken(tokenString string) []string {
	var parts []string
	current := ""
	inPart := true

	for _, char := range tokenString {
		if char == '.' {
			parts = append(parts, current)
			current = ""
			inPart = false
		} else {
			if !inPart {
				inPart = true
			}
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	return parts
}