package testutils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// AdminTestHelper provides utilities for testing admin functionality
type AdminTestHelper struct {
	AuthTestHelper *AuthTestHelper
}

// NewAdminTestHelper creates a new AdminTestHelper
func NewAdminTestHelper() *AdminTestHelper {
	return &AdminTestHelper{
		AuthTestHelper: NewAuthTestHelper(),
	}
}

// NewAdminTestHelperWithSecret creates a new AdminTestHelper with a custom secret
func NewAdminTestHelperWithSecret(secret string) *AdminTestHelper {
	return &AdminTestHelper{
		AuthTestHelper: NewAuthTestHelperWithSecret(secret),
	}
}

// CreateValidAdminToken creates a valid JWT token for a test admin
func (a *AdminTestHelper) CreateValidAdminToken(userID uint, email, name, provider string) (string, error) {
	return a.AuthTestHelper.CreateValidAdminToken(userID, email, name, provider)
}

// CreateValidModeratorToken creates a valid JWT token for a test moderator
func (a *AdminTestHelper) CreateValidModeratorToken(userID uint, email, name, provider string) (string, error) {
	return CreateValidToken(userID, email, "moderator", a.AuthTestHelper.Secret)
}

// CheckAdminPermissions verifies if a token has admin-level permissions
func (a *AdminTestHelper) CheckAdminPermissions(token string) (bool, error) {
	claims, err := ValidateToken(token, a.AuthTestHelper.Secret)
	if err != nil {
		return false, err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, nil
	}

	return role == "admin", nil
}

// CheckModeratorPermissions verifies if a token has moderator-level permissions
func (a *AdminTestHelper) CheckModeratorPermissions(token string) (bool, error) {
	claims, err := ValidateToken(token, a.AuthTestHelper.Secret)
	if err != nil {
		return false, err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, nil
	}

	return role == "admin" || role == "moderator", nil
}

// CreateAdminTestServer creates a test server with admin-specific routes
func CreateAdminTestServer() (*TestServer, *AdminTestHelper, *MockAuthProvider) {
	// Create a new router for the admin test server
	gin.SetMode(gin.TestMode)
	router := gin.New()

	authHelper := NewAuthTestHelper()
	mockProvider := NewMockAuthProvider()

	SetupMockAdminRoutes(router, &AdminTestHelper{AuthTestHelper: authHelper}, mockProvider)

	server := httptest.NewServer(router)

	// Create a TestServer wrapper
	testServer := &TestServer{
		Server: server,
		Router: router,
		Config: GetTestConfig(),
	}

	return testServer, &AdminTestHelper{AuthTestHelper: authHelper}, mockProvider
}

// SetupMockAdminRoutes sets up mock admin routes for testing
func SetupMockAdminRoutes(router *gin.Engine, adminHelper *AdminTestHelper, mockProvider *MockAuthProvider) {
	// This function would set up mock routes specific to admin functionality
	// For this implementation, we'll reuse the auth setup since admin functions
	// typically build upon authentication
	SetupMockAuthRoutes(router, adminHelper.AuthTestHelper, mockProvider)

	// Add any additional admin-specific routes here
	router.GET("/admin/dashboard", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := authHeader[7:]
		isAdmin, err := adminHelper.CheckAdminPermissions(token)
		if err != nil || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden - admin permissions required"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to admin dashboard",
			"metrics": gin.H{
				"total_users":      1000,
				"total_activities": 500,
				"pending_reviews":  25,
			},
		})
	})
}
