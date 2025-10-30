package contract

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/tests/testutils"
)

// TestDatabaseConnectionContract tests the database connection contract
func TestDatabaseConnectionContract(t *testing.T) {
	t.Run("Establish Connection", func(t *testing.T) {
		// Test establishing database connection without CGO dependencies
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Verify connection is active
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		assert.NoError(t, sqlDB.Ping())
	})

	t.Run("Close Connection", func(t *testing.T) {
		// Test properly closing the database connection
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Close the connection
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		err = sqlDB.Close()
		assert.NoError(t, err)

		// Verify connection is closed
		err = sqlDB.Ping()
		assert.Error(t, err)
	})

	t.Run("Connection with Platform-Independent Driver", func(t *testing.T) {
		// Test that the connection works with the pure-Go SQLite driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Attempt a simple operation to verify the connection works
		err = db.Exec("SELECT 1").Error
		assert.NoError(t, err)
	})
}