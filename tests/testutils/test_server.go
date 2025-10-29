package testutils

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

// TestServer wraps httptest.Server with additional utilities
type TestServer struct {
	Server *httptest.Server
	Router *gin.Engine
	Config TestConfig
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

	return &TestServer{
		Server: server,
		Router: router,
		Config: config,
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
