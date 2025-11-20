package integration

import (
	"testing"

	"free2free/models"
	"free2free/tests/testutils"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestRegressionTesting tests for any regressions in database functionality
func TestRegressionTesting(t *testing.T) {
	t.Run("Full User Flow Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Migrate all relevant models
		modelsToMigrate := []interface{}{
			&models.User{}, &models.Activity{}, &models.Location{},
			&models.Match{}, &models.MatchParticipant{}, &models.Review{},
			&models.ReviewLike{}, &models.Admin{}, &models.RefreshToken{},
		}
		for _, model := range modelsToMigrate {
			err := testutils.MigrateTestDB(db, model)
			assert.NoError(t, err)
		}

		// Simulate a complete user flow that was working before
		// 1. Create user
		user := &models.User{
			Name:           "Regression Test User",
			Email:          "regression@example.com",
			SocialID:       "reg_123",
			SocialProvider: "facebook",
			AvatarURL:      "https://example.com/regression_avatar.jpg",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// 2. Create a location
		location := &models.Location{
			Name:      "Regression Test Location",
			Address:   "123 Regression St",
			Latitude:  25.0,
			Longitude: 121.0,
		}
		result = db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// 3. Create an activity
		activity := &models.Activity{
			Title:       "Regression Test Activity",
			Description: "Activity for regression testing",
			TargetCount: 2,
			LocationID:  location.ID,
			CreatedBy:   user.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// 4. Admin creates activity (no status field in Activity model)
		result = db.Save(activity)
		assert.NoError(t, result.Error)

		// 5. Create a match for the activity
		match := &models.Match{
			ActivityID: activity.ID,
			Status:     "open",
		}
		result = db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// 6. Add user as participant
		participant := &models.MatchParticipant{
			MatchID: match.ID,
			UserID:  user.ID,
			Status:  "confirmed",
		}
		result = db.Create(participant)
		assert.NoError(t, result.Error)
		assert.NotZero(t, participant.ID)

		// 7. Complete the match and create a review
		match.Status = "completed"
		result = db.Save(match)
		assert.NoError(t, result.Error)

		review := &models.Review{
			MatchID:    match.ID,
			ReviewerID: user.ID,
			RevieweeID: user.ID,
			Score:      5,
			Comment:    "Great regression test experience!",
		}
		result = db.Create(review)
		assert.NoError(t, result.Error)
		assert.NotZero(t, review.ID)

		// 8. Someone likes the review
		like := &models.ReviewLike{
			ReviewID: review.ID,
			UserID:   user.ID,
			IsLike:   true,
		}
		result = db.Create(like)
		assert.NoError(t, result.Error)
		assert.NotZero(t, like.ID)

		// 9. Verify all relationships are maintained
		var verifiedMatch models.Match
		result = db.Preload("Activity").First(&verifiedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "completed", verifiedMatch.Status)
		assert.Equal(t, "Regression Test Activity", verifiedMatch.Activity.Title)

		// Check participants separately
		var participants []models.MatchParticipant
		result = db.Where("match_id = ?", match.ID).Find(&participants)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(participants))

		var verifiedReview models.Review
		result = db.First(&verifiedReview, review.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Great regression test experience!", verifiedReview.Comment)

		// Check review likes separately
		var reviewLikes []models.ReviewLike
		result = db.Where("review_id = ?", review.ID).Find(&reviewLikes)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(reviewLikes))
	})

	t.Run("Transaction Integrity Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// Test transaction integrity to ensure no regression
		err = db.Transaction(func(tx *gorm.DB) error {
			// Create location in transaction
			location := &models.Location{
				Name:      "Transaction Test Location",
				Address:   "456 Transaction Ave",
				Latitude:  24.5,
				Longitude: 121.5,
			}
			result := tx.Create(location)
			assert.NoError(t, result.Error)
			assert.NotZero(t, location.ID)

			// Create activity in transaction
			activity := &models.Activity{
				Title:       "Transaction Test Activity",
				Description: "Activity for transaction regression test",
				TargetCount: 2,
				LocationID:  location.ID,
				CreatedBy:   1, // Mock admin ID
			}
			result = tx.Create(activity)
			assert.NoError(t, result.Error)
			assert.NotZero(t, activity.ID)

			// Query within transaction to ensure consistency
			var count int64
			err := tx.Model(&models.Location{}).Count(&count).Error
			assert.NoError(t, err)
			assert.Equal(t, int64(1), count)

			return nil // Commit transaction
		})

		assert.NoError(t, err)

		// Verify that both records exist after transaction
		var location models.Location
		result := db.First(&location)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Transaction Test Location", location.Name)

		var activity models.Activity
		result = db.First(&activity)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Transaction Test Activity", activity.Title)
		assert.Equal(t, location.ID, activity.LocationID)
	})

	t.Run("Complex Query Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create multiple users for query testing
		users := []models.User{
			{Name: "Regression Alice", Email: "reg.alice@example.com", SocialProvider: "facebook", SocialID: "ra1"},
			{Name: "Regression Bob", Email: "reg.bob@example.com", SocialProvider: "instagram", SocialID: "rb2"},
			{Name: "Regression Charlie", Email: "reg.charlie@example.com", SocialProvider: "facebook", SocialID: "rc3"},
			{Name: "Regression Diana", Email: "reg.diana@example.com", SocialProvider: "facebook", SocialID: "rd4"},
		}

		for _, user := range users {
			result := db.Create(&user)
			assert.NoError(t, result.Error)
			assert.NotZero(t, user.ID)
		}

		// Test complex queries that should work as before
		var facebookUsers []models.User
		result := db.Where("social_provider = ?", "facebook").Order("name ASC").Find(&facebookUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 3, len(facebookUsers))
		assert.Equal(t, "Regression Alice", facebookUsers[0].Name)

		var emailPatternUsers []models.User
		result = db.Where("email LIKE ?", "reg.%@example.com").Find(&emailPatternUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 4, len(emailPatternUsers))

		var limitedUsers []models.User
		result = db.Limit(2).Offset(1).Order("name ASC").Find(&limitedUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(limitedUsers))

		var allUsers []models.User
		result = db.Order("name DESC").Find(&allUsers)
		assert.NoError(t, result.Error)
		assert.Equal(t, 4, len(allUsers))
		assert.Equal(t, "Regression Diana", allUsers[0].Name)
	})

	t.Run("Migration and Schema Changes Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		// Test migrations work properly without regression
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Insert data with all expected fields
		user := &models.User{
			Name:           "Migration Regression User",
			Email:          "migreg@example.com",
			SocialProvider: "facebook",
			SocialID:       "migreg_789",
			AvatarURL:      "https://example.com/migreg_avatar.jpg",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Verify all fields are correctly stored and retrieved
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Migration Regression User", retrievedUser.Name)
		assert.Equal(t, "migreg@example.com", retrievedUser.Email)
		assert.Equal(t, "facebook", retrievedUser.SocialProvider)
		assert.Equal(t, "migreg_789", retrievedUser.SocialID)
		assert.Equal(t, "https://example.com/migreg_avatar.jpg", retrievedUser.AvatarURL)

		// Test updates work correctly
		retrievedUser.Name = "Updated Migration Regression User"
		result = db.Save(&retrievedUser)
		assert.NoError(t, result.Error)

		var updatedUser models.User
		result = db.First(&updatedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated Migration Regression User", updatedUser.Name)
	})

	t.Run("Data Integrity Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{})
		assert.NoError(t, err)

		// Test data integrity with various data types
		activity := &models.Activity{
			Title:       "Data Integrity Test Activity with Special Characters: !@#$%^&*()",
			Description: "Testing data integrity with various characters and ensuring no data corruption: 0123456789, abcdefghijklmnopqrstuvwxyz, ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			TargetCount: 2,
			LocationID:  1,
			CreatedBy:   1,
		}
		result := db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Retrieve and verify data integrity
		var retrievedActivity models.Activity
		result = db.First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, activity.Title, retrievedActivity.Title)
		assert.Equal(t, activity.Description, retrievedActivity.Description)

		// Test with longer text
		longDescription := "Very long description. "
		for i := 0; i < 50; i++ {
			longDescription += "This is a repeated sentence to make the text longer. "
		}
		activity.Description = longDescription
		result = db.Save(activity)
		assert.NoError(t, result.Error)

		var updatedActivity models.Activity
		result = db.First(&updatedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, longDescription, updatedActivity.Description)
	})

	t.Run("Relationship Integrity Regression", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{}, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// Create records with relationships
		location := &models.Location{
			Name:      "Relationship Test Location",
			Address:   "789 Relationship Blvd",
			Latitude:  24.0,
			Longitude: 120.0,
		}
		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		user := &models.User{
			Name:           "Relationship Test User",
			Email:          "reltest@example.com",
			SocialProvider: "facebook",
			SocialID:       "reltest_101",
		}
		result = db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		activity := &models.Activity{
			Title:       "Relationship Test Activity",
			Description: "Activity to test relationship integrity",
			TargetCount: 2,
			LocationID:  location.ID,
			CreatedBy:   user.ID,
		}
		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Test relationship queries
		var activitiesWithLocations []models.Activity
		result = db.Preload("Location").Find(&activitiesWithLocations)
		assert.NoError(t, result.Error)
		assert.Greater(t, len(activitiesWithLocations), 0)
		for _, act := range activitiesWithLocations {
			if act.ID == activity.ID {
				assert.Equal(t, location.Name, act.Location.Name)
				assert.Equal(t, location.Address, act.Location.Address)
			}
		}

		// Test joins
		var joinedResults []struct {
			ActivityID   uint
			Title        string
			LocationID   uint
			LocationName string
		}
		result = db.Raw("SELECT a.id as activity_id, a.title, a.location_id, l.name as location_name FROM activities a JOIN locations l ON a.location_id = l.id WHERE a.id = ?", activity.ID).Scan(&joinedResults)
		assert.NoError(t, result.Error)
		assert.Equal(t, 1, len(joinedResults))
		assert.Equal(t, activity.Title, joinedResults[0].Title)
		assert.Equal(t, location.Name, joinedResults[0].LocationName)
	})
}
