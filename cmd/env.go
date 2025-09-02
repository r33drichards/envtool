package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/username/envtool/pkg/envfile"
)

const (
	// Key for tracking managed environment variables
	ManagedEnvVarsKey = "ENVTOOL_MANAGED_ENV_VARS"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env [shell]",
	Short: "Generate shell commands to set environment variables",
	Long: `Generate shell commands to set environment variables from a .env file.
The output should be evaluated by the shell to apply the changes.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get shell type (bash, zsh, etc.) if provided
		shellType := "bash" // Default
		if len(args) > 0 {
			shellType = args[0]
		}
		
		// Get path to .env file
		envFilePath := viper.GetString("env-file")
		
		// Parse .env file
		parser := &envfile.DefaultParser{}
		envVars, err := parser.Parse(envFilePath)
		if err != nil {
			// Just return empty if file doesn't exist or can't be parsed
			return nil
		}
		
		// Get currently managed env vars
		managedEnvVarsStr := os.Getenv(ManagedEnvVarsKey)
		managedEnvVars := []string{}
		if managedEnvVarsStr != "" {
			managedEnvVars = strings.Split(managedEnvVarsStr, ",")
		}
		
		// Generate export commands
		output := generateExportCommands(managedEnvVars, envVars, shellType)
		
		// Print to stdout (will be captured by eval in the shell)
		fmt.Print(output)
		return nil
	},
}

// generateExportCommands generates shell commands to export/unset env vars
func generateExportCommands(currentVars []string, newVars map[string]string, shellType string) string {
	commands := []string{}
	
	// Generate unset commands for variables that are no longer present
	for _, key := range currentVars {
		if _, exists := newVars[key]; !exists && key != "" {
			commands = append(commands, fmt.Sprintf("unset %s", key))
		}
	}
	
	// Generate export commands for new variables
	newVarKeys := make([]string, 0, len(newVars))
	for key := range newVars {
		newVarKeys = append(newVarKeys, key)
	}
	
	// Sort keys for consistent output
	sort.Strings(newVarKeys)
	
	for _, key := range newVarKeys {
		value := newVars[key]
		// Quote value if not already quoted
		if !strings.HasPrefix(value, "'") {
			value = shellescape.Quote(value)
		}
		commands = append(commands, fmt.Sprintf("export %s=%s", key, value))
	}
	
	// Add command to update the managed vars list
	if len(newVarKeys) > 0 {
		newManagedVars := strings.Join(newVarKeys, ",")
		commands = append(commands, fmt.Sprintf("export %s=%s", ManagedEnvVarsKey, newManagedVars))
	}
	
	return strings.Join(commands, "\n")
}

func init() {
	rootCmd.AddCommand(envCmd)
}