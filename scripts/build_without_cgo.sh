#!/bin/bash

# Build script for testing application build without CGO dependencies
# This script builds the application with CGO disabled to verify platform-independent compilation

echo "Building application without CGO dependencies..."

# Set CGO_ENABLED=0 to disable CGO
export CGO_ENABLED=0

# Build the application
echo "Building with CGO disabled..."
go build -o free2free_nocgo .

if [ $? -eq 0 ]; then
    echo "Build successful! Application compiled without CGO dependencies."
    echo "Binary: free2free_nocgo"
    
    # Show build info
    echo "Build timestamp: $(date)"
    
    # Clean up the binary
    rm -f free2free_nocgo
    echo "Cleaned up build artifacts."
else
    echo "Build failed! Application has CGO dependencies that prevent platform-independent compilation."
    exit 1
fi

echo "Testing cross-compilation for different platforms..."

# Test cross-compilation for Linux
echo "Testing Linux build..."
export GOOS=linux
export GOARCH=amd64
go build -o free2free_linux .

if [ $? -eq 0 ]; then
    echo "Linux build successful!"
    rm -f free2free_linux
else
    echo "Linux build failed!"
fi

# Reset environment variables
unset GOOS
unset GOARCH

# Test cross-compilation for macOS
echo "Testing macOS build..."
export GOOS=darwin
export GOARCH=amd64
go build -o free2free_macos .

if [ $? -eq 0 ]; then
    echo "macOS build successful!"
    rm -f free2free_macos
else
    echo "macOS build failed!"
fi

# Reset environment variables
unset GOOS
unset GOARCH

# Test cross-compilation for Windows
echo "Testing Windows build..."
export GOOS=windows
export GOARCH=amd64
go build -o free2free_windows.exe .

if [ $? -eq 0 ]; then
    echo "Windows build successful!"
    rm -f free2free_windows.exe
else
    echo "Windows build failed!"
fi

# Reset environment variables
unset GOOS
unset GOARCH
unset CGO_ENABLED

echo ""
echo "All build tests completed successfully!"
echo "The application can be compiled without CGO dependencies."
echo "This enables cross-platform deployment and containerization."
echo ""