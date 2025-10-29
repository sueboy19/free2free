package testutils

import (
	"fmt"
	"os"
	"time"
)

// Use modernc.org/sqlite as the underlying driver (no CGO required)
import _ "modernc.org/sqlite"

// TestConfig holds configuration for tests
type TestConfig struct {
	DatabaseURL    string
	ServerPort     string
	JWTSecret      string
	TestTimeout    time.Duration
	MaxConnections int
}

// GetTestConfig loads test configuration from environment variables or defaults
func GetTestConfig() TestConfig {
	return TestConfig{
		DatabaseURL:    getEnvOrDefault("TEST_DATABASE_URL", "sqlite://test.db"),
		ServerPort:     getEnvOrDefault("TEST_SERVER_PORT", "8081"),
		JWTSecret:      getEnvOrDefault("TEST_JWT_SECRET", "test-secret-for-development"),
		TestTimeout:    time.Duration(getEnvIntOrDefault("TEST_TIMEOUT", 30)) * time.Second,
		MaxConnections: getEnvIntOrDefault("TEST_MAX_CONNECTIONS", 10),
	}
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := parseInt(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func parseInt(s string) (int, error) {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0, fmt.Errorf("invalid integer: %s", s)
		}
	}
	return result, nil
}
