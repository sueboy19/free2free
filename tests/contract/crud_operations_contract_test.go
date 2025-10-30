package contract

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"free2free/models"
	"free2free/tests/testutils"
)

// TestCRUDOperationsContract tests the CRUD operations contract
func TestCRUDOperationsContract(t *testing.T) {
	t.Run("Create Operation (INSERT)", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Migrate the User model for testing
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Test creating a new user record
		user := &models.User{
			Name:           "Test User",
			Email:          "test@example.com",
			SocialProvider: "facebook",
			SocialID:       "123456",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)
	})

	t.Run("Read Operation (SELECT)", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Migrate the User model for testing
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a test user
		user := &models.User{
			Name:           "Test User",
			Email:          "test@example.com",
			SocialProvider: "facebook",
			SocialID:       "123456",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Test reading the user record
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("Update Operation (UPDATE)", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Migrate the User model for testing
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a test user
		user := &models.User{
			Name:           "Test User",
			Email:          "test@example.com",
			SocialProvider: "facebook",
			SocialID:       "123456",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Test updating the user record
		user.Name = "Updated User"
		result = db.Save(user)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedUser models.User
		result = db.First(&updatedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, "Updated User", updatedUser.Name)
	})

	t.Run("Delete Operation (DELETE)", func(t *testing.T) {
		db, err := testutils.CreateTestDB()
		if err != nil {
			t.Logf("Database connection error: %v", err)
			if err != nil && (err.Error() == "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub" ||
				err.Error() == "failed to connect database") {
				t.Skip("Skipping test due to CGO dependency issue - this is expected in some environments")
			}
			assert.NoError(t, err)
			return
		}
		assert.NotNil(t, db)

		// Migrate the User model for testing
		err = testutils.MigrateTestDB(db, &models.User{})
		assert.NoError(t, err)

		// Create a test user
		user := &models.User{
			Name:           "Test User",
			Email:          "test@example.com",
			SocialProvider: "facebook",
			SocialID:       "123456",
			AvatarURL:      "https://example.com/avatar.jpg",
		}

		result := db.Create(user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Test deleting the user record
		result = db.Delete(&models.User{}, user.ID)
		assert.NoError(t, result.Error)

		// Verify the deletion
		var deletedUser models.User
		result = db.First(&deletedUser, user.ID)
		assert.Error(t, result.Error)
	})
}