package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"free2free/database"
	"free2free/handlers"
	"free2free/models"
	"free2free/routes"
)

// TestServer holds the test server instance and related components
type TestServer struct {
	Server *httptest.Server
	Router *gin.Engine
	DB     *gorm.DB
	Config TestConfig
}

// NewTestServer creates a new test server instance
func NewTestServer() *TestServer {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Get test configuration
	config := GetTestConfig()

	// Initialize database connection for tests
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", 
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database for testing: " + err.Error())
	}

	// Set global DB for handlers
	database.GlobalDB = &database.ActualGormDB{Conn: db}

	// Create the router
	router := gin.Default()

	// Setup routes
	routes.SetupAdminRoutes(router)
	routes.SetupUserRoutes(router)
	routes.SetupOrganizerRoutes(router)
	routes.SetupReviewRoutes(router)
	routes.SetupReviewLikeRoutes(router)

	// Add OAuth routes
	router.GET("/auth/:provider", handlers.OauthBegin)
	router.GET("/auth/:provider/callback", handlers.OauthCallback)
	router.GET("/logout", handlers.Logout)
	router.GET("/auth/token", handlers.ExchangeToken)
	router.POST("/auth/refresh", handlers.RefreshTokenHandler)
	router.GET("/profile", handlers.Profile)

	// Create test server
	server := httptest.NewServer(router)

	return &TestServer{
		Server: server,
		Router: router,
		DB:     db,
		Config: config,
	}
}

// Close closes the test server and cleans up resources
func (ts *TestServer) Close() {
	ts.Server.Close()
}

// GetURL returns the test server URL
func (ts *TestServer) GetURL(path string) string {
	return ts.Server.URL + path
}

// DoRequest makes an HTTP request to the test server
func (ts *TestServer) DoRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, ts.GetURL(path), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set content type if not already set and we have a body
	if body != nil && headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// CreateTestUser creates a test user in the database
func (ts *TestServer) CreateTestUser() (*models.User, error) {
	user := &models.User{
		SocialID:       "test_user_social_id",
		SocialProvider: "facebook",
		Name:           "Test User",
		Email:          "test@example.com",
		AvatarURL:      "https://example.com/avatar.jpg",
		IsAdmin:        false,
	}

	result := ts.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// SetupTestDatabase runs migrations for test models
func (ts *TestServer) SetupTestDatabase() error {
	return ts.DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Location{},
		&models.Activity{},
		&models.Match{},
		&models.MatchParticipant{},
		&models.Review{},
		&models.ReviewLike{},
		&models.RefreshToken{},
	)
}

// ClearTestData clears test data from the database
func (ts *TestServer) ClearTestData() error {
	tables := []interface{}{
		&models.ReviewLike{}, &models.Review{}, &models.MatchParticipant{},
		&models.Match{}, &models.Activity{}, &models.Location{},
		&models.Admin{}, &models.User{}, &models.RefreshToken{},
	}

	for _, table := range tables {
		if err := ts.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			// If table doesn't exist, continue
			if !strings.Contains(err.Error(), "doesn't exist") {
				return err
			}
		}
	}

	return nil
}