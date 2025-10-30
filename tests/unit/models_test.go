package unit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Use modernc.org/sqlite as the underlying driver (no CGO required)
import _ "modernc.org/sqlite"

// Simple test models for GORM CRUD operations
type TestUser struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

type TestLocation struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

type TestActivity struct {
	ID         int64  `gorm:"primaryKey"`
	Title      string `gorm:"size:255"`
	LocationID int64
}

type TestMatch struct {
	ID         int64 `gorm:"primaryKey"`
	ActivityID int64
	Status     string `gorm:"size:50"`
}

type TestMatchParticipant struct {
	ID      int64 `gorm:"primaryKey"`
	MatchID int64
	Status  string `gorm:"size:50"`
}

type TestReview struct {
	ID      int64 `gorm:"primaryKey"`
	MatchID int64
	Score   int
	Comment string `gorm:"size:1000"`
}

type TestDB struct {
	conn *gorm.DB
}

func (d *TestDB) AutoMigrate(dst ...interface{}) error {
	return d.conn.AutoMigrate(dst...)
}

func (d *TestDB) Create(value interface{}) *gorm.DB {
	return d.conn.Create(value)
}

func (d *TestDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.First(dest, conds...)
}

func (d *TestDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return d.conn.Where(query, args...)
}

func (d *TestDB) Save(value interface{}) *gorm.DB {
	return d.conn.Save(value)
}

func (d *TestDB) WithContext(ctx context.Context) *gorm.DB {
	return d.conn.WithContext(ctx)
}

func (d *TestDB) Preload(query string, args ...interface{}) *gorm.DB {
	return d.conn.Preload(query, args...)
}

func (d *TestDB) Model(value interface{}) *gorm.DB {
	return d.conn.Model(value)
}

func (d *TestDB) Update(column string, value interface{}) *gorm.DB {
	return d.conn.Model(nil).Update(column, value)
}

func (d *TestDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.Find(dest, conds...)
}

func (d *TestDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.Delete(value, conds...)
}

func (d *TestDB) Joins(query string, args ...interface{}) *gorm.DB {
	return d.conn.Joins(query, args...)
}

func (d *TestDB) Raw(sql string, values ...interface{}) *gorm.DB {
	return d.conn.Raw(sql, values...)
}

func (d *TestDB) Order(value interface{}) *gorm.DB {
	return d.conn.Order(value)
}

func (d *TestDB) DB() (*sql.DB, error) {
	return d.conn.DB()
}

func setupTestDB(t *testing.T) *TestDB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Logf("Database connection error: %v", err)
		if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
			err.Error() == "failed to connect database") {
			t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
		}
		assert.NoError(t, err)
		return nil
	}

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxIdleConns(10)

	testDB := &TestDB{conn: db}
	err = testDB.AutoMigrate(
		&TestUser{},
		&TestLocation{},
		&TestActivity{},
		&TestMatch{},
		&TestMatchParticipant{},
		&TestReview{},
	)
	if err != nil {
		t.Logf("Migration error: %v", err)
		if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
			err.Error() == "Error 101 (HY000): failed to open database") {
			t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
		}
		assert.NoError(t, err)
		return nil
	}

	return testDB
}

func TestGORMCreate(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Test Create
	user := &TestUser{Name: "Test User"}
	result := db.Create(user)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), user.ID)

	location := &TestLocation{Name: "Test Location"}
	result = db.Create(location)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), location.ID)

	activity := &TestActivity{Title: "Test Activity", LocationID: 1}
	result = db.Create(activity)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), activity.ID)

	match := &TestMatch{ActivityID: 1, Status: "open"}
	result = db.Create(match)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), match.ID)

	participant := &TestMatchParticipant{MatchID: 1, Status: "pending"}
	result = db.Create(participant)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), participant.ID)

	review := &TestReview{MatchID: 1, Score: 5, Comment: "Great!"}
	result = db.Create(review)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), review.ID)
}

func TestGORMFind(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := &TestUser{Name: "Test User"}
	result := db.Create(user)
	assert.NoError(t, result.Error)

	// Test First
	var foundUser TestUser
	result = db.First(&foundUser, 1)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Test User", foundUser.Name)

	// Test Find
	var users []TestUser
	result = db.Where("id = ?", 1).Find(&users)
	assert.NoError(t, result.Error)
	assert.Len(t, users, 1)
	assert.Equal(t, "Test User", users[0].Name)
}

func TestGORMUpdate(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := &TestUser{Name: "Original"}
	result := db.Create(user)
	assert.NoError(t, result.Error)

	// Test Save
	user.Name = "Updated"
	result = db.Save(user)
	assert.NoError(t, result.Error)

	var updatedUser TestUser
	result = db.First(&updatedUser, 1)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated", updatedUser.Name)

	// Test Updates
	result = db.Model(&TestUser{}).Where("id = ?", 1).Update("name", "Updated Again")
	assert.NoError(t, result.Error)

	result = db.First(&updatedUser, 1)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated Again", updatedUser.Name)
}

func TestGORMDelete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := &TestUser{Name: "Test User"}
	result := db.Create(user)
	assert.NoError(t, result.Error)

	// Test Delete
	result = db.Delete(&TestUser{}, 1)
	assert.NoError(t, result.Error)

	var deletedUser TestUser
	result = db.First(&deletedUser, 1)
	assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))

	// Test Delete with Where
	anotherUser := &TestUser{Name: "Another"}
	result = db.Create(anotherUser)
	assert.NoError(t, result.Error)
	
	result = db.Where("id = ?", 2).Delete(&TestUser{})
	assert.NoError(t, result.Error)

	result = db.First(&deletedUser, 2)
	assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func TestGORMPreload(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	location := &TestLocation{Name: "Test Loc"}
	result := db.Create(location)
	assert.NoError(t, result.Error)

	activity := &TestActivity{Title: "Test Act", LocationID: 1}
	result = db.Create(activity)
	assert.NoError(t, result.Error)

	var foundActivity TestActivity
	result = db.Preload("Location").First(&foundActivity, 1)
	assert.NoError(t, result.Error)
	// Note: Preload test is basic; in real models with associations, it would load related data
}

func TestGORMJoins(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// If setup failed due to database connection issues, skip the test
		t.Skip("Skipping test due to database setup failure")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Create test data
	user := &TestUser{Name: "Test User"}
	result := db.Create(user)
	assert.NoError(t, result.Error)

	var users []TestUser
	result = db.Joins("JOIN some_table ON some_condition").Find(&users)
	assert.NoError(t, result.Error) // Basic test for Joins method
}
