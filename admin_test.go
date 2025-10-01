package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockDB struct {
	calledMethods map[string]interface{}
	data          map[string]interface{}
}

func (m *MockDB) AutoMigrate(dst ...interface{}) error {
	return nil
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	// Simulate create
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	if users, ok := value.(*[]User); ok {
		for i, u := range *users {
			newU := u
			newU.ID = int64(i + 1)
			if existing, ok := m.data["users"].([]User); ok {
				m.data["users"] = append(existing, newU)
			} else {
				m.data["users"] = []User{newU}
			}
		}
	}
	if activity, ok := value.(*Activity); ok {
		activity.ID = 1
		if existing, ok := m.data["activities"].([]Activity); ok {
			m.data["activities"] = append(existing, *activity)
		} else {
			m.data["activities"] = []Activity{*activity}
		}
	}
	if location, ok := value.(*Location); ok {
		location.ID = 1
		if existing, ok := m.data["locations"].([]Location); ok {
			m.data["locations"] = append(existing, *location)
		} else {
			m.data["locations"] = []Location{*location}
		}
	}
	return &gorm.DB{Error: nil}
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	if user, ok := dest.(*User); ok {
		if users, ok := m.data["users"].([]User); ok && len(users) > 0 {
			*user = users[0]
			return &gorm.DB{Error: nil}
		}
		return &gorm.DB{Error: gorm.ErrRecordNotFound}
	}
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Preload(query string, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Order(value interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Joins(query string, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Raw(sql string, values ...interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) WithContext(ctx context.Context) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) DB() (*sql.DB, error) {
	return nil, nil
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	if activities, ok := dest.(*[]Activity); ok {
		*activities = []Activity{
			{ID: 1, Title: "Test Activity", TargetCount: 4, LocationID: 1, Description: "Test", CreatedBy: 1},
		}
		return &gorm.DB{Error: nil}
	}
	if locations, ok := dest.(*[]Location); ok {
		*locations = []Location{
			{ID: 1, Name: "Test Location", Address: "Test Addr", Latitude: 25.0330, Longitude: 121.5654},
		}
		return &gorm.DB{Error: nil}
	}
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return &gorm.DB{Error: nil}
}

func mockAuthenticatedUser(c *gin.Context) (*User, error) {
	return &User{ID: 1, Name: "Admin User", IsAdmin: true}, nil
}

func TestListActivities(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	listActivities(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var activities []Activity
	err := json.Unmarshal(w.Body.Bytes(), &activities)
	assert.NoError(t, err)
	assert.Len(t, activities, 1)
	assert.Equal(t, "Test Activity", activities[0].Title)
}

func TestCreateActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"title":"New Activity","target_count":4,"location_id":1,"description":"New Test"}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	createActivity(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdActivity Activity
	err := json.Unmarshal(w.Body.Bytes(), &createdActivity)
	assert.NoError(t, err)
	assert.Equal(t, "New Activity", createdActivity.Title)
}

func TestUpdateActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer([]byte(`{"title":"Updated Activity","target_count":5,"location_id":1,"description":"Updated Test"}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	updateActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var updatedActivity Activity
	err := json.Unmarshal(w.Body.Bytes(), &updatedActivity)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Activity", updatedActivity.Title)
}

func TestDeleteActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	deleteActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "活動已刪除", response["message"])
}

func TestListLocations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	listLocations(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var locations []Location
	err := json.Unmarshal(w.Body.Bytes(), &locations)
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Equal(t, "Test Location", locations[0].Name)
}

func TestCreateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"name":"New Location","address":"New Addr","latitude":25.0,"longitude":121.0}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	createLocation(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdLocation Location
	err := json.Unmarshal(w.Body.Bytes(), &createdLocation)
	assert.NoError(t, err)
	assert.Equal(t, "New Location", createdLocation.Name)
}

func TestUpdateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer([]byte(`{"name":"Updated Location","address":"Updated Addr","latitude":26.0,"longitude":122.0}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	updateLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var updatedLocation Location
	err := json.Unmarshal(w.Body.Bytes(), &updatedLocation)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Location", updatedLocation.Name)
}

func TestDeleteLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Mock adminDB
	oldAdminDB := adminDB
	adminDB = &MockDB{}
	defer func() {
		adminDB = oldAdminDB
	}()

	// Mock getAuthenticatedUser
	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = mockAuthenticatedUser
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	deleteLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "地點已刪除", response["message"])
}

func TestAdminAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test unauthorized
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	oldGetAuth := getAuthenticatedUser
	getAuthenticatedUser = func(c *gin.Context) (*User, error) {
		return nil, errors.New("unauthorized")
	}
	defer func() {
		getAuthenticatedUser = oldGetAuth
	}()

	middleware := AdminAuthMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test authorized
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	getAuthenticatedUser = mockAuthenticatedUser
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code) // No error
}
