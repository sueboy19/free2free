package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var Store sessions.Store

func init() {
	// Initialize session store with configuration from environment
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		// For testing, use a default key if not set
		sessionKey = "test-session-key-for-testing-environment"
	}

	var authKey, encryptionKey []byte
	if len(sessionKey) >= 64 {
		authKey = []byte(sessionKey[:32])
		encryptionKey = []byte(sessionKey[32:64])
	} else if len(sessionKey) >= 32 {
		authKey = []byte(sessionKey[:32])
		// If key is shorter than 64 but longer than 32, repeat or pad
		encryptionKey = make([]byte, 32)
		for i := 0; i < 32; i++ {
			encryptionKey[i] = sessionKey[i%len(sessionKey)]
		}
	} else {
		// If key is too short, pad it to appropriate lengths
		authKey = make([]byte, 32)
		encryptionKey = make([]byte, 32)
		for i := 0; i < 32; i++ {
			authKey[i] = sessionKey[i%len(sessionKey)]
			encryptionKey[i] = sessionKey[i%len(sessionKey)]
		}
	}

	Store = sessions.NewCookieStore(authKey, encryptionKey)
}

// SessionMiddleware handles session creation and management
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or create a session
		session, err := Store.Get(c.Request, "free2free-session")
		if err != nil {
			// If there's an error getting the session, log it but don't panic
			// Just create a new empty session
			session, _ = Store.New(c.Request, "free2free-session")
		}

		// Make sure session is never nil
		if session == nil {
			session, _ = Store.New(c.Request, "free2free-session")
		}

		// Set the session in the context
		c.Set("session", session)
		
		// Continue with the request
		c.Next()
		
		// Save the session if it was modified
		err = Store.Save(c.Request, c.Writer, session)
		if err != nil {
			// Log the error but don't fail the request
			fmt.Printf("Error saving session: %v\n", err)
		}
	}
}

// GetSession safely retrieves the session from context
func GetSession(c *gin.Context) (*sessions.Session, error) {
	sessionVal, exists := c.Get("session")
	if !exists {
		return nil, fmt.Errorf("session not found in context")
	}
	
	session, ok := sessionVal.(*sessions.Session)
	if !ok {
		return nil, fmt.Errorf("session type assertion failed")
	}
	
	return session, nil
}