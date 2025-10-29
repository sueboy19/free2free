# API Contract: Complete API Testing

## Overview
This document defines the API contracts for testing the complete system workflow: login, creating free2free items, management, and approval processes.

## Authentication Endpoints

### POST /auth/token
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

### GET /auth/:provider
Initiate OAuth login flow

**Request**:
- Method: GET
- Path Parameters: provider (facebook|instagram)
- Headers: None required
- Query Parameters: None
- Body: None

**Response**:
- Redirect (307): To OAuth provider
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

## Free2Free Item Endpoints

### POST /api/activities
Create a new free2free item

**Request**:
- Method: POST
- Headers: 
  - Authorization: Bearer <JWT token>
  - Content-Type: application/json
- Body:
  ```json
  {
    "title": "Activity title",
    "description": "Activity description",
    "location_id": "ID of location"
  }
  ```

**Response**:
- Success (201):
  ```json
  {
    "id": "new activity ID",
    "title": "Activity title",
    "description": "Activity description",
    "location_id": "ID of location",
    "status": "pending",
    "creator_id": "ID of creator",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Unauthorized (401):
  ```json
  {
    "error": "unauthorized"
  }
  ```
- Validation Error (400):
  ```json
  {
    "error": "validation failed",
    "details": ["title is required", "description too short"]
  }
  ```

### GET /api/activities/:id
Get a specific free2free item

**Request**:
- Method: GET
- Headers: 
  - Authorization: Bearer <JWT token>
- Path Parameters: id (activity ID)

**Response**:
- Success (200):
  ```json
  {
    "id": "activity ID",
    "title": "Activity title",
    "description": "Activity description",
    "location_id": "ID of location",
    "status": "activity status",
    "creator_id": "ID of creator",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```

### PUT /api/activities/:id
Update a specific free2free item

**Request**:
- Method: PUT
- Headers: 
  - Authorization: Bearer <JWT token>
  - Content-Type: application/json
- Path Parameters: id (activity ID)
- Body:
  ```json
  {
    "title": "Updated activity title",
    "description": "Updated activity description",
    "location_id": "Updated location ID"
  }
  ```

**Response**:
- Success (200):
  ```json
  {
    "id": "activity ID",
    "title": "Updated activity title",
    "description": "Updated activity description",
    "location_id": "Updated location ID",
    "status": "activity status",
    "creator_id": "ID of creator",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Unauthorized (401):
  ```json
  {
    "error": "unauthorized"
  }
  ```
- Forbidden (403):
  ```json
  {
    "error": "forbidden - cannot modify this activity"
  }
  ```

## Administrative Endpoints

### GET /admin/activities
Get all free2free items (for admin management)

**Request**:
- Method: GET
- Headers: 
  - Authorization: Bearer <JWT token>
  - Requires admin role
- Query Parameters: 
  - status (optional): filter by status
  - page (optional): pagination page
  - limit (optional): items per page

**Response**:
- Success (200):
  ```json
  {
    "activities": [
      {
        "id": "activity ID",
        "title": "Activity title",
        "description": "Activity description",
        "location_id": "ID of location",
        "status": "activity status",
        "creator_id": "ID of creator",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "total": "total count",
    "page": "current page",
    "limit": "items per page"
  }
  ```

### PUT /admin/activities/:id/approve
Approve a free2free item

**Request**:
- Method: PUT
- Headers: 
  - Authorization: Bearer <JWT token>
  - Requires admin role
- Path Parameters: id (activity ID)
- Body: None

**Response**:
- Success (200):
  ```json
  {
    "id": "activity ID",
    "status": "approved",
    "approved_at": "timestamp",
    "approved_by": "admin ID"
  }
  ```
- Forbidden (403):
  ```json
  {
    "error": "forbidden - admin permissions required"
  }
  ```

### PUT /admin/activities/:id/reject
Reject a free2free item

**Request**:
- Method: PUT
- Headers: 
  - Authorization: Bearer <JWT token>
  - Requires admin role
- Path Parameters: id (activity ID)
- Body:
  ```json
  {
    "reason": "reason for rejection (optional)"
  }
  ```

**Response**:
- Success (200):
  ```json
  {
    "id": "activity ID",
    "status": "rejected",
    "rejection_reason": "rejection reason if provided"
  }
  ```
- Forbidden (403):
  ```json
  {
    "error": "forbidden - admin permissions required"
  }
  ```