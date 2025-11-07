package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewImagesCommand creates the images command
func NewImagesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List images",
		Long:  `List all container images available in NebulaBox.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listImages()
		},
	}

	return cmd
}

// NewRmiCommand creates the rmi (remove image) command
func NewRmiCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "rmi IMAGE",
		Short: "Remove an image",
		Long:  `Remove one or more images from NebulaBox.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeImages(args, force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force removal of image")

	return cmd
}

func listImages() error {
	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// List images
	images, err := client.ListImages()
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(images) == 0 {
		fmt.Println("No images found")
		return nil
	}

	// Display images in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "IMAGE ID\tREPOSITORY\tTAG\tSIZE\tCREATED")
	fmt.Fprintln(w, "--------\t----------\t---\t----\t-------")

	for _, img := range images {
		imageID := img.ID
		if len(imageID) > 12 {
			imageID = imageID[:12]
		}

		repo := img.Name
		tag := img.Tag
		if tag == "" {
			tag = "latest"
		}

		size := formatSize(img.Size)
		created := img.CreatedAt.Format("2006-01-02 15:04:05")

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", imageID, repo, tag, size, created)
	}

	w.Flush()
	return nil
}

func removeImages(images []string, force bool) error {
	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	for _, imageRef := range images {
		fmt.Printf("🗑️  Removing image: %s\n", imageRef)

		if err := client.DeleteImage(imageRef); err != nil {
			if force {
				fmt.Printf("⚠️  Warning: %v (continuing due to --force)\n", err)
				continue
			}
			return fmt.Errorf("failed to remove image %s: %w", imageRef, err)
		}

		fmt.Printf("✅ Image removed: %s\n", imageRef)
	}

	return nil
}

func formatSize(bytes int64) string {
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

