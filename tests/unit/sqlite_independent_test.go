package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/tests/testutils"
)

// TestSQLiteOperationsWithoutCGO tests SQLite operations without requiring CGO
func TestSQLiteOperationsWithoutCGO(t *testing.T) {
	t.Run("In-Memory Database Creation", func(t *testing.T) {
		// Test creating an in-memory database without CGO
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// Verify that the database works by executing a simple query
		err = db.Exec("SELECT 1").Error
		assert.NoError(t, err)
	})

	t.Run("Basic CRUD Operations", func(t *testing.T) {
		// Test basic CRUD operations work with platform-independent database
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Create a simple table
		type TestRecord struct {
			ID   uint
			Name string
		}

		err = testutils.MigrateTestDB(db, &TestRecord{})
		assert.NoError(t, err)

		// Create
		record := &TestRecord{Name: "Test Record"}
		result := db.Create(record)
		assert.NoError(t, result.Error)
		assert.NotZero(t, record.ID)

		// Read
		var retrievedRecord TestRecord
		result = db.First(&retrievedRecord, record.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, record.Name, retrievedRecord.Name)

		// Update
		retrievedRecord.Name = "Updated Record"
		result = db.Save(&retrievedRecord)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedRecord TestRecord
		result = db.First(&updatedRecord, record.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated Record", updatedRecord.Name)

		// Delete
		result = db.Delete(&TestRecord{}, record.ID)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedRecord TestRecord
		result = db.First(&deletedRecord, record.ID)
		assert.Error(t, result.Error)
	})

	t.Run("Transaction Operations", func(t *testing.T) {
		// Test transaction operations work with platform-independent database
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Create a simple table
		type TestRecord struct {
			ID   uint
			Name string
		}

		err = testutils.MigrateTestDB(db, &TestRecord{})
		assert.NoError(t, err)

		// Test transaction
		err = db.Transaction(func(tx *gorm.DB) error {
			// Create a record inside transaction
			record := &TestRecord{Name: "Transaction Record"}
			result := tx.Create(record)
			assert.NoError(t, result.Error)
			assert.NotZero(t, record.ID)

			// Read the record inside transaction
			var retrievedRecord TestRecord
			result = tx.First(&retrievedRecord, record.ID)
			assert.NoError(t, result.Error)
			assert.Equal(t, record.Name, retrievedRecord.Name)

			return nil // Commit the transaction
		})

		assert.NoError(t, err)

		// Verify the record exists after transaction
		var finalRecord TestRecord
		result := db.First(&finalRecord)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Transaction Record", finalRecord.Name)
	})

	t.Run("Multiple Connections", func(t *testing.T) {
		// Test multiple database connections work with platform-independent database
		db1, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		db2, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Create a simple table in both databases
		type TestRecord struct {
			ID   uint
			Name string
		}

		err = testutils.MigrateTestDB(db1, &TestRecord{})
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db2, &TestRecord{})
		assert.NoError(t, err)

		// Create records in both databases
		record1 := &TestRecord{Name: "DB1 Record"}
		result1 := db1.Create(record1)
		assert.NoError(t, result1.Error)
		assert.NotZero(t, record1.ID)

		record2 := &TestRecord{Name: "DB2 Record"}
		result2 := db2.Create(record2)
		assert.NoError(t, result2.Error)
		assert.NotZero(t, record2.ID)

		// Verify records exist in respective databases
		var checkRecord1 TestRecord
		result := db1.First(&checkRecord1, record1.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "DB1 Record", checkRecord1.Name)

		var checkRecord2 TestRecord
		result = db2.First(&checkRecord2, record2.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "DB2 Record", checkRecord2.Name)
	})
}