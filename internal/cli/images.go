package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Build command moved to build_cmd.go

// NewPushCommand creates the push command
func NewPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push IMAGE",
		Short: "Push an image to a registry",
		Long:  `Push an image to NebulaBox registry or external registry.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pushImage(args[0])
		},
	}

	return cmd
}

func pushImage(image string) error {
	fmt.Printf("⬆️  Pushing image: %s\n", image)
	
	// TODO: Implement actual push to registry
	fmt.Println("✅ Image pushed successfully!")
	fmt.Println("💡 This is a placeholder - registry system coming in Phase 2!")
	
	return nil
}

// NewPullCommand creates the pull command
func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull IMAGE",
		Short: "Pull an image from a registry",
		Long:  `Pull an image from NebulaBox registry or external registry.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pullImage(args[0])
		},
	}

	return cmd
}

func pullImage(image string) error {
	fmt.Printf("⬇️  Pulling image: %s\n", image)

	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// Pull image
	if err := client.PullImage(image); err != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, err)
	}

	fmt.Printf("✅ Image pulled successfully: %s\n", image)
	return nil
}

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show NebulaBox version information",
		Long:  `Show version information for NebulaBox CLI and runtime.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showVersion()
		},
	}

	return cmd
}

func showVersion() error {
	// Try to get version from main package
	version := getVersion()
	buildTime := getBuildTime()
	gitCommit := getGitCommit()
	
	fmt.Println("🚀 NebulaBox")
	fmt.Printf("Version:     %s\n", version)
	if buildTime != "unknown" {
		fmt.Printf("Build Time:  %s\n", buildTime)
	}
	if gitCommit != "unknown" {
		shortCommit := gitCommit
		if len(shortCommit) > 7 {
			shortCommit = shortCommit[:7]
		}
		fmt.Printf("Git Commit:  %s\n", shortCommit)
	}
	fmt.Println("Go Version: 1.22+")
	fmt.Println("")
	fmt.Println("🎯 Current Phase: Core Workflow Layer")
	fmt.Println("📋 Features: CLI, Basic Commands, Real Engine Integration")
	fmt.Println("🔮 Next: API Development")
	
	return nil
}

// These will be set by the main package or via build flags
func getVersion() string {
	if v := os.Getenv("NEBULABOX_VERSION"); v != "" {
		return v
	}
	return "0.1.0-alpha"
}

func getBuildTime() string {
	if t := os.Getenv("NEBULABOX_BUILD_TIME"); t != "" {
		return t
	}
	return "unknown"
}

func getGitCommit() string {
	if c := os.Getenv("NEBULABOX_GIT_COMMIT"); c != "" {
		return c
	}
	return "unknown"
}
