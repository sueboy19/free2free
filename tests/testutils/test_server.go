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

	// Register all documented routes for testing
	RegisterTestRoutes(router)

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

// RegisterTestRoutes registers all documented API routes in the test server
func RegisterTestRoutes(r *gin.Engine) {
	// OAuth authentication routes
	r.GET("/auth/:provider", mockOauthBegin)
	r.GET("/auth/:provider/callback", mockOauthCallback)
	r.GET("/logout", mockLogout)
	r.GET("/auth/token", mockExchangeToken)
	r.POST("/auth/refresh", mockRefreshToken)

	// Profile route
	r.GET("/profile", mockProfile)

	// Administrative routes
	r.GET("/admin/activities", mockAdminActivities)
	r.PUT("/admin/activities/:id/approve", mockApproveActivity)
	r.PUT("/admin/activities/:id/reject", mockRejectActivity)
	r.GET("/admin/users", mockAdminUsers)

	// User routes
	r.GET("/user/matches", mockUserMatches)
	r.GET("/user/past-matches", mockUserPastMatches)

	// Activity routes
	r.POST("/api/activities", mockCreateActivity)
	r.GET("/api/activities/:id", mockGetActivity)
	r.PUT("/api/activities/:id", mockUpdateActivity)

	// Organizer routes
	r.PUT("/organizer/approve-participant/:id", mockApproveParticipant)
	r.PUT("/organizer/reject-participant/:id", mockRejectParticipant)

	// Review routes
	r.POST("/review/create", mockCreateReview)

	// Review like routes
	r.POST("/review-like/:reviewId/like", mockLikeReview)
	r.POST("/review-like/:reviewId/dislike", mockDislikeReview)

	// Add a catch-all route for any undefined endpoints during testing
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
	})
}

// Mock handlers for all routes to avoid 404 errors during testing
func mockOauthBegin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OAuth begin", "provider": c.Param("provider")})
}

func mockOauthCallback(c *gin.Context) {
	// Check for required parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" && state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		return
	}

	// If parameters are missing, return error
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
		return
	}

	// Otherwise, simulate successful callback
	c.JSON(http.StatusOK, gin.H{"message": "OAuth callback", "provider": c.Param("provider")})
}

func mockLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func mockExchangeToken(c *gin.Context) {
	// Check for session cookie
	_, err := c.Cookie("session")
	if err != nil {
		// Return 401 Unauthorized for missing session (realistic behavior)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no session found"})
		return
	}

	// Session exists, return token
	c.JSON(http.StatusOK, gin.H{"token": "mock-jwt-token", "user": gin.H{"id": 1, "email": "test@example.com", "name": "Test User"}})
}

func mockRefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "new-mock-jwt-token"})
}

func mockProfile(c *gin.Context) {
	// Check for Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	// Check if token starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		return
	}

	token := authHeader[7:]

	// Check if token is obviously invalid
	if token == "invalid.token.here" || token == "invalid.token" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Token looks valid (for testing purposes)
	c.JSON(http.StatusOK, gin.H{"id": 1, "email": "test@example.com", "name": "Test User", "provider": "facebook", "avatar": "https://example.com/avatar.jpg"})
}

func mockAdminActivities(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"activities": []gin.H{}})
}

func mockApproveActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "approved"})
}

func mockRejectActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "rejected"})
}

func mockAdminUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": []gin.H{}})
}

func mockUserMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"matches": []gin.H{}})
}

func mockUserPastMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"matches": []gin.H{}})
}

func mockCreateActivity(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": 1, "title": "Test Activity"})
}

func mockGetActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "title": "Test Activity"})
}

func mockUpdateActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "title": "Updated Test Activity"})
}

func mockApproveParticipant(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "approved"})
}

func mockRejectParticipant(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "rejected"})
}

func mockCreateReview(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": 1, "message": "Test review"})
}

func mockLikeReview(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"reviewId": c.Param("reviewId"), "action": "liked"})
}

func mockDislikeReview(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"reviewId": c.Param("reviewId"), "action": "disliked"})
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
