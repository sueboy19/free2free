# Data Model: Complete API Testing

## Overview
This document describes the key data entities and their relationships relevant to API testing for the complete system workflow. The models are derived from the existing database models in the application and the test requirements.

## Key Entities

### User Session
Represents an authenticated user's interaction with the system, including their permissions and session token

- **Entity**: User Session (Runtime concept, not persisted)
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
  - ID: int (primary key)
  - Email: string (user's email)
  - Name: string (user's display name)
  - Provider: string (OAuth provider)
  - ProviderID: string (ID from OAuth provider)
  - Avatar: string (URL to user's avatar)
  - Role: string (user role: "user", "admin", etc.)
  - CreatedAt: time.Time
  - UpdatedAt: time.Time

### Free2Free Item (Activity)
Represents the core content in the system - a free2free item that can be created, managed, and approved

- **Entity**: Activity
- **Fields**:
  - ID: int (primary key)
  - Title: string (title of the activity)
  - Description: string (detailed description)
  - Status: string (status: "pending", "approved", "rejected", "active")
  - CreatorID: int (foreign key to User)
  - LocationID: int (foreign key to Location)
  - CreatedAt: time.Time
  - UpdatedAt: time.Time
  - ApprovedAt: *time.Time (nullable, when approved)
  - ApprovedBy: *int (nullable, admin who approved)

### Admin Review
Represents the moderation workflow for free2free items, including status changes and approval/rejection actions

- **Entity**: Admin (for administrative actions)
- **Fields**:
  - ID: int (primary key)
  - UserID: int (foreign key to User)
  - Permissions: string (admin permissions)
  - CreatedAt: time.Time
  - UpdatedAt: time.Time
  
Also includes review-related data stored in associated entities:
- Activity status changes
- Approval/rejection notes
- Audit trail of admin actions

### Location
Represents locations associated with free2free items

- **Entity**: Location
- **Fields**:
  - ID: int (primary key)
  - Name: string
  - Address: string
  - Latitude: float64
  - Longitude: float64
  - CreatedAt: time.Time
  - UpdatedAt: time.Time

## Relationships

1. **User → Activity** (One-to-Many)
   - A user can create multiple free2free items
   - Foreign key: Activity.CreatorID → User.ID

2. **Activity → Location** (Many-to-One)
   - Multiple activities can be associated with one location
   - Foreign key: Activity.LocationID → Location.ID

3. **Admin → Activity** (One-to-Many for approvals)
   - An admin can approve multiple activities
   - Foreign key: Activity.ApprovedBy → Admin.ID (nullable)

## State Transitions (for Activity model)

### Status Transitions
- **pending** → **approved** (when admin approves)
- **pending** → **rejected** (when admin rejects)
- **rejected** → **pending** (when admin reconsiders)
- **approved** → **active** (when activity goes live)
- **active** → **inactive** (when activity ends)

## Validation Rules from Requirements

1. **Activity Creation**:
   - Title must not be empty
   - Description must be at least 10 characters
   - Location must exist
   - Creator must be authenticated

2. **Activity Management**:
   - Only admins can change status from pending to approved/rejected
   - Only the creator can modify their own activities in pending state
   - Approved activities cannot be modified by creators

3. **Authentication**:
   - All API calls require valid JWT token (except login endpoints)
   - Admin endpoints require admin role

4. **Data Validation**:
   - All input fields must be validated for security (XSS, SQL injection)
   - Email format validation for user accounts

## Test Data Requirements

For comprehensive API testing, the following test data patterns are required:

1. **Valid Data Sets**: Complete, valid data for successful operations
2. **Invalid Data Sets**: Invalid inputs to test validation
3. **Edge Cases**: Boundary conditions, special characters, large inputs
4. **Security Test Data**: Malicious inputs to test security measures
5. **Permission-Specific Data**: Different user roles to test authorization
6. **Historical Data**: Past activities, expired sessions for state testing