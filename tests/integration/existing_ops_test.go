package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/models"
	"free2free/tests/testutils"
)

// TestExistingDatabaseOperations tests that all existing database operations work identically
func TestExistingDatabaseOperations(t *testing.T) {
	t.Run("User Authentication Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Test creating a user with all fields (as would happen in auth flow)
		user := &models.User{
			Name:       "Existing Op Test User",
			Email:      "existing@example.com",
			Provider:   "facebook",
			ProviderID: "existing_123",
			Avatar:     "https://example.com/existing_avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Test finding user by provider and provider ID (common auth operation)
		var foundUser models.User
		result = db.Where("provider = ? AND provider_id = ?", "facebook", "existing_123").First(&foundUser)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Existing Op Test User", foundUser.Name)
		assert.Equal(t, "existing@example.com", foundUser.Email)

		// Test updating user information (common in profile updates)
		foundUser.Name = "Updated Existing Op Test User"
		result = db.Save(&foundUser)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedUser models.User
		result = db.First(&updatedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated Existing Op Test User", updatedUser.Name)
	})

	t.Run("Activity Lifecycle Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{}, &models.Location{})
		assert.NoError(t, err)

		// Create a location first
		location := &models.Location{
			Name:      "Test Location",
			Address:   "123 Test St",
			Latitude:  25.0,
			Longitude: 121.0,
		}
		result := db.Create(location)
		assert.NoError(t, result.Error)
		assert.NotZero(t, location.ID)

		// Create an activity (normal creation flow)
		activity := &models.Activity{
			Title:       "Existing Op Test Activity",
			Description: "Activity for existing operations test",
			Status:      "pending",
			LocationID:  location.ID,
		}

		result = db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Test reading activity with location (using preload)
		var retrievedActivity models.Activity
		result = db.Preload("Location").First(&retrievedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Existing Op Test Activity", retrievedActivity.Title)
		assert.Equal(t, location.ID, retrievedActivity.LocationID)
		assert.Equal(t, "Test Location", retrievedActivity.Location.Name)

		// Test updating activity status (admin workflow)
		retrievedActivity.Status = "approved"
		result = db.Save(&retrievedActivity)
		assert.NoError(t, result.Error)

		// Verify status update
		var updatedActivity models.Activity
		result = db.First(&updatedActivity, activity.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "approved", updatedActivity.Status)
	})

	t.Run("Admin Management Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.User{}, &models.Admin{})
		assert.NoError(t, err)

		// Create a user
		user := &models.User{
			Name:       "Admin Test User",
			Email:      "admin@example.com",
			Provider:   "facebook",
			ProviderID: "admin_456",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Create an admin entry
		admin := &models.Admin{
			UserID:      user.ID,
			Permissions: "read,write,delete",
		}

		result = db.Create(admin)
		assert.NoError(t, result.Error)
		assert.NotZero(t, admin.ID)

		// Test finding admin by user ID
		var foundAdmin models.Admin
		result = db.Where("user_id = ?", user.ID).First(&foundAdmin)
		assert.NoError(t, result.Error)
		assert.Equal(t, user.ID, foundAdmin.UserID)

		// Test updating admin permissions
		foundAdmin.Permissions = "read,write"
		result = db.Save(&foundAdmin)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedAdmin models.Admin
		result = db.First(&updatedAdmin, admin.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "read,write", updatedAdmin.Permissions)
	})

	t.Run("Match Creation and Management", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Activity{}, &models.Match{}, &models.MatchParticipant{})
		assert.NoError(t, err)

		// Create an activity first
		activity := &models.Activity{
			Title:       "Match Test Activity",
			Description: "Activity for match creation test",
			Status:      "active",
		}
		result := db.Create(activity)
		assert.NoError(t, result.Error)
		assert.NotZero(t, activity.ID)

		// Create a match for the activity
		match := &models.Match{
			ActivityID: activity.ID,
			Status:     "open",
		}
		result = db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// Add a participant to the match
		participant := &models.MatchParticipant{
			MatchID: match.ID,
			Status:  "confirmed",
		}
		result = db.Create(participant)
		assert.NoError(t, result.Error)
		assert.NotZero(t, participant.ID)

		// Test retrieving match with participants
		var retrievedMatch models.Match
		result = db.Preload("Participants").First(&retrievedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, activity.ID, retrievedMatch.ActivityID)
		assert.Equal(t, "open", retrievedMatch.Status)
		assert.Equal(t, 1, len(retrievedMatch.Participants))

		// Test updating match status
		retrievedMatch.Status = "closed"
		result = db.Save(&retrievedMatch)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedMatch models.Match
		result = db.First(&updatedMatch, match.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "closed", updatedMatch.Status)
	})

	t.Run("Review and ReviewLike Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.Match{}, &models.Review{}, &models.ReviewLike{})
		assert.NoError(t, err)

		// Create a match first
		match := &models.Match{
			ActivityID: 1,
			Status:     "completed",
		}
		result := db.Create(match)
		assert.NoError(t, result.Error)
		assert.NotZero(t, match.ID)

		// Create a review for the match
		review := &models.Review{
			MatchID:   match.ID,
			Score:     5,
			Comment:   "Great experience!",
		}
		result = db.Create(review)
		assert.NoError(t, result.Error)
		assert.NotZero(t, review.ID)

		// Create a like for the review
		like := &models.ReviewLike{
			ReviewID: review.ID,
			Score:    1, // Like
		}
		result = db.Create(like)
		assert.NoError(t, result.Error)
		assert.NotZero(t, like.ID)

		// Test retrieving review with likes count
		var retrievedReview models.Review
		result = db.Preload("ReviewLikes").First(&retrievedReview, review.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, match.ID, retrievedReview.MatchID)
		assert.Equal(t, "Great experience!", retrievedReview.Comment)
		assert.Equal(t, 1, len(retrievedReview.ReviewLikes))

		// Test updating review
		retrievedReview.Comment = "Updated review comment"
		result = db.Save(&retrievedReview)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedReview models.Review
		result = db.First(&updatedReview, review.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated review comment", updatedReview.Comment)
	})

	t.Run("Refresh Token Operations", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		assert.NoError(t, err)

		err = testutils.MigrateTestDB(db, &models.RefreshToken{})
		assert.NoError(t, err)

		// Create a refresh token
		token := &models.RefreshToken{
			UserID: 1,
			Token:  "refresh_token_value",
		}
		result := db.Create(token)
		assert.NoError(t, result.Error)
		assert.NotZero(t, token.ID)

		// Test finding by token value
		var foundToken models.RefreshToken
		result = db.Where("token = ?", "refresh_token_value").First(&foundToken)
		assert.NoError(t, result.Error)
		assert.Equal(t, uint(1), foundToken.UserID)

		// Test updating token (rotation)
		foundToken.Token = "new_refresh_token_value"
		result = db.Save(&foundToken)
		assert.NoError(t, result.Error)

		// Verify update
		var updatedToken models.RefreshToken
		result = db.Where("token = ?", "new_refresh_token_value").First(&updatedToken)
		assert.NoError(t, result.Error)
		assert.Equal(t, "new_refresh_token_value", updatedToken.Token)

		// Test deletion (when token is revoked)
		result = db.Delete(&models.RefreshToken{}, updatedToken.ID)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedToken models.RefreshToken
		result = db.First(&deletedToken, updatedToken.ID)
		assert.Error(t, result.Error)
	})
}