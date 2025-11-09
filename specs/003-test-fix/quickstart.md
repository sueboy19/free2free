# Quickstart Guide: Test Fix

## Overview
This guide provides instructions for implementing fixes to identified test issues in the API testing functionality, including session handling problems, missing API routes, build errors in tests, and OAuth flow issues.

## Prerequisites

1. Go 1.25.0 or higher installed
2. Project dependencies installed (`go mod download`)
3. Database setup and running (MariaDB/MySQL)
4. Environment variables configured (use `.env.example` as template)
5. Basic understanding of the Gin framework and Go testing

## Understanding the Issues

Before implementing fixes, it's important to understand the specific test issues that need to be addressed:

### 1. Session Handling Problems
- Runtime panics: `interface conversion: interface {} is nil, not *sessions.Session`
- Occurs in authentication handlers when accessing session data
- Root cause: Session middleware not properly initializing session in test context

### 2. Missing API Routes
- 404 errors for documented endpoints like `/profile`
- Root cause: Routes not registered in test server setup

### 3. Build Errors in Tests
- Compilation errors in performance and integration tests
- Examples: JWT claims structure mismatches, model field mismatches
- Root cause: Test code not aligned with actual implementation

### 4. OAuth Flow Issues
- SESSION_SECRET warnings in test environment
- Root cause: Missing environment configuration for tests

## Fix Implementation Strategy

### 1. Fix Session Handling
Create proper session initialization for test environments:

```go
// In test utilities
func NewTestServer() *TestServer {
    r := gin.Default()
    
    // Initialize session store similar to main.go
    sessionKey := os.Getenv("SESSION_KEY") 
    if sessionKey == "" {
        // Use test-specific key for testing
        sessionKey = "test-session-key-for-testing-environment"
    }
    
    var authKey, encryptionKey []byte
    // ... key setup logic similar to main.go
    
    store := sessions.NewCookieStore(authKey, encryptionKey)
    store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   86400 * 7,
        HttpOnly: true,
        Secure:   false, // Set to false for testing
        SameSite: http.SameSiteLaxMode,
    }
    
    // Add session middleware to router
    r.Use(sessions.Sessions("free2free-session", store))
    
    // Register routes
    RegisterRoutes(r)
    
    return &TestServer{Router: r}
}
```

### 2. Fix Missing API Routes
Ensure all documented routes are registered in the test server:

```go
// In test utilities
func RegisterRoutes(r *gin.Engine) {
    // OAuth authentication routes
    r.GET("/auth/:provider", handlers.OauthBegin)
    r.GET("/auth/:provider/callback", handlers.OauthCallback)
    r.GET("/logout", handlers.Logout)
    r.GET("/auth/token", handlers.ExchangeToken)
    r.POST("/auth/refresh", handlers.RefreshTokenHandler)
    
    // Profile route
    r.GET("/profile", middleware.AuthRequired(), handlers.Profile)
    
    // Other routes as documented
    routes.SetupAdminRoutes(r)
    routes.SetupUserRoutes(r)
    routes.SetupOrganizerRoutes(r)
    routes.SetupReviewRoutes(r)
    routes.SetupReviewLikeRoutes(r)
}
```

### 3. Fix Build Errors
Update test code to match actual implementation:

```go
// Update JWT claims handling
func validateJWTToken(tokenString string) (*jwt.Token, error) {
    // Use the same claims structure as the actual implementation
    type Claims struct {
        UserID uint   `json:"user_id"`
        Email  string `json:"email"`
        jwt.RegisteredClaims
    }
    
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return token, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### 4. Fix OAuth Flow Issues
Ensure proper environment configuration in tests:

```go
// In test setup
func setupTestEnvironment() {
    // Set required environment variables for testing
    if os.Getenv("SESSION_KEY") == "" {
        os.Setenv("SESSION_KEY", "test-key-for-oauth-flow")
    }
    if os.Getenv("JWT_SECRET") == "" {
        os.Setenv("JWT_SECRET", "test-jwt-secret-for-testing")
    }
    
    // Initialize OAuth providers for testing if needed
    // ...
}
```

## Running Tests After Fixes

### 1. Run All Tests
```bash
cd free2free
go test ./tests/... -v
```

### 2. Run Specific Test Types
```bash
# Unit tests
go test ./tests/unit/... -v

# Integration tests
go test ./tests/integration/... -v

# API tests
go test ./tests/api/... -v

# End-to-end tests
go test ./tests/e2e/... -v

# Security tests
go test ./tests/security/... -v
```

### 3. Run with Coverage
```bash
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices for Test Fixes

### 1. Environment Consistency
- Ensure test environment mirrors production as closely as possible
- Use appropriate configuration for different test scenarios
- Avoid hardcoding values in tests

### 2. Error Handling
- Properly handle session initialization failures
- Provide meaningful error messages without exposing sensitive information
- Test both success and failure paths

### 3. Security Validation
- Ensure session handling maintains security properties
- Validate JWT tokens appropriately
- Test OAuth flow with both valid and invalid credentials

### 4. Performance Validation
- Verify that fixes don't introduce performance regressions
- Ensure OAuth flow completes within 10 seconds as required
- Maintain API response times under 500ms

## Validation Checklist

Before considering the fixes complete:

- [ ] All authentication tests pass without runtime panics
- [ ] All documented API endpoints return appropriate responses
- [ ] All test files compile without errors
- [ ] OAuth flow completes successfully in test environment
- [ ] Session handling works correctly across requests
- [ ] JWT validation works without panics
- [ ] Performance requirements are maintained
- [ ] Security requirements are maintained