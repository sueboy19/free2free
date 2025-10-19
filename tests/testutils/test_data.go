package testutils

import (
	"fmt"

	"free2free/models"
	"gorm.io/gorm"
)

// TestUser represents a test user for testing purposes
type TestUser struct {
	ID           int64
	Name         string
	Email        string
	IsAdmin      bool
	SocialID     string
	SocialProvider string
	AvatarURL    string
}

// TestDataGenerator provides utilities to generate test data
type TestDataGenerator struct {
	DB *gorm.DB
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator(db *gorm.DB) *TestDataGenerator {
	return &TestDataGenerator{
		DB: db,
	}
}

// CreateTestUser creates a test user in the database
func (gen *TestDataGenerator) CreateTestUser(name, email string, isAdmin bool) (*models.User, error) {
	user := &models.User{
		SocialID:       fmt.Sprintf("test_%s_social_id", name),
		SocialProvider: "facebook", // Default to Facebook for test
		Name:           name,
		Email:          email,
		AvatarURL:      fmt.Sprintf("https://example.com/%s-avatar.jpg", name),
		IsAdmin:        isAdmin,
	}

	result := gen.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// CreateTestAdmin creates a test admin user
func (gen *TestDataGenerator) CreateTestAdmin() (*models.User, error) {
	return gen.CreateTestUser("Test Admin", "admin@test.com", true)
}

// CreateTestRegularUser creates a test regular user
func (gen *TestDataGenerator) CreateTestRegularUser() (*models.User, error) {
	return gen.CreateTestUser("Test User", "user@test.com", false)
}

// CreateTestLocation creates a test location
func (gen *TestDataGenerator) CreateTestLocation(name, address string, lat, lng float64) (*models.Location, error) {
	location := &models.Location{
		Name:      name,
		Address:   address,
		Latitude:  lat,
		Longitude: lng,
	}

	result := gen.DB.Create(location)
	if result.Error != nil {
		return nil, result.Error
	}

	return location, nil
}

// CreateTestActivity creates a test activity
func (gen *TestDataGenerator) CreateTestActivity(title, description string, locationID uint, targetCount int) (*models.Activity, error) {
	activity := &models.Activity{
		Title:       title,
		Description: description,
		LocationID:  locationID,
		TargetCount: targetCount,
		CreatedBy:   1, // Default to first user
	}

	result := gen.DB.Create(activity)
	if result.Error != nil {
		return nil, result.Error
	}

	return activity, nil
}

// CreateTestMatch creates a test match
func (gen *TestDataGenerator) CreateTestMatch(activityID, organizerID uint, matchTime string) (*models.Match, error) {
	match := &models.Match{
		ActivityID:  activityID,
		OrganizerID: organizerID,
		MatchTime:   matchTime,
		Status:      "open", // Default status
	}

	result := gen.DB.Create(match)
	if result.Error != nil {
		return nil, result.Error
	}

	return match, nil
}

// CreateTestReview creates a test review
func (gen *TestDataGenerator) CreateTestReview(matchID, reviewerID, revieweeID uint, score int, comment string) (*models.Review, error) {
	review := &models.Review{
		MatchID:     matchID,
		ReviewerID:  reviewerID,
		RevieweeID:  revieweeID,
		Score:       score,
		Comment:     comment,
	}

	result := gen.DB.Create(review)
	if result.Error != nil {
		return nil, result.Error
	}

	return review, nil
}

// CreateTestMatchParticipant creates a test match participant
func (gen *TestDataGenerator) CreateTestMatchParticipant(matchID, userID uint, status string) (*models.MatchParticipant, error) {
	participant := &models.MatchParticipant{
		MatchID: matchID,
		UserID:  userID,
		Status:  status, // pending, approved, rejected
	}

	result := gen.DB.Create(participant)
	if result.Error != nil {
		return nil, result.Error
	}

	return participant, nil
}

// CreateTestReviewLike creates a test review like/dislike
func (gen *TestDataGenerator) CreateTestReviewLike(reviewID, userID uint, isLike bool) (*models.ReviewLike, error) {
	like := &models.ReviewLike{
		ReviewID: reviewID,
		UserID:   userID,
		IsLike:   isLike,
	}

	result := gen.DB.Create(like)
	if result.Error != nil {
		return nil, result.Error
	}

	return like, nil
}

// ClearAllTestData removes all test data from the database
func (gen *TestDataGenerator) ClearAllTestData() error {
	// Delete in reverse order to respect foreign key constraints
	tables := []interface{}{
		&models.ReviewLike{}, &models.Review{}, &models.MatchParticipant{},
		&models.Match{}, &models.Activity{}, &models.Location{},
		&models.Admin{}, &models.User{}, &models.RefreshToken{},
	}

	for _, table := range tables {
		if err := gen.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			// If table doesn't exist, continue
			return err
		}
	}

	return nil
}

// SetupSampleData creates sample data for testing
func (gen *TestDataGenerator) SetupSampleData() error {
	// Create a test admin user
	admin, err := gen.CreateTestAdmin()
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	// Create a regular user
	user, err := gen.CreateTestRegularUser()
	if err != nil {
		return fmt.Errorf("failed to create regular user: %w", err)
	}

	// Create a location
	location, err := gen.CreateTestLocation("Test Location", "123 Test St", 25.0, 121.0)
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	// Create an activity
	activity, err := gen.CreateTestActivity("Test Activity", "Sample activity for testing", location.ID, 2)
	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	// Create a match organized by the admin
	match, err := gen.CreateTestMatch(activity.ID, admin.ID, "2025-12-31T10:00:00Z")
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	// Add the organizer as a participant
	_, err = gen.CreateTestMatchParticipant(match.ID, admin.ID, "approved")
	if err != nil {
		return fmt.Errorf("failed to create organizer participant: %w", err)
	}

	// Add the regular user as a participant
	_, err = gen.CreateTestMatchParticipant(match.ID, user.ID, "approved")
	if err != nil {
		return fmt.Errorf("failed to create regular user participant: %w", err)
	}

	// Create a review from regular user to admin (after match completion)
	_, err = gen.CreateTestReview(match.ID, user.ID, admin.ID, 5, "Great organizer!")
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}