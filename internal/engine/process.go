package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// ProcessManager manages container processes
type ProcessManager struct {
	processes map[string]*ProcessInfo // containerID -> ProcessInfo
}

// ProcessInfo contains information about a container process
type ProcessInfo struct {
	Pid        int
	ContainerID string
	Command    []string
	Args       []string
	StartedAt  time.Time
	ExitCode   int
	State      ProcessState
}

// ProcessState represents the state of a process
type ProcessState string

const (
	ProcessStateRunning ProcessState = "running"
	ProcessStateStopped ProcessState = "stopped"
	ProcessStateExited  ProcessState = "exited"
)

// NewProcessManager creates a new process manager
func NewProcessManager() (*ProcessManager, error) {
	return &ProcessManager{
		processes: make(map[string]*ProcessInfo),
	}, nil
}

// StartContainer starts a container process
func (pm *ProcessManager) StartContainer(ctx context.Context, container *Container) (int, error) {
	// Create command
	var cmd *exec.Cmd
	if len(container.Spec.Command) > 0 {
		cmd = exec.CommandContext(ctx, container.Spec.Command[0], container.Spec.Command[1:]...)
	} else {
		// Default command
		cmd = exec.CommandContext(ctx, "/bin/sh")
	}

	// Set working directory
	if container.Spec.WorkingDir != "" {
		cmd.Dir = container.Spec.WorkingDir
	}

	// Set environment variables
	env := os.Environ()
	for k, v := range container.Spec.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// Setup namespaces
	// This would use syscall.Clone or unshare
	// For POC, we'll run without full isolation

	// Start the process
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start process: %w", err)
	}

	pid := cmd.Process.Pid

	// Store process info
	pm.processes[container.ID] = &ProcessInfo{
		Pid:        pid,
		ContainerID: container.ID,
		Command:    container.Spec.Command,
		Args:       container.Spec.Args,
		StartedAt:  time.Now(),
		State:      ProcessStateRunning,
	}

	return pid, nil
}

// StopContainer stops a container process
func (pm *ProcessManager) StopContainer(ctx context.Context, container *Container) (int, error) {
	procInfo, exists := pm.processes[container.ID]
	if !exists {
		return 0, fmt.Errorf("process for container %s not found", container.ID)
	}

	// Find the process
	process, err := os.FindProcess(procInfo.Pid)
	if err != nil {
		return 0, fmt.Errorf("failed to find process: %w", err)
	}

	// Send SIGTERM first
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// Process might already be dead
	}

	// Wait a bit for graceful shutdown
	time.Sleep(2 * time.Second)

	// Check if process is still running
	if process, err := os.FindProcess(procInfo.Pid); err == nil {
		// Send SIGKILL if still running
		process.Signal(syscall.SIGKILL)
	}

	// Wait for process to exit
	state, err := process.Wait()
	exitCode := 0
	if err == nil {
		if status, ok := state.Sys().(syscall.WaitStatus); ok {
			exitCode = status.ExitStatus()
		}
	}

	procInfo.State = ProcessStateExited
	procInfo.ExitCode = exitCode

	return exitCode, nil
}

// GetProcessInfo returns process information
func (pm *ProcessManager) GetProcessInfo(containerID string) (*ProcessInfo, error) {
	procInfo, exists := pm.processes[containerID]
	if !exists {
		return nil, fmt.Errorf("process info for container %s not found", containerID)
	}

	return procInfo, nil
}

// IsProcessRunning checks if a process is running
func (pm *ProcessManager) IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// CollectLogs collects logs from a container process
func (pm *ProcessManager) CollectLogs(containerID string) ([]byte, error) {
	// In a full implementation, we'd capture stdout/stderr
	// For now, return empty
	return []byte{}, nil
}

// WaitForProcess waits for a process to exit
func (pm *ProcessManager) WaitForProcess(pid int) (int, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return 0, fmt.Errorf("failed to find process: %w", err)
	}

	state, err := process.Wait()
	if err != nil {
		return 0, err
	}

	exitCode := 0
	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		exitCode = status.ExitStatus()
	}

	return exitCode, nil
}

