package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nebulabox/nebulabox/internal/engine"
)

// NewHierarchyCommand creates the hierarchy command
func NewHierarchyCommand() *cobra.Command {
	hierarchyCmd := &cobra.Command{
		Use:   "hierarchy",
		Short: "Manage hierarchical containers",
		Long:  "Manage containers that can contain other containers (infinitely nested)",
	}

	hierarchyCmd.AddCommand(NewHierarchyCreateCommand())
	hierarchyCmd.AddCommand(NewHierarchyListCommand())
	hierarchyCmd.AddCommand(NewHierarchyTreeCommand())
	hierarchyCmd.AddCommand(NewHierarchyAddGroupCommand())

	return hierarchyCmd
}

// NewHierarchyCreateCommand creates a nested container
func NewHierarchyCreateCommand() *cobra.Command {
	var parentID string
	var specFile string

	cmd := &cobra.Command{
		Use:   "create [container-id]",
		Short: "Create a nested container within another container",
		Long: `Create a container within another container. This enables infinite nesting.
Containers can contain other containers, and those can contain more containers.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if parentID == "" && len(args) == 0 {
				return fmt.Errorf("parent container ID required (--parent or argument)")
			}

			if specFile == "" {
				return fmt.Errorf("container specification file required (--file)")
			}

			// Read container spec
			data, err := os.ReadFile(specFile)
			if err != nil {
				return fmt.Errorf("failed to read container spec: %w", err)
			}

			var spec engine.ContainerSpec
			if err := json.Unmarshal(data, &spec); err != nil {
				return fmt.Errorf("failed to parse container spec: %w", err)
			}

			// Use parent from args if provided
			if len(args) > 0 {
				parentID = args[0]
			}

			// Create engine client
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			// Create nested container
			container, err := client.CreateNestedContainer(parentID, &spec)
			if err != nil {
				return fmt.Errorf("failed to create nested container: %w", err)
			}

			fmt.Printf("✅ Nested container created: %s (ID: %s)\n", container.Name, container.ID)
			fmt.Printf("   Parent: %s\n", parentID)
			fmt.Printf("   Depth: %d\n", getContainerDepth(client, container.ID))

			return nil
		},
	}

	cmd.Flags().StringVar(&parentID, "parent", "", "Parent container ID")
	cmd.Flags().StringVarP(&specFile, "file", "f", "", "Container specification file (JSON)")

	return cmd
}

// NewHierarchyListCommand lists containers in a hierarchy
func NewHierarchyListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list [root-container-id]",
		Short: "List containers in a hierarchy",
		Long:  "List all containers in a hierarchy tree starting from a root container",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			if len(args) > 0 {
				// List containers in specific hierarchy
				containers, err := client.ListContainersInHierarchy(args[0])
				if err != nil {
					return fmt.Errorf("failed to list hierarchy: %w", err)
				}

				fmt.Printf("📋 Containers in hierarchy (root: %s):\n", args[0])
				fmt.Printf("%-20s %-15s %-10s %-8s\n", "CONTAINER ID", "NAME", "IMAGE", "DEPTH")
				fmt.Println(strings.Repeat("-", 60))
				for _, container := range containers {
					depth := getContainerDepth(client, container.ID)
					id := container.ID
					if len(id) > 12 {
						id = id[:12]
					}
					fmt.Printf("%-20s %-15s %-10s %-8d\n", id, container.Name, container.Image, depth)
				}
			} else {
				// List all hierarchies
				trees, err := client.GetFullHierarchy()
				if err != nil {
					return fmt.Errorf("failed to get full hierarchy: %w", err)
				}

				if len(trees) == 0 {
					fmt.Println("No hierarchical containers found")
					return nil
				}

				fmt.Println("📋 Container Hierarchies:")
				for _, tree := range trees {
					printTree(tree, 0)
				}
			}

			return nil
		},
	}
}

// NewHierarchyTreeCommand shows hierarchy tree
func NewHierarchyTreeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tree [container-id]",
		Short: "Show hierarchy tree",
		Long:  "Display the complete hierarchy tree for a container",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			if len(args) > 0 {
				// Show tree for specific container
				tree, err := client.GetHierarchy(args[0])
				if err != nil {
					return fmt.Errorf("failed to get hierarchy: %w", err)
				}

				fmt.Printf("🌳 Hierarchy Tree for: %s\n", args[0])
				printTree(tree, 0)
			} else {
				// Show all trees
				trees, err := client.GetFullHierarchy()
				if err != nil {
					return fmt.Errorf("failed to get full hierarchy: %w", err)
				}

				if len(trees) == 0 {
					fmt.Println("No hierarchical containers found")
					return nil
				}

				fmt.Println("🌳 All Container Hierarchies:")
				for i, tree := range trees {
					fmt.Printf("\n[%d] Root: %s\n", i+1, tree.Container.ID)
					printTree(tree, 0)
				}
			}

			return nil
		},
	}
}

// NewHierarchyAddGroupCommand adds a container to a group
func NewHierarchyAddGroupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add-group CONTAINER GROUP",
		Short: "Add container to group",
		Long:  "Add a container (from any hierarchy level) to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerID := args[0]
			groupID := args[1]

			client, err := NewEngineClient()
			if err != nil {
				return fmt.Errorf("failed to create engine client: %w", err)
			}

			if err := client.AddContainerToGroup(containerID, groupID); err != nil {
				return fmt.Errorf("failed to add container to group: %w", err)
			}

			fmt.Printf("✅ Container %s added to group %s\n", containerID, groupID)
			return nil
		},
	}
}

// Helper functions

func printTree(tree *engine.HierarchyTree, indent int) {
	prefix := strings.Repeat("  ", indent)
	if indent > 0 {
		prefix += "└─ "
	}

	id := tree.Container.ID
	if len(id) > 12 {
		id = id[:12]
	}

	fmt.Printf("%s%s (%s) [depth: %d]\n", prefix, tree.Container.Name, id, tree.Depth)

	// Show groups
	if len(tree.Groups) > 0 {
		for _, group := range tree.Groups {
			fmt.Printf("%s  └─ 📦 Group: %s\n", strings.Repeat("  ", indent), group.Name)
		}
	}

	// Show nested groups
	if len(tree.NestedGroups) > 0 {
		for _, group := range tree.NestedGroups {
			fmt.Printf("%s  └─ 📦 Nested Group: %s\n", strings.Repeat("  ", indent), group.Name)
		}
	}

	// Print children
	for _, child := range tree.Children {
		printTree(child, indent+1)
	}
}

func getContainerDepth(client *EngineClient, containerID string) int {
	// This would query the hierarchy manager
	// For now, return 0 as placeholder
	return 0
}

