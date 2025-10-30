package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContainerBuild tests that the application can be built in a minimal container environment
func TestContainerBuild(t *testing.T) {
	// This test verifies that the app can build in a container environment without CGO
	// For this test, we'll simulate the container build environment by ensuring CGO is disabled
	
	// Save original CGO_ENABLED value to restore later
	originalCgoEnabled := os.Getenv("CGO_ENABLED")
	
	// Ensure CGO is disabled
	err := os.Setenv("CGO_ENABLED", "0")
	assert.NoError(t, err)
	
	// Restore original value after test
	defer os.Setenv("CGO_ENABLED", originalCgoEnabled)
	
	// Test building the application
	cmd := exec.Command("go", "build", "-v", ".")
	cmd.Dir = ".." // Go up to project root
	output, err := cmd.CombinedOutput()
	
	t.Logf("Container build output: %s", string(output))
	
	// The build should succeed without CGO
	assert.NoError(t, err, "Application should build successfully in container-like environment")
}

// TestContainerRun tests that the application can run in a container environment
func TestContainerRun(t *testing.T) {
	// For this test we'll create a basic execution test that ensures the app can run
	// without requiring native dependencies
	
	// Save original CGO_ENABLED value
	originalCgoEnabled := os.Getenv("CGO_ENABLED")
	
	// Ensure CGO is disabled
	err := os.Setenv("CGO_ENABLED", "0")
	assert.NoError(t, err)
	
	// Restore original value after test
	defer os.Setenv("CGO_ENABLED", originalCgoEnabled)
	
	// Test that our Go code can run by executing a simple test
	cmd := exec.Command("go", "test", "-run", "TestSimpleInMemoryDB", "./tests/unit/")
	cmd.Dir = ".." // Go up to project root
	output, err := cmd.CombinedOutput()
	
	t.Logf("Container run output: %s", string(output))
	
	// The test should succeed
	assert.NoError(t, err, "Application should run successfully in container-like environment")
}