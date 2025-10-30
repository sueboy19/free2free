package testutils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Use modernc.org/sqlite as the underlying driver (no CGO required)
// This import should be processed before any other SQLite driver imports
// It registers the pure-Go SQLite driver to replace the CGO-based one
import _ "modernc.org/sqlite"

// CreateTestDB creates an in-memory SQLite database for testing
// Uses modernc.org/sqlite pure-Go driver (no CGO required)
func CreateTestDB() (*gorm.DB, error) {
	// We use gorm.io/driver/sqlite but with the modernc.org/sqlite driver registered,
	// which should now be used instead of the CGO-based github.com/mattn/go-sqlite3
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
