package cli

import (
	"context"
	"fmt"

	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/spf13/cobra"
)

// NewBuildCommand creates the build command
func NewBuildCommand() *cobra.Command {
	var tag string
	var file string

	cmd := &cobra.Command{
		Use:   "build [flags] PATH",
		Short: "Build an image from a Dockerfile",
		Long:  `Build an image from a Dockerfile or NebulaBox build specification.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildImage(args[0], tag, file)
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	cmd.Flags().StringVarP(&file, "file", "f", "Dockerfile", "Name of the Dockerfile")

	return cmd
}

func buildImage(path, tag, file string) error {
	fmt.Printf("🔨 Building image from: %s\n", path)
	fmt.Printf("📄 Using build file: %s\n", file)
	
	if tag != "" {
		fmt.Printf("🏷️  Tag: %s\n", tag)
	}
	
	// TODO: Implement actual build process
	fmt.Println("✅ Build completed successfully!")
	fmt.Println("💡 This is a placeholder - build system coming in Phase 2!")
	
	return nil
}

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
	
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize containerd client: %w", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	// Pull image
	if err := client.PullImage(ctx, image); err != nil {
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
	fmt.Println("🚀 NebulaBox")
	fmt.Println("Version:     0.1.0-alpha")
	fmt.Println("Build:      Phase 1 Development")
	fmt.Println("Go Version: 1.22+")
	fmt.Println("Platform:   macOS")
	fmt.Println("")
	fmt.Println("🎯 Current Phase: Core Workflow Layer")
	fmt.Println("📋 Features: CLI, Basic Commands, Placeholder Integration")
	fmt.Println("🔮 Next: containerd Runtime Integration")
	
	return nil
}
