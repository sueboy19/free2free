package performance

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"free2free/models"
	"free2free/tests/testutils"
)

// TestPerformanceComparison tests performance comparison between implementations
func TestPerformanceComparison(t *testing.T) {
	t.Run("Single Operation Performance", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Measure create operation performance
		startTime := time.Now()
		user := &models.User{
			Name:       "Perf Test User",
			Email:      "perf@example.com",
			Provider:   "facebook",
			ProviderID: "perf_123",
		}
		result := db.Create(user)
		createDuration := time.Since(startTime)

		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Measure read operation performance
		startTime = time.Now()
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		readDuration := time.Since(startTime)

		assert.NoError(t, result.Error)
		assert.Equal(t, "Perf Test User", retrievedUser.Name)

		// Measure update operation performance
		startTime = time.Now()
		user.Name = "Updated Perf Test User"
		result = db.Save(user)
		updateDuration := time.Since(startTime)

		assert.NoError(t, result.Error)

		// Measure delete operation performance
		startTime = time.Now()
		result = db.Delete(&models.User{}, user.ID)
		deleteDuration := time.Since(startTime)

		assert.NoError(t, result.Error)

		// Log performance metrics
		t.Logf("Single Operation Performance:")
		t.Logf("  Create: %v", createDuration)
		t.Logf("  Read: %v", readDuration)
		t.Logf("  Update: %v", updateDuration)
		t.Logf("  Delete: %v", deleteDuration)

		// Verify performance is within acceptable limits (100ms for each operation as an example)
		assert.Less(t, createDuration, 100*time.Millisecond, "Create operation should be faster than 100ms")
		assert.Less(t, readDuration, 100*time.Millisecond, "Read operation should be faster than 100ms")
		assert.Less(t, updateDuration, 100*time.Millisecond, "Update operation should be faster than 100ms")
		assert.Less(t, deleteDuration, 100*time.Millisecond, "Delete operation should be faster than 100ms")
	})

	t.Run("Batch Operation Performance", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create multiple users in batch
		users := make([]*models.User, 100)
		for i := 0; i < 100; i++ {
			users[i] = &models.User{
				Name:       "Batch User " + string(rune('0'+i)),
				Email:      "batch" + string(rune('0'+i)) + "@example.com",
				Provider:   "facebook",
				ProviderID: "batch_" + string(rune('0'+i)),
			}
		}

		startTime := time.Now()
		for _, user := range users {
			result := db.Create(user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)
		}
		batchCreateDuration := time.Since(startTime)

		// Read all users back
		startTime = time.Now()
		var allUsers []models.User
		result := db.Find(&allUsers)
		batchReadDuration := time.Since(startTime)

		assert.NoError(t, result.Error)
		assert.Equal(t, 100, len(allUsers))

		// Update all users
		startTime = time.Now()
		for i := range allUsers {
			allUsers[i].Name = "Updated " + allUsers[i].Name
			result := db.Save(&allUsers[i])
			assert.NoError(t, result.Error)
		}
		batchUpdateDuration := time.Since(startTime)

		// Delete all users
		startTime = time.Now()
		for _, user := range users {
			result := db.Delete(&models.User{}, user.ID)
			assert.NoError(t, result.Error)
		}
		batchDeleteDuration := time.Since(startTime)

		// Log performance metrics
		t.Logf("Batch Operation Performance (100 records):")
		t.Logf("  Batch Create: %v", batchCreateDuration)
		t.Logf("  Batch Read: %v", batchReadDuration)
		t.Logf("  Batch Update: %v", batchUpdateDuration)
		t.Logf("  Batch Delete: %v", batchDeleteDuration)

		// Verify batch performance is within acceptable limits
		assert.Less(t, batchCreateDuration, 2*time.Second, "Batch create should be faster than 2 seconds")
		assert.Less(t, batchReadDuration, 500*time.Millisecond, "Batch read should be faster than 500ms")
		assert.Less(t, batchUpdateDuration, 2*time.Second, "Batch update should be faster than 2 seconds")
		assert.Less(t, batchDeleteDuration, 1*time.Second, "Batch delete should be faster than 1 second")
	})

	t.Run("Transaction Performance", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{}, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// Measure transaction performance with multiple operations
		startTime := time.Now()
		err = db.Transaction(func(tx *gorm.DB) error {
			// Create a location
			location := &models.Location{
				Name:      "Transaction Perf Test Location",
				Address:   "123 Transaction Perf St",
				Latitude:  25.0,
				Longitude: 121.0,
			}
			result := tx.Create(location)
			assert.NoError(t, result.Error)
			assert.NotZero(t, location.ID)

			// Create an activity
			activity := &models.Activity{
				Title:       "Transaction Perf Test Activity",
				Description: "Activity for transaction performance test",
				Status:      "pending",
				LocationID:  location.ID,
			}
			result = tx.Create(activity)
			assert.NoError(t, result.Error)
			assert.NotZero(t, activity.ID)

			// Create a user
			user := &models.User{
				Name:       "Transaction Perf Test User",
				Email:      "txperf@example.com",
				Provider:   "facebook",
				ProviderID: "txperf_456",
			}
			result = tx.Create(user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)

			return nil // Commit transaction
		})
		transactionDuration := time.Since(startTime)

		assert.NoError(t, err)

		// Log transaction performance
		t.Logf("Transaction Performance:")
		t.Logf("  Multi-operation transaction: %v", transactionDuration)

		// Verify transaction performance is within acceptable limits
		assert.Less(t, transactionDuration, 500*time.Millisecond, "Transaction should be faster than 500ms")
	})

	t.Run("Complex Query Performance", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create test data
		for i := 0; i < 50; i++ {
			user := &models.User{
				Name:       "Query Perf Test User " + string(rune('0'+(i%10))),
				Email:      "queryperf" + string(rune('0'+i)) + "@example.com",
				Provider:   "facebook",
				ProviderID: "qperf_" + string(rune('0'+i)),
			}
			result := db.Create(user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)
		}

		// Measure complex query performance
		startTime := time.Now()
		var facebookUsers []models.User
		result := db.Where("provider = ?", "facebook").Order("name ASC").Find(&facebookUsers)
		queryDuration := time.Since(startTime)

		assert.NoError(t, result.Error)
		assert.Greater(t, len(facebookUsers), 0)

		// Measure join query performance
		startTime = time.Now()
		var joinedResults []struct {
			UserID   uint
			UserName string
		}
		result = db.Raw("SELECT id as user_id, name as user_name FROM users WHERE provider = ? ORDER BY name", "facebook").Scan(&joinedResults)
		joinQueryDuration := time.Since(startTime)

		assert.NoError(t, result.Error)
		assert.Greater(t, len(joinedResults), 0)

		// Log query performance
		t.Logf("Complex Query Performance:")
		t.Logf("  Filtered query: %v", queryDuration)
		t.Logf("  Raw join query: %v", joinQueryDuration)

		// Verify query performance is within acceptable limits
		assert.Less(t, queryDuration, 200*time.Millisecond, "Filtered query should be faster than 200ms")
		assert.Less(t, joinQueryDuration, 200*time.Millisecond, "Raw join query should be faster than 200ms")
	})

	t.Run("Concurrent Operation Performance", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Measure concurrent operation performance
		const numGoroutines = 10
		const operationsPerGoroutine = 10

		startTime := time.Now()

		// Run multiple goroutines performing operations
		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					user := &models.User{
						Name:       "Concurrent User " + string(rune('0'+goroutineID)) + "-" + string(rune('0'+j)),
						Email:      "conc" + string(rune('0'+goroutineID)) + "-" + string(rune('0'+j)) + "@example.com",
						Provider:   "facebook",
						ProviderID: "conc_" + string(rune('0'+goroutineID)) + "_" + string(rune('0'+j)),
					}
					result := db.Create(user)
					assert.NoError(t, result.Error)
					assert.NotZero(t, user.ID)

					var retrievedUser models.User
					result = db.First(&retrievedUser, user.ID)
					assert.NoError(t, result.Error)
				}
			}(i)
		}

		wg.Wait()
		concurrentDuration := time.Since(startTime)

		// Verify all records were created
		var allUsers []models.User
		result := db.Find(&allUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, numGoroutines*operationsPerGoroutine, len(allUsers))

		// Log concurrent performance
		t.Logf("Concurrent Operation Performance:")
		t.Logf("  %d goroutines x %d operations each = %d operations total: %v", 
			numGoroutines, operationsPerGoroutine, numGoroutines*operationsPerGoroutine, concurrentDuration)

		// Verify concurrent performance is within acceptable limits
		assert.Less(t, concurrentDuration, 3*time.Second, "Concurrent operations should be faster than 3 seconds")
	})

	t.Run("Memory Usage During Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a moderate amount of data to measure memory usage
		const numRecords = 1000
		
		// Measure memory before operations
		var memStatsBefore runtime.MemStats
		runtime.ReadMemStats(&memStatsBefore)

		// Create records
		for i := 0; i < numRecords; i++ {
			user := &models.User{
				Name:       "Memory Test User " + string(rune('0'+(i%100))),
				Email:      "memtest" + string(rune('0'+i)) + "@example.com",
				Provider:   "facebook",
				ProviderID: "memtest_" + string(rune('0'+i)),
			}
			result := db.Create(user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)
		}

		// Measure memory after creation
		var memStatsAfterCreate runtime.MemStats
		runtime.ReadMemStats(&memStatsAfterCreate)

		// Read all records
		var allUsers []models.User
		result := db.Find(&allUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, numRecords, len(allUsers))

		// Measure memory after read
		var memStatsAfterRead runtime.MemStats
		runtime.ReadMemStats(&memStatsAfterRead)

		// Log memory usage
		t.Logf("Memory Usage:")
		t.Logf("  Before operations: %d KB", memStatsBefore.Alloc/1024)
		t.Logf("  After create (%d records): %d KB", numRecords, memStatsAfterCreate.Alloc/1024)
		t.Logf("  After read (%d records): %d KB", numRecords, memStatsAfterRead.Alloc/1024)
		t.Logf("  Memory growth from create: %d KB", (memStatsAfterCreate.Alloc-memStatsBefore.Alloc)/1024)
		t.Logf("  Memory growth from read: %d KB", (memStatsAfterRead.Alloc-memStatsAfterCreate.Alloc)/1024)

		// Verify memory usage is reasonable (less than 10MB growth for 1000 records)
		assert.Less(t, memStatsAfterRead.Alloc-memStatsBefore.Alloc, uint64(10*1024*1024), 
			"Memory usage should grow less than 10MB for 1000 records")
	})
}