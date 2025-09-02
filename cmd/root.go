package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	envFile string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "envtool",
	Short: "A tool for managing environment variables",
	Long: `A tool for managing environment variables in shell environments.
It can initialize shell configurations and dynamically load variables from .env files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envtool.yaml)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", ".env", "path to .env file")

	// Bind flags to viper
	viper.BindPFlag("env-file", rootCmd.PersistentFlags().Lookup("env-file"))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			viper.AddConfigPath(home)
			viper.SetConfigName(".envtool")
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
