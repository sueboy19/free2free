package testutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestServer wraps httptest.Server with additional utilities
type TestServer struct {
	Server *httptest.Server
	Router *gin.Engine
	Config TestConfig
	DB     *gorm.DB
}

// NewTestServer creates a new test server with a configured router
func NewTestServer() *TestServer {
	config := GetTestConfig()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Add middlewares for testing
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS for testing
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	server := httptest.NewServer(router)

	// Create test database instance
	db, err := CreateTestDB()
	if err != nil {
		// If we can't create a test DB, log it but continue
		// Tests that require the DB will handle the nil check
		db = nil
	}

	return &TestServer{
		Server: server,
		Router: router,
		Config: config,
		DB:     db,
	}
}

// NewTestServerWithTimeout creates a new test server with a specified timeout
func NewTestServerWithTimeout(timeout time.Duration) *TestServer {
	ts := NewTestServer()
	ts.Config.TestTimeout = timeout
	return ts
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.Server.Close()
}

// ClearTestData clears test data from the database
func (ts *TestServer) ClearTestData() error {
	// This method is a placeholder for clearing test data
	// Implementation would depend on specific test database setup
	return nil
}

// SetupTestDatabase sets up the test database
func (ts *TestServer) SetupTestDatabase() error {
	// This method is a placeholder for setting up the test database
	// Implementation would depend on specific test database setup
	return nil
}

// CreateTestUser creates a test user
func (ts *TestServer) CreateTestUser() (*TestUser, error) {
	// Return a standard test user
	user := &TestUser{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "facebook",
		Role:     "user",
		IsAdmin:  false, // Adding IsAdmin field as it's referenced in integration tests
	}
	return user, nil
}

// GetURL returns the full URL for an endpoint
func (ts *TestServer) GetURL(endpoint string) string {
	return ts.Server.URL + endpoint
}

// DoRequest sends an HTTP request to the test server
func (ts *TestServer) DoRequest(method, endpoint string, headers map[string]string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	
	url := ts.GetURL(endpoint)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	
	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	return client.Do(req)
}
