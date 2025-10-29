package unit

import (
	"strings"
	"testing"
	"time"

	"free2free/tests/testutils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestJWTTokenGeneration tests the JWT token generation functionality
func TestJWTTokenGeneration(t *testing.T) {
	secret := "test-secret"

	t.Run("Valid Token Creation", func(t *testing.T) {
		userID := uint(123)
		email := "test@example.com"
		role := "user"

		tokenString, err := testutils.CreateValidToken(userID, email, role, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Verify the token can be parsed and has correct claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)

		assert.Equal(t, float64(userID), claims["user_id"])
		assert.Equal(t, email, claims["email"])
		assert.Equal(t, role, claims["role"])

		// Check that expiration is in the future
		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.True(t, exp > float64(time.Now().Unix()))
	})

	t.Run("Expired Token Creation", func(t *testing.T) {
		userID := uint(123)
		email := "test@example.com"
		role := "user"

		tokenString, err := testutils.CreateExpiredToken(userID, email, role, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Verify the token is expired when validated
		// jwt.ParseWithOptions can handle validation errors differently
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		// The error would be returned if the token is expired
		// So we check both error and token validity
		if err != nil {
			// If there's an error, it might be a validation error (like expired)
			// Check if it's a validation error related to expiration
			if !strings.Contains(err.Error(), "token is expired") {
				assert.NoError(t, err) // This would fail if it's some other error
			}
		}
		
		// If we get a token back, it might be invalid
		if token != nil {
			assert.False(t, token.Valid)
		}

		// The token should be expired, so try parsing without verification to check claims
		tokenUnverified, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't validate signatures or claims for this check
			return []byte(secret), nil
		})
		
		if claims, ok := tokenUnverified.Claims.(jwt.MapClaims); ok {
			exp, ok := claims["exp"].(float64)
			assert.True(t, ok)
			assert.True(t, exp < float64(time.Now().Unix()))
		}
	})

	t.Run("Token Validation", func(t *testing.T) {
		userID := uint(456)
		email := "validate@example.com"
		role := "admin"

		tokenString, err := testutils.CreateValidToken(userID, email, role, secret)
		assert.NoError(t, err)

		claims, err := testutils.ValidateToken(tokenString, secret)
		assert.NoError(t, err)

		assert.Equal(t, float64(userID), claims["user_id"])
		assert.Equal(t, email, claims["email"])
		assert.Equal(t, role, claims["role"])
	})

	t.Run("Invalid Token Validation", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		_, err := testutils.ValidateToken(invalidToken, secret)
		assert.Error(t, err)
	})

	t.Run("Token with Wrong Secret", func(t *testing.T) {
		userID := uint(789)
		email := "wrong-secret@example.com"
		role := "user"

		tokenString, err := testutils.CreateValidToken(userID, email, role, secret)
		assert.NoError(t, err)

		// Try to validate with wrong secret
		_, err = testutils.ValidateToken(tokenString, "wrong-secret")
		assert.Error(t, err)
	})
}

// TestMockJWTValidator tests the mock JWT validator functionality
func TestMockJWTValidator(t *testing.T) {
	validator := testutils.NewMockJWTValidator()

	t.Run("Add and Validate Token", func(t *testing.T) {
		token := "test-token-123"
		userID := uint(123)

		validator.AddValidToken(token, userID)

		returnedUserID, err := validator.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, returnedUserID)
	})

	t.Run("Invalid Token Rejection", func(t *testing.T) {
		token := "invalid-token-456"

		validator.AddInvalidToken(token)

		_, err := validator.ValidateToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Unknown Token Rejection", func(t *testing.T) {
		token := "unknown-token-789"

		_, err := validator.ValidateToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})
}

// TestTokenExpiration tests token expiration behavior
func TestTokenExpiration(t *testing.T) {
	secret := "test-secret"

	t.Run("Token Expiration Validation", func(t *testing.T) {
		// Create an expired token
		tokenString, err := testutils.CreateExpiredToken(111, "expired@example.com", "user", secret)
		assert.NoError(t, err)

		// Try to validate the expired token
		_, err = testutils.ValidateToken(tokenString, secret)
		assert.Error(t, err)
	})

	t.Run("Valid Token Doesn't Expire Immediately", func(t *testing.T) {
		userID := uint(222)
		email := "valid@example.com"
		role := "user"

		tokenString, err := testutils.CreateValidToken(userID, email, role, secret)
		assert.NoError(t, err)

		// Immediately validate the token - should still be valid
		claims, err := testutils.ValidateToken(tokenString, secret)
		assert.NoError(t, err)

		assert.Equal(t, float64(userID), claims["user_id"])
		assert.Equal(t, email, claims["email"])
		assert.Equal(t, role, claims["role"])
	})
}

// BenchmarkTokenCreation benchmarks JWT token creation performance
func BenchmarkTokenCreation(b *testing.B) {
	secret := "benchmark-secret"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testutils.CreateValidToken(uint(i), "test@example.com", "user", secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTokenValidation benchmarks JWT token validation performance
func BenchmarkTokenValidation(b *testing.B) {
	secret := "benchmark-secret"
	tokenString, err := testutils.CreateValidToken(999, "benchmark@example.com", "user", secret)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testutils.ValidateToken(tokenString, secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}
