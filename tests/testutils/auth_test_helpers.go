package testutils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	StateManager   *OAuthStateManager
	PKCEManager    *PKCEManager
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
		StateManager:   NewOAuthStateManager(),
		PKCEManager:    NewPKCEManager(),
	}
}

// AddValidAuthCode adds a valid authorization code to the mock
func (m *MockAuthProvider) AddValidAuthCode(code string, user MockUser) {
	m.ValidAuthCodes[code] = true
	m.ValidTokens[code] = user
}

// AddValidAuthCodeWithStateAndPKCE adds a valid auth code with state and PKCE
func (m *MockAuthProvider) AddValidAuthCodeWithStateAndPKCE(code string, user MockUser, state string, verifier string) {
	m.ValidAuthCodes[code] = true
	m.ValidTokens[code] = user
	m.StateManager.ValidStates[state] = true
	m.PKCEManager.Codes[code] = verifier
}

// ValidateState validates an OAuth state parameter
func (m *MockAuthProvider) ValidateState(state string) bool {
	return m.StateManager.ValidateState(state)
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

// CreateTokenWithClaims creates a JWT token with the specified claims
func CreateTokenWithClaims(claims map[string]interface{}, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		// If there's an error, return an empty string
		return ""
	}

	return tokenString
}

// CreateTokenWithCustomClaims creates a JWT token with custom claims
func CreateTokenWithCustomClaims(userID uint, email, name, role, secret string, customClaims map[string]interface{}) string {
	// Create base claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Add custom claims
	for key, value := range customClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		// If there's an error, return an empty string
		return ""
	}

	return tokenString
}

// CreateTestUser creates a test user for testing purposes
func CreateTestUser() *TestUser {
	return &TestUser{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "facebook",
		Role:     "user",
		IsAdmin:  false,
	}
}

// OAuthStateManager manages OAuth state parameters for CSRF protection
type OAuthStateManager struct {
	ValidStates map[string]bool
}

// NewOAuthStateManager creates a new OAuth state manager
func NewOAuthStateManager() *OAuthStateManager {
	return &OAuthStateManager{
		ValidStates: make(map[string]bool),
	}
}

// GenerateState generates a random state parameter for OAuth flow
func (m *OAuthStateManager) GenerateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	m.ValidStates[state] = true
	return state
}

// ValidateState validates an OAuth state parameter and returns true if valid
func (m *OAuthStateManager) ValidateState(state string) bool {
	valid, exists := m.ValidStates[state]
	if exists {
		delete(m.ValidStates, state) // Mark state as used (one-time use)
	}
	return valid
}

// PKCEManager manages PKCE (Proof Key for Code Exchange) for OAuth
type PKCEManager struct {
	Codes map[string]string
}

// NewPKCEManager creates a new PKCE manager
func NewPKCEManager() *PKCEManager {
	return &PKCEManager{
		Codes: make(map[string]string),
	}
}

// GenerateCodeChallenge generates a code verifier and challenge for PKCE
func (m *PKCEManager) GenerateCodeChallenge() (verifier string, challenge string, err error) {
	// Generate a random code verifier (43 bytes for URL-safe base64)
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Create code verifier (base64url encoded)
	verifier = base64.URLEncoding.EncodeToString(b)

	// Create code challenge (SHA256 of verifier, base64url encoded)
	hash := sha256.Sum256([]byte(verifier))
	challenge = base64.URLEncoding.EncodeToString(hash[:])

	return verifier, challenge, nil
}

// AddValidAuthCodeWithPKCE adds a valid auth code with PKCE verifier to PKCE manager
func (m *PKCEManager) AddValidAuthCodeWithPKCE(code string, verifier string) {
	m.Codes[code] = verifier
}

// ValidateCodeVerifier validates a PKCE code verifier
func (m *PKCEManager) ValidateCodeVerifier(verifier string) bool {
	for storedVerifier := range m.Codes {
		if storedVerifier == verifier {
			delete(m.Codes, verifier) // Mark as used (one-time use)
			return true
		}
	}
	return false
}
