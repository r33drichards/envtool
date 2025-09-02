package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration test for the full workflow
// To run: INTEGRATION_TEST=true go test -v ./tests
func TestFullWorkflow(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test; set INTEGRATION_TEST=true to run")
	}
	
	// Create temporary test directory
	tempDir, err := ioutil.TempDir("", "envtool-integration")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Build the binary
	binaryPath := filepath.Join(tempDir, "envtool")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../main.go")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}
	
	// Create test shell config files
	bashrcPath := filepath.Join(tempDir, "bashrc")
	zshrcPath := filepath.Join(tempDir, "zshrc")
	
	// Create test .env file
	envFilePath := filepath.Join(tempDir, ".env")
	err = ioutil.WriteFile(envFilePath, []byte("TEST_VAR=test_value\nANOTHER_VAR=another_value"), 0644)
	assert.NoError(t, err)
	
	// Run init command with custom paths
	initCmd := exec.Command(binaryPath, "init", "--bashrc", bashrcPath, "--zshrc", zshrcPath)
	output, err = initCmd.CombinedOutput()
	assert.NoError(t, err, "Init command failed: %s", output)
	
	// Verify shell config files were created and contain hooks
	bashContent, err := ioutil.ReadFile(bashrcPath)
	assert.NoError(t, err)
	assert.Contains(t, string(bashContent), "_envtool_hook")
	
	zshContent, err := ioutil.ReadFile(zshrcPath)
	assert.NoError(t, err)
	assert.Contains(t, string(zshContent), "_envtool_hook")
	
	// Run env command with custom env file
	envCmd := exec.Command(binaryPath, "env", "bash", "--env-file", envFilePath)
	output, err = envCmd.CombinedOutput()
	assert.NoError(t, err, "Env command failed: %s", output)
	
	// Verify env command output contains expected exports
	outputStr := string(output)
	assert.Contains(t, outputStr, "export TEST_VAR='test_value'")
	assert.Contains(t, outputStr, "export ANOTHER_VAR='another_value'")
	assert.Contains(t, outputStr, "export ENVTOOL_MANAGED_ENV_VARS=")
	
	// Simulate a change in the .env file
	err = ioutil.WriteFile(envFilePath, []byte("TEST_VAR=updated_value\nNEW_VAR=new_value"), 0644)
	assert.NoError(t, err)
	
	// Set environment variable to simulate previous run
	os.Setenv("ENVTOOL_MANAGED_ENV_VARS", "TEST_VAR,ANOTHER_VAR")
	defer os.Unsetenv("ENVTOOL_MANAGED_ENV_VARS")
	
	// Run env command again
	envCmd = exec.Command(binaryPath, "env", "bash", "--env-file", envFilePath)
	output, err = envCmd.CombinedOutput()
	assert.NoError(t, err, "Second env command failed: %s", output)
	
	// Verify updated output
	outputStr = string(output)
	assert.Contains(t, outputStr, "export TEST_VAR='updated_value'")
	assert.Contains(t, outputStr, "export NEW_VAR='new_value'")
	assert.Contains(t, outputStr, "unset ANOTHER_VAR")
}