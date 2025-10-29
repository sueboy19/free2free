# Quickstart Guide: Cross-Platform Testing Without Native Dependencies

## Overview
This guide provides instructions for implementing platform-independent database functionality that works without native compilation dependencies, enabling consistent testing across different environments.

## Prerequisites

1. Go 1.25.0 or higher installed
2. Project dependencies installed (`go mod download`)
3. Standard Go environment without CGO requirements (CGO_ENABLED=0)
4. Understanding of the current database testing patterns

## Setting Up the Environment

### 1. Install Required Dependencies
```bash
go get modernc.org/sqlite
go get gorm.io/driver/sqlite
go get gorm.io/gorm
```

### 2. Update Imports for Pure-Go SQLite
Add the following import to register the pure-Go SQLite driver:
```go
import _ "modernc.org/sqlite"
```

This import should be placed in files that initialize the database, such as:
- Main application entry point
- Test initialization files
- Database utility files

## Implementation Approach

### 1. Database Connection Setup
The connection setup remains largely the same, but with the driver replacement:

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    // Register pure-Go SQLite driver to replace CGO-based implementation
    _ "modernc.org/sqlite"
)

// Create database connection (no code changes needed in this function)
func createTestDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

### 2. Testing with Platform-Independent Database
Update your test utilities to ensure they work with the new driver:

```go
// In your test utils
package testutils

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "modernc.org/sqlite"
)

// Register pure-Go SQLite driver
import _ "modernc.org/sqlite"

func CreateTestDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    // Additional configuration as needed
    return db, err
}
```

## Running Tests

### 1. Run Tests with CGO Disabled
```bash
CGO_ENABLED=0 go test ./tests/... -v
```

### 2. Run Tests in Standard Environment
```bash
go test ./tests/... -v
```

### 3. Run with Coverage
```bash
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Migration Steps

### 1. Add Import Replacement
Add the import `_ "modernc.org/sqlite"` in the main package or in test initialization code to register the pure-Go SQLite driver.

### 2. Verify Compatibility
- Run existing test suite to ensure all functionality works as before
- Check all CRUD operations continue to work identically
- Verify migration operations function correctly

### 3. Test in CGO-Disabled Environment
- Run `CGO_ENABLED=0 go test` to verify tests work without CGO dependencies
- Test in minimal container environments if applicable

## Validation Checklist

### Before Migration
- [ ] All existing tests pass with current implementation
- [ ] Document current performance metrics
- [ ] Identify all uses of database functionality

### After Migration 
- [ ] All tests pass with new implementation
- [ ] Tests run successfully with CGO_ENABLED=0
- [ ] Performance degradation is within acceptable thresholds (≤20%)
- [ ] All CRUD operations work identically
- [ ] Database migration operations function correctly
- [ ] In-memory database functionality preserved

## Best Practices

### 1. Import Order Matters
- Ensure the import replacement (`_ "modernc.org/sqlite"`) happens early
- Place the import in the main package or in the entry point of your application

### 2. Testing Consistency
- Maintain the same test coverage as before
- Verify all database operations behave identically to previous implementation
- Test edge cases that might behave differently between implementations

### 3. Performance Monitoring
- Compare performance metrics between old and new implementations
- Monitor for any significant performance degradation
- Consider caching strategies if performance is impacted

### 4. Error Handling
- Ensure error handling patterns remain consistent
- Verify error messages are appropriate and helpful
- Check that all error conditions are properly handled

## Troubleshooting

### Common Issues
1. **Import Registration**: Ensure the import `_ "modernc.org/sqlite"` is registered before database initialization
2. **Feature Differences**: Some SQLite-specific features might behave differently between implementations
3. **Performance Variations**: Pure-Go implementations may have different performance characteristics

### Verification Commands
```bash
# Verify CGO is disabled
go env CGO_ENABLED

# Run tests with CGO explicitly disabled
CGO_ENABLED=0 go test ./...

# Check for CGO dependencies in the build
go build -a -v .
```