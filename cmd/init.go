package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/username/envtool/pkg/shell"
)

// Default paths for shell configuration files
const (
	defaultBashrcPath = "/etc/bash.bashrc"
	defaultZshrcPath  = "/etc/zsh/zshrc"
)

var (
	bashrcPath string
	zshrcPath  string
	userOnly   bool
	bashOnly   bool
	zshOnly    bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize shell configuration",
	Long: `Initialize shell configuration by adding hooks to shell rc files.
This allows envtool to automatically update environment variables when
the shell prompt is displayed.

By default, it modifies system-wide configuration files. Use --user flag
to modify user-specific configuration files instead.

You can also specify custom paths for bash and zsh configuration files
using the --bashrc and --zshrc flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileManager := &shell.DefaultFileManager{}
		
		// Determine which shells to update
		updateBash := true
		updateZsh := true
		if bashOnly && !zshOnly {
			updateZsh = false
		}
		if zshOnly && !bashOnly {
			updateBash = false
		}
		
		// If user passed positional args but selected both shells, return error
		if len(args) > 0 && updateBash && updateZsh {
			return fmt.Errorf("positional rc/env paths are supported only when selecting exactly one shell with --bash or --zsh")
		}
		
		// Optional env-file to embed in hook
		envPathFromArgs := ""
		
		// If user-only flag is set, use user-specific config files
		if userOnly {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get user home directory: %w", err)
			}
			
			// Only override if not explicitly set via flags
			if bashrcPath == defaultBashrcPath {
				bashrcPath = filepath.Join(homeDir, ".bashrc")
			}
			if zshrcPath == defaultZshrcPath {
				zshrcPath = filepath.Join(homeDir, ".zshrc")
			}
		}
		
		// Handle positional args for exactly one shell
		if updateBash && !updateZsh {
			if len(args) >= 1 {
				bashrcPath = args[0]
			}
			if len(args) >= 2 {
				envPathFromArgs = args[1]
			}
		} else if updateZsh && !updateBash {
			if len(args) >= 1 {
				zshrcPath = args[0]
			}
			if len(args) >= 2 {
				envPathFromArgs = args[1]
			}
		}
		
		// Determine env-file to embed in hook
		envPathForHook := strings.TrimSpace(envPathFromArgs)
		if envPathForHook == "" {
			// fall back to flag/config if set
			if value := strings.TrimSpace(viper.GetString("env-file")); value != "" && value != ".env" {
				envPathForHook = value
			}
		}
		envFlag := ""
		if envPathForHook != "" {
			envFlag = fmt.Sprintf(" --env-file %s", envPathForHook)
		}
		
		// Build hook contents dynamically
		bashHook := fmt.Sprintf(`
_envtool_hook() {
  local previous_exit_status=$?;
  trap -- '' SIGINT;
  eval "$(envtool env bash%s)";
  trap - SIGINT;
  return $previous_exit_status;
};
if ! [[ "${PROMPT_COMMAND:-}" =~ _envtool_hook ]]; then
  PROMPT_COMMAND="_envtool_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi
`, envFlag)
		zshHook := fmt.Sprintf(`
_envtool_hook() {
  trap -- '' SIGINT;
  eval "$(envtool env zsh%s)";
  trap - SIGINT;
}
typeset -ag precmd_functions;
if [[ -z "${precmd_functions[(r)_envtool_hook]+1}" ]]; then
  precmd_functions=( _envtool_hook ${precmd_functions[@]} )
fi
typeset -ag chpwd_functions;
if [[ -z "${chpwd_functions[(r)_envtool_hook]+1}" ]]; then
  chpwd_functions=( _envtool_hook ${chpwd_functions[@]} )
fi
`, envFlag)
		
		// Setup for bash
		if updateBash {
			exists, err := fileManager.FileExists(bashrcPath)
			if err != nil {
				return fmt.Errorf("failed to check bash configuration: %w", err)
			}
			
			// Create parent directory if it doesn't exist
			if !exists {
				dir := filepath.Dir(bashrcPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory for bash configuration: %w", err)
				}
				
				if err := fileManager.WriteFile(bashrcPath, ""); err != nil {
					return fmt.Errorf("failed to create bash configuration: %w", err)
				}
			}
			
			if err := fileManager.AppendToFile(bashrcPath, bashHook); err != nil {
				return fmt.Errorf("failed to update bash configuration: %w", err)
			}
		}
		
		// Setup for zsh
		if updateZsh {
			exists, err := fileManager.FileExists(zshrcPath)
			if err != nil {
				return fmt.Errorf("failed to check zsh configuration: %w", err)
			}
			
			if !exists {
				dir := filepath.Dir(zshrcPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory for zsh configuration: %w", err)
				}
				
				if err := fileManager.WriteFile(zshrcPath, ""); err != nil {
					return fmt.Errorf("failed to create zsh configuration: %w", err)
				}
			}
			
			if err := fileManager.AppendToFile(zshrcPath, zshHook); err != nil {
				return fmt.Errorf("failed to update zsh configuration: %w", err)
			}
		}
		
		fmt.Printf("Shell configurations updated successfully:\n")
		if updateBash {
			fmt.Printf("- Bash: %s\n", bashrcPath)
		}
		if updateZsh {
			fmt.Printf("- Zsh: %s\n", zshrcPath)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	
	// Add flags for customizing configuration paths
	initCmd.Flags().StringVar(&bashrcPath, "bashrc", defaultBashrcPath, "Path to bash configuration file")
	initCmd.Flags().StringVar(&zshrcPath, "zshrc", defaultZshrcPath, "Path to zsh configuration file")
	initCmd.Flags().BoolVar(&userOnly, "user", false, "Modify user-specific configuration files instead of system-wide")
	initCmd.Flags().BoolVar(&bashOnly, "bash", false, "Only update bash configuration (default: both shells)")
	initCmd.Flags().BoolVar(&zshOnly, "zsh", false, "Only update zsh configuration (default: both shells)")
	
	// Bind to viper for config file support
	viper.BindPFlag("init.bashrc", initCmd.Flags().Lookup("bashrc"))
	viper.BindPFlag("init.zshrc", initCmd.Flags().Lookup("zshrc"))
	viper.BindPFlag("init.user", initCmd.Flags().Lookup("user"))
	viper.BindPFlag("init.bash", initCmd.Flags().Lookup("bash"))
	viper.BindPFlag("init.zsh", initCmd.Flags().Lookup("zsh"))
}