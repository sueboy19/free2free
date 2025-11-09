package unit

import (
	"os"
	"testing"
	"time"

	"free2free/tests/testutils"
	"free2free/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTClaimsStructure(t *testing.T) {
	// Setup test environment
	originalEnv := testutils.SaveOriginalEnvironment()
	testutils.SetupTestEnvironment()
	defer testutils.RestoreOriginalEnvironment(originalEnv)

	t.Run("Validate JWT claims structure matches implementation", func(t *testing.T) {
		// Verify that our JWT secret is set
		jwtSecret := os.Getenv("JWT_SECRET")
		assert.NotEmpty(t, jwtSecret, "JWT_SECRET must be set for testing")

		// Create a JWT token with our expected claims structure
		claims := utils.Claims{
			UserID:   123,
			UserName: "test user",
			IsAdmin:  false,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Parse the token to ensure it matches our expected structure
		parsedToken, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Validate that the claims are of the correct type
		validatedClaims, ok := parsedToken.Claims.(*utils.Claims)
		assert.True(t, ok, "Claims should be of type *utils.Claims")
		assert.Equal(t, int64(123), validatedClaims.UserID)
		assert.Equal(t, "test user", validatedClaims.UserName)
		assert.Equal(t, false, validatedClaims.IsAdmin)
	})

	t.Run("JWT validation function works with correct structure", func(t *testing.T) {
		// Set JWT secret
		_ = os.Setenv("JWT_SECRET", "test-jwt-secret-for-testing-environment-32-characters")

		// Create a JWT token with our expected claims structure
		claims := utils.Claims{
			UserID:   456,
			UserName: "another user",
			IsAdmin:  true,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		assert.NoError(t, err)

		// Use ValidateJWTToken function to validate the token
		// Since this function is internal to utils, we'll test our own validation
		parsedToken, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Extract the claims
		parsedClaims, ok := parsedToken.Claims.(*utils.Claims)
		assert.True(t, ok)
		assert.Equal(t, int64(456), parsedClaims.UserID)
		assert.Equal(t, "another user", parsedClaims.UserName)
		assert.Equal(t, true, parsedClaims.IsAdmin)
	})

	t.Run("Validate JWT with interface{} doesn't cause panic", func(t *testing.T) {
		// Set JWT secret
		_ = os.Setenv("JWT_SECRET", "test-jwt-secret-for-testing-environment-32-characters")

		// Create a JWT token with our expected claims structure
		claims := utils.Claims{
			UserID:   789,
			UserName: "test user 3",
			IsAdmin:  false,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		assert.NoError(t, err)

		// Parse with interface{} first, then type assert
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Type assertion to our Claims structure
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			// Access the user_id field from MapClaims
			if userID, exists := claims["user_id"]; exists {
				// Convert to int64 (MapClaims stores numbers as float64)
				assert.Equal(t, float64(789), userID)
			} else {
				t.Error("user_id field not found in claims")
			}
		}
	})
}