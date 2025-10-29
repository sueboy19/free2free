package testutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateValidToken creates a valid JWT token for testing purposes
func CreateValidToken(userID uint, email, role string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		"iat":     time.Now().Unix(),                     // Issued at time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateExpiredToken creates an expired JWT token for testing purposes
func CreateExpiredToken(userID uint, email, role string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(-time.Hour).Unix(),     // Token expired 1 hour ago
		"iat":     time.Now().Add(-time.Hour * 2).Unix(), // Issued 2 hours ago
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateAndExtractClaims parses and validates a JWT token, extracting common claims
func ValidateAndExtractClaims(tokenString, secret string) (userID uint, email, role string, err error) {
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return 0, "", "", err
	}

	// Extract user ID
	if idFloat, ok := claims["user_id"].(float64); ok {
		userID = uint(idFloat)
	} else {
		return 0, "", "", fmt.Errorf("user_id not found in token")
	}

	// Extract email
	if emailStr, ok := claims["email"].(string); ok {
		email = emailStr
	} else {
		return 0, "", "", fmt.Errorf("email not found in token")
	}

	// Extract role
	if roleStr, ok := claims["role"].(string); ok {
		role = roleStr
	} else {
		return 0, "", "", fmt.Errorf("role not found in token")
	}

	return userID, email, role, nil
}

// ValidateTokenAndCheckRole validates a JWT token and checks user role
func ValidateTokenAndCheckRole(tokenString, secret, requiredRole string) (uint, error) {
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return 0, err
	}

	// Extract role
	if role, ok := claims["role"].(string); ok {
		if role != requiredRole {
			return 0, fmt.Errorf("insufficient permissions: required %s, got %s", requiredRole, role)
		}
	} else {
		return 0, fmt.Errorf("role not found in token")
	}

	// Extract user ID
	if idFloat, ok := claims["user_id"].(float64); ok {
		return uint(idFloat), nil
	}

	return 0, fmt.Errorf("user_id not found in token")
}

// MockJWTValidator provides a mock implementation for testing
type MockJWTValidator struct {
	ValidTokens   map[string]bool
	InvalidTokens map[string]bool
	TokenUserMap  map[string]uint
}

// NewMockJWTValidator creates a new mock JWT validator
func NewMockJWTValidator() *MockJWTValidator {
	return &MockJWTValidator{
		ValidTokens:   make(map[string]bool),
		InvalidTokens: make(map[string]bool),
		TokenUserMap:  make(map[string]uint),
	}
}

// ValidateToken checks if the token is valid in the mock
func (m *MockJWTValidator) ValidateToken(tokenString string) (uint, error) {
	if m.InvalidTokens[tokenString] {
		return 0, fmt.Errorf("invalid token")
	}

	if userID, exists := m.TokenUserMap[tokenString]; exists {
		return userID, nil
	}

	return 0, fmt.Errorf("token not found")
}

// AddValidToken adds a valid token to the mock
func (m *MockJWTValidator) AddValidToken(token string, userID uint) {
	m.ValidTokens[token] = true
	m.TokenUserMap[token] = userID
}

// AddInvalidToken adds an invalid token to the mock
func (m *MockJWTValidator) AddInvalidToken(token string) {
	m.InvalidTokens[token] = true
}
