package shell

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test-file")
	err = ioutil.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	// Test file exists
	fileManager := &DefaultFileManager{}
	exists, err := fileManager.FileExists(testFile)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test file does not exist
	exists, err = fileManager.FileExists(filepath.Join(tempDir, "nonexistent"))
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestReadWriteFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test file path
	testFile := filepath.Join(tempDir, "test-file")
	
	// Test content
	content := "test content"

	// Write file
	fileManager := &DefaultFileManager{}
	err = fileManager.WriteFile(testFile, content)
	assert.NoError(t, err)

	// Read file
	readContent, err := fileManager.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, content, readContent)
}

func TestAppendToFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test file path
	testFile := filepath.Join(tempDir, "test-file")
	
	// Initial content
	initialContent := "initial content\n"
	
	// Content to append
	appendContent := "appended content"

	// Create file with initial content
	fileManager := &DefaultFileManager{}
	err = fileManager.WriteFile(testFile, initialContent)
	assert.NoError(t, err)

	// Append to file
	err = fileManager.AppendToFile(testFile, appendContent)
	assert.NoError(t, err)

	// Read file
	readContent, err := fileManager.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, initialContent+appendContent, readContent)

	// Append again (should NOT duplicate because ContainsContent prevents duplicates)
	err = fileManager.AppendToFile(testFile, appendContent)
	assert.NoError(t, err)

	// Read file again
	readContent, err = fileManager.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, initialContent+appendContent, readContent)
}

func TestContainsContent(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test file path
	testFile := filepath.Join(tempDir, "test-file")
	
	// File content
	content := "line 1\nline 2\nline 3\n"

	// Create file
	fileManager := &DefaultFileManager{}
	err = fileManager.WriteFile(testFile, content)
	assert.NoError(t, err)

	// Test contains content
	contains, err := fileManager.ContainsContent(testFile, "line 2")
	assert.NoError(t, err)
	assert.True(t, contains)

	// Test does not contain content
	contains, err = fileManager.ContainsContent(testFile, "line 4")
	assert.NoError(t, err)
	assert.False(t, contains)

	// Test file does not exist
	contains, err = fileManager.ContainsContent(filepath.Join(tempDir, "nonexistent"), "content")
	assert.NoError(t, err)
	assert.False(t, contains)
}