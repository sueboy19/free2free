package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/markbates/goth"
)

// MockFacebookProvider simulates a Facebook OAuth provider for testing
type MockFacebookProvider struct {
	Server *httptest.Server
	AppID  string
}

// NewMockFacebookProvider creates a new mock Facebook provider
func NewMockFacebookProvider() *MockFacebookProvider {
	mfp := &MockFacebookProvider{
		AppID: "test-facebook-app-id",
	}

	// Create a test server that simulates Facebook OAuth endpoints
	server := httptest.NewServer(http.HandlerFunc(mfp.handler))
	mfp.Server = server

	return mfp
}

// Close closes the mock provider server
func (mfp *MockFacebookProvider) Close() {
	mfp.Server.Close()
}

// GetAuthURL returns the authorization URL for the mock provider
func (mfp *MockFacebookProvider) GetAuthURL(redirectURL string) string {
	return fmt.Sprintf("%s/auth?redirect_uri=%s", mfp.Server.URL, redirectURL)
}

// handler implements the HTTP handler for mock Facebook endpoints
func (mfp *MockFacebookProvider) handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth":
		mfp.handleAuth(w, r)
	case "/callback":
		mfp.handleCallback(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleAuth handles the authorization request
func (mfp *MockFacebookProvider) handleAuth(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	
	if redirectURI == "" {
		http.Error(w, "Missing redirect_uri", http.StatusBadRequest)
		return
	}

	// Redirect to callback with a mock code
	callbackURL := fmt.Sprintf("%s?code=mock_auth_code&state=", redirectURI)
	http.Redirect(w, r, callbackURL, http.StatusFound)
}

// handleCallback handles the OAuth callback and returns user data
func (mfp *MockFacebookProvider) handleCallback(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would verify the code and return an access token
	// For testing, we directly return user information
	
	user := goth.User{
		UserID:         "mock_facebook_user_id_12345",
		Email:          "mockuser@example.com",
		Name:           "Mock Test User",
		FirstName:      "Mock",
		LastName:       "User",
		NickName:       "mockuser",
		Description:    "A mock Facebook user for testing",
		AvatarURL:      "https://example.com/mock-avatar.jpg",
		Provider:       "facebook",
		AccessToken:    "mock_access_token_12345",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
		RefreshToken:   "mock_refresh_token_12345",
	}
	
	// Return user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetMockUser returns a mock user object for testing
func GetMockUser() goth.User {
	return goth.User{
		UserID:         "mock_facebook_user_id_12345",
		Email:          "mockuser@example.com",
		Name:           "Mock Test User",
		FirstName:      "Mock",
		LastName:       "User",
		NickName:       "mockuser",
		Description:    "A mock Facebook user for testing",
		AvatarURL:      "https://example.com/mock-avatar.jpg",
		Provider:       "facebook",
		AccessToken:    "mock_access_token_12345",
		ExpiresAt:      time.Now().Add(1 * time.Hour),
		RefreshToken:   "mock_refresh_token_12345",
	}
}