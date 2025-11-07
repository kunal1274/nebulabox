package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/nebulabox/nebulabox/internal/engine"
)

// NewListCommand creates the list command
func NewListCommand() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List containers",
		Long:  `List all running containers managed by NebulaBox.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listContainers(all)
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (including stopped)")

	return cmd
}

func listContainers(all bool) error {
	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// List containers
	containers, err := client.ListContainers()
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		fmt.Println("No containers found")
		return nil
	}

	// Display containers in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "CONTAINER ID\tIMAGE\tSTATUS\tNAME\tCREATED")
	fmt.Fprintln(w, "------------\t-----\t------\t----\t-------")

	for _, container := range containers {
		state := string(container.State)
		if !all && container.State == engine.StateStopped {
			continue // Skip stopped containers unless --all is specified
		}

		containerID := container.ID
		if len(containerID) > 12 {
			containerID = containerID[:12]
		}

		created := container.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			containerID,
			container.Image,
			state,
			container.Name,
			created)
	}

	w.Flush()
	return nil
}

// NewStopCommand creates the stop command
func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop CONTAINER",
		Short: "Stop a running container",
		Long:  `Stop one or more running containers.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stopContainers(args)
		},
	}

	return cmd
}

func stopContainers(containers []string) error {
	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	for _, containerID := range containers {
		fmt.Printf("🛑 Stopping container: %s\n", containerID)

		if err := client.StopContainer(containerID); err != nil {
			fmt.Printf("❌ Failed to stop container %s: %v\n", containerID, err)
			continue
		}

		fmt.Printf("✅ Container %s stopped\n", containerID)
	}

	return nil
}

// NewLogsCommand creates the logs command
func NewLogsCommand() *cobra.Command {
	var follow bool
	var tail string

	cmd := &cobra.Command{
		Use:   "logs CONTAINER",
		Short: "Fetch the logs of a container",
		Long:  `Fetch the logs of a container.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showLogs(args[0], follow, tail)
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().StringVar(&tail, "tail", "all", "Number of lines to show from the end of the logs")

	return cmd
}

func showLogs(container string, follow bool, tail string) error {
	fmt.Printf("📜 Logs for container: %s\n", container)
	fmt.Println("----------------------------------------------------")

	// Create engine client
	client, err := NewEngineClient()
	if err != nil {
		return fmt.Errorf("failed to create engine client: %w", err)
	}

	// Get container
	containerObj, err := client.GetContainer(container)
	if err != nil {
		return fmt.Errorf("container %s not found: %w", container, err)
	}

	// Get logs from engine
	logs, err := client.GetContainerLogs(container)
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}

	if len(logs) == 0 {
		fmt.Printf("No logs available for container %s\n", containerObj.Name)
		fmt.Println("(Container may not be running or logs not yet collected)")
		return nil
	}

	// Display logs
	fmt.Println(string(logs))

	if follow {
		fmt.Println("\n💡 Follow mode not yet fully implemented - logs shown above")
	}

	return nil
}
