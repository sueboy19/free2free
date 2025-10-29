# Free2Free API Testing Documentation

## Overview
This document describes the testing setup and strategies implemented for the Free2Free API. The tests cover comprehensive workflows from user login to creating free2free items, management, and approval processes.

## Test Structure

The tests are organized in the following directory structure:

```
tests/
├── unit/               # Pure unit tests for individual functions
├── integration/        # Integration tests for multiple components
├── contract/           # API contract validation tests
├── api/                # Tests for API endpoints
├── e2e/                # End-to-end workflow tests
├── performance/        # Performance and load tests
└── testutils/          # Test utilities and helpers
```

## Running Tests

### All Tests
```bash
go test ./tests/... -v
```

### Specific Test Types
```bash
# Unit tests
go test ./tests/unit/... -v

# Integration tests
go test ./tests/integration/... -v

# Contract tests
go test ./tests/contract/... -v

# End-to-end tests
go test ./tests/e2e/... -v

# Performance tests
go test ./tests/performance/... -v
```

### With Coverage
```bash
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Utilities

The `testutils/` directory contains several helper packages:

- `config.go`: Test configuration management
- `test_server.go`: Test server setup utilities
- `mock_db.go`: Mock database functionality for tests
- `jwt_validator.go`: JWT token utilities for testing
- `api_helpers.go`: Common API test helpers
- `auth_test_helpers.go`: Authentication-specific test helpers
- `admin_test_helpers.go`: Admin-specific test helpers
- `test_data.go`: Test data generation utilities

## Authentication Testing

The authentication testing approach includes:

1. **Contract Tests**: Validate the API contracts for authentication endpoints
2. **Integration Tests**: Test the OAuth login flow and token generation
3. **Unit Tests**: Test JWT token generation and validation in isolation
4. **End-to-End Tests**: Test the complete login flow with valid and invalid credentials

## Activities Management Testing

The activities management testing includes:

1. **Contract Tests**: Validate the API contracts for activities endpoints
2. **Integration Tests**: Test the create, update, and retrieve functionality
3. **Validation Tests**: Test data validation for activity creation
4. **End-to-End Tests**: Test the complete activity creation workflow

## Admin Functionality Testing

The admin functionality testing includes:

1. **Contract Tests**: Validate the API contracts for admin endpoints
2. **Integration Tests**: Test the approval and rejection workflows
3. **Permission Tests**: Test role-based access controls
4. **Management Tests**: Test admin-specific management operations

## Testing Best Practices

### Test Isolation
- Each test is independent and does not rely on state from other tests
- Use setup/teardown functions to create clean state for each test
- Avoid shared state between tests

### Data Validation
- Validate all user inputs according to defined rules
- Test both valid and invalid input scenarios
- Ensure security measures are in place against injection attacks

### Performance Testing
- Ensure API endpoints respond within 500ms under normal load
- Test authentication and token operations for efficiency
- Validate database query performance

### Security Testing
- Ensure all endpoints require proper authentication
- Validate authorization checks for different user roles
- Verify that sensitive information is not exposed in error messages

## API Endpoints Covered

### Authentication Endpoints
- `POST /auth/token` - Exchange OAuth session for JWT token
- `GET /auth/:provider` - Initiate OAuth login flow
- `GET /auth/:provider/callback` - OAuth callback endpoint
- `POST /auth/refresh` - Refresh expired JWT token
- `GET /logout` - Logout user and clear session

### User Endpoints
- `GET /profile` - Get user profile information

### Activities Endpoints
- `POST /api/activities` - Create a new free2free item
- `GET /api/activities/:id` - Get a specific free2free item
- `PUT /api/activities/:id` - Update a specific free2free item

### Administrative Endpoints
- `GET /admin/activities` - Get all free2free items (for admin management)
- `PUT /admin/activities/:id/approve` - Approve a free2free item
- `PUT /admin/activities/:id/reject` - Reject a free2free item

## Quality Assurance

The testing strategy ensures:
- All major user workflows (login, create, manage, approve) can be completed through API calls with 99% success rate
- All API endpoints respond within 500ms under normal load conditions
- Comprehensive API test coverage of 95% of all endpoints and workflows
- All authentication and authorization requirements are validated through API tests with zero security vulnerabilities
- Users can successfully complete the entire process from login to creating a free2free item in under 3 minutes with appropriate feedback
- System administrators can manage and approve free2free items with 99% success rate through the API
- All API error conditions are properly handled and tested with appropriate user feedback