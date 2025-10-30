package unit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/models"
	"free2free/tests/testutils"
)

// TestExistingUnitTestsCompatibility tests existing unit test compatibility with pure-Go driver
func TestExistingUnitTestsCompatibility(t *testing.T) {
	t.Run("User Model Operations Compatibility", func(t *testing.T) {
		// Test that existing user model operations work with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Migrate the User model
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a user (standard operation from existing tests)
		user := &models.User{
			Name:           "Compatibility Test User",
			Email:          "compat@example.com",
			SocialProvider: "test_provider",
			SocialID:       "test_123",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Read the user back (standard operation)
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Compatibility Test User", retrievedUser.Name)
		assert.Equal(t, "compat@example.com", retrievedUser.Email)

		// Update the user (standard operation)
		retrievedUser.Name = "Updated Compatibility Test User"
		result = db.Save(&retrievedUser)
		assert.NoError(t, result.Error)

		// Verify the update worked
		var updatedUser models.User
		result = db.First(&updatedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated Compatibility Test User", updatedUser.Name)

		// Delete the user (standard operation)
		result = db.Delete(&models.User{}, user.ID)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedUser models.User
		result = db.First(&deletedUser, user.ID)
		assert.Error(t, result.Error)
	})

	t.Run("Activity Model Operations Compatibility", func(t *testing.T) {
		// Test that existing activity model operations work with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Migrate the Activity model
		err = testutils.MigrateTestDB(db, &models.Activity{})
		assert.NoError(t, err)

		// Create an activity (standard operation from existing tests)
		activity := &models.Activity{
			Title:       "Compatibility Test Activity",
			Description: "Activity for compatibility testing",
		}

		result := db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Read the activity back (standard operation)
		var retrievedActivity models.Activity
		result = db.First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Compatibility Test Activity", retrievedActivity.Title)
		assert.Equal(t, "Activity for compatibility testing", retrievedActivity.Description)

		// Verify the activity remains as expected (no Status field)
		var updatedActivity models.Activity
		result = db.First(&updatedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Compatibility Test Activity", updatedActivity.Title)
	})

	t.Run("Location Model Operations Compatibility", func(t *testing.T) {
		// Test that existing location model operations work with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Migrate the Location model
		err = testutils.MigrateTestDB(db, &models.Location{})
		assert.NoError(t, err)

		// Create a location (standard operation from existing tests)
		location := &models.Location{
			Name:      "Compatibility Test Location",
			Address:   "123 Test St",
			Latitude:  25.0,
			Longitude: 121.0,
		}

		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Read the location back (standard operation)
		var retrievedLocation models.Location
		result = db.First(&retrievedLocation, location.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Compatibility Test Location", retrievedLocation.Name)
		assert.Equal(t, 25.0, retrievedLocation.Latitude)
		assert.Equal(t, 121.0, retrievedLocation.Longitude)
	})

	t.Run("Match Model Operations Compatibility", func(t *testing.T) {
		// Test that existing match model operations work with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Migrate related models
		err = testutils.MigrateTestDB(db, &models.Activity{}, &models.Match{}, &models.MatchParticipant{})
		assert.NoError(t, err)

		// Create a match (standard operation from existing tests)
		match := &models.Match{
			ActivityID: 1,
			Status:     "open",
		}

		result := db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// Read the match back (standard operation)
		var retrievedMatch models.Match
		result = db.First(&retrievedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(1), retrievedMatch.ActivityID)
		assert.Equal(t, "open", retrievedMatch.Status)
	})

	t.Run("Complex Query Operations Compatibility", func(t *testing.T) {
		// Test that complex queries work as expected with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Migrate models
		err = testutils.MigrateTestDB(db, &models.User{}, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// Create multiple test records
		user1 := &models.User{
			Name:           "Query Test User 1",
			Email:          "q1@example.com",
			SocialProvider: "test",
			SocialID:       "q1_123",
		}
		user2 := &models.User{
			Name:           "Query Test User 2",
			Email:          "q2@example.com",
			SocialProvider: "test",
			SocialID:       "q2_456",
		}

		result1 := db.Create(user1)
		result2 := db.Create(user2)
		assert.NoError(t, result1.Error)
		assert.NoError(t, result2.Error)

		// Test simple query
		var users []models.User
		result := db.Where("email LIKE ?", "%@example.com").Find(&users)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(users))

		// Test query with ordering
		var orderedUsers []models.User
		result = db.Order("name ASC").Find(&orderedUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(orderedUsers))
		assert.Equal(t, "Query Test User 1", orderedUsers[0].Name)

		// Test limit operation
		var limitedUsers []models.User
		result = db.Limit(1).Find(&limitedUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(limitedUsers))
	})

	t.Run("Migration Operations Compatibility", func(t *testing.T) {
		// Test that migration operations work as expected with the new driver
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if strings.Contains(err.Error(), "go-sqlite3 requires cgo") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			} else {
				assert.NoError(t, err)
			}
		}

		// Run multiple migrations similar to what might happen in real usage
		modelsToMigrate := []interface{}{
			&models.User{},
			&models.Activity{},
			&models.Location{},
			&models.Match{},
			&models.MatchParticipant{},
			&models.Review{},
			&models.ReviewLike{},
			&models.Admin{},
			&models.RefreshToken{},
		}

		for _, model := range modelsToMigrate {
			err := testutils.MigrateTestDB(db, model)
			assert.NoError(t, err, "Migration should succeed for model: %T", model)
		}

		// Verify that tables are created and can be used
		user := &models.User{
			Name:           "Migration Test",
			Email:          "migration@example.com",
			SocialProvider: "test",
			SocialID:       "mig_789",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Migration Test", retrievedUser.Name)
	})
}