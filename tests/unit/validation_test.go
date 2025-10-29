package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestActivityDataValidation tests validation of activity data
func TestActivityDataValidation(t *testing.T) {
	t.Run("Valid Activity Data", func(t *testing.T) {
		// Test valid activity data that should pass validation
		title := "Valid Activity Title"
		description := "This is a valid activity description that meets the minimum length requirement."
		locationID := uint(1)

		titleValid := validateTitle(title)
		descriptionValid := validateDescription(description)
		locationValid := validateLocationID(locationID)

		assert.True(t, titleValid, "Valid title should pass validation")
		assert.True(t, descriptionValid, "Valid description should pass validation")
		assert.True(t, locationValid, "Valid location ID should pass validation")
	})

	t.Run("Invalid Title Validation", func(t *testing.T) {
		// Test invalid titles that should fail validation
		invalidTitles := []string{
			"",  // Empty title
			"X", // Too short (less than 1 char min)
			"Very long title that exceeds the maximum allowed length by quite a bit and is definitely over the limit", // Too long (over 100 chars max)
		}

		for _, title := range invalidTitles {
			t.Run("Invalid title: "+title, func(t *testing.T) {
				titleValid := validateTitle(title)
				assert.False(t, titleValid, "Invalid title should fail validation")
			})
		}
	})

	t.Run("Invalid Description Validation", func(t *testing.T) {
		// Test invalid descriptions that should fail validation
		invalidDescriptions := []string{
			"",           // Empty description
			"Short",      // Too short (less than 10 chars min)
			"Very short", // Still too short
			"A very long description that definitely exceeds the maximum allowed length of 500 characters. " +
				"This description continues with additional text to ensure it surpasses the 500 character limit. " +
				"It contains multiple sentences and various words to achieve the required length. " +
				"Additional text is included to ensure it is well over the specified limit. " +
				"More text continues here to make sure the character count is exceeded. " +
				"Even more text is added to guarantee the limit has been surpassed. " +
				"Yet more text follows to ensure the validation will fail as expected. " +
				"More content is included. The count continues to grow. Adding more words. " +
				"Still adding more text. More text. More text. More text. More text. More text.", // Too long (over 500 chars max)
		}

		for _, description := range invalidDescriptions {
			t.Run("Invalid description length: "+string(rune(len(description))), func(t *testing.T) {
				descriptionValid := validateDescription(description)
				assert.False(t, descriptionValid, "Invalid description should fail validation")
			})
		}
	})

	t.Run("Invalid Location ID Validation", func(t *testing.T) {
		// Test invalid location IDs that should fail validation
		// For this test, we'll consider 0 as invalid if not allowed
		invalidLocationID := uint(0) // Assuming 0 is not a valid location ID

		locationValid := validateLocationID(invalidLocationID)
		assert.False(t, locationValid, "Invalid location ID should fail validation")
	})

	t.Run("Boundary Values for Title", func(t *testing.T) {
		// Test boundary values for title length
		validMinTitle := "A" // 1 character - minimum valid length
		invalidMaxTitle := ""
		for i := 0; i < 101; i++ {
			invalidMaxTitle += "A"
		} // 101 characters - exceeds max length

		minTitleValid := validateTitle(validMinTitle)
		maxTitleValid := validateTitle(invalidMaxTitle)

		assert.True(t, minTitleValid, "Minimum valid title length should pass")
		assert.False(t, maxTitleValid, "Maximum invalid title length should fail")
	})

	t.Run("Boundary Values for Description", func(t *testing.T) {
		// Test boundary values for description length
		validMinDescription := "Ten Chars!" // 10 characters - minimum valid length
		invalidMaxDescription := ""
		for i := 0; i < 501; i++ {
			invalidMaxDescription += "A"
		} // 501 characters - exceeds max length

		minDescriptionValid := validateDescription(validMinDescription)
		maxDescriptionValid := validateDescription(invalidMaxDescription)

		assert.True(t, minDescriptionValid, "Minimum valid description length should pass")
		assert.False(t, maxDescriptionValid, "Maximum invalid description length should fail")
	})
}

// TestUserDataValidation tests validation of user data
func TestUserDataValidation(t *testing.T) {
	t.Run("Valid Email Format", func(t *testing.T) {
		validEmails := []string{
			"user@example.com",
			"test.email+tag@domain.co.uk",
			"user123@sub.domain.org",
		}

		for _, email := range validEmails {
			t.Run("Valid email: "+email, func(t *testing.T) {
				emailValid := validateEmailFormat(email)
				assert.True(t, emailValid, "Valid email format should pass validation")
			})
		}
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		invalidEmails := []string{
			"invalid-email",
			"@invalid.com",
			"user@",
			"user..name@example.com",
			"user@domain",
			"",
		}

		for _, email := range invalidEmails {
			t.Run("Invalid email: "+email, func(t *testing.T) {
				emailValid := validateEmailFormat(email)
				assert.False(t, emailValid, "Invalid email format should fail validation")
			})
		}
	})

	t.Run("Valid User Name", func(t *testing.T) {
		validNames := []string{
			"John Doe",
			"Jane",
			"User Name With Spaces",
		}

		for _, name := range validNames {
			t.Run("Valid name: "+name, func(t *testing.T) {
				nameValid := validateUserName(name)
				assert.True(t, nameValid, "Valid user name should pass validation")
			})
		}
	})

	t.Run("Invalid User Name", func(t *testing.T) {
		invalidNames := []string{
			"",   // Empty name
			"  ", // Whitespace only
			"A",  // Too short (if we have min length)
			"Very long name that might exceed limits if we set any",
		}

		for _, name := range invalidNames {
			t.Run("Invalid name: "+name, func(t *testing.T) {
				nameValid := validateUserName(name)
				assert.False(t, nameValid, "Invalid user name should fail validation")
			})
		}
	})
}

// validateTitle validates the title of an activity
func validateTitle(title string) bool {
	if len(title) < 1 || len(title) > 100 {
		return false
	}
	return true
}

// validateDescription validates the description of an activity
func validateDescription(description string) bool {
	if len(description) < 10 || len(description) > 500 {
		return false
	}
	return true
}

// validateLocationID validates the location ID of an activity
func validateLocationID(locationID uint) bool {
	if locationID == 0 {
		return false
	}
	return true
}

// validateEmailFormat validates the format of an email address
func validateEmailFormat(email string) bool {
	// Simple email validation: check for @ and basic format
	if len(email) == 0 {
		return false
	}

	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				// Multiple @ symbols
				return false
			}
			atIndex = i
		}
	}

	if atIndex <= 0 || atIndex == len(email)-1 {
		// @ at beginning or end
		return false
	}

	// Check for dot after @
	domainPart := email[atIndex+1:]
	if len(domainPart) == 0 {
		return false
	}

	hasDot := false
	for _, char := range domainPart {
		if char == '.' {
			hasDot = true
			break
		}
	}

	return hasDot
}

// validateUserName validates a user name
func validateUserName(name string) bool {
	// Check if name is not empty and has reasonable length
	nameLen := len(name)
	if nameLen < 1 || nameLen > 100 {
		return false
	}

	// Check if it's just whitespace
	allSpaces := true
	for _, char := range name {
		if char != ' ' {
			allSpaces = false
			break
		}
	}

	if allSpaces {
		return false
	}

	return true
}
