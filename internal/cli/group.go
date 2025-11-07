package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nebulabox/nebulabox/internal/engine"
)

// NewGroupCommand creates the group command
func NewGroupCommand() *cobra.Command {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage container groups",
		Long:  "Create and manage container groups for flexible architecture testing",
	}

	groupCmd.AddCommand(NewGroupCreateCommand())
	groupCmd.AddCommand(NewGroupStartCommand())
	groupCmd.AddCommand(NewGroupStopCommand())
	groupCmd.AddCommand(NewGroupListCommand())
	groupCmd.AddCommand(NewGroupStatusCommand())

	return groupCmd
}

// NewGroupCreateCommand creates the group create command
func NewGroupCreateCommand() *cobra.Command {
	var specFile string

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a container group",
		Long: `Create a container group from a group specification file.
Groups allow you to test different architectures (monolithic, microservices, etc.)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if specFile == "" {
				return fmt.Errorf("group specification file required (--file)")
			}

			// Read group spec
			data, err := os.ReadFile(specFile)
			if err != nil {
				return fmt.Errorf("failed to read group spec: %w", err)
			}

			var spec engine.GroupSpec
			if err := json.Unmarshal(data, &spec); err != nil {
				return fmt.Errorf("failed to parse group spec: %w", err)
			}

			// Use group name from args if provided
			if len(args) > 0 {
				spec.Name = args[0]
			}

			// Create engine client
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			// Create group
			group, err := client.CreateGroup(&spec)
			if err != nil {
				return fmt.Errorf("failed to create group: %w", err)
			}

			fmt.Printf("✅ Group created: %s (ID: %s)\n", group.Name, group.ID)
			fmt.Printf("   Strategy: %s\n", group.Strategy)
			fmt.Printf("   Containers: %d\n", len(group.Containers))

			return nil
		},
	}

	cmd.Flags().StringVarP(&specFile, "file", "f", "", "Group specification file (JSON)")

	return cmd
}

// NewGroupStartCommand creates the group start command
func NewGroupStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start [group-id]",
		Short: "Start a container group",
		Long:  "Start all containers in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			if err := client.StartGroup(args[0]); err != nil {
				return fmt.Errorf("failed to start group: %w", err)
			}

			fmt.Printf("✅ Group started: %s\n", args[0])
			return nil
		},
	}
}

// NewGroupStopCommand creates the group stop command
func NewGroupStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [group-id]",
		Short: "Stop a container group",
		Long:  "Stop all containers in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			if err := client.StopGroup(args[0]); err != nil {
				return fmt.Errorf("failed to stop group: %w", err)
			}

			fmt.Printf("✅ Group stopped: %s\n", args[0])
			return nil
		},
	}
}

// NewGroupListCommand creates the group list command
func NewGroupListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List container groups",
		Long:  "List all container groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			groups, err := client.ListGroups()
			if err != nil {
				return fmt.Errorf("failed to list groups: %w", err)
			}

			if len(groups) == 0 {
				fmt.Println("No groups found")
				return nil
			}

			fmt.Printf("%-20s %-15s %-10s %-8s\n", "GROUP ID", "NAME", "STRATEGY", "STATE")
			fmt.Println(strings.Repeat("-", 60))
			for _, group := range groups {
				id := group.ID
				if len(id) > 12 {
					id = id[:12]
				}
				fmt.Printf("%-20s %-15s %-10s %-8s\n",
					id, group.Name, string(group.Strategy), string(group.State))
			}

			return nil
		},
	}
}

// NewGroupStatusCommand creates the group status command
func NewGroupStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status [group-id]",
		Short: "Show group status",
		Long:  "Show detailed status of a container group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			group, err := client.GetGroup(args[0])
			if err != nil {
				return fmt.Errorf("failed to get group: %w", err)
			}

			fmt.Printf("Group: %s (ID: %s)\n", group.Name, group.ID)
			fmt.Printf("Strategy: %s\n", group.Strategy)
			fmt.Printf("State: %s\n", group.State)
			fmt.Printf("Containers: %d\n", len(group.Containers))
			for i, containerID := range group.Containers {
				fmt.Printf("  [%d] %s\n", i+1, containerID)
			}

			return nil
		},
	}
}

