package e2e

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"free2free/tests/testutils"
)

// TestEnvironmentSetupValidation tests that the test environment is properly set up
func TestEnvironmentSetupValidation(t *testing.T) {
	t.Run("Required environment variables are set", func(t *testing.T) {
		// Check that required environment variables are available
		envVars := []string{
			"TEST_DB_HOST",
			"TEST_DB_PORT",
			"TEST_DB_USER",
			"TEST_DB_PASSWORD",
			"TEST_DB_NAME",
			"TEST_JWT_SECRET",
			"TEST_FACEBOOK_KEY",
			"TEST_FACEBOOK_SECRET",
		}

		for _, envVar := range envVars {
			value := os.Getenv(envVar)
			// For this test, we'll use default values if env vars aren't set
			// In real implementation, these would be required
			assert.NotEmpty(t, value, "Environment variable %s should be set", envVar)
		}
	})

	t.Run("Database connection is available", func(t *testing.T) {
		// Initialize test server which connects to database
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Verify that we can ping the database
		db, err := testServer.DB.DB()
		assert.NoError(t, err)

		err = db.Ping()
		assert.NoError(t, err, "Should be able to connect to database")
	})

	t.Run("Test database schema exists", func(t *testing.T) {
		// Initialize test server
		testServer := testutils.NewTestServer()
		defer testServer.Close()

		// Verify that the database has the required tables
		// We'll try to access one of the models to ensure the schema exists
		err := testServer.SetupTestDatabase()
		assert.NoError(t, err, "Should be able to setup test database schema")
	})

	t.Run("Test configuration loads properly", func(t *testing.T) {
		config := testutils.GetTestConfig()

		assert.NotEmpty(t, config.DBHost, "DBHost should not be empty")
		assert.NotEmpty(t, config.DBPort, "DBPort should not be empty")
		assert.NotEmpty(t, config.JWTSecret, "JWTSecret should not be empty")

		// For JWT secret, ensure it's of sufficient length
		assert.GreaterOrEqual(t, len(config.JWTSecret), 32, "JWTSecret should be at least 32 characters")
	})

	t.Run("Test server starts without errors", func(t *testing.T) {
		// This test verifies that the test server can be initialized without errors
		testServer := testutils.NewTestServer()

		// If we got here without panicking, the server initialization worked
		assert.NotNil(t, testServer, "Test server should initialize without errors")
		assert.NotNil(t, testServer.Server, "Test server instance should exist")
		assert.NotNil(t, testServer.Router, "Test router should exist")

		// Clean up
		testServer.Close()
	})

	t.Run("Required dependencies are available", func(t *testing.T) {
		// This test checks if all required Go dependencies are available
		// In Go, dependencies are resolved at compile time, so we're testing
		// that we can import and use the required packages

		// We've already used testutils, which depends on gin, gorm, and other packages
		config := testutils.GetTestConfig()
		assert.NotEmpty(t, config.BaseURL, "Configuration should load properly")

		// Initialize test server (which uses gin and gorm)
		testServer := testutils.NewTestServer()
		assert.NotNil(t, testServer)
		testServer.Close()
	})
}

// TestEnvironmentConsistency tests that the test environment is consistent
func TestEnvironmentConsistency(t *testing.T) {
	t.Run("Multiple server instances have consistent configuration", func(t *testing.T) {
		// Start first server
		server1 := testutils.NewTestServer()
		config1 := server1.Config

		// Start second server
		server2 := testutils.NewTestServer()
		config2 := server2.Config

		// Both servers should have the same base configuration
		// (Though they run on different ports)
		assert.Equal(t, config1.DBHost, config2.DBHost, "DB host should be consistent")
		assert.Equal(t, config1.DBPort, config2.DBPort, "DB port should be consistent")
		assert.Equal(t, config1.JWTSecret, config2.JWTSecret, "JWT secret should be consistent")

		server1.Close()
		server2.Close()
	})

	t.Run("Environment configuration does not change during test run", func(t *testing.T) {
		// Get config at start
		configStart := testutils.GetTestConfig()

		// Simulate some test operations
		testServer := testutils.NewTestServer()
		testServer.Close()

		// Get config at end
		configEnd := testutils.GetTestConfig()

		// Configuration should remain the same
		assert.Equal(t, configStart.DBHost, configEnd.DBHost)
		assert.Equal(t, configStart.JWTSecret, configEnd.JWTSecret)
		assert.Equal(t, configStart.FacebookAppID, configEnd.FacebookAppID)
	})
}
