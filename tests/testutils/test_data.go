package testutils

import (
	"time"

	"free2free/models"
)



// CreateTestActivity creates a test activity with valid data structure
func CreateTestActivity() models.Activity {
	return models.Activity{
		ID:          1,
		Title:       "Test Activity",
		TargetCount: 2,
		LocationID:  1,
		Description: "This is a test activity",
		CreatedBy:   1,
	}
}

// CreateTestLocation creates a test location with valid data structure
func CreateTestLocation() models.Location {
	return models.Location{
		ID:        1,
		Name:      "Test Location",
		Address:   "123 Test Street",
		Latitude:  25.03,
		Longitude: 121.5,
	}
}

// CreateTestMatch creates a test match with valid data structure
func CreateTestMatch() models.Match {
	return models.Match{
		ID:          1,
		ActivityID:  1,
		OrganizerID: 1,
		MatchTime:   time.Now().Add(24 * time.Hour), // Tomorrow
		Status:      "open",
	}
}

// CreateTestMatchParticipant creates a test match participant with valid data structure
func CreateTestMatchParticipant() models.MatchParticipant {
	return models.MatchParticipant{
		ID:       1,
		MatchID:  1,
		UserID:   1,
		Status:   "pending",
		JoinedAt: time.Now(),
	}
}

// CreateTestReview creates a test review with valid data structure
func CreateTestReview() models.Review {
	return models.Review{
		ID:         1,
		MatchID:    1,
		ReviewerID: 1,
		RevieweeID: 2,
		Score:      5,
		Comment:    "Great experience!",
		CreatedAt:  time.Now(),
	}
}

// CreateTestReviewLike creates a test review like with valid data structure
func CreateTestReviewLike() models.ReviewLike {
	return models.ReviewLike{
		ID:       1,
		ReviewID: 1,
		UserID:   1,
		IsLike:   true,
	}
}

// CreateTestRefreshToken creates a test refresh token with valid data structure
func CreateTestRefreshToken() models.RefreshToken {
	return models.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     "refresh-token-string",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

// CreateTestAdmin creates a test admin with valid data structure
func CreateTestAdmin() models.Admin {
	return models.Admin{
		ID:       1,
		Username: "test_admin",
		Email:    "admin@example.com",
	}
}