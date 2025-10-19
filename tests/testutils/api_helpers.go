package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIResponse represents a generic API response
type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ParseResponse parses an HTTP response into an APIResponse struct
func ParseResponse(resp *http.Response, target interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// For error responses, try to parse the error message
		var apiResp APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return fmt.Errorf("request failed with status: %d", resp.StatusCode)
		}
		if apiResp.Error != "" {
			return fmt.Errorf("API error: %s", apiResp.Error)
		}
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}

// MakeAuthenticatedRequest makes an HTTP request with an authorization header
func MakeAuthenticatedRequest(testServer *TestServer, method, path, token string, body interface{}) (*http.Response, error) {
	headers := make(map[string]string)
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}

	return testServer.DoRequest(method, path, body, headers)
}

// GetUserProfile makes a request to get the user profile
func GetUserProfile(testServer *TestServer, token string) (*APIResponse, error) {
	resp, err := MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	err = ParseResponse(resp, &apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// CheckResponseStatus checks if the response has the expected status code
func CheckResponseStatus(resp *http.Response, expectedCode int) error {
	if resp.StatusCode != expectedCode {
		return fmt.Errorf("expected status code %d, got %d", expectedCode, resp.StatusCode)
	}
	return nil
}

// ValidateAPIResponse validates that an API response has the expected structure
func ValidateAPIResponse(resp *http.Response, expectedStatusCode int) error {
	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status %d, got %d", expectedStatusCode, resp.StatusCode)
	}

	// Try to parse the response to ensure it's valid JSON
	var apiResp APIResponse
	err := json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return fmt.Errorf("invalid JSON response: %w", err)
	}

	return nil
}