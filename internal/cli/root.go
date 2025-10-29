package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// NewRootCommand creates the root command for NebulaBox CLI
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nebulabox",
		Short: "NebulaBox - A modern DevOps platform",
		Long: `NebulaBox is a modern DevOps platform that simplifies containerization, 
orchestration, and deployment workflows. It provides Docker-like functionality 
with enhanced simplicity and intelligence.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize configuration
			initConfig()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nebulabox.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(NewRunCommand())
	rootCmd.AddCommand(NewListCommand())
	rootCmd.AddCommand(NewStopCommand())
	rootCmd.AddCommand(NewLogsCommand())
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewPushCommand())
	rootCmd.AddCommand(NewPullCommand())
	rootCmd.AddCommand(NewVersionCommand())

	return rootCmd
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".nebulabox" (without extension)
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			return
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nebulabox")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
