package cli

import (
	"context"
	"fmt"

	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/spf13/cobra"
)

// NewRunCommand creates the run command
func NewRunCommand() *cobra.Command {
	var (
		detach   bool
		port     string
		name     string
		env      []string
		volume   []string
	)

	cmd := &cobra.Command{
		Use:   "run [flags] IMAGE [COMMAND] [ARG...]",
		Short: "Run a container from an image",
		Long: `Run a container from the specified image. This is the main command 
for starting containers with NebulaBox.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			image := args[0]
			var command []string
			if len(args) > 1 {
				command = args[1:]
			}

			return runContainer(image, command, runOptions{
				detach: detach,
				port:   port,
				name:   name,
				env:    env,
				volume: volume,
			})
		},
	}

	cmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run container in background")
	cmd.Flags().StringVarP(&port, "port", "p", "", "Publish container port to host")
	cmd.Flags().StringVar(&name, "name", "", "Assign a name to the container")
	cmd.Flags().StringArrayVarP(&env, "env", "e", []string{}, "Set environment variables")
	cmd.Flags().StringArrayVar(&volume, "volume", []string{}, "Bind mount a volume")

	return cmd
}

type runOptions struct {
	detach bool
	port   string
	name   string
	env    []string
	volume []string
}

func runContainer(image string, command []string, opts runOptions) error {
	fmt.Printf("🚀 NebulaBox: Starting container from image '%s'\n", image)
	
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize containerd client: %w", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	// Pull image first
	fmt.Printf("⬇️  Pulling image: %s\n", image)
	if err := client.PullImage(ctx, image); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	
	// Create container options
	containerOpts := &containerd.ContainerOptions{
		Name:        opts.name,
		Ports:       make(map[string]string),
		Environment: make(map[string]string),
		Volumes:     make(map[string]string),
		Detach:      opts.detach,
	}
	
	// Parse port mapping
	if opts.port != "" {
		containerOpts.Ports["80"] = opts.port // Default to port 80
	}
	
	// Parse environment variables
	for _, env := range opts.env {
		containerOpts.Environment[env] = env // Simplified for now
	}
	
	// Parse volume mounts
	for _, vol := range opts.volume {
		containerOpts.Volumes[vol] = vol // Simplified for now
	}
	
	// Create container
	fmt.Printf("📦 Creating container...\n")
	container, err := client.CreateContainer(ctx, image, opts.name, containerOpts)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	
	// Start container
	fmt.Printf("🔄 Starting container: %s\n", container.ID)
	if err := client.StartContainer(ctx, container.ID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	
	fmt.Printf("✅ Container started successfully!\n")
	fmt.Printf("   ID: %s\n", container.ID)
	fmt.Printf("   Name: %s\n", container.Name)
	fmt.Printf("   Image: %s\n", container.Image)
	fmt.Printf("   Status: %s\n", container.Status)
	
	if opts.detach {
		fmt.Println("🔄 Running in background...")
	} else {
		fmt.Println("🔄 Running in foreground...")
	}
	
	return nil
}
