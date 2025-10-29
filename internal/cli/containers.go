package cli

import (
	"context"
	"fmt"

	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/spf13/cobra"
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
	fmt.Println("📋 NebulaBox Containers:")
	fmt.Println("CONTAINER ID    IMAGE           STATUS          NAMES")
	fmt.Println("----------------------------------------------------")
	
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize containerd client: %w", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	// Get containers
	containers, err := client.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}
	
	if len(containers) == 0 {
		fmt.Println("No containers found")
		return nil
	}
	
	// Display containers
	for _, container := range containers {
		if !all && container.Status == "stopped" {
			continue // Skip stopped containers unless --all is specified
		}
		
		fmt.Printf("%-15s %-15s %-15s %s\n", 
			container.ID, 
			container.Image, 
			container.Status, 
			container.Name)
	}
	
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
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize containerd client: %w", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	for _, container := range containers {
		fmt.Printf("🛑 Stopping container: %s\n", container)
		
		if err := client.StopContainer(ctx, container); err != nil {
			fmt.Printf("❌ Failed to stop container %s: %v\n", container, err)
			continue
		}
		
		fmt.Printf("✅ Container %s stopped\n", container)
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
	
	// Initialize containerd client
	client, err := containerd.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize containerd client: %w", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	// Get container logs
	logs, err := client.GetContainerLogs(ctx, container)
	if err != nil {
		return fmt.Errorf("failed to get logs for container %s: %w", container, err)
	}
	
	if len(logs) == 0 {
		fmt.Println("No logs available")
		return nil
	}
	
	// Display logs
	for _, log := range logs {
		fmt.Println(log)
	}
	
	return nil
}
