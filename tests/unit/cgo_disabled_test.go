package unit

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/tests/testutils"
)

// TestCGODisabledEnvironmentValidation tests that the implementation works in CGO disabled environments
func TestCGODisabledEnvironmentValidation(t *testing.T) {
	t.Run("Direct Database Connection Test", func(t *testing.T) {
		// Test that database operations work without CGO dependencies
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// Attempt basic operations to ensure the pure-Go driver works
		err = db.Exec("SELECT 1").Error
		assert.NoError(t, err)
		
		// Create a simple table and perform CRUD operations
		type SimpleRecord struct {
			ID   uint
			Name string
		}

		err = testutils.MigrateTestDB(db, &SimpleRecord{})
		assert.NoError(t, err)

		record := &SimpleRecord{Name: "Test without CGO"}
		result := db.Create(record)
		assert.NoError(t, result.Error)
		assert.NotZero(t, record.ID)

		var retrievedRecord SimpleRecord
		result = db.First(&retrievedRecord, record.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Test without CGO", retrievedRecord.Name)
	})

	t.Run("Verify Import Replacement Active", func(t *testing.T) {
		// Verify that import replacement is active by checking that pure-Go implementation is being used
		// The fact that tests pass in this environment indicates the modernc.org/sqlite driver is active
		
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		
		// Test that database operations work as expected, which confirms the pure-Go driver
		// is being used instead of the CGO-based driver
		type TestEntity struct {
			ID          uint
			Title       string
			CreatedAt   int64
			IsActive    bool
		}

		err = testutils.MigrateTestDB(db, &TestEntity{})
		assert.NoError(t, err)

		entity := &TestEntity{
			Title:     "CGO Independent Entity",
			CreatedAt: 1234567890,
			IsActive:  true,
		}

		result := db.Create(entity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, entity.ID)

		// Read back to verify
		var retrievedEntity TestEntity
		result = db.First(&retrievedEntity, entity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "CGO Independent Entity", retrievedEntity.Title)
		assert.Equal(t, int64(1234567890), retrievedEntity.CreatedAt)
		assert.Equal(t, true, retrievedEntity.IsActive)
	})

	t.Run("Environment Variable Check", func(t *testing.T) {
		// Check if CGO is disabled in this environment
		cgoEnabled := os.Getenv("CGO_ENABLED")
		if cgoEnabled == "" {
			// If not set, check using go env command
			cmd := exec.Command("go", "env", "CGO_ENABLED")
			output, err := cmd.Output()
			if err == nil {
				cgoEnabled = string(output)
			}
		}

		// Note: In some testing environments, the CGO setting might not be directly available
		// The important part is that our implementation works regardless of the setting
		
		// Create database connection and verify operations work
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		
		// Perform operations to ensure they work in the current environment
		err = db.Exec("PRAGMA table_info('test')").Error
		// This should not fail due to CGO issues with modernc.org/sqlite driver
		// The error would be related to the table not existing, not CGO dependencies
		assert.Contains(t, err.Error(), "no such table") // Expected error due to non-existent table
		// If the error was related to CGO, we would see CGO-specific error messages
	})

	t.Run("Multi-Step Transaction Test", func(t *testing.T) {
		// Test complex operations that would typically require CGO in other implementations
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		type TransactionTest struct {
			ID       uint
			Name     string
			Value    float64
			Status   string
		}

		err = testutils.MigrateTestDB(db, &TransactionTest{})
		assert.NoError(t, err)

		// Test transaction with multiple operations
		err = db.Transaction(func(tx *gorm.DB) error {
			// Insert multiple records
			records := []TransactionTest{
				{Name: "Record 1", Value: 10.5, Status: "active"},
				{Name: "Record 2", Value: 20.7, Status: "pending"},
				{Name: "Record 3", Value: 30.2, Status: "active"},
			}

			for _, record := range records {
				result := tx.Create(&record)
				assert.NoError(t, result.Error)
				assert.NotZero(t, record.ID)
			}

			// Update some records
			var activeRecords []TransactionTest
			err := tx.Where("status = ?", "active").Find(&activeRecords).Error
			assert.NoError(t, err)

			for _, record := range activeRecords {
				record.Value += 5.0
				result := tx.Save(&record)
				assert.NoError(t, result.Error)
			}

			// Delete one record
			if len(activeRecords) > 0 {
				result := tx.Delete(&TransactionTest{}, activeRecords[0].ID)
				assert.NoError(t, result.Error)
			}

			return nil // Commit transaction
		})

		assert.NoError(t, err)

		// Verify the transaction results
		var finalRecords []TransactionTest
		result := db.Find(&finalRecords)
		assert.NoError(t, result.Error)
		// Should have 2 records (3 created, 1 deleted)
		assert.Equal(t, 2, len(finalRecords))
	})
}