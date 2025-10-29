# Quickstart Guide: Complete API Testing

## Overview
This guide provides instructions for implementing comprehensive API tests that cover the complete system workflow: login, creating free2free items, management, and approval processes.

## Prerequisites

1. Go 1.25.0 or higher installed
2. Project dependencies installed (`go mod download`)
3. Database setup and running (MariaDB/MySQL)
4. Environment variables configured (use `.env.example` as template)
5. Basic understanding of the Gin framework and Go testing

## Setting Up the Environment

### 1. Clone and Setup the Project
```bash
git clone <repository-url>
cd free2free
cp .env.example .env
# Update .env with your configuration
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Run Database Migrations
Set `AUTO_MIGRATE=true` in your environment to create the required tables.

## Testing Architecture

The project follows a modular testing approach:

- `tests/unit/` - Pure unit tests for individual functions
- `tests/integration/` - Integration tests for multiple components
- `tests/contract/` - API contract validation tests
- `tests/api/` - Tests for API endpoints
- `tests/e2e/` - End-to-end workflow tests
- `tests/testutils/` - Test utilities and helpers

## Implementing Complete API Workflow Tests

### 1. End-to-End Workflow Test
Create a test that covers the complete user journey:

```go
// tests/e2e/complete_api_workflow_test.go
package e2e

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "free2free/main" // adjust import path as needed
)

func TestCompleteAPIWorkflow(t *testing.T) {
    gin.SetMode(gin.TestMode)

    // Setup test router
    router := setupRouter()

    // 1. Test login/oauth flow
    // 2. Create free2free item
    // 3. Admin approval
    // 4. Verify end result
}

func setupRouter() *gin.Engine {
    // Setup your router with middleware and routes
    r := gin.Default()
    // Add your routes here
    return r
}
```

### 2. Authentication Flow Tests
Test the complete authentication flow:

```go
// tests/integration/auth_flow_test.go
package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAuthenticationFlow(t *testing.T) {
    gin.SetMode(gin.TestMode)

    router := gin.Default()
    
    // Add authentication routes to router
    // router.GET("/auth/:provider", handlers.OauthBegin)
    // router.GET("/auth/:provider/callback", handlers.OauthCallback)
    // router.GET("/auth/token", handlers.ExchangeToken)
    // router.GET("/logout", handlers.Logout)

    // Create test request
    req, _ := http.NewRequest("GET", "/auth/facebook", nil)
    w := httptest.NewRecorder()
    
    // Execute request
    router.ServeHTTP(w, req)

    // Assert expected response
    assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}
```

### 3. API Contract Tests
Ensure API endpoints meet the contract specifications:

```go
// tests/contract/api_contract_test.go
package contract

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

type ActivityResponse struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
}

func TestActivityCreationContract(t *testing.T) {
    gin.SetMode(gin.TestMode)

    router := gin.Default()
    // Add your routes

    // Test creating an activity
    req, _ := http.NewRequest("POST", "/api/activities", 
        json.RawMessage(`{"title":"Test Activity", "description":"Test Description", "location_id":1}`))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer <valid-jwt-token>")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var response ActivityResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Test Activity", response.Title)
    assert.Equal(t, "pending", response.Status)
}
```

### 4. Performance Tests
Test that endpoints meet performance requirements:

```go
// tests/performance/api_performance_test.go
package performance

import (
    "testing"
    "time"
)

func TestAPITimeout(t *testing.T) {
    timeout := 500 * time.Millisecond

    // Measure API call duration
    start := time.Now()
    // Make API call
    // ...
    duration := time.Since(start)

    if duration > timeout {
        t.Errorf("API call took %v, expected less than %v", duration, timeout)
    }
}
```

## Running Tests

### 1. Run All Tests
```bash
go test ./tests/... -v
```

### 2. Run Specific Test Type
```bash
# Unit tests
go test ./tests/unit/... -v

# Integration tests
go test ./tests/integration/... -v

# API tests
go test ./tests/api/... -v

# End-to-end tests
go test ./tests/e2e/... -v
```

### 3. Run with Coverage
```bash
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices for API Testing

### 1. Test All User Scenarios
Implement tests for all user scenarios defined in the feature specification:
- Login flow (valid/invalid credentials)
- Creating free2free items (valid/invalid data)
- Management and approval (admin/non-admin access)

### 2. Handle Edge Cases
Test edge cases like:
- Expired tokens
- Concurrent requests
- Invalid data inputs
- Missing permissions

### 3. Maintain Test Independence
Each test should be independent and not rely on the state from other tests:
- Use setup/teardown functions to create clean state
- Avoid shared state between tests

### 4. Use Test Utilities
Leverage the existing test utilities in `tests/testutils/`:
- Test servers
- Mock data
- Helper functions

### 5. Validate Security Requirements
Ensure tests validate:
- Authentication requirements
- Authorization checks
- Input validation
- Session management

## Adding New Test Cases

1. Identify the new functionality or workflow to be tested
2. Determine which test type is appropriate (unit, integration, contract, API, e2e)
3. Create the test file in the appropriate directory
4. Follow existing code patterns and conventions
5. Ensure the test validates the feature requirements
6. Run the test and verify it passes