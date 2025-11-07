package engine

import (
	"context"
	"fmt"
	"time"
)

// ContainerLifecycle manages container lifecycle operations
type ContainerLifecycle struct {
	runtime *NebulaRuntime
}

// NewContainerLifecycle creates a new container lifecycle manager
func NewContainerLifecycle(runtime *NebulaRuntime) *ContainerLifecycle {
	return &ContainerLifecycle{
		runtime: runtime,
	}
}

// Create creates a new container
func (cl *ContainerLifecycle) Create(ctx context.Context, spec *ContainerSpec) (*Container, error) {
	return cl.runtime.CreateContainer(ctx, spec)
}

// Start starts a container
func (cl *ContainerLifecycle) Start(ctx context.Context, id string) error {
	return cl.runtime.StartContainer(ctx, id)
}

// Stop stops a container
func (cl *ContainerLifecycle) Stop(ctx context.Context, id string) error {
	return cl.runtime.StopContainer(ctx, id)
}

// Restart restarts a container
func (cl *ContainerLifecycle) Restart(ctx context.Context, id string) error {
	// Stop the container
	if err := cl.Stop(ctx, id); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Wait a bit
	time.Sleep(1 * time.Second)

	// Start the container
	if err := cl.Start(ctx, id); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// Pause pauses a container
func (cl *ContainerLifecycle) Pause(ctx context.Context, id string) error {
	container, err := cl.runtime.GetContainer(ctx, id)
	if err != nil {
		return err
	}

	if container.State != StateRunning {
		return fmt.Errorf("container %s is not running", id)
	}

	// Send SIGSTOP to pause
	// In full implementation, use cgroup freezer
	container.State = StatePaused
	return nil
}

// Unpause unpauses a container
func (cl *ContainerLifecycle) Unpause(ctx context.Context, id string) error {
	container, err := cl.runtime.GetContainer(ctx, id)
	if err != nil {
		return err
	}

	if container.State != StatePaused {
		return fmt.Errorf("container %s is not paused", id)
	}

	// Send SIGCONT to unpause
	container.State = StateRunning
	return nil
}

// Remove removes a container
func (cl *ContainerLifecycle) Remove(ctx context.Context, id string) error {
	return cl.runtime.DeleteContainer(ctx, id)
}

// GetStatus returns container status
func (cl *ContainerLifecycle) GetStatus(ctx context.Context, id string) (*ContainerStatus, error) {
	container, err := cl.runtime.GetContainer(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if process is still running
	isRunning := false
	if container.State == StateRunning && container.Pid > 0 {
		// Check process
		isRunning = cl.runtime.process.IsProcessRunning(container.Pid)
		if !isRunning {
			// Process died, update state
			container.State = StateStopped
		}
	}

	return &ContainerStatus{
		ID:          container.ID,
		Name:        container.Name,
		State:       container.State,
		Pid:         container.Pid,
		IsRunning:   isRunning,
		CreatedAt:   container.CreatedAt,
		StartedAt:   container.StartedAt,
		StoppedAt:   container.StoppedAt,
		ExitCode:    container.ExitCode,
		Image:       container.Image,
		Labels:      container.Labels,
		GroupID:     container.GroupID,
	}, nil
}

// ContainerStatus represents the status of a container
type ContainerStatus struct {
	ID        string
	Name      string
	State     ContainerState
	Pid       int
	IsRunning bool
	CreatedAt time.Time
	StartedAt *time.Time
	StoppedAt *time.Time
	ExitCode  int
	Image     string
	Labels    map[string]string
	GroupID   string
}

// Wait waits for a container to stop
func (cl *ContainerLifecycle) Wait(ctx context.Context, id string) (int, error) {
	container, err := cl.runtime.GetContainer(ctx, id)
	if err != nil {
		return 0, err
	}

	if container.State != StateRunning {
		return container.ExitCode, nil
	}

	// Wait for process to exit
	exitCode, err := cl.runtime.process.WaitForProcess(container.Pid)
	if err != nil {
		return 0, err
	}

	// Update container state
	now := time.Now()
	container.State = StateStopped
	container.StoppedAt = &now
	container.ExitCode = exitCode

	return exitCode, nil
}

// GetStats returns container statistics
func (cl *ContainerLifecycle) GetStats(ctx context.Context, id string) (*ContainerStats, error) {
	container, err := cl.runtime.GetContainer(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get CPU and memory usage from cgroup
	cpuUsage := 0.0
	memoryUsage := int64(0)

	if container.CGroupPath != "" {
		if usage, err := cl.runtime.cgroups.GetCPUUsage(container.CGroupPath); err == nil {
			cpuUsage = usage
		}
		if usage, err := cl.runtime.cgroups.GetMemoryUsage(container.CGroupPath); err == nil {
			memoryUsage = usage
		}
	}

	return &ContainerStats{
		ID:          container.ID,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		Pid:         container.Pid,
		State:       container.State,
	}, nil
}

// ContainerStats represents container statistics
type ContainerStats struct {
	ID          string
	CPUUsage    float64
	MemoryUsage int64
	Pid         int
	State       ContainerState
}

