package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nebulabox/nebulabox/internal/engine"
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

	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// Pull image first (if not exists)
	fmt.Printf("⬇️  Pulling image: %s\n", image)
	if err := client.PullImage(image); err != nil {
		// Image might already exist, continue
		fmt.Printf("⚠️  Image pull note: %v (continuing...)\n", err)
	}

	// Generate container ID if name not provided
	containerID := opts.name
	if containerID == "" {
		containerID = generateContainerID(image)
	}

	// Parse port mappings
	ports := make(map[string]string)
	if opts.port != "" {
		// Parse "host:container" or "container"
		parts := strings.Split(opts.port, ":")
		if len(parts) == 2 {
			ports[parts[1]] = parts[0]
		} else {
			ports[opts.port] = opts.port
		}
	}

	// Parse environment variables
	env := make(map[string]string)
	for _, e := range opts.env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		} else {
			env[e] = ""
		}
	}

	// Parse volume mounts
	volumes := make(map[string]string)
	for _, vol := range opts.volume {
		parts := strings.SplitN(vol, ":", 2)
		if len(parts) == 2 {
			volumes[parts[1]] = parts[0]
		} else {
			volumes[vol] = vol
		}
	}

	// Create container spec
	spec := &engine.ContainerSpec{
		ID:      containerID,
		Name:    opts.name,
		Image:   image,
		Command: command,
		Ports:   ports,
		Env:     env,
		Volumes: volumes,
	}

	// Create container
	fmt.Printf("📦 Creating container...\n")
	container, err := client.CreateContainer(spec)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	fmt.Printf("🔄 Starting container: %s\n", container.ID)
	if err := client.StartContainer(container.ID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("✅ Container started successfully!\n")
	fmt.Printf("   ID: %s\n", container.ID)
	fmt.Printf("   Name: %s\n", container.Name)
	fmt.Printf("   Image: %s\n", container.Image)
	fmt.Printf("   Status: %s\n", container.State)

	if opts.detach {
		fmt.Println("🔄 Running in background...")
	} else {
		fmt.Println("🔄 Running in foreground...")
	}

	return nil
}

func generateContainerID(image string) string {
	// Simple ID generation from image name
	parts := strings.Split(image, ":")
	name := parts[0]
	if len(name) > 8 {
		name = name[:8]
	}
	return fmt.Sprintf("%s-%d", name, getTimestamp())
}

func getTimestamp() int64 {
	return time.Now().UnixNano()
}
