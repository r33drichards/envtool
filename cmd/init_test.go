package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmd_CustomPaths(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "envtool-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test file paths
	bashrc := filepath.Join(tempDir, "bashrc")
	zshrc := filepath.Join(tempDir, "zshrc")
	
	// Save original flags
	origBashrc := bashrcPath
	origZshrc := zshrcPath
	origUserOnly := userOnly
	origBashOnly := bashOnly
	origZshOnly := zshOnly
	
	// Restore originals after test
	defer func() {
		bashrcPath = origBashrc
		zshrcPath = origZshrc
		userOnly = origUserOnly
		bashOnly = origBashOnly
		zshOnly = origZshOnly
	}()
	
	// Configure for test: both shells with custom paths
	bashrcPath = bashrc
	zshrcPath = zshrc
	userOnly = false
	bashOnly = false
	zshOnly = false
	
	// Execute the command
	cmd := initCmd
	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
	
	// Verify bash file contents
	bashContent, err := ioutil.ReadFile(bashrc)
	assert.NoError(t, err)
	assert.Contains(t, string(bashContent), "_envtool_hook")
	assert.Contains(t, string(bashContent), "envtool env bash")
	
	// Verify zsh file contents
	zshContent, err := ioutil.ReadFile(zshrc)
	assert.NoError(t, err)
	assert.Contains(t, string(zshContent), "_envtool_hook")
	assert.Contains(t, string(zshContent), "envtool env zsh")
}

func TestInitCmd_BashOnly_PositionalPaths(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "envtool-test-bash")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	rcPath := filepath.Join(tempDir, "bashrc")
	envPath := filepath.Join(tempDir, "custom.env")
	if err := ioutil.WriteFile(envPath, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatalf("Failed to write env file: %v", err)
	}
	
	// Save originals
	origBashrc := bashrcPath
	origZshrc := zshrcPath
	origUserOnly := userOnly
	origBashOnly := bashOnly
	origZshOnly := zshOnly
	defer func() {
		bashrcPath = origBashrc
		zshrcPath = origZshrc
		userOnly = origUserOnly
		bashOnly = origBashOnly
		zshOnly = origZshOnly
	}()
	
	// Configure: bash only
	userOnly = false
	bashOnly = true
	zshOnly = false
	// Paths will be set via positional args
	
	cmd := initCmd
	err = cmd.RunE(cmd, []string{rcPath, envPath})
	assert.NoError(t, err)
	
	content, err := ioutil.ReadFile(rcPath)
	assert.NoError(t, err)
	text := string(content)
	assert.Contains(t, text, "envtool env bash")
	assert.Contains(t, text, "--env-file "+envPath)
	assert.NotContains(t, text, "envtool env zsh")
}

func TestInitCmd_ZshOnly_PositionalPaths(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "envtool-test-zsh")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	rcPath := filepath.Join(tempDir, "zshrc")
	envPath := filepath.Join(tempDir, "custom.env")
	if err := ioutil.WriteFile(envPath, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatalf("Failed to write env file: %v", err)
	}
	
	// Save originals
	origBashrc := bashrcPath
	origZshrc := zshrcPath
	origUserOnly := userOnly
	origBashOnly := bashOnly
	origZshOnly := zshOnly
	defer func() {
		bashrcPath = origBashrc
		zshrcPath = origZshrc
		userOnly = origUserOnly
		bashOnly = origBashOnly
		zshOnly = origZshOnly
	}()
	
	// Configure: zsh only
	userOnly = false
	bashOnly = false
	zshOnly = true
	
	cmd := initCmd
	err = cmd.RunE(cmd, []string{rcPath, envPath})
	assert.NoError(t, err)
	
	content, err := ioutil.ReadFile(rcPath)
	assert.NoError(t, err)
	text := string(content)
	assert.Contains(t, text, "envtool env zsh")
	assert.Contains(t, text, "--env-file "+envPath)
	assert.NotContains(t, text, "envtool env bash")
}