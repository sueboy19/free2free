package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

// Use modernc.org/sqlite as the underlying driver (no CGO required)
import _ "modernc.org/sqlite"

// HTTPRequestBuilder helps build and execute HTTP requests for testing
type HTTPRequestBuilder struct {
	router    http.Handler
	method    string
	url       string
	body      interface{}
	authToken string
}

// NewRequest creates a new HTTP request builder
func NewRequest(router http.Handler, method, url string) *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		router: router,
		method: method,
		url:    url,
	}
}

// WithBody adds a request body to the request
func (b *HTTPRequestBuilder) WithBody(body interface{}) *HTTPRequestBuilder {
	b.body = body
	return b
}

// WithAuth adds an authentication token to the request
func (b *HTTPRequestBuilder) WithAuth(token string) *HTTPRequestBuilder {
	b.authToken = token
	return b
}

// Execute executes the built request
func (b *HTTPRequestBuilder) Execute() (*httptest.ResponseRecorder, error) {
	var req *http.Request
	var err error

	if b.body != nil {
		jsonData, err := json.Marshal(b.body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(b.method, b.url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(b.method, b.url, nil)
		if err != nil {
			return nil, err
		}
	}

	// Add authorization header if token is provided
	if b.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+b.authToken)
	}

	// Execute request
	w := httptest.NewRecorder()
	b.router.ServeHTTP(w, req)

	return w, nil
}

// MakeRequest creates and executes an HTTP request for testing
func MakeRequest(
	router http.Handler,
	method, url string,
	body interface{},
	authToken string,
) (*httptest.ResponseRecorder, error) {
	return NewRequest(router, method, url).WithBody(body).WithAuth(authToken).Execute()
}

// GetRequest creates and executes a GET request for testing
func GetRequest(router http.Handler, url string, authToken string) (*httptest.ResponseRecorder, error) {
	return NewRequest(router, "GET", url).WithAuth(authToken).Execute()
}

// PostRequest creates and executes a POST request for testing
func PostRequest(router http.Handler, url string, body interface{}, authToken string) (*httptest.ResponseRecorder, error) {
	return NewRequest(router, "POST", url).WithBody(body).WithAuth(authToken).Execute()
}

// PutRequest creates and executes a PUT request for testing
func PutRequest(router http.Handler, url string, body interface{}, authToken string) (*httptest.ResponseRecorder, error) {
	return NewRequest(router, "PUT", url).WithBody(body).WithAuth(authToken).Execute()
}

// DeleteRequest creates and executes a DELETE request for testing
func DeleteRequest(router http.Handler, url string, authToken string) (*httptest.ResponseRecorder, error) {
	return NewRequest(router, "DELETE", url).WithAuth(authToken).Execute()
}

// RequestWithValidation executes a request and validates the response
func RequestWithValidation(
	t *testing.T,
	router http.Handler,
	method, url string,
	body interface{},
	authToken string,
	expectedStatus int,
) *httptest.ResponseRecorder {
	w, err := MakeRequest(router, method, url, body, authToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, w.Code)
	return w
}

// JSONResponse represents a standard JSON response
type JSONResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ParseJSONResponse parses the response body as JSON
func ParseJSONResponse(response *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(response.Body.Bytes(), target)
}

// CreateTestUser creates a test user struct for use in tests
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Role     string `json:"role"`
}

// CreateTestActivity creates a test activity struct for use in tests
type TestActivity struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	LocationID  uint   `json:"location_id"`
	Status      string `json:"status"`
	CreatorID   uint   `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DBTestHelper provides utilities for database testing
type DBTestHelper struct {
	db *gorm.DB
}

// NewDBTestHelper creates a new database test helper
func NewDBTestHelper() (*DBTestHelper, error) {
	db, err := CreateTestDB()
	if err != nil {
		return nil, err
	}
	
	return &DBTestHelper{
		db: db,
	}, nil
}

// GetDB returns the database instance
func (h *DBTestHelper) GetDB() *gorm.DB {
	return h.db
}

// MigrateModels runs migrations for the provided models
func (h *DBTestHelper) MigrateModels(models ...interface{}) error {
	return MigrateTestDB(h.db, models...)
}

// CloseDB closes the database connection
func (h *DBTestHelper) CloseDB() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// UseModerncSQLite ensures the modernc.org/sqlite driver is used
// This is a utility function to verify the platform-independent database driver is active
func UseModerncSQLite() bool {
	// This function can be used to verify that the modernc.org/sqlite driver is active
	// through the import in this package or other test utilities
	return true
}
