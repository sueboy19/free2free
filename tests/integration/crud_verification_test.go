package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/models"
	"free2free/tests/testutils"
)

// TestCRUDOperationVerification tests CRUD operations with the platform-independent database
func TestCRUDOperationVerification(t *testing.T) {
	t.Run("User Model CRUD Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// CREATE
		user := &models.User{
			Name:       "CRUD Test User",
			Email:      "crud@example.com",
			Provider:   "facebook",
			ProviderID: "crud_123",
			Avatar:     "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// READ (single)
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "CRUD Test User", retrievedUser.Name)
		assert.Equal(t, "crud@example.com", retrievedUser.Email)

		// READ (multiple)
		var allUsers []models.User
		result = db.Find(&allUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(allUsers))

		// UPDATE
		user.Name = "Updated CRUD Test User"
		result = db.Save(user)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedUser models.User
		result = db.First(&updatedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated CRUD Test User", updatedUser.Name)

		// DELETE
		result = db.Delete(&models.User{}, user.ID)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedUser models.User
		result = db.First(&deletedUser, user.ID)
		assert.Error(t, result.Error)
	})

	t.Run("Activity Model CRUD Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{})
		assert.NoError(t, err)

		// CREATE
		activity := &models.Activity{
			Title:       "CRUD Test Activity",
			Description: "Activity for CRUD testing",
			Status:      "pending",
		}

		result := db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// READ (single)
		var retrievedActivity models.Activity
		result = db.First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "CRUD Test Activity", retrievedActivity.Title)

		// UPDATE
		activity.Status = "approved"
		result = db.Save(activity)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedActivity models.Activity
		result = db.First(&updatedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "approved", updatedActivity.Status)

		// DELETE
		result = db.Delete(&models.Activity{}, activity.ID)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedActivity models.Activity
		result = db.First(&deletedActivity, activity.ID)
		assert.Error(t, result.Error)
	})

	t.Run("Location Model CRUD Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Location{})
		assert.NoError(t, err)

		// CREATE
		location := &models.Location{
			Name:      "CRUD Test Location",
			Address:   "456 CRUD Ave",
			Latitude:  24.5,
			Longitude: 121.5,
		}

		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// READ
		var retrievedLocation models.Location
		result = db.First(&retrievedLocation, location.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "CRUD Test Location", retrievedLocation.Name)

		// UPDATE
		location.Latitude = 25.0
		result = db.Save(location)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedLocation models.Location
		result = db.First(&updatedLocation, location.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, 25.0, updatedLocation.Latitude)

		// DELETE
		result = db.Delete(&models.Location{}, location.ID)
		assert.NoError(t, result.Error)
	})

	t.Run("Match Model CRUD Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Match{})
		assert.NoError(t, err)

		// CREATE
		match := &models.Match{
			ActivityID: 1,
			Status:     "open",
		}

		result := db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// READ
		var retrievedMatch models.Match
		result = db.First(&retrievedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(1), retrievedMatch.ActivityID)

		// UPDATE
		match.Status = "closed"
		result = db.Save(match)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedMatch models.Match
		result = db.First(&updatedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "closed", updatedMatch.Status)

		// DELETE
		result = db.Delete(&models.Match{}, match.ID)
		assert.NoError(t, result.Error)
	})

	t.Run("Multiple Model Relationships", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Migrate related models
		err = testutils.MigrateTestDB(db, &models.User{}, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// CREATE related records
		// Create a location
		location := &models.Location{
			Name:      "Test Location for Relationships",
			Address:   "789 Relationship Blvd",
			Latitude:  24.0,
			Longitude: 120.0,
		}
		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Create a user
		user := &models.User{
			Name:       "Test User for Relationships",
			Email:      "rel@example.com",
			Provider:   "facebook",
			ProviderID: "rel_456",
		}
		result = db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Create an activity with the location
		activity := &models.Activity{
			Title:       "Test Activity with Relationships",
			Description: "Activity for relationship testing",
			Status:      "pending",
			LocationID:  location.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// READ with preloading to verify relationships
		var retrievedActivity models.Activity
		result = db.Preload("Location").First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Test Activity with Relationships", retrievedActivity.Title)
		assert.Equal(t, location.ID, retrievedActivity.LocationID)
		assert.Equal(t, "Test Location for Relationships", retrievedActivity.Location.Name)

		// UPDATE with relationships
		activity.Title = "Updated Test Activity with Relationships"
		result = db.Save(activity)
		assert.NoError(t, result.Error)

		var updatedActivity models.Activity
		result = db.Preload("Location").First(&updatedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated Test Activity with Relationships", updatedActivity.Title)

		// DELETE related records
		result = db.Delete(&models.Activity{}, activity.ID)
		assert.NoError(t, result.Error)
		result = db.Delete(&models.User{}, user.ID)
		assert.NoError(t, result.Error)
		result = db.Delete(&models.Location{}, location.ID)
		assert.NoError(t, result.Error)
	})

	t.Run("Complex Query Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create multiple records for complex query testing
		users := []models.User{
			{Name: "Alice", Email: "alice@example.com", Provider: "facebook", ProviderID: "a1"},
			{Name: "Bob", Email: "bob@example.com", Provider: "facebook", ProviderID: "b2"},
			{Name: "Charlie", Email: "charlie@example.com", Provider: "instagram", ProviderID: "c3"},
			{Name: "Diana", Email: "diana@example.com", Provider: "facebook", ProviderID: "d4"},
		}

		for _, user := range users {
			result := db.Create(&user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)
		}

		// Test complex queries
		var facebookUsers []models.User
		result := db.Where("provider = ?", "facebook").Find(&facebookUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 3, len(facebookUsers))

		var orderedUsers []models.User
		result = db.Order("name ASC").Find(&orderedUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 4, len(orderedUsers))
		assert.Equal(t, "Alice", orderedUsers[0].Name)
		assert.Equal(t, "Diana", orderedUsers[3].Name)

		var limitedUsers []models.User
		result = db.Limit(2).Offset(1).Find(&limitedUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(limitedUsers))

		var emailPatternUsers []models.User
		result = db.Where("email LIKE ?", "%@example.com").Find(&emailPatternUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 4, len(emailPatternUsers))
	})
}