package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nebulabox/nebulabox/internal/engine"
)

// EngineClient provides direct integration with the custom engine
type EngineClient struct {
	runtime *engine.NebulaRuntime
	ctx     context.Context
}

// NewEngineClient creates a new engine client
func NewEngineClient() (*EngineClient, error) {
	// Determine storage paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	storagePath := filepath.Join(homeDir, ".nebulabox", "storage")
	statePath := filepath.Join(homeDir, ".nebulabox", "state")

	// Create runtime
	runtime, err := engine.NewRuntime(storagePath, statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	return &EngineClient{
		runtime: runtime,
		ctx:     context.Background(),
	}, nil
}

// BuildImage builds an image from a BuildSpec
func (ec *EngineClient) BuildImage(spec *engine.BuildSpec) (*engine.Image, error) {
	return ec.runtime.BuildImage(ec.ctx, spec)
}

// CreateContainer creates a container
func (ec *EngineClient) CreateContainer(spec *engine.ContainerSpec) (*engine.Container, error) {
	return ec.runtime.CreateContainer(ec.ctx, spec)
}

// StartContainer starts a container
func (ec *EngineClient) StartContainer(id string) error {
	return ec.runtime.StartContainer(ec.ctx, id)
}

// StopContainer stops a container
func (ec *EngineClient) StopContainer(id string) error {
	return ec.runtime.StopContainer(ec.ctx, id)
}

// ListContainers lists all containers
func (ec *EngineClient) ListContainers() ([]*engine.Container, error) {
	return ec.runtime.ListContainers(ec.ctx)
}

// GetContainer gets a container by ID
func (ec *EngineClient) GetContainer(id string) (*engine.Container, error) {
	return ec.runtime.GetContainer(ec.ctx, id)
}

// DeleteContainer deletes a container
func (ec *EngineClient) DeleteContainer(id string) error {
	return ec.runtime.DeleteContainer(ec.ctx, id)
}

// PullImage pulls an image
func (ec *EngineClient) PullImage(image string) error {
	return ec.runtime.PullImage(ec.ctx, image)
}

// ListImages lists all images
func (ec *EngineClient) ListImages() ([]*engine.Image, error) {
	return ec.runtime.ListImages(ec.ctx)
}

// DeleteImage deletes an image
func (ec *EngineClient) DeleteImage(imageRef string) error {
	return ec.runtime.DeleteImage(ec.ctx, imageRef)
}

// CreateGroup creates a container group
func (ec *EngineClient) CreateGroup(spec *engine.GroupSpec) (*engine.Group, error) {
	return ec.runtime.CreateGroup(ec.ctx, spec)
}

// StartGroup starts a container group
func (ec *EngineClient) StartGroup(groupID string) error {
	return ec.runtime.StartGroup(ec.ctx, groupID)
}

// StopGroup stops a container group
func (ec *EngineClient) StopGroup(groupID string) error {
	return ec.runtime.StopGroup(ec.ctx, groupID)
}

// GetRuntimeInfo returns runtime information
func (ec *EngineClient) GetRuntimeInfo() (*engine.RuntimeInfo, error) {
	return ec.runtime.Info(ec.ctx)
}

// ListGroups lists all container groups
func (ec *EngineClient) ListGroups() ([]*engine.Group, error) {
	return ec.runtime.ListGroups(ec.ctx)
}

// GetGroup gets a group by ID
func (ec *EngineClient) GetGroup(groupID string) (*engine.Group, error) {
	return ec.runtime.GetGroup(ec.ctx, groupID)
}

// CreateNestedContainer creates a container within another container
func (ec *EngineClient) CreateNestedContainer(parentID string, spec *engine.ContainerSpec) (*engine.Container, error) {
	return ec.runtime.Hierarchy.CreateNestedContainer(ec.ctx, parentID, spec)
}

// CreateNestedGroup creates a group within a container
func (ec *EngineClient) CreateNestedGroup(containerID string, spec *engine.GroupSpec) (*engine.Group, error) {
	return ec.runtime.Hierarchy.CreateNestedGroup(ec.ctx, containerID, spec)
}

// AddContainerToGroup adds a container to a group
func (ec *EngineClient) AddContainerToGroup(containerID, groupID string) error {
	return ec.runtime.Hierarchy.AddContainerToGroup(ec.ctx, containerID, groupID)
}

// GetHierarchy gets the hierarchy tree for a container
func (ec *EngineClient) GetHierarchy(containerID string) (*engine.HierarchyTree, error) {
	return ec.runtime.Hierarchy.GetHierarchy(ec.ctx, containerID)
}

// GetFullHierarchy gets all hierarchy trees
func (ec *EngineClient) GetFullHierarchy() ([]*engine.HierarchyTree, error) {
	return ec.runtime.Hierarchy.GetFullHierarchy(ec.ctx)
}

// ListContainersInHierarchy lists all containers in a hierarchy
func (ec *EngineClient) ListContainersInHierarchy(rootID string) ([]*engine.Container, error) {
	return ec.runtime.GetHierarchyContainers(ec.ctx, rootID)
}

// GetContainerLogs gets logs for a container
func (ec *EngineClient) GetContainerLogs(containerID string) ([]byte, error) {
	return ec.runtime.GetContainerLogs(ec.ctx, containerID)
}

