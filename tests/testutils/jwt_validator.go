package testutils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// ValidateJWTToken validates a JWT token and returns the claims
func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Get JWT secret from config or environment
	jwtSecret := os.Getenv("TEST_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-jwt-secret-key-32-chars-long-enough!!"
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Validate the token itself
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract and return the claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// CreateMockJWTToken creates a mock JWT token for testing
func CreateMockJWTToken(userID int64, userName string, isAdmin bool) (string, error) {
	// Get JWT secret from config or environment
	jwtSecret := os.Getenv("TEST_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-jwt-secret-key-32-chars-long-enough!!"
	}

	// Create the claims
	claims := &JWTClaims{
		UserID:   userID,
		UserName: userName,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "free2free-test",
			Subject:   fmt.Sprintf("user:%d", userID),
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(tokenString string) (bool, error) {
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return true, err
	}

	// Check if the token is expired
	return time.Now().After(claims.ExpiresAt.Time), nil
}

// GetUserIDFromToken extracts the user ID from a JWT token
func GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// GetUserNameFromToken extracts the user name from a JWT token
func GetUserNameFromToken(tokenString string) (string, error) {
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserName, nil
}

// IsUserAdminFromToken checks if the user is an admin based on the token
func IsUserAdminFromToken(tokenString string) (bool, error) {
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return false, err
	}

	return claims.IsAdmin, nil
}