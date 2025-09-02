package envfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "envfile-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name: "Basic env file",
			content: `
FOO=bar
BAZ=qux
`,
			expected: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
		},
		{
			name: "Comments and empty lines",
			content: `
# This is a comment
FOO=bar

# Another comment
BAZ=qux
`,
			expected: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
		},
		{
			name: "Quoted values",
			content: `
FOO="quoted value"
BAR='single quoted'
BAZ="nested 'quotes'"
`,
			expected: map[string]string{
				"FOO": "quoted value",
				"BAR": "single quoted",
				"BAZ": "nested 'quotes'",
			},
		},
		{
			name: "Spaces in values",
			content: `
FOO=bar baz
SPACE_BEFORE= value with space before
SPACE_AFTER=value with space after 
`,
			expected: map[string]string{
				"FOO":          "bar baz",
				"SPACE_BEFORE": "value with space before",
				"SPACE_AFTER":  "value with space after",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test .env file
			envPath := filepath.Join(tempDir, ".env-"+tc.name)
			err := ioutil.WriteFile(envPath, []byte(tc.content), 0644)
			assert.NoError(t, err)

			// Parse the file
			parser := &DefaultParser{}
			result, err := parser.Parse(envPath)
			assert.NoError(t, err)

			// Verify the result
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParse_FileNotExist(t *testing.T) {
	parser := &DefaultParser{}
	result, err := parser.Parse("/nonexistent/file.env")
	assert.Error(t, err)
	assert.Nil(t, result)
}