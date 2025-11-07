package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nebulabox/nebulabox/internal/engine"
)

// NewBuildCommand creates the build command
func NewBuildCommand() *cobra.Command {
	var tag string
	var file string

	cmd := &cobra.Command{
		Use:   "build [flags] PATH",
		Short: "Build an image from a BuildSpec",
		Long: `Build an image from a NebulaBox BuildSpec (JSON) or Dockerfile.
BuildSpec is NebulaBox's structured JSON format for defining container images.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return buildImage(path, tag, file)
		},
	}

	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	cmd.Flags().StringVarP(&file, "file", "f", "buildspec.json", "Name of the BuildSpec file (default: buildspec.json)")

	return cmd
}

func buildImage(path, tag, file string) error {
	buildspecPath := filepath.Join(path, file)

	// Check if file exists
	if _, err := os.Stat(buildspecPath); os.IsNotExist(err) {
		return fmt.Errorf("BuildSpec file not found: %s", buildspecPath)
	}

	fmt.Printf("🔨 Building image from BuildSpec: %s\n", buildspecPath)

	// Read BuildSpec
	data, err := os.ReadFile(buildspecPath)
	if err != nil {
		return fmt.Errorf("failed to read BuildSpec: %w", err)
	}

	// Parse BuildSpec
	var spec engine.BuildSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("failed to parse BuildSpec: %w", err)
	}

	// Override tag if provided
	if tag != "" {
		parts := strings.Split(tag, ":")
		spec.Name = parts[0]
		if len(parts) > 1 {
			spec.Tag = parts[1]
		} else {
			spec.Tag = "latest"
		}
	}

	// Ensure tag is set
	if spec.Tag == "" {
		spec.Tag = "latest"
	}

	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// Build image
	fmt.Printf("📦 Building image: %s:%s\n", spec.Name, spec.Tag)
	image, err := client.BuildImage(&spec)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	fmt.Printf("✅ Image built successfully!\n")
	fmt.Printf("   ID: %s\n", image.ID)
	fmt.Printf("   Name: %s:%s\n", image.Name, image.Tag)
	fmt.Printf("   Size: %s\n", formatImageSize(image.Size))
	fmt.Printf("   Digest: %s\n", image.Digest)

	return nil
}

func formatImageSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

