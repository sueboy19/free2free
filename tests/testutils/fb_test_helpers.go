package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// FacebookAuthTestHelper provides utilities for testing Facebook OAuth flow
type FacebookAuthTestHelper struct {
	BaseURL string
	Client  *http.Client
}

// NewFacebookAuthTestHelper creates a new helper instance
func NewFacebookAuthTestHelper(baseURL string) *FacebookAuthTestHelper {
	return &FacebookAuthTestHelper{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// StartFacebookAuth initiates the Facebook OAuth flow
func (fath *FacebookAuthTestHelper) StartFacebookAuth() (*http.Response, error) {
	authURL := fmt.Sprintf("%s/auth/facebook", fath.BaseURL)
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return nil, err
	}

	return fath.Client.Do(req)
}

// SimulateFacebookCallback simulates the Facebook OAuth callback
// In real testing, this would be called by Facebook after user authentication
func (fath *FacebookAuthTestHelper) SimulateFacebookCallback(code, state string) (*http.Response, error) {
	callbackURL := fmt.Sprintf("%s/auth/facebook/callback", fath.BaseURL)
	// Add query parameters
	callbackURL += fmt.Sprintf("?code=%s&state=%s", code, state)

	req, err := http.NewRequest("GET", callbackURL, nil)
	if err != nil {
		return nil, err
	}

	return fath.Client.Do(req)
}

// GetFacebookAuthResponse extracts the response from the Facebook OAuth callback
func (fath *FacebookAuthTestHelper) GetFacebookAuthResponse(resp *http.Response) (map[string]interface{}, error) {
	var authResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		// If it's not JSON, it might be a redirect (status code 302)
		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect {
			// In this case, return the redirect info
			return map[string]interface{}{
				"redirect_url": resp.Header.Get("Location"),
				"status":       resp.StatusCode,
			}, nil
		}
		return nil, err
	}
	return authResponse, nil
}

// ParseAuthURL extracts the redirect URL from the initial auth request
// This would point to Facebook's authentication page in a real implementation
func (fath *FacebookAuthTestHelper) ParseAuthURL(resp *http.Response) (string, error) {
	// If the server redirects to Facebook, the Location header will contain the Facebook auth URL
	location := resp.Header.Get("Location")
	if location == "" {
		// If no redirect happens, the server might return the URL in the response body
		// This depends on the specific implementation
		return "", fmt.Errorf("no redirect URL found in response")
	}
	return location, nil
}

// ValidateFacebookAuthSuccess validates that the Facebook auth was successful
// and extracts the JWT token from the response
func (fath *FacebookAuthTestHelper) ValidateFacebookAuthSuccess(response map[string]interface{}) (string, error) {
	// Check if the response contains access_token (JWT token in our implementation)
	if token, exists := response["access_token"]; exists {
		if tokenStr, ok := token.(string); ok && tokenStr != "" {
			return tokenStr, nil
		}
	}

	// Alternative: Check if the response structure is as expected
	// Our handlers return: {user, access_token, refresh_token, expires_in}
	if _, exists := response["user"]; exists {
		if token, exists := response["access_token"]; exists {
			if tokenStr, ok := token.(string); ok {
				return tokenStr, nil
			}
		}
	}

	return "", fmt.Errorf("access_token not found in auth response")
}

// CheckFacebookAuthError checks if the Facebook auth resulted in an error
func (fath *FacebookAuthTestHelper) CheckFacebookAuthError(resp *http.Response, response map[string]interface{}) (bool, string) {
	// Check the status code
	if resp.StatusCode >= 400 {
		return true, fmt.Sprintf("HTTP error: %d", resp.StatusCode)
	}

	// Check if the response contains an error field
	if errorField, exists := response["error"]; exists {
		if errorStr, ok := errorField.(string); ok {
			return true, errorStr
		}
	}

	// Check if the response contains error as a top-level string
	if errorStr, ok := response["error"].(string); ok && errorStr != "" {
		return true, errorStr
	}

	return false, ""
}

// PrepareFacebookCallbackURL constructs the callback URL with required parameters
func (fath *FacebookAuthTestHelper) PrepareFacebookCallbackURL(params map[string]string) string {
	baseURL := fmt.Sprintf("%s/auth/facebook/callback", fath.BaseURL)

	// Create query string
	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
	}

	if len(queryParams) > 0 {
		return baseURL + "?" + strings.Join(queryParams, "&")
	}

	return baseURL
}
