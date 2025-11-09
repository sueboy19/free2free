package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSessionInitialization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid session initialization", func(t *testing.T) {
		// Create a test session handler
		testSession := testutils.NewTestSession()

		// Create a gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Simulate a request
		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req

		// Since we can't access the unexported store field, create the session in a different way
		// The AddSession method should create the session properly
		err := testSession.AddSession(c, "free2free-session", map[string]interface{}{
			"user_id": 1,
			"email":   "test@example.com",
		})
		assert.NoError(t, err)

		// Test that we can access the session without panic
		retrievedSession, err := testSession.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedSession)
	})

	t.Run("Session with values", func(t *testing.T) {
		// Create a test session handler
		testSession := testutils.NewTestSession()

		// Create a gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Add session with values
		err := testSession.AddSession(c, "free2free-session", map[string]interface{}{
			"user_id": 1,
			"email":   "test@example.com",
		})
		assert.NoError(t, err)

		// Simulate a request
		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req

		// Test that we can access the session without panic and has correct values
		session, err := testSession.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check that the session has the expected values
		userID, exists := session.Values["user_id"]
		assert.True(t, exists)
		assert.Equal(t, 1, userID)

		email, exists := session.Values["email"]
		assert.True(t, exists)
		assert.Equal(t, "test@example.com", email)
	})

	t.Run("Session retrieval without initialization", func(t *testing.T) {
		// Create a gin context without session
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req

		// Create a test session handler
		ts := testutils.NewTestSession()

		// This should not panic and should create a new session
		session, err := ts.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, session)
	})
}