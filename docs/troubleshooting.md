# Troubleshooting Guide

This document provides solutions to common issues encountered when using the free2free application, particularly those related to the platform-independent database implementation.

## Common Issues

### 1. Build Errors with CGO Dependencies

**Problem**: Application fails to build with CGO-related errors.

**Symptoms**: 
- Build fails with errors about missing CGO tools
- Error messages mentioning `github.com/mattn/go-sqlite3`

**Solution**:
1. Ensure you're using the pure-Go SQLite driver:
   ```bash
   go get modernc.org/sqlite
   ```

2. Set the environment variable during build:
   ```bash
   CGO_ENABLED=0 go build .
   ```

3. Verify your imports include the pure-Go driver:
   ```go
   import _ "modernc.org/sqlite"  // This registers the pure-Go driver
   ```

### 2. Database Connection Failures

**Problem**: Application can't connect to the SQLite database.

**Symptoms**:
- Error messages about database driver not found
- Connection errors during application startup

**Solution**:
1. Verify the SQLite driver is properly imported in your code:
   ```go
   import (
       "gorm.io/driver/sqlite"
       "gorm.io/gorm"
       _ "modernc.org/sqlite"  // Import for side effects only
   )
   ```

2. Check the database connection string is valid:
   - For file-based: `sqlite.Open("data.db")`
   - For in-memory: `sqlite.Open("file::memory:?cache=shared")`

3. Ensure the application has write permissions to the database directory.

### 3. Test Failures in Minimal Environments

**Problem**: Tests fail in containerized or minimal environments.

**Symptoms**:
- Tests pass locally but fail in Docker containers
- Tests work with CGO enabled but fail with CGO disabled

**Solution**:
1. Run tests with CGO explicitly disabled:
   ```bash
   CGO_ENABLED=0 go test ./tests/... -v
   ```

2. Ensure your test setup uses the platform-independent database:
   - Use in-memory databases for tests: `"file::memory:?cache=shared"`
   - Import the pure-Go driver in test files

3. Check that all database operations work identically with the pure-Go implementation.

### 4. Performance Differences

**Problem**: Database operations are slower than expected.

**Symptoms**:
- Increased latency for database operations
- Slower test execution times

**Solution**:
1. Note that pure-Go implementations may have slightly different performance characteristics (within 20% of CGO version).
2. For performance-critical operations, consider optimizing queries or migrating to a server-based database in production.
3. Use database connection pooling to improve performance:
   ```go
   sqlDB, err := db.DB()
   if err != nil {
       return err
   }
   sqlDB.SetMaxOpenConns(25)
   sqlDB.SetMaxIdleConns(25)
   sqlDB.SetConnMaxLifetime(5 * time.Minute)
   ```

### 5. Missing SQLite Features

**Problem**: Specific SQLite features don't work as expected.

**Symptoms**:
- Certain SQL functions not supported
- Extensions not available

**Solution**:
1. Note that the pure-Go implementation may not support all SQLite extensions.
2. Verify that your SQL is compatible with standard SQLite syntax.
3. Check the [modernc.org/sqlite documentation](https://pkg.go.dev/modernc.org/sqlite) for supported features.
4. If you require specific extensions, consider using a server-based database instead.

### 6. Import Registration Issues

**Problem**: Database driver not registered properly.

**Symptoms**:
- Driver not found errors
- Application panics on database operations

**Solution**:
1. Ensure the import is in the main package or in an imported package:
   ```go
   import _ "modernc.org/sqlite"  // Import for side effects only
   ```

2. Verify there are no import conflicts with other SQLite drivers.

3. Make sure the import happens early in the package initialization sequence.

## Verification Steps

### Verify Pure-Go Implementation
To verify your application is using the pure-Go SQLite implementation:

```bash
# Check if build works without CGO
CGO_ENABLED=0 go build .

# Run tests without CGO
CGO_ENABLED=0 go test ./tests/... -v

# Check for CGO dependencies in the build
go build -a -v .
```

### Check Database Driver
To verify the correct driver is being used:

```go
// In your application or test
import (
    "database/sql"
    _ "modernc.org/sqlite"  // Pure-Go driver
)

var registeredDrivers []string
for i := 0; ; i++ {
    driver := sql.Drivers()
    if i >= len(driver) {
        break
    }
    registeredDrivers = append(registeredDrivers, driver[i])
}
// Check that "sqlite" is in the list
```

## Best Practices

### 1. Import Registration
- Always import the pure-Go driver with `_ "modernc.org/sqlite"`
- Place the import early in your package initialization
- Avoid conflicts with other SQLite drivers

### 2. Testing Environments
- Test with `CGO_ENABLED=0` to ensure platform independence
- Use in-memory databases for faster test execution
- Verify all database operations work identically

### 3. Error Handling
- Handle database errors consistently across different implementations
- Provide meaningful error messages for debugging
- Log database operations for troubleshooting

### 4. Performance Monitoring
- Monitor performance differences between implementations
- Use connection pooling to improve performance
- Consider query optimization if needed