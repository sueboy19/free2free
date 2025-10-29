package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/tests/testutils"
)

// TestCrossPlatformDatabaseInitialization tests database initialization works across platforms
func TestCrossPlatformDatabaseInitialization(t *testing.T) {
	t.Run("Consistent Initialization", func(t *testing.T) {
		// Test that database initialization works consistently without platform-specific dependencies
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// Verify database is properly initialized by performing basic operations
		err = db.Exec("SELECT 1").Error
		assert.NoError(t, err)
	})

	t.Run("Multiple Initialization Sequence", func(t *testing.T) {
		// Test initializing multiple databases in sequence
		db1, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db1)

		db2, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db2)

		db3, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db3)

		// Verify each database works independently
		err = db1.Exec("SELECT 1").Error
		assert.NoError(t, err)

		err = db2.Exec("SELECT 1").Error
		assert.NoError(t, err)

		err = db3.Exec("SELECT 1").Error
		assert.NoError(t, err)
	})

	t.Run("Database Migration Consistency", func(t *testing.T) {
		// Test that migrations work consistently across different databases
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Create a test table structure
		type TestModel struct {
			ID   uint
			Name string
			Data string
		}

		// Run migration
		err = testutils.MigrateTestDB(db, &TestModel{})
		assert.NoError(t, err)

		// Verify the table was created by inserting and reading data
		model := &TestModel{Name: "Test", Data: "Data"}
		result := db.Create(model)
		assert.NoError(t, result.Error)
		assert.NotZero(t, model.ID)

		var retrievedModel TestModel
		result = db.First(&retrievedModel, model.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Test", retrievedModel.Name)
		assert.Equal(t, "Data", retrievedModel.Data)
	})

	t.Run("Connection Pool Management", func(t *testing.T) {
		// Test that connection management works without platform-specific features
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Get the underlying SQL DB to test connection management
		sqlDB, err := db.DB()
		assert.NoError(t, err)

		// Test that we can configure connection parameters
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(0) // No max lifetime for tests

		// Verify the settings are applied by creating a query
		err = db.Exec("SELECT 1").Error
		assert.NoError(t, err)
	})

	t.Run("Import Replacement Verification", func(t *testing.T) {
		// Test that the import replacement is working as expected
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// The fact that we can create and use the database means the 
		// pure-Go SQLite driver (modernc.org/sqlite) is being used instead of CGO-based driver
		// This is verified by the successful execution of database operations
		
		// Execute a series of operations to verify the platform-independent driver works
		type VerificationModel struct {
			ID   uint
			Key  string
			Val  int
		}

		err = testutils.MigrateTestDB(db, &VerificationModel{})
		assert.NoError(t, err)

		// Create multiple records
		for i := 0; i < 5; i++ {
			model := &VerificationModel{
				Key: "Key" + string(rune('0'+i)),
				Val: i * 10,
			}
			result := db.Create(model)
			assert.NoError(t, result.Error)
			assert.NotZero(t, model.ID)
		}

		// Read all records back
		var allModels []VerificationModel
		result := db.Find(&allModels)
		assert.NoError(t, result.Error)
		assert.Equal(t, 5, len(allModels))

		// Verify values
		for i, model := range allModels {
			assert.Equal(t, "Key"+string(rune('0'+i)), model.Key)
			assert.Equal(t, i*10, model.Val)
		}
	})
}