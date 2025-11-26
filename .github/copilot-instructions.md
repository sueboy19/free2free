# Free2Free AI Coding Agent Instructions

## Project Overview
This is a Go-based "buy one get one free" matching web application using Gin framework with OAuth authentication (Facebook/Instagram), JWT tokens, and MariaDB backend. The codebase follows hexagonal architecture with comprehensive testing.

## Key Architecture Patterns

### Hexagonal Architecture Structure
- **Core**: Business logic in `models/` and `utils/`
- **Adapters**: `database/` (GORM), `handlers/` (HTTP), `routes/` (routing)
- **Entry Points**: HTTP API endpoints in `routes/` packages

### Authentication Flow
```go
// OAuth → Session → JWT pattern
GET /auth/:provider → Facebook/Instagram OAuth → Session storage → JWT generation
```
- Uses Goth library for OAuth 2.0 flows
- Session management with Gorilla Sessions
- JWT tokens with 15-minute validity + refresh tokens
- Admin users flagged with `IsAdmin` field

### Error Handling Pattern
```go
// Standardized error structure
type AppError struct {
    Code      int    `json:"code"`
    Message   string `json:"error"`
    ErrorCode string `json:"code_error,omitempty"`
}
// Usage: c.Error(apperrors.NewValidationError("invalid input"))
```

## Critical Developer Workflows

### Environment Setup
```bash
# Required environment variables
cp .env.example .env
# Set: SESSION_KEY, DB_USER, DB_PASSWORD, DB_NAME, FACEBOOK_KEY, FACEBOOK_SECRET, BASE_URL

# Database setup
docker-compose up -d  # MariaDB on port 3306
go mod tidy           # Install dependencies
```

### Build & Test Commands
```bash
# Platform-independent build (no CGO)
CGO_ENABLED=0 go build -o free2free-staging .

# Docker build
docker build -t free2free .

# Comprehensive test suite
go test -v ./tests/...
go test -v ./tests/e2e/  # End-to-end tests
go test -v ./tests/contract/  # API contract tests

# Run with auto-reload (development)
air  # Requires air.toml config
```

### Database Operations
```bash
# Auto-migrate with environment flag
AUTO_MIGRATE=true go run .
# Tables: users, admins, locations, activities, matches, match_participants, reviews, review_likes
```

## Project-Specific Conventions

### API Response Patterns
```go
// Standard success response
c.JSON(http.StatusOK, gin.H{
    "data": result,
    "total": count,
    "page": page,
    "limit": limit,
})

// Error response
c.Error(apperrors.NewValidationError("invalid input"))
c.Abort()
```

### Route Organization
- `routes/admin.go` - Admin endpoints (requires IsAdmin flag)
- `routes/organizer.go` - Match organizer operations
- `routes/review.go` - Review and rating system
- `routes/user.go` - User-specific operations
- `routes/auth_handlers.go` - OAuth authentication

### Testing Structure
- `tests/unit/` - Unit tests for individual functions
- `tests/integration/` - Database integration tests
- `tests/contract/` - API contract validation tests
- `tests/e2e/` - Full application workflow tests
- `tests/testutils/` - Shared test utilities

### Model Validation
```go
// Using go-playground/validator
type User struct {
    Name string `validate:"required,min=1,max=100"`
    Email string `validate:"required,email"`
    SocialProvider string `validate:"required,oneof=facebook instagram"`
}
```

## Integration Points

### External Dependencies
- **OAuth**: Facebook/Instagram via Goth library
- **Database**: MariaDB (production) + SQLite (testing, pure-Go)
- **JWT**: golang-jwt/jwt/v5 with HMAC signing
- **CORS**: gin-contrib/cors for cross-origin requests
- **Swagger**: swaggo for API documentation

### Cross-Component Communication
- **Session → JWT**: Authentication state flows from OAuth session to JWT tokens
- **Middleware Chain**: CORS → Sessions → Authentication → Route handlers
- **Database Abstraction**: `database.GlobalDB` interface for testability

### Key Files for Understanding
- `main.go` - Application entry point with middleware setup
- `database/db.go` - Database interface and implementations
- `utils/auth.go` - Authentication utilities and JWT handling
- `errors/errors.go` - Centralized error handling
- `models/models.go` - Data models with validation tags

## Performance Considerations
- **Response Time**: Target <500ms for API endpoints
- **JWT Validation**: Target <10ms token verification
- **OAuth Flow**: Target <30 seconds for Facebook login
- **Database**: GORM with connection pooling, auto-migration enabled

## Security Patterns
- **Input Validation**: Struct-level validation with go-playground/validator
- **Password Handling**: bcrypt hashing for admin passwords
- **Session Security**: Encrypted cookie storage with rotating keys
- **CSRF Protection**: Built into Gorilla Sessions
- **Rate Limiting**: Consider adding for production

## Development Tips
- Use `CGO_ENABLED=0` for cross-platform compatibility
- Run `go test -v ./tests/contract/` to verify API contracts
- Check `swagger.yaml` for complete API documentation
- Use `tests/testutils/` for shared test data and helpers
- Follow the hexagonal architecture when adding new features