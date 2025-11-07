package engine

import (
	"context"
	"fmt"
	"time"
)

// GroupManager manages container groups
type GroupManager struct {
	runtime *NebulaRuntime
}

// NewGroupManager creates a new group manager
func NewGroupManager(runtime *NebulaRuntime) *GroupManager {
	return &GroupManager{
		runtime: runtime,
	}
}

// CreateGroup creates a new container group
func (gm *GroupManager) CreateGroup(ctx context.Context, spec *GroupSpec) (*Group, error) {
	groupID := generateGroupID(spec.Name)

	// Create containers for the group
	containerIDs := make([]string, 0, len(spec.Containers))
	for _, containerSpec := range spec.Containers {
		// Create container spec
		cspec := &ContainerSpec{
			ID:          generateContainerID(containerSpec.Name),
			Name:        containerSpec.Name,
			Image:       containerSpec.Image,
			Ports:       parsePorts(containerSpec.Ports),
			Env:         containerSpec.Env,
			Volumes:     containerSpec.Volumes,
			Network:     spec.Networking.Bridge,
			NetworkMode: "bridge",
			GroupID:     groupID,
		}

		// Create container
		container, err := gm.runtime.CreateContainer(ctx, cspec)
		if err != nil {
			return nil, fmt.Errorf("failed to create container %s: %w", containerSpec.Name, err)
		}

		containerIDs = append(containerIDs, container.ID)
	}

	// Store container specs in group for dependency resolution
	// (In full implementation, we'd store this in the Group struct)

	// Create group
	group := &Group{
		ID:          groupID,
		Name:        spec.Name,
		Strategy:    spec.Strategy,
		Containers:  containerIDs,
		State:       GroupStateCreated,
		Networking:  spec.Networking,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	gm.runtime.mu.Lock()
	gm.runtime.groups[groupID] = group
	gm.runtime.mu.Unlock()

	return group, nil
}

// StartGroup starts all containers in a group
func (gm *GroupManager) StartGroup(ctx context.Context, groupID string) error {
	group, err := gm.runtime.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// Start containers in dependency order
	started := make(map[string]bool)
	for _, containerID := range group.Containers {
		_, err := gm.runtime.GetContainer(ctx, containerID)
		if err != nil {
			return fmt.Errorf("failed to get container %s: %w", containerID, err)
		}

		// Check dependencies
		// In a full implementation, we'd store container specs in the group
		// For now, we'll start containers in order
		// Dependencies would be resolved by checking container names

		// Start container
		if err := gm.runtime.StartContainer(ctx, containerID); err != nil {
			return fmt.Errorf("failed to start container %s: %w", containerID, err)
		}
		started[containerID] = true
	}

	// Update group state
	group.State = GroupStateRunning
	group.UpdatedAt = time.Now()

	return nil
}

// StopGroup stops all containers in a group
func (gm *GroupManager) StopGroup(ctx context.Context, groupID string) error {
	group, err := gm.runtime.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// Stop all containers
	for _, containerID := range group.Containers {
		if err := gm.runtime.StopContainer(ctx, containerID); err != nil {
			// Log error but continue stopping other containers
			fmt.Printf("Warning: failed to stop container %s: %v\n", containerID, err)
		}
	}

	// Update group state
	group.State = GroupStateStopped
	group.UpdatedAt = time.Now()

	return nil
}

// RemoveGroup removes a container group
func (gm *GroupManager) RemoveGroup(ctx context.Context, groupID string) error {
	group, err := gm.runtime.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// Stop group first if running
	if group.State == GroupStateRunning {
		if err := gm.StopGroup(ctx, groupID); err != nil {
			return fmt.Errorf("failed to stop group: %w", err)
		}
	}

	// Remove all containers
	for _, containerID := range group.Containers {
		if err := gm.runtime.DeleteContainer(ctx, containerID); err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to remove container %s: %v\n", containerID, err)
		}
	}

	// Remove group
	gm.runtime.mu.Lock()
	delete(gm.runtime.groups, groupID)
	gm.runtime.mu.Unlock()

	return nil
}

// GetGroupStatus returns the status of all containers in a group
func (gm *GroupManager) GetGroupStatus(ctx context.Context, groupID string) (*GroupStatus, error) {
	group, err := gm.runtime.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Get status of all containers
	containerStatuses := make([]*ContainerStatus, 0, len(group.Containers))
	for _, containerID := range group.Containers {
		container, err := gm.runtime.GetContainer(ctx, containerID)
		if err != nil {
			continue
		}

		status := &ContainerStatus{
			ID:      container.ID,
			Name:    container.Name,
			State:   container.State,
			Pid:     container.Pid,
			Image:   container.Image,
			GroupID: container.GroupID,
		}
		containerStatuses = append(containerStatuses, status)
	}

	return &GroupStatus{
		ID:               group.ID,
		Name:             group.Name,
		Strategy:         group.Strategy,
		State:            group.State,
		Containers:       containerStatuses,
		ContainerCount:   len(group.Containers),
		RunningCount:     countRunning(containerStatuses),
		StoppedCount:     countStopped(containerStatuses),
		CreatedAt:        group.CreatedAt,
		UpdatedAt:        group.UpdatedAt,
	}, nil
}

// GroupStatus represents the status of a container group
type GroupStatus struct {
	ID             string
	Name           string
	Strategy       GroupStrategy
	State          GroupState
	Containers     []*ContainerStatus
	ContainerCount int
	RunningCount   int
	StoppedCount   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Helper functions

func generateGroupID(name string) string {
	return fmt.Sprintf("group-%s-%d", name, time.Now().UnixNano())
}

func generateContainerID(name string) string {
	return fmt.Sprintf("container-%s-%d", name, time.Now().UnixNano())
}

func parsePorts(ports []string) map[string]string {
	result := make(map[string]string)
	for _, port := range ports {
		// Parse "host:container" or "container"
		parts := splitPort(port)
		if len(parts) == 2 {
			result[parts[1]] = parts[0]
		} else {
			result[parts[0]] = parts[0]
		}
	}
	return result
}

func splitPort(port string) []string {
	// Simple port parsing
	// In full implementation, handle all formats
	return []string{port}
}

func countRunning(statuses []*ContainerStatus) int {
	count := 0
	for _, status := range statuses {
		if status.State == StateRunning {
			count++
		}
	}
	return count
}

func countStopped(statuses []*ContainerStatus) int {
	count := 0
	for _, status := range statuses {
		if status.State == StateStopped {
			count++
		}
	}
	return count
}

