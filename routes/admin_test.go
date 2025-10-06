package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"free2free/database"
	"free2free/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDatabase() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Set up global database instance for tests
	database.SetGlobalDB(&database.ActualGormDB{Conn: db})
	
	// Auto migrate the tables for testing
	db.AutoMigrate(&models.Activity{}, &models.Location{})
}

func mockAuthenticatedUser(c *gin.Context) (*models.User, error) {
	return &models.User{ID: 1, Name: "Admin User", IsAdmin: true}, nil
}

func TestListActivities(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	listActivities(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var activities []models.Activity
	err := json.Unmarshal(w.Body.Bytes(), &activities)
	assert.NoError(t, err)
	// Should return an empty array since no activities have been created
	assert.Len(t, activities, 0)
}

func TestCreateActivity(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"title":"New Activity","target_count":4,"location_id":1,"description":"New Test"}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	createActivity(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdActivity models.Activity
	err := json.Unmarshal(w.Body.Bytes(), &createdActivity)
	assert.NoError(t, err)
	assert.Equal(t, "New Activity", createdActivity.Title)
}

func TestUpdateActivity(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer([]byte(`{"title":"Updated Activity","target_count":5,"location_id":1,"description":"Updated Test"}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	updateActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var updatedActivity models.Activity
	err := json.Unmarshal(w.Body.Bytes(), &updatedActivity)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Activity", updatedActivity.Title)
}

func TestDeleteActivity(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	deleteActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "活動已刪除", response["message"])
}

func TestListLocations(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	listLocations(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var locations []models.Location
	err := json.Unmarshal(w.Body.Bytes(), &locations)
	assert.NoError(t, err)
	// Should return empty array since no locations have been created
	assert.Len(t, locations, 0)
}

func TestCreateLocation(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"name":"New Location","address":"New Addr","latitude":25.0,"longitude":121.0}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	createLocation(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdLocation models.Location
	err := json.Unmarshal(w.Body.Bytes(), &createdLocation)
	assert.NoError(t, err)
	assert.Equal(t, "New Location", createdLocation.Name)
}

func TestUpdateLocation(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer([]byte(`{"name":"Updated Location","address":"Updated Addr","latitude":26.0,"longitude":122.0}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	updateLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var updatedLocation models.Location
	err := json.Unmarshal(w.Body.Bytes(), &updatedLocation)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Location", updatedLocation.Name)
}

func TestDeleteLocation(t *testing.T) {
	setupTestDatabase()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock the context to simulate an authenticated user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})

	deleteLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "地點已刪除", response["message"])
}

func TestAdminAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test unauthorized - simulate a non-admin user
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock a non-admin user
	c.Set("user", &models.User{ID: 2, Name: "Regular User", IsAdmin: false})

	middleware := AdminAuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// Test authorized - simulate an admin user
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	// Mock an admin user
	c.Set("user", &models.User{ID: 1, Name: "Admin User", IsAdmin: true})
	middleware(c)

	// If the user is admin, the middleware should call c.Next() and not write an error
	// The status code would remain 0 (default) if the next handler is not called
	assert.Equal(t, 0, w.Code) // Middleware didn't write an error response, so status remains 0
}
