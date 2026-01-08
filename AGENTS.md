# AGENTS.md - Free2Free AI Agent Guidelines

## Project Overview

Free2Free is a Go-based "buy one get one free" matching web application using Gin framework with OAuth authentication (Facebook/Instagram), JWT tokens, and MariaDB backend. The project follows hexagonal architecture with comprehensive testing.

**Tech Stack:**
- Backend: Go 1.25 + Gin + GORM + MariaDB (production) + SQLite (testing, pure-Go via modernc.org/sqlite)
- Frontend: Vue 3 + TypeScript + Vite + Tailwind CSS
- Authentication: OAuth (Facebook/Instagram) + JWT tokens with 15min expiry + refresh tokens
- Testing: Go testing + testify + mock database (in-memory SQLite)

**Architecture:**
- Hexagonal architecture with Core (models/utils), Adapter (database/handlers/routes), Entry Points (HTTP API)

## Build Commands

### Go Backend

```bash
# Standard build
go build -o free2free.exe .

# Cross-platform build (no CGO - required for testing)
CGO_ENABLED=0 go build -o free2free.exe .

# Development with hot reload
air

# Docker build
docker build -t free2free .
```

### Frontend (cd frontend)

```bash
# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview
```

## Test Commands

### All Tests

```bash
# Run all tests
go test ./tests/... -v

# Platform-independent tests (no CGO - recommended)
CGO_ENABLED=0 go test ./tests/... -v
```

### Specific Test Types

```bash
# Integration tests
go test ./tests/integration/... -v

# Unit tests
go test ./tests/unit/... -v

# Contract tests
go test ./tests/contract/... -v

# End-to-end tests
go test ./tests/e2e/... -v

# Performance tests
go test ./tests/performance/... -v
```

### Single Test

```bash
# Run specific test function
go test -run TestFunctionName ./tests/integration/... -v

# Run tests in specific file
go test ./tests/integration/auth_integration_test.go -v -run TestOAuthLoginFlow
```

### Coverage

```bash
# Generate coverage report
go test ./tests/... -v -coverprofile=coverage.out

# View HTML coverage
go tool cover -html=coverage.out -o coverage.html

# View terminal coverage
go tool cover -func=coverage.out
```

## Code Style Guidelines

### Import Organization

```go
// Order: Standard library → Third-party → Local packages
import (
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/gorilla/sessions"
    "gorm.io/gorm"

    apperrors "free2free/errors"
    "free2free/database"
    "free2free/models"
)
```

**Note:** Use `apperrors` alias to avoid conflict with stdlib `errors` package.

### Naming Conventions

- **Exported functions/types**: CamelCase (e.g., `CreateUser`, `GetAuthenticatedUser`)
- **Private functions/types**: camelCase (e.g., `validateUser`, `parseToken`)
- **Files**: snake_case (e.g., `auth_handlers.go`, `user_routes.go`)
- **Constants**: UPPER_SNAKE_CASE or CamelCase depending on scope
- **Comments**: English or Chinese acceptable, maintain consistency

### Error Handling Pattern

```go
// ALWAYS use centralized error types from errors package
c.Error(apperrors.NewValidationError("invalid input"))
c.Error(apperrors.NewUnauthorizedError("authentication required"))
c.Error(apperrors.NewForbiddenError("admin only"))
c.Error(apperrors.NewInternalError("database error"))

// ALWAYS call c.Error() and return - NEVER use c.JSON() for errors
if err != nil {
    c.Error(apperrors.MapGORMError(err))
    return
}

// Middleware handles errors automatically
// Standardized error structure: {code, error, code_error}
```

**Critical:** Never use `c.JSON()` to return error responses. Always use `c.Error()` with `apperrors` types so that error middleware can handle them consistently.

### Architecture: Hexagonal Pattern

```
Core Layer (Business Logic & Domain Models):
  - models/        : Data models with validation tags
  - utils/         : Business utilities (auth, helpers)

Adapter Layer (External Dependencies):
  - database/       : GORM interface, DB connection
  - handlers/       : HTTP handlers, OAuth processing
  - middleware/     : Session, error handling, CORS

Entry Points (HTTP API):
  - routes/         : Route definitions, path organization
  - main.go         : App entry, middleware setup, server config
```

**Key Pattern:** Use `database.GlobalDB.Conn` interface for testability - never import database driver directly in handlers.

### Database Operations

```go
// Use GORM with GlobalDB interface
if err := database.GlobalDB.Conn.Create(&model).Error; err != nil {
    c.Error(apperrors.MapGORMError(err))
    return
}

// Preload related data for efficiency (avoid N+1 queries)
database.GlobalDB.Conn.
    Preload("Activity").
    Preload("Organizer").
    Preload("Location").
    Where("status = ? AND match_time > ?", "open", time.Now()).
    Find(&matches)

// Always map GORM errors using MapGORMError()
```

**Testing:** Use `modernc.org/sqlite` (pure-Go) with `CGO_ENABLED=0` for platform-independent testing.

### Authentication & Authorization

```go
// Get authenticated user from session or JWT
user, err := utils.GetAuthenticatedUser(c)
if err != nil {
    c.Error(apperrors.NewUnauthorizedError("not logged in"))
    return
}

// Check admin status
if !user.IsAdmin {
    c.Error(apperrors.NewForbiddenError("admin only"))
    return
}

// OAuth → Session → JWT flow
// 1. OAuth provider authenticates user
// 2. Session stores user_id in session
// 3. JWT generated (15min access + 7day refresh)
// 4. Token rotation on refresh
```

**Token Management:**
- Access token: 15 minutes expiry
- Refresh token: 7 days expiry
- Session supports 1000+ concurrent requests

### Validation Pattern

```go
// Use validator tags on models
type User struct {
    Name           string `validate:"required,min=1,max=100"`
    Email          string `validate:"required,email"`
    SocialProvider string `validate:"required,oneof=facebook instagram"`
}

// Validate in handlers
v := validator.New()
if err := v.Struct(&input); err != nil {
    c.Error(apperrors.NewValidationError(err.Error()))
    return
}

// Validator tags: required, min, max, email, url, oneof
```

### Testing Patterns

```go
// Use test utilities
ts := testutils.NewTestServer()
defer ts.Close()

db, err := testutils.CreateTestDB()
assert.NoError(t, err)

// Create authenticated requests
jwtToken := testutils.CreateMockJWTToken(userID, userName, isAdmin)
resp, err := testutils.MakeAuthenticatedRequest(ts, "GET", "/user/matches", jwtToken, nil)

// Test structure
func TestFunctionName(t *testing.T) {
    t.Run("Subtest description", func(t *testing.T) {
        // Setup
        ts := testutils.NewTestServer()
        defer ts.Close()

        // Execute
        w := httptest.NewRecorder()
        req := httptest.NewRequest("GET", "/endpoint", nil)

        // Assert
        assert.Equal(t, http.StatusOK, w.Code)
    })
}
```

**Key Utilities:**
- `testutils.NewTestServer()` - Create test server with Gin router
- `testutils.CreateTestDB()` - Create in-memory SQLite DB
- `testutils.CreateMockJWTToken()` - Generate JWT for testing
- `testutils.MakeAuthenticatedRequest()` - Make authenticated HTTP requests

### API Response Pattern

```go
// Success response
c.JSON(http.StatusOK, gin.H{
    "data": result,
    "total": count,
    "page": page,
})

// Error response (handled by middleware)
c.Error(apperrors.NewAppError(http.StatusBadRequest, "message"))
c.Abort()

// DO NOT use c.JSON() for errors - middleware will format correctly
```

**Standard Success Format:**
```json
{
  "data": {...},
  "total": 100,
  "page": 1
}
```

**Standard Error Format:**
```json
{
  "code": 400,
  "error": "error message",
  "code_error": "ERROR_CODE"
}
```

### Frontend Guidelines (Vue 3)

```typescript
// Use Composition API with <script setup>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
</script>

// TypeScript with strict mode
// Use proper typing for all props and reactive data

// State management with Pinia
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()

// API calls via services/api.ts
import api from '@/services/api'
await api.get('/user/matches')

// Route lazy loading
const Home = () => import('@/views/Home.vue')
```

## Performance Requirements

From specs/002-complete-api-testing and specs/001-fb-login-test-suite:

- **API Response Time**: Target < 500ms under normal load
- **JWT Validation**: Target < 10ms for token verification
- **OAuth Flow**: Target < 30 seconds for Facebook login complete
- **Session Management**: Support 1000+ concurrent requests
- **Test Coverage**: Target 95%+ of all endpoints and workflows
- **Success Rate**: Target 99% for login, create, management operations

## Critical Files

- `main.go` - App entry point, middleware setup, OAuth provider config
- `models/models.go` - Data models with validation tags
- `errors/errors.go` - Centralized error handling and types
- `utils/auth.go` - Authentication utilities, JWT validation
- `database/db.go` - Database interface and GlobalDB wrapper
- `handlers/auth_handlers.go` - OAuth flow, token generation, refresh
- `routes/*.go` - API route definitions (admin.go, user.go, organizer.go, review.go)
- `middleware/error_handler.go` - Global error handling middleware
- `middleware/session.go` - Session management middleware
- `tests/testutils/` - Shared test helpers (test_server.go, mock_db.go, jwt_validator.go)
- `tests/README.md` - Test structure and execution guide

## Important Constraints & Best Practices

### CGO Dependency

```bash
# ALWAYS use CGO_ENABLED=0 for testing to enable cross-platform compatibility
CGO_ENABLED=0 go test ./tests/... -v

# Use modernc.org/sqlite (pure-Go) instead of mattn/go-sqlite3
import _ "modernc.org/sqlite"
```

**Why:** Enables platform-independent testing and compilation without native build tools.

### Error Handling

- **ALWAYS** use `c.Error()` with `apperrors` types
- **NEVER** use `c.JSON()` to return error responses
- **ALWAYS** call `c.Abort()` after `c.Error()` in handlers
- Let middleware handle error formatting and HTTP status codes

### Import Aliases

- Use `apperrors` alias to avoid conflict with stdlib `errors`
- Use package aliases only when necessary for clarity

### Testing Environment

- Use in-memory SQLite for complete data isolation
- Each test should be independent with setup/teardown
- Use `defer ts.Close()` for cleanup
- Use `defer db.Close()` for database cleanup

### Authentication Flow

1. OAuth provider (Facebook/Instagram) authenticates user
2. Session stores `user_id` in session
3. JWT generated: 15-minute access token + 7-day refresh token
4. Refresh tokens stored in database, rotated on refresh
5. Session supports both session cookies and JWT tokens

### Security Considerations

- Input validation on all endpoints (validator tags)
- SQL injection prevention via GORM parameterized queries
- CSRF protection via Gorilla Sessions
- Token rotation on refresh (delete old, create new)
- Rate limiting consideration for production
- Secure cookies with HttpOnly, SameSite settings

## Development Workflow

### Making Changes

1. Write tests first (TDD approach recommended)
2. Implement following hexagonal architecture
3. Use `database.GlobalDB.Conn` interface
4. Use centralized error handling
5. Run tests with `CGO_ENABLED=0`
6. Verify no regression: `go test ./tests/... -v`

### Adding New Features

1. Define models in `models/models.go` with validation tags
2. Create handlers in `handlers/` or routes in `routes/`
3. Add routes in appropriate routes file
4. Write integration tests in `tests/integration/`
5. Update API documentation (Swagger)
6. Run full test suite: `go test ./tests/... -v`

### Common Issues & Solutions

- **CGO build fails**: Use `CGO_ENABLED=0` for cross-platform builds
- **Test import errors**: Ensure using `modernc.org/sqlite` not `mattn/go-sqlite3`
- **Session nil panic**: Check session middleware setup, ensure session always exists in context
- **JWT validation fails**: Verify JWT_SECRET environment variable is set and >=32 characters
- **Database connection fails**: Check MariaDB is running, verify .env DB_HOST setting

## Additional Resources

- **API Documentation**: `http://localhost:8080/swagger/index.html`
- **Test Guide**: `tests/README.md`
- **Database Design**: `database_design.md`
- **Security Design**: `security_design.md`
- **Frontend README**: `frontend/README.md`
- **Deployment**: `docs/deployment.md`
