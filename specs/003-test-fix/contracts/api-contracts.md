# API Contract: Test Fix

## Overview
This document defines the API contracts for the fixed functionality in the API testing system. The contracts ensure that all documented endpoints are properly implemented and accessible, authentication flows work correctly, and session handling is consistent.

## Authentication Endpoints

### GET /auth/:provider
Initiate OAuth login flow

**Request**:
- Method: GET
- Path Parameters: provider (facebook|instagram)
- Headers: None required
- Query Parameters: None
- Body: None

**Response**:
- Success (307): Redirect to OAuth provider
- Error (400):
  ```json
  {
    "error": "invalid provider"
  }
  ```

### GET /auth/:provider/callback
OAuth callback endpoint

**Request**:
- Method: GET
- Path Parameters: provider (facebook|instagram)
- Query Parameters: OAuth callback parameters
- Headers: None required
- Body: None

**Response**:
- Success (307): Redirect to frontend with session
- Error (400):
  ```json
  {
    "error": "oauth failed"
  }
  ```

### GET /auth/token
Exchange OAuth session for JWT token

**Request**:
- Method: GET
- Headers: None required
- Query Parameters: None
- Body: None

**Response**:
- Success (200):
  ```json
  {
    "token": "JWT token string",
    "user": {
      "id": "user ID",
      "email": "user email",
      "name": "user name",
      "provider": "oauth provider",
      "role": "user role"
    }
  }
  ```
- Unauthorized (401):
  ```json
  {
    "error": "authentication failed"
  }
  ```

### POST /auth/refresh
Refresh expired JWT token

**Request**:
- Method: POST
- Headers: 
  - Content-Type: application/json
- Body:
  ```json
  {
    "refresh_token": "refresh token string"
  }
  ```

**Response**:
- Success (200):
  ```json
  {
    "token": "new JWT token",
    "refresh_token": "new refresh token"
  }
  ```
- Unauthorized (401):
  ```json
  {
    "error": "invalid refresh token"
  }
  ```

### GET /logout
Logout user and clear session

**Request**:
- Method: GET
- Headers: None required
- Query Parameters: None
- Body: None

**Response**:
- Success (200):
  ```json
  {
    "message": "logged out"
  }
  ```

## User Endpoints

### GET /profile
Get user profile information

**Request**:
- Method: GET
- Headers: 
  - Authorization: Bearer <JWT token>
- Query Parameters: None
- Body: None

**Response**:
- Success (200):
  ```json
  {
    "id": "user ID",
    "email": "user email",
    "name": "user name",
    "provider": "oauth provider",
    "avatar": "avatar URL"
  }
  ```
- Unauthorized (401):
  ```json
  {
    "error": "unauthorized"
  }
  ```

## Error Response Format

For all error conditions, the API should return consistent error responses:

- 400 Bad Request:
  ```json
  {
    "error": "descriptive error message",
    "code": "ERROR_CODE"
  }
  ```

- 401 Unauthorized:
  ```json
  {
    "error": "authentication required",
    "code": "AUTH_REQUIRED"
  }
  ```

- 403 Forbidden:
  ```json
  {
    "error": "insufficient permissions",
    "code": "FORBIDDEN"
  }
  ```

- 404 Not Found:
  ```json
  {
    "error": "resource not found",
    "code": "NOT_FOUND"
  }
  ```

- 500 Internal Server Error:
  ```json
  {
    "error": "internal server error",
    "code": "INTERNAL_ERROR"
  }
  ```

## Session Handling Contract

### Session Initialization
- Sessions must be properly initialized before access
- Accessing uninitialized session must return 500 error, not cause runtime panic
- Session data must be validated before use

### JWT Token Validation
- JWT tokens must be validated for expiration
- Invalid tokens must return 401 error
- Token claims structure must be consistent with implementation
- Token validation must not cause runtime panics

## Test Environment Contract

### Test Server Setup
- All documented API endpoints must be registered in test server
- Test server must initialize session handling similar to production
- Environment variables required for authentication must be available in test context

### Test Execution 
- All test files must compile without errors
- Tests must execute without build or runtime errors
- Test completion rate must meet success criteria defined in spec