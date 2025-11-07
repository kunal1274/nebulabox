package engine

import (
	"context"
	"fmt"
)

// GetContainerLogs gets logs for a container
func (r *NebulaRuntime) GetContainerLogs(ctx context.Context, containerID string) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.containers[containerID]
	if !exists {
		return nil, fmt.Errorf("container not found: %s", containerID)
	}

	// Get logs from process manager
	logs, err := r.process.CollectLogs(containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	return logs, nil
}

// GetHierarchyContainers lists all containers in a hierarchy
func (r *NebulaRuntime) GetHierarchyContainers(ctx context.Context, rootID string) ([]*Container, error) {
	if r.Hierarchy == nil {
		return nil, fmt.Errorf("hierarchy manager not initialized")
	}
	return r.Hierarchy.ListContainersInHierarchy(ctx, rootID)
}

