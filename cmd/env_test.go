package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateExportCommands(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		currentVars []string
		newVars     map[string]string
		shellType   string
		expected    []string
	}{
		{
			name:        "New variables only",
			currentVars: []string{},
			newVars: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
			shellType: "bash",
			expected: []string{
				"export BAZ=qux",
				"export FOO=bar",
				"export ENVTOOL_MANAGED_ENV_VARS=BAZ,FOO",
			},
		},
		{
			name:        "Unset variables",
			currentVars: []string{"FOO", "BAR", "BAZ"},
			newVars: map[string]string{
				"BAZ": "updated",
				"QUX": "new",
			},
			shellType: "bash",
			expected: []string{
				"unset FOO",
				"unset BAR",
				"export BAZ=updated",
				"export QUX=new",
				"export ENVTOOL_MANAGED_ENV_VARS=BAZ,QUX",
			},
		},
		{
			name:        "No changes",
			currentVars: []string{"FOO", "BAR"},
			newVars: map[string]string{
				"FOO": "unchanged",
				"BAR": "unchanged",
			},
			shellType: "bash",
			expected: []string{
				"export BAR=unchanged",
				"export FOO=unchanged",
				"export ENVTOOL_MANAGED_ENV_VARS=BAR,FOO",
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := generateExportCommands(tc.currentVars, tc.newVars, tc.shellType)
			lines := strings.Split(strings.TrimSpace(output), "\n")
			
			assert.Equal(t, len(tc.expected), len(lines), "Number of output lines doesn't match expected")
			
			for i, expected := range tc.expected {
				assert.Equal(t, expected, lines[i], "Output line doesn't match expected")
			}
		})
	}
}