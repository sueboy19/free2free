package testutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
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

// IsTokenExpired checks if a JWT token is expired (does not validate signature)
func IsTokenExpired(tokenString string) (bool, error) {
	// Use ParseUnverified to parse token without signature validation
	// This allows us to check expiration without needing the secret key
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			return time.Now().Unix() > int64(exp), nil
		}
	}

	return false, fmt.Errorf("could not parse expiration time")
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

// TamperWithJWTToken modifies the payload of a JWT token for testing tampering detection
func TamperWithJWTToken(tokenString string, modifications map[string]interface{}) (string, error) {
	// JWT format: header.payload.signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode payload: %w", err)
	}

	// Parse the payload into a map
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Apply modifications
	for key, value := range modifications {
		claims[key] = value
	}

	// Marshal the modified claims
	modifiedPayload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal modified payload: %w", err)
	}

	// Encode the modified payload
	encodedPayload := base64.RawURLEncoding.EncodeToString(modifiedPayload)

	// Reassemble the token (keep original header and signature)
	tamperedToken := parts[0] + "." + encodedPayload + "." + parts[2]

	return tamperedToken, nil
}

// ValidateTokenSignature validates only the signature of a JWT token (without checking expiration)
func ValidateTokenSignature(tokenString, secret string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithoutClaimsValidation())

	// For signature validation specifically, return false if there's any error
	// This is to support the security test that expects false for invalid signatures
	if err != nil {
		return false, nil
	}

	// Check if the signature is valid
	return token.Valid, nil
}
