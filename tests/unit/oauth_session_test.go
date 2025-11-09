package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"free2free/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOAuthSessionEstablishment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Session validator creates valid sessions", func(t *testing.T) {
		// Create a session validator
		sessionValidator := testutils.NewSessionValidator()

		// Create a gin router and add test handler to access the context after middleware runs
		r := gin.New()
		sessionValidator.AddSessionMiddleware(r)

		var sessionVal interface{}
		var exists bool

		r.GET("/", func(c *gin.Context) {
			// Capture session info in closure variables after middleware has run
			sessionVal, exists = c.Get("session")
		})

		// Create a request
		req, _ := http.NewRequest("GET", "/", nil)

		// Test the middleware by sending a request through the router
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Now check the captured session values
		assert.True(t, exists, "Session should exist in context")
		assert.NotNil(t, sessionVal, "Session should not be nil")
	})

	t.Run("Adding values to session works correctly", func(t *testing.T) {
		// Create a test session handler
		ts := testutils.NewTestSession()

		// Create a gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a request
		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req

		// Add session with values
		sessionValues := map[string]interface{}{
			"user_id": int64(123),
			"name":    "Test User",
		}

		err := ts.AddSession(c, "free2free-session", sessionValues)
		assert.NoError(t, err)

		// Verify we can retrieve the session and values
		session, err := ts.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check that values were set correctly
		assert.Equal(t, int64(123), session.Values["user_id"])
		assert.Equal(t, "Test User", session.Values["name"])
	})

	t.Run("Session validator handles missing session gracefully", func(t *testing.T) {
		// Create a gin context without session
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a request
		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req

		// Create a test session handler
		ts := testutils.NewTestSession()

		// This should create a new session instead of panicking
		session, err := ts.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, session)
	})

	t.Run("Session with user authentication data", func(t *testing.T) {
		// Create a test session handler
		ts := testutils.NewTestSession()

		// Create a gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a request
		req, _ := http.NewRequest("GET", "/profile", nil)
		c.Request = req

		// Add authentication session data
		authSessionData := map[string]interface{}{
			"user_id": 456,
			"email":   "user@example.com",
		}

		err := ts.AddSession(c, "free2free-session", authSessionData)
		assert.NoError(t, err)

		// Verify session was created with correct data
		session, err := ts.GetSessionFromContext(c)
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check authentication data
		assert.Equal(t, 456, session.Values["user_id"])
		assert.Equal(t, "user@example.com", session.Values["email"])
	})

	t.Run("Creating session request with values", func(t *testing.T) {
		// Create a test session handler
		ts := testutils.NewTestSession()

		// Create a request with pre-filled session
		sessionData := map[string]interface{}{
			"user_id":    789,
			"is_admin":   false,
			"user_name":  "Admin User",
		}

		req, w, err := ts.CreateSessionRequest("GET", "/admin", sessionData)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.NotNil(t, w)

		// Verify that the request has session cookie set
		cookies := req.Cookies()
		assert.NotEmpty(t, cookies)

		// Check cookie names contain session information
		hasSessionCookie := false
		for _, cookie := range cookies {
			if cookie.Name == "free2free-session" {
				hasSessionCookie = true
				break
			}
		}
		assert.True(t, hasSessionCookie, "Request should have session cookie")
	})
}