package main

import (
	"fmt"
	"os"

	"github.com/nebulabox/nebulabox/internal/cli"
	"github.com/sirupsen/logrus"
)

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
