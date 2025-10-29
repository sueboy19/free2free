package testutils

import (
	"math/rand"
	"time"
)

// TestDataGenerator provides utilities for generating test data
type TestDataGenerator struct {
	rand *rand.Rand
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateUser creates a mock user for testing
func (g *TestDataGenerator) GenerateUser(id uint, email, name, provider string) TestUser {
	if email == "" {
		email = g.GenerateEmail()
	}
	if name == "" {
		name = g.GenerateName()
	}
	if provider == "" {
		provider = g.GenerateProvider()
	}

	return TestUser{
		ID:       id,
		Email:    email,
		Name:     name,
		Provider: provider,
		Role:     "user",
	}
}

// GenerateActivity creates a mock activity for testing
func (g *TestDataGenerator) GenerateActivity(id uint, title, description string, locationID, creatorID uint) TestActivity {
	if title == "" {
		title = g.GenerateTitle()
	}
	if description == "" {
		description = g.GenerateDescription()
	}

	statuses := []string{"pending", "approved", "rejected", "active"}
	status := statuses[g.rand.Intn(len(statuses))]

	return TestActivity{
		ID:          id,
		Title:       title,
		Description: description,
		LocationID:  locationID,
		Status:      status,
		CreatorID:   creatorID,
	}
}

// GenerateEmail creates a random email for testing
func (g *TestDataGenerator) GenerateEmail() string {
	adjectives := []string{"happy", "quick", "bright", "clever", "friendly", "smart", "cool", "sunny", "brave", "calm"}
	nouns := []string{"cat", "dog", "bird", "fish", "lion", "tiger", "bear", "wolf", "fox", "rabbit"}
	domains := []string{"example.com", "test.com", "mock.org", "fakemail.net", "tempmail.co"}

	adjective := adjectives[g.rand.Intn(len(adjectives))]
	noun := nouns[g.rand.Intn(len(nouns))]
	number := g.rand.Intn(1000)
	domain := domains[g.rand.Intn(len(domains))]

	return adjective + noun + string(rune('0'+number/100)) + string(rune('0'+(number/10)%10)) + string(rune('0'+number%10)) + "@" + domain
}

// GenerateName creates a random name for testing
func (g *TestDataGenerator) GenerateName() string {
	firstNames := []string{"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", "William", "Elizabeth"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez"}

	firstName := firstNames[g.rand.Intn(len(firstNames))]
	lastName := lastNames[g.rand.Intn(len(lastNames))]

	return firstName + " " + lastName
}

// GenerateTitle creates a random title for testing
func (g *TestDataGenerator) GenerateTitle() string {
	adjectives := []string{"Amazing", "Fantastic", "Wonderful", "Incredible", "Superb", "Excellent", "Outstanding", "Impressive", "Remarkable", "Phenomenal"}
	nouns := []string{"Event", "Activity", "Experience", "Opportunity", "Gathering", "Meeting", "Session", "Workshop", "Seminar", "Conference"}

	adjective := adjectives[g.rand.Intn(len(adjectives))]
	noun := nouns[g.rand.Intn(len(nouns))]

	return adjective + " " + noun
}

// GenerateDescription creates a random description for testing
func (g *TestDataGenerator) GenerateDescription() string {
	sentences := []string{
		"This is a wonderful opportunity to learn and grow.",
		"Don't miss this chance to connect with like-minded individuals.",
		"A perfect event for expanding your knowledge and network.",
		"An excellent way to spend your time with meaningful activities.",
		"Join us for an unforgettable experience filled with fun.",
		"This activity promises to be both educational and entertaining.",
		"A great opportunity to meet new people and have fun.",
		"An event designed to bring the community together.",
		"Experience something new and exciting with this activity.",
		"A chance to develop new skills in a supportive environment.",
	}

	// Generate a description with 2-4 sentences
	numSentences := 2 + g.rand.Intn(3)
	description := ""

	for i := 0; i < numSentences; i++ {
		description += sentences[g.rand.Intn(len(sentences))]
		if i < numSentences-1 {
			description += " "
		}
	}

	return description
}

// GenerateProvider creates a random provider for testing
func (g *TestDataGenerator) GenerateProvider() string {
	providers := []string{"facebook", "instagram", "google", "twitter"}
	return providers[g.rand.Intn(len(providers))]
}

// GenerateTestUsers creates multiple test users
func (g *TestDataGenerator) GenerateTestUsers(count int) []TestUser {
	users := make([]TestUser, count)
	for i := 0; i < count; i++ {
		users[i] = g.GenerateUser(uint(i+1), "", "", "")
	}
	return users
}

// GenerateTestActivities creates multiple test activities
func (g *TestDataGenerator) GenerateTestActivities(count int, creatorID uint) []TestActivity {
	activities := make([]TestActivity, count)
	for i := 0; i < count; i++ {
		activities[i] = g.GenerateActivity(uint(i+1), "", "", uint(i+1), creatorID)
	}
	return activities
}

// Predefined test data for consistent testing
var (
	// Standard test user
	StandardTestUser = TestUser{
		ID:       999,
		Email:    "standard@test.com",
		Name:     "Standard Test User",
		Provider: "facebook",
		Role:     "user",
	}

	// Standard test activity
	StandardTestActivity = TestActivity{
		ID:          999,
		Title:       "Standard Test Activity",
		Description: "This is a standard test activity for consistent testing.",
		LocationID:  1,
		Status:      "pending",
		CreatorID:   999,
	}
)
