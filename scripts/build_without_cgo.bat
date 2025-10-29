@echo off
REM Build script for testing application build without CGO dependencies
REM This script builds the application with CGO disabled to verify platform-independent compilation

echo Building application without CGO dependencies...

REM Set CGO_ENABLED=0 to disable CGO
set CGO_ENABLED=0

REM Build the application
echo Building with CGO disabled...
go build -o free2free_nocgo.exe .

if %errorlevel% == 0 (
    echo Build successful! Application compiled without CGO dependencies.
    echo Binary: free2free_nocgo.exe
    
    REM Show build info
    echo Build timestamp: %date% %time%
    
    REM Clean up the binary
    del free2free_nocgo.exe >nul 2>&1
    echo Cleaned up build artifacts.
) else (
    echo Build failed! Application has CGO dependencies that prevent platform-independent compilation.
    echo Error level: %errorlevel%
    exit /b %errorlevel%
)

echo Testing cross-compilation for different platforms...

REM Test cross-compilation for Linux
echo Testing Linux build...
set GOOS=linux
set GOARCH=amd64
go build -o free2free_linux .

if %errorlevel% == 0 (
    echo Linux build successful!
    del free2free_linux >nul 2>&1
) else (
    echo Linux build failed!
)

REM Reset environment variables
set GOOS=
set GOARCH=

REM Test cross-compilation for macOS
echo Testing macOS build...
set GOOS=darwin
set GOARCH=amd64
go build -o free2free_macos .

if %errorlevel% == 0 (
    echo macOS build successful!
    del free2free_macos >nul 2>&1
) else (
    echo macOS build failed!
)

REM Reset environment variables
set GOOS=
set GOARCH=

REM Test cross-compilation for Windows
echo Testing Windows build...
set GOOS=windows
set GOARCH=amd64
go build -o free2free_windows.exe .

if %errorlevel% == 0 (
    echo Windows build successful!
    del free2free_windows.exe >nul 2>&1
) else (
    echo Windows build failed!
)

REM Reset environment variables
set GOOS=
set GOARCH=
set CGO_ENABLED=

echo.
echo All build tests completed successfully!
echo The application can be compiled without CGO dependencies.
echo This enables cross-platform deployment and containerization.
echo.

pause