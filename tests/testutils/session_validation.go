package testutils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// SessionValidator provides utilities for validating session handling in tests
type SessionValidator struct {
	store sessions.Store
}

// NewSessionValidator creates a new session validator with a test session store
func NewSessionValidator() *SessionValidator {
	// Create a test session store with proper key sizes (32 bytes)
	authKey := make([]byte, 32)
	copy(authKey, "test-auth-key-32-characters-lon") // exactly 32 characters
	encryptionKey := make([]byte, 32)
	copy(encryptionKey, "test-encrypt-key-32-characters!") // exactly 32 characters
	store := sessions.NewCookieStore(authKey, encryptionKey)
	
	return &SessionValidator{
		store: store,
	}
}

// AddSessionMiddleware adds proper session handling to a gin router for testing
func (sv *SessionValidator) AddSessionMiddleware(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		session, err := sv.store.Get(c.Request, "free2free-session")
		if err != nil {
			// If there's an error getting the session, log it but don't panic
			// Just create a new empty session
			session, _ = sv.store.New(c.Request, "free2free-session")
		}

		// Make sure session is never nil
		if session == nil {
			session, _ = sv.store.New(c.Request, "free2free-session")
		}

		// Set the session in the context
		c.Set("session", session)
		
		// Continue with the request
		c.Next()
		
		// Save the session if needed
		sv.store.Save(c.Request, c.Writer, session)
	})
}

// AddSessionToRequest adds a session to an HTTP request for testing
func (sv *SessionValidator) AddSessionToRequest(req *http.Request, sessionName string, values map[string]interface{}) error {
	session, err := sv.store.Get(req, sessionName)
	if err != nil {
		return err
	}

	// Set the session values
	for k, v := range values {
		session.Values[k] = v
	}

	// Save the session to the request
	return session.Save(req, httptest.NewRecorder())
}

// ValidateSessionExists checks if a session exists and is accessible without panics
func (sv *SessionValidator) ValidateSessionExists(c *gin.Context, sessionName string) (bool, error) {
	session, err := sv.store.Get(c.Request, sessionName)
	if err != nil {
		return false, err
	}

	// Check if session is properly initialized
	if session == nil {
		return false, nil
	}

	return true, nil
}