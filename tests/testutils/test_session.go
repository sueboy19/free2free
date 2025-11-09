package testutils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// TestSession provides utilities for handling sessions in test environments
type TestSession struct {
	store sessions.Store
}

// NewTestSession creates a new test session handler
func NewTestSession() *TestSession {
	// Create a test session store with proper key sizes (32 bytes)
	authKey := make([]byte, 32)
	copy(authKey, "test-auth-key-32-characters-lon") // exactly 32 characters
	encryptionKey := make([]byte, 32)
	copy(encryptionKey, "test-encrypt-key-32-characters!") // exactly 32 characters
	store := sessions.NewCookieStore(authKey, encryptionKey)
	
	// Set appropriate options for testing
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to false for testing
		SameSite: http.SameSiteLaxMode,
	}
	
	return &TestSession{
		store: store,
	}
}

// AddSession adds a session to a gin context for testing
func (ts *TestSession) AddSession(c *gin.Context, sessionName string, values map[string]interface{}) error {
	// Ensure c.Request is not nil before using it
	if c.Request == nil {
		req, _ := http.NewRequest("GET", "/", nil)
		c.Request = req
	}
	
	session, err := ts.store.Get(c.Request, sessionName)
	if err != nil {
		return err
	}

	// Set the session values
	for k, v := range values {
		session.Values[k] = v
	}

	// Save the session
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return err
	}

	// Add the session to the context
	c.Set("session", session)

	return nil
}

// CreateSessionRequest creates an HTTP request with a pre-filled session
func (ts *TestSession) CreateSessionRequest(method, url string, sessionValues map[string]interface{}) (*http.Request, *httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Create a temporary response recorder
	w := httptest.NewRecorder()

	// Add session to the request
	session, err := ts.store.Get(req, "free2free-session")
	if err != nil {
		return nil, nil, err
	}

	// Set the session values
	for k, v := range sessionValues {
		session.Values[k] = v
	}

	// Save the session to the response recorder
	err = session.Save(req, w)
	if err != nil {
		return nil, nil, err
	}

	// Add session cookie to the request
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	return req, w, nil
}

// GetSessionFromContext safely retrieves the session from gin context
func (ts *TestSession) GetSessionFromContext(c *gin.Context) (*sessions.Session, error) {
	sessionVal, exists := c.Get("session")
	if !exists {
		// Ensure c.Request is not nil before using it
		if c.Request == nil {
			req, _ := http.NewRequest("GET", "/", nil)
			c.Request = req
		}
		
		// Create a new empty session instead of panicking
		session, err := ts.store.New(c.Request, "free2free-session")
		if err != nil {
			return nil, err
		}
		c.Set("session", session)
		return session, nil
	}
	
	session, ok := sessionVal.(*sessions.Session)
	if !ok {
		// Ensure c.Request is not nil before using it
		if c.Request == nil {
			req, _ := http.NewRequest("GET", "/", nil)
			c.Request = req
		}
		
		// Create a new empty session instead of panicking
		session, err := ts.store.New(c.Request, "free2free-session")
		if err != nil {
			return nil, err
		}
		c.Set("session", session)
		return session, nil
	}
	
	return session, nil
}