package testutils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Use modernc.org/sqlite as the underlying driver (no CGO required)
// This import should be processed before any other SQLite driver imports
import _ "modernc.org/sqlite"

// CreateTestDB creates an in-memory SQLite database for testing
func CreateTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table name
		},
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// MigrateTestDB applies all migrations to the test database
func MigrateTestDB(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}
