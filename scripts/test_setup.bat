@echo off
REM Test setup script for Facebook Login and API Test Suite
REM This script sets up the test environment for local testing

echo Setting up test environment for Facebook Login and API tests...

REM Set test environment variables
echo Setting test environment variables...
set TEST_DB_HOST=localhost
set TEST_DB_PORT=3306
set TEST_DB_USER=root
set TEST_DB_PASSWORD=password
set TEST_DB_NAME=free2free_test
set TEST_JWT_SECRET=test-jwt-secret-key-32-chars-long-enough!!
set TEST_FACEBOOK_KEY=test-facebook-app-id
set TEST_FACEBOOK_SECRET=test-facebook-app-secret
set TEST_BASE_URL=http://localhost:8080
set TEST_SERVER_PORT=8080

echo Environment variables set successfully.

REM Check if Docker is available
docker --version >nul 2>&1
if %errorlevel% == 0 (
    echo Docker is available, checking if MariaDB container is running...
    
    REM Check if the MariaDB container is already running
    docker ps --format "table {{.Names}}\t{{.Status}}" | findstr -i mariadb >nul
    if %errorlevel% == 0 (
        echo MariaDB container is already running.
    ) else (
        echo MariaDB container not running, checking if it exists...
        
        REM Check if container exists but is stopped
        docker ps -a --format "table {{.Names}}\t{{.Status}}" | findstr -i mariadb >nul
        if %errorlevel% == 0 (
            echo MariaDB container exists but is stopped. Starting it...
            docker start mariadb_container
        ) else (
            echo MariaDB container does not exist. Please start your MariaDB service or container separately.
        )
    )
) else (
    echo Docker not found. Please ensure MariaDB service is running separately.
)

REM Wait a moment for services to be ready
timeout /t 2 /nobreak >nul

echo Test environment setup complete.
echo You can now run tests with: go test ./tests/...
echo For example:
echo   go test ./tests/unit/...
echo   go test ./tests/integration/...
echo   go test ./tests/e2e/...

pause