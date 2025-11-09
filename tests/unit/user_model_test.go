package unit

import (
	"testing"

	"free2free/models"

	"github.com/stretchr/testify/assert"
)

func TestUserModelFields(t *testing.T) {
	t.Run("User model has expected fields", func(t *testing.T) {
		// Create a user instance
		user := models.User{}

		// Verify that the model has the expected fields by attempting to access them
		// This test ensures the structure matches what's expected in the implementation
		
		// The model should have these fields:
		// - ID (primary key)
		// - SocialID
		// - SocialProvider
		// - Name
		// - Email
		// - AvatarURL
		// - IsAdmin
		
		// We can't directly access unexported fields, but we can check the structure exists
		// For Go models, if we can create an instance, the structure exists
		
		// Check that we can assign values to basic fields
		user.ID = 1
		user.SocialID = "social123"
		user.SocialProvider = "facebook"
		user.Name = "Test User"
		user.Email = "test@example.com"
		user.AvatarURL = "https://example.com/avatar.jpg"
		user.IsAdmin = false
		
		// Verify fields were set correctly
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "social123", user.SocialID)
		assert.Equal(t, "facebook", user.SocialProvider)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
		assert.Equal(t, false, user.IsAdmin)
	})

	t.Run("User model with provider fields exists", func(t *testing.T) {
		// Create a user with values that match how it's used in the actual implementation
		user := models.User{
			ID:             1,
			SocialID:       "facebook123",
			SocialProvider: "facebook",
			Name:           "John Doe",
			Email:          "john@example.com",
			AvatarURL:      "https://graph.facebook.com/123/picture",
			IsAdmin:        false,
		}

		// Verify all fields have been set properly
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "facebook123", user.SocialID)
		assert.Equal(t, "facebook", user.SocialProvider)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "https://graph.facebook.com/123/picture", user.AvatarURL)
		assert.Equal(t, false, user.IsAdmin)
	})

	t.Run("User model admin field works correctly", func(t *testing.T) {
		// Test that the IsAdmin field works as expected
		adminUser := models.User{IsAdmin: true}
		normalUser := models.User{IsAdmin: false}

		assert.True(t, adminUser.IsAdmin)
		assert.False(t, normalUser.IsAdmin)
	})

	t.Run("User model has required fields for OAuth", func(t *testing.T) {
		// Test that the model has the fields required for OAuth flow
		user := models.User{
			SocialID:       "test123",
			SocialProvider: "test_provider",
			Name:           "Test User",
			Email:          "test@example.com",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		// Verify OAuth-related fields are present
		assert.NotEmpty(t, user.SocialID)
		assert.NotEmpty(t, user.SocialProvider)
		assert.NotEmpty(t, user.Name)
		assert.NotEmpty(t, user.Email)
		assert.NotEmpty(t, user.AvatarURL)
	})
}