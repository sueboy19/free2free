# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Free2Free** is a "Buy One Get One Free" (BOGO) matching website built with Go (Gin framework) backend and Vue.js frontend. Users can create and join matching groups for various promotions via OAuth authentication (Facebook/Instagram).

**Tech Stack:**
- Backend: Go 1.25 + Gin + GORM + MariaDB (SQLite for testing)
- Frontend: Vue 3 + TypeScript + Vite + Tailwind CSS
- Authentication: OAuth 2.0 (Facebook/Instagram) via Goth + JWT tokens
- Documentation: Swagger UI at `/swagger/index.html`

## Common Commands

### Development
```bash
# Run with hot-reload (recommended for development)
air

# Run directly
go run .

# Build
go build
CGO_ENABLED=0 go build -o free2free-staging .

# Database (MariaDB via Docker)
docker-compose up -d    # Start database
docker-compose down     # Stop database
```

### Testing
```bash
# Run all tests
go test ./tests/... -v

# Platform-independent tests (no CGO dependency)
CGO_ENABLED=0 go test ./tests/... -v

# Specific test types
go test ./tests/unit/... -v          # Unit tests
go test ./tests/integration/... -v   # Integration tests
go test ./tests/contract/... -v      # API contract tests
go test ./tests/e2e/... -v           # End-to-end tests
go test ./tests/performance/... -v   # Performance tests

# With coverage
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Frontend
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Start dev server
npm run build        # Build for production
```

## Architecture

### Hexagonal Architecture
The project follows hexagonal architecture with clear separation of concerns:

- **Core Layer**: `models/` (domain models with GORM), `utils/` (business logic)
- **Adapter Layer**: `database/` (GORM interface), `handlers/` (HTTP handlers)
- **Entry Layer**: `routes/` (API routing by domain)

### Authentication Flow
```
OAuth Login (Facebook/Instagram)
    ↓
Session Storage (Gorilla Sessions)
    ↓
JWT Token Exchange (GET /auth/token)
    ↓
API Access (Bearer token in Authorization header)
```

- JWT tokens are valid for 15 minutes
- Admin users flagged with `IsAdmin` field in `Admin` model

### Database Design
Eight main entities managed by GORM auto-migration:
- `users` - OAuth users (social_id, provider: facebook/instagram)
- `admins` - System administrators (password hashed with bcrypt)
- `locations` - Physical locations for activities
- `activities` - BOGO promotion activities
- `matches` - Individual match instances
- `match_participants` - Participants with approval status (pending/approved/rejected)
- `reviews` - Ratings and comments between matched users
- `review_likes` - Like/dislike functionality for reviews

### Route Organization
- `routes/auth_handlers.go` - OAuth endpoints (`/auth/:provider`, `/auth/:provider/callback`, `/logout`, `/auth/token`)
- `routes/user.go` - User profile (`/profile`)
- `routes/admin.go` - Admin endpoints (CRUD for activities/locations, requires `IsAdmin`)
- `routes/organizer.go` - Match organizer operations (`/api/matches`)
- `routes/review.go` - Review system (`/api/reviews`)
- `routes/review_like.go` - Review likes (`/api/review-likes`)

## Standard Patterns

### Error Handling
```go
// Use centralized error types from errors/errors.go
c.Error(apperrors.NewValidationError("invalid input"))
c.Error(apperrors.NewUnauthorizedError("authentication required"))
c.Error(apperrors.NewNotFoundError("resource not found"))
c.Abort()
```

### API Response Pattern
```go
// Success response with pagination
c.JSON(http.StatusOK, gin.H{
    "data": result,
    "total": count,
    "page": page,
    "limit": limit,
})

// Simple success response
c.JSON(http.StatusOK, gin.H{
    "message": "操作成功",
    "id": createdID,
})
```

### Model Validation
```go
// Uses go-playground/validator tags
type User struct {
    Name string `validate:"required,min=1,max=100"`
    Email string `validate:"required,email"`
    SocialProvider string `validate:"required,oneof=facebook instagram"`
}
```

## Environment Setup

Copy `.env.example` to `.env` and configure:
- `SESSION_KEY` - 32+ character session encryption key
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - MariaDB credentials
- `FACEBOOK_KEY`, `FACEBOOK_SECRET` - Facebook OAuth app credentials
- `INSTAGRAM_KEY`, `INSTAGRAM_SECRET` - Instagram OAuth app credentials
- `BASE_URL` - Application base URL (e.g., `http://localhost:8080`)
- `JWT_SECRET` - JWT signing secret

## Testing Strategy

- **Platform Independence**: Tests use SQLite (pure Go via `modernc.org/sqlite`) instead of MariaDB to avoid CGO dependency
- **Test Organization**: `tests/unit/`, `tests/integration/`, `tests/contract/`, `tests/e2e/`, `tests/performance/`
- **Test Utilities**: Shared helpers in `tests/testutils/` (mock DB, JWT helpers, test data generators)
- **Coverage Target**: 95% API endpoint coverage

## Key Files

- `main.go` - Application entry point with middleware setup (CORS, sessions, auth, error handling)
- `database/db.go` - Database interface abstraction for testability
- `models/models.go` - GORM models with validation tags
- `handlers/auth_handlers.go` - OAuth authentication flow handlers
- `middleware/error_handler.go` - Global error handling middleware
- `middleware/session.go` - Session management middleware
- `utils/auth.go` - JWT token generation and validation
- `errors/errors.go` - Centralized error types

## OAuth Testing in Swagger UI

1. In browser: `http://localhost:8080/auth/facebook` (login via OAuth)
2. After login, visit `http://localhost:8080/auth/token` to get JWT token
3. In Swagger UI, click "Authorize" button
4. Enter: `Bearer <your_token>` (note the space after "Bearer")
5. Execute authenticated endpoints

## Frontend Development

The Vue.js frontend in `frontend/` uses:
- **State Management**: Pinia stores
- **Routing**: Vue Router 4
- **HTTP Client**: Axios with interceptors for JWT token handling
- **Styling**: Tailwind CSS
- **Notifications**: Vue Toastification

Backend API base URL: `http://localhost:8080`
