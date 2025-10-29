package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of token to create
type TokenType string

const (
	UserToken  TokenType = "user"
	AdminToken TokenType = "admin"
)

// AuthTestHelper provides utilities for testing authentication functionality
type AuthTestHelper struct {
	Secret string
}

// NewAuthTestHelper creates a new AuthTestHelper with a default secret
func NewAuthTestHelper() *AuthTestHelper {
	return &AuthTestHelper{
		Secret: "test-secret-change-in-production",
	}
}

// NewAuthTestHelperWithSecret creates a new AuthTestHelper with a custom secret
func NewAuthTestHelperWithSecret(secret string) *AuthTestHelper {
	return &AuthTestHelper{
		Secret: secret,
	}
}

// CreateValidToken creates a valid JWT token for the specified type
func (a *AuthTestHelper) CreateValidToken(userID uint, email, name, provider string, tokenType TokenType) (string, error) {
	return CreateValidToken(userID, email, string(tokenType), a.Secret)
}

// CreateValidUserToken creates a valid JWT token for a test user
func (a *AuthTestHelper) CreateValidUserToken(userID uint, email, name, provider string) (string, error) {
	return a.CreateValidToken(userID, email, name, provider, UserToken)
}

// CreateValidAdminToken creates a valid JWT token for a test admin
func (a *AuthTestHelper) CreateValidAdminToken(userID uint, email, name, provider string) (string, error) {
	return a.CreateValidToken(userID, email, name, provider, AdminToken)
}

// CreateExpiredUserToken creates an expired JWT token for testing
func (a *AuthTestHelper) CreateExpiredUserToken(userID uint, email, name, provider string) (string, error) {
	return CreateExpiredToken(userID, email, "user", a.Secret)
}

// AddAuthHeader adds authentication header to a request
func (a *AuthTestHelper) AddAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// CreateTestAuthRequest creates an HTTP request with authentication token
func CreateTestAuthRequest(method, url string, body interface{}, token string) (*http.Request, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	// Add authorization header if token is provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

// MockAuthProvider simulates an OAuth provider for testing
type MockAuthProvider struct {
	ValidAuthCodes map[string]bool
	ValidTokens    map[string]MockUser
}

// MockUser represents a user in the mock OAuth provider
type MockUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Avatar   string `json:"avatar"`
}

// NewMockAuthProvider creates a new mock OAuth provider
func NewMockAuthProvider() *MockAuthProvider {
	return &MockAuthProvider{
		ValidAuthCodes: make(map[string]bool),
		ValidTokens:    make(map[string]MockUser),
	}
}

// AddValidAuthCode adds a valid authorization code to the mock
func (m *MockAuthProvider) AddValidAuthCode(code string, user MockUser) {
	m.ValidAuthCodes[code] = true
	m.ValidTokens[code] = user
}

// ValidateAuthCode validates an authorization code and returns user info
func (m *MockAuthProvider) ValidateAuthCode(code string) (MockUser, bool) {
	user, exists := m.ValidTokens[code]
	if !exists {
		return MockUser{}, false
	}

	// Mark the code as used so it can't be reused
	delete(m.ValidAuthCodes, code)

	return user, true
}

// SetupMockAuthRoutes sets up mock authentication routes for testing
func SetupMockAuthRoutes(router *gin.Engine, authHelper *AuthTestHelper, mockProvider *MockAuthProvider) {
	// Mock OAuth provider begin
	router.GET("/auth/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		if provider != "facebook" && provider != "instagram" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
			return
		}

		// Generate a mock auth code
		authCode := "mock-auth-code-" + provider
		mockProvider.AddValidAuthCode(authCode, MockUser{
			ID:       "123456",
			Email:    "test@example.com",
			Name:     "Test User",
			Provider: provider,
			Avatar:   "https://example.com/avatar.jpg",
		})

		// Redirect with auth code
		c.Redirect(http.StatusTemporaryRedirect, "/auth/"+provider+"/callback?code="+authCode)
	})

	// Mock OAuth provider callback
	router.GET("/auth/:provider/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}

		user, valid := mockProvider.ValidateAuthCode(code)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
			return
		}

		// In a real implementation, this would create a session or token
		// For testing, we'll just return user data
		c.JSON(http.StatusOK, gin.H{
			"message": "authenticated",
			"user":    user,
		})
	})

	// Mock token exchange
	router.GET("/auth/token", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"token": "mock-jwt-token",
			"user": gin.H{
				"id":       1,
				"email":    "test@example.com",
				"name":     "Test User",
				"provider": "facebook",
				"role":     "user",
			},
		})
	})
}

// CreateMockAuthServer creates a test server with mock authentication routes
func CreateMockAuthServer() (*httptest.Server, *AuthTestHelper, *MockAuthProvider) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	authHelper := NewAuthTestHelper()
	mockProvider := NewMockAuthProvider()

	SetupMockAuthRoutes(router, authHelper, mockProvider)

	server := httptest.NewServer(router)

	return server, authHelper, mockProvider
}

// ExtractUserIDFromToken extracts the user ID from a JWT token
func (a *AuthTestHelper) ExtractUserIDFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return uint(userIDFloat), nil
		}
	}

	return 0, nil
}

// ValidateTokenWithTimeout validates a token with a performance constraint
func (a *AuthTestHelper) ValidateTokenWithTimeout(tokenString string, timeout time.Duration) (jwt.MapClaims, error) {
	type result struct {
		claims jwt.MapClaims
		err    error
	}

	resultChan := make(chan result, 1)

	// Run validation in a goroutine to enforce timeout
	go func() {
		claims, err := ValidateToken(tokenString, a.Secret)
		resultChan <- result{claims, err}
	}()

	// Wait for result or timeout
	select {
	case res := <-resultChan:
		return res.claims, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("token validation exceeded timeout of %v", timeout)
	}
}
