package testutils

import (
	"os"
)

// SetupTestEnvironment sets up the required environment variables for OAuth flow testing
func SetupTestEnvironment() {
	// Set default environment variables required for OAuth flow
	if os.Getenv("SESSION_KEY") == "" {
		os.Setenv("SESSION_KEY", "test-session-key-for-oauth-flow-32-characters!")
	}
	
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-jwt-secret-for-testing-environment")
	}
	
	if os.Getenv("FACEBOOK_KEY") == "" {
		os.Setenv("FACEBOOK_KEY", "test-facebook-key")
	}
	
	if os.Getenv("FACEBOOK_SECRET") == "" {
		os.Setenv("FACEBOOK_SECRET", "test-facebook-secret")
	}
	
	if os.Getenv("INSTAGRAM_KEY") == "" {
		os.Setenv("INSTAGRAM_KEY", "test-instagram-key")
	}
	
	if os.Getenv("INSTAGRAM_SECRET") == "" {
		os.Setenv("INSTAGRAM_SECRET", "test-instagram-secret")
	}
	
	if os.Getenv("BASE_URL") == "" {
		os.Setenv("BASE_URL", "http://localhost:8080")
	}
	
	// Set secure cookie to false for testing
	if os.Getenv("SECURE_COOKIE") == "" {
		os.Setenv("SECURE_COOKIE", "false")
	}
}

// CleanupTestEnvironment removes the test environment variables
func CleanupTestEnvironment() {
	os.Unsetenv("SESSION_KEY")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("FACEBOOK_KEY")
	os.Unsetenv("FACEBOOK_SECRET")
	os.Unsetenv("INSTAGRAM_KEY")
	os.Unsetenv("INSTAGRAM_SECRET")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("SECURE_COOKIE")
}

// SaveOriginalEnvironment saves the current environment variables
func SaveOriginalEnvironment() map[string]string {
	originalEnv := make(map[string]string)
	
	originalEnv["SESSION_KEY"] = os.Getenv("SESSION_KEY")
	originalEnv["JWT_SECRET"] = os.Getenv("JWT_SECRET")
	originalEnv["FACEBOOK_KEY"] = os.Getenv("FACEBOOK_KEY")
	originalEnv["FACEBOOK_SECRET"] = os.Getenv("FACEBOOK_SECRET")
	originalEnv["INSTAGRAM_KEY"] = os.Getenv("INSTAGRAM_KEY")
	originalEnv["INSTAGRAM_SECRET"] = os.Getenv("INSTAGRAM_SECRET")
	originalEnv["BASE_URL"] = os.Getenv("BASE_URL")
	originalEnv["SECURE_COOKIE"] = os.Getenv("SECURE_COOKIE")
	
	return originalEnv
}

// RestoreOriginalEnvironment restores the original environment variables
func RestoreOriginalEnvironment(originalEnv map[string]string) {
	for key, value := range originalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}