package main

import (
	"fmt"
	"os"

	"github.com/nebulabox/nebulabox/internal/cli"
	"github.com/sirupsen/logrus"
)

// Version and build information (set via ldflags during build)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func init() {
	// Set environment variables for version info
	if Version != "dev" {
		os.Setenv("NEBULABOX_VERSION", Version)
	}
	if BuildTime != "unknown" {
		os.Setenv("NEBULABOX_BUILD_TIME", BuildTime)
	}
	if GitCommit != "unknown" {
		os.Setenv("NEBULABOX_GIT_COMMIT", GitCommit)
	}
}

func main() {
	// Set up logging
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create and execute the root command
	rootCmd := cli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
