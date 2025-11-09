# Data Model: Test Fix

## Overview
This document describes the data model considerations for fixing test issues in the API testing functionality. The models are derived from the existing database models in the application and the test requirements for proper validation.

## Key Entities

### User Session (Runtime concept, not persisted)
Represents the user's authenticated state across requests, containing user identity and permissions

- **Entity**: User Session
- **Fields**:
  - sessionID: string (unique identifier for the session)
  - userID: int (reference to user ID)
  - token: string (JWT token)
  - permissions: []string (list of permissions)
  - createdAt: time.Time (when session was created)
  - expiresAt: time.Time (when session expires)
  - provider: string (OAuth provider, e.g., "facebook", "instagram")

### User
Represents a registered user in the system

- **Entity**: User
- **Fields**:
  - ID: uint (primary key)
  - Email: string (user's email)
  - Name: string (user's display name)
  - Provider: string (OAuth provider)
  - ProviderID: string (ID from OAuth provider)
  - Avatar: string (URL to user's avatar)
  - Role: string (user role: "user", "admin", etc.)
  - CreatedAt: time.Time
  - UpdatedAt: time.Time

### Authentication Token
Represents authenticated user identity that can be validated across requests

- **Entity**: Authentication Token
- **Fields**:
  - token: string (the JWT token string)
  - userID: int (user identifier)
  - issuedAt: time.Time (when token was issued)
  - expiresAt: time.Time (when token expires)
  - claims: map[string]interface{} (token claims)

### OAuth Provider
Represents external authentication service that validates user identity

- **Entity**: OAuth Provider
- **Fields**:
  - name: string (provider name: "facebook", "instagram")
  - clientID: string (OAuth client ID)
  - clientSecret: string (OAuth client secret)
  - callbackURL: string (OAuth callback URL)

## Relationships

1. **User → User Session** (One-to-Many)
   - A user can have multiple session records (though typically only one active)
   - Foreign key: Session.userID → User.ID

2. **OAuth Provider → User** (Many-to-Many for authentication history)
   - A user can authenticate through multiple OAuth providers
   - A provider can authenticate multiple users

## Validation Rules from Requirements

1. **Session Validation**:
   - Session must be properly initialized before use
   - Session data must not be nil when accessed
   - Session must have valid user ID

2. **Token Validation**:
   - JWT tokens must be properly formatted
   - Tokens must not be expired
   - Token claims must contain required fields (userID)

3. **User Data Validation**:
   - User data must be consistent between test and implementation
   - Required fields (Provider, ProviderID) must be present in User model

## Test Data Requirements

For comprehensive testing of the fixes, the following test data patterns are required:

1. **Valid Session Data**: Complete, valid session data for successful operations
2. **Invalid Session Data**: Nil or malformed session data to test error handling
3. **Expired Token Data**: Tokens that have exceeded their validity period
4. **Invalid Token Data**: Malformed or unsigned JWTs for validation testing
5. **Missing User Data**: User records with missing required fields for validation
6. **OAuth Provider Data**: Complete provider configurations for OAuth testing
7. **Environment Configuration**: Proper environment variables for session handling