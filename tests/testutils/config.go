package testutils

import (
	"os"
)

// TestConfig holds the configuration for tests
type TestConfig struct {
	FacebookAppID     string
	FacebookAppSecret string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	ServerPort       string
	BaseURL          string
}

// GetTestConfig retrieves test configuration from environment variables
func GetTestConfig() TestConfig {
	return TestConfig{
		FacebookAppID:     getEnv("TEST_FACEBOOK_KEY", "test-facebook-app-id"),
		FacebookAppSecret: getEnv("TEST_FACEBOOK_SECRET", "test-facebook-app-secret"),
		DBHost:           getEnv("TEST_DB_HOST", "localhost"),
		DBPort:           getEnv("TEST_DB_PORT", "3306"),
		DBUser:           getEnv("TEST_DB_USER", "root"),
		DBPassword:       getEnv("TEST_DB_PASSWORD", "password"),
		DBName:           getEnv("TEST_DB_NAME", "free2free_test"),
		JWTSecret:        getEnv("TEST_JWT_SECRET", "test-jwt-secret-key-32-chars-long-enough!!"),
		ServerPort:       getEnv("TEST_SERVER_PORT", "8080"),
		BaseURL:          getEnv("TEST_BASE_URL", "http://localhost:8080"),
	}
}

// getEnv retrieves environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}