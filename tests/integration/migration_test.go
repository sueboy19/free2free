package integration

import (
	"testing"
	"time"

	"free2free/models"
	"free2free/tests/testutils"

	"github.com/stretchr/testify/assert"
)

// TestSchemaMigrationOperations tests schema migration operations with the new driver
func TestSchemaMigrationOperations(t *testing.T) {
	t.Run("Single Model Migration", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Test migration of a single model
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Verify the table was created by attempting to use it
		user := &models.User{
			Name:           "Migration Test User",
			Email:          "mig@example.com",
			SocialProvider: "facebook",
			SocialID:       "mig_123",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Verify the record was properly stored
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Migration Test User", retrievedUser.Name)
		assert.Equal(t, "mig@example.com", retrievedUser.Email)
	})

	t.Run("Multiple Model Migration", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Test migration of multiple models at once
		modelsToMigrate := []interface{}{
			&models.User{}, &models.Activity{}, &models.Location{},
		}

		for _, model := range modelsToMigrate {
			err := testutils.MigrateTestDB(db, model)
			assert.NoError(t, err)
		}

		// Test that all tables are functional
		// Create a location
		location := &models.Location{
			Name:      "Migration Test Location",
			Address:   "123 Migration St",
			Latitude:  25.0,
			Longitude: 121.0,
		}
		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Create a user
		user := &models.User{
			Name:           "Migration Test User",
			Email:          "mig2@example.com",
			SocialProvider: "facebook",
			SocialID:       "mig2_456",
		}
		result = db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Create an activity
		activity := &models.Activity{
			Title:       "Migration Test Activity",
			Description: "Activity for migration test",
			LocationID:  location.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Verify relationships work
		var retrievedActivity models.Activity
		result = db.Preload("Location").First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, location.ID, retrievedActivity.LocationID)
		assert.Equal(t, "Migration Test Location", retrievedActivity.Location.Name)
	})

	t.Run("Complex Model Migration", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Test migration of all application models
		allModels := []interface{}{
			&models.User{}, &models.Activity{}, &models.Location{},
			&models.Match{}, &models.MatchParticipant{}, &models.Review{},
			&models.ReviewLike{}, &models.Admin{}, &models.RefreshToken{},
		}

		for i, model := range allModels {
			err := testutils.MigrateTestDB(db, model)
			assert.NoError(t, err, "Migration should succeed for model %d: %T", i, model)
		}

		// Test operations across all models to ensure they're properly created
		// Create a user
		user := &models.User{
			Name:           "Complex Migration Test User",
			Email:          "complexmig@example.com",
			SocialProvider: "facebook",
			SocialID:       "cmig_789",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Create an admin
		admin := &models.Admin{
			Username: "complex_mig_admin",
			Email:    "admin@complexmig.com",
		}
		result = db.Create(admin)
		assert.NoError(t, result.Error)
		assert.NotZero(t, admin.ID)

		// Create a location
		location := &models.Location{
			Name:      "Complex Migration Test Location",
			Address:   "456 Complex Migration Ave",
			Latitude:  24.5,
			Longitude: 120.5,
		}
		result = db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Create an activity
		activity := &models.Activity{
			Title:       "Complex Migration Test Activity",
			Description: "Activity for complex migration test",
			LocationID:  location.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Create a match
		match := &models.Match{
			ActivityID: activity.ID,
			Status:     "open",
		}
		result = db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// Create a match participant
		participant := &models.MatchParticipant{
			MatchID:  match.ID,
			UserID:   user.ID,
			Status:   "approved",
			JoinedAt: time.Now(),
		}
		result = db.Create(participant)
		assert.NoError(t, result.Error)
		assert.NotZero(t, participant.ID)

		// Create a review
		review := &models.Review{
			MatchID:    match.ID,
			ReviewerID: user.ID,
			RevieweeID: user.ID,
			Score:      5,
			Comment:    "Great complex migration test!",
			CreatedAt:  time.Now(),
		}
		result = db.Create(review)
		assert.NoError(t, result.Error)
		assert.NotZero(t, review.ID)

		// Create a review like
		like := &models.ReviewLike{
			ReviewID: review.ID,
			UserID:   user.ID,
			IsLike:   true,
		}
		result = db.Create(like)
		assert.NoError(t, result.Error)
		assert.NotZero(t, like.ID)

		// Create a refresh token
		token := &models.RefreshToken{
			UserID:    uint(user.ID),
			Token:     "complex_migration_token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		}
		result = db.Create(token)
		assert.NoError(t, result.Error)
		assert.NotZero(t, token.ID)

		// Verify all relationships work
		var fullMatch models.Match
		result = db.Preload("Activity").Preload("Activity.Location").First(&fullMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Complex Migration Test Activity", fullMatch.Activity.Title)
		assert.Equal(t, "Complex Migration Test Location", fullMatch.Activity.Location.Name)

		var fullReview models.Review
		result = db.Preload("User").First(&fullReview, review.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Great complex migration test!", fullReview.Comment)
	})

	t.Run("Migration Idempotency", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Run migration for User model
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a user to verify the table works
		user := &models.User{
			Name:           "Idempotency Test User",
			Email:          "idem@example.com",
			SocialProvider: "facebook",
			SocialID:       "idem_001",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Run the same migration again - should be idempotent
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// The existing user should still be there
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Idempotency Test User", retrievedUser.Name)

		// Create another user to verify the table still works after re-migration
		user2 := &models.User{
			Name:           "Idempotency Test User 2",
			Email:          "idem2@example.com",
			SocialProvider: "instagram",
			SocialID:       "idem_002",
		}
		result = db.Create(user2)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user2.ID)

		// Verify both users exist
		var allUsers []models.User
		result = db.Find(&allUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(allUsers))
	})

	t.Run("Migration with Data Preservation", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Initial migration
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create some data
		user := &models.User{
			Name:           "Preservation Test User",
			Email:          "pres@example.com",
			SocialProvider: "facebook",
			SocialID:       "pres_001",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Re-run migration (which should preserve data in this context)
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Verify the original data is preserved
		var preservedUser models.User
		result = db.First(&preservedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Preservation Test User", preservedUser.Name)
		assert.Equal(t, "pres@example.com", preservedUser.Email)
		assert.Equal(t, "facebook", preservedUser.SocialProvider)
	})

	t.Run("Sequential Migration of Related Models", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Migrate related models in sequence (parent before child)
		err = testutils.MigrateTestDB(db, &models.Location{})
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{})
		assert.NoError(t, err)

		// Create a location
		location := &models.Location{
			Name:      "Sequential Migration Location",
			Address:   "789 Sequential St",
			Latitude:  23.5,
			Longitude: 119.5,
		}
		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Create an activity linked to the location
		activity := &models.Activity{
			Title:       "Sequential Migration Activity",
			Description: "Activity for sequential migration test",
			LocationID:  location.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Verify the relationship works
		var retrievedActivity models.Activity
		result = db.Preload("Location").First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, location.ID, retrievedActivity.LocationID)
		assert.Equal(t, "Sequential Migration Location", retrievedActivity.Location.Name)

		// Add more models in sequence
		err = testutils.MigrateTestDB(db, &models.Match{})
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.MatchParticipant{})
		assert.NoError(t, err)

		// Create match and participant
		match := &models.Match{
			ActivityID: activity.ID,
			Status:     "open",
		}
		result = db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		participant := &models.MatchParticipant{
			MatchID:  match.ID,
			UserID:   1, // Simple test user
			Status:   "approved",
			JoinedAt: time.Now(),
		}
		result = db.Create(participant)
		assert.NoError(t, result.Error)
		assert.NotZero(t, participant.ID)

		// Verify full chain works
		var fullMatch models.Match
		result = db.Preload("Activity").Preload("Activity.Location").First(&fullMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, activity.ID, fullMatch.ActivityID)
		assert.Equal(t, location.ID, fullMatch.Activity.LocationID)
	})
}
