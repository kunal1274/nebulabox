package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Runtime represents the NebulaBox custom container runtime engine
// This is a completely independent runtime that uses Linux primitives directly
type Runtime interface {
	// Container lifecycle
	CreateContainer(ctx context.Context, spec *ContainerSpec) (*Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	DeleteContainer(ctx context.Context, id string) error
	GetContainer(ctx context.Context, id string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)

	// Image management
	BuildImage(ctx context.Context, spec *BuildSpec) (*Image, error)
	PullImage(ctx context.Context, image string) error
	DeleteImage(ctx context.Context, image string) error
	ListImages(ctx context.Context) ([]*Image, error)
	GetImage(ctx context.Context, image string) (*Image, error)

	// Container groups
	CreateGroup(ctx context.Context, spec *GroupSpec) (*Group, error)
	StartGroup(ctx context.Context, groupID string) error
	StopGroup(ctx context.Context, groupID string) error
	GetGroup(ctx context.Context, groupID string) (*Group, error)
	ListGroups(ctx context.Context) ([]*Group, error)

	// Runtime info
	Info(ctx context.Context) (*RuntimeInfo, error)
	Version(ctx context.Context) (string, error)
}

// NebulaRuntime is the concrete implementation of the Runtime interface
type NebulaRuntime struct {
	containers    map[string]*Container
	images        map[string]*Image
	groups        map[string]*Group
	namespaces    *NamespaceManager
	cgroups       *CGroupManager
	filesystem    *FilesystemManager
	network       *NetworkManager
	process       *ProcessManager
	imageManager  *ImageManager
	Hierarchy     *HierarchyManager
	mu            sync.RWMutex
	storagePath   string
	statePath     string
}

// NewRuntime creates a new NebulaBox runtime instance
func NewRuntime(storagePath, statePath string) (*NebulaRuntime, error) {
	// Initialize namespace manager
	nsMgr, err := NewNamespaceManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace manager: %w", err)
	}

	// Initialize cgroup manager
	cgMgr, err := NewCGroupManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create cgroup manager: %w", err)
	}

	// Initialize filesystem manager
	fsMgr, err := NewFilesystemManager(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem manager: %w", err)
	}

	// Initialize network manager
	netMgr, err := NewNetworkManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create network manager: %w", err)
	}

	// Initialize process manager
	procMgr, err := NewProcessManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create process manager: %w", err)
	}

	// Initialize image manager
	imgMgr, err := NewImageManager(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create image manager: %w", err)
	}

	runtime := &NebulaRuntime{
		containers:  make(map[string]*Container),
		images:      make(map[string]*Image),
		groups:      make(map[string]*Group),
		namespaces:  nsMgr,
		cgroups:     cgMgr,
		filesystem:  fsMgr,
		network:     netMgr,
		process:     procMgr,
		imageManager: imgMgr,
		storagePath: storagePath,
		statePath:   statePath,
	}

	// Initialize hierarchy manager
	runtime.Hierarchy = NewHierarchyManager(runtime)

	return runtime, nil
}

// ContainerSpec defines how to create a container
type ContainerSpec struct {
	ID          string
	Name        string
	Image       string
	Command     []string
	Args        []string
	Env         map[string]string
	WorkingDir  string
	User        string
	NetworkMode string
	Network     string
	Ports       map[string]string // containerPort -> hostPort
	Volumes     map[string]string // containerPath -> hostPath
	Labels      map[string]string
	Resources   *ResourceLimits
	Security    *SecurityConfig
	GroupID     string // Optional: container group this belongs to
}

// ResourceLimits defines container resource constraints
type ResourceLimits struct {
	CPUShares    int64
	CPUQuota     int64
	MemoryLimit  int64 // bytes
	MemorySwap   int64 // bytes
	PidsLimit    int64
	BlkioWeight  uint16
	BlkioReadBps int64
	BlkioWriteBps int64
}

// SecurityConfig defines security settings for a container
type SecurityConfig struct {
	ReadonlyRootfs bool
	NoNewPrivs     bool
	AppArmorProfile string
	SELinuxLabel   string
	Capabilities   []string
	DropCapabilities []string
}

// Container represents a running or stopped container
type Container struct {
	ID          string
	Name        string
	Image       string
	State       ContainerState
	CreatedAt   time.Time
	StartedAt   *time.Time
	StoppedAt   *time.Time
	ExitCode    int
	Pid         int
	Spec        *ContainerSpec
	Namespaces  *NamespaceSet
	CGroupPath  string
	Network     *NetworkConfig
	Volumes     map[string]string
	Labels      map[string]string
	GroupID     string
}

// ContainerState represents the state of a container
type ContainerState string

const (
	StateCreated    ContainerState = "created"
	StateRunning    ContainerState = "running"
	StatePaused     ContainerState = "paused"
	StateStopped    ContainerState = "stopped"
	StateRemoved    ContainerState = "removed"
)

// Image represents a container image
type Image struct {
	ID          string
	Name        string
	Tag         string
	Digest      string
	Size        int64
	Layers      []ImageLayer
	Config      *ImageConfig
	CreatedAt   time.Time
	Manifest    *ImageManifest
}

// ImageLayer represents a single layer in an image
type ImageLayer struct {
	ID       string
	Digest   string
	Size     int64
	MediaType string
	Path     string
}

// ImageConfig contains image configuration
type ImageConfig struct {
	Architecture string
	OS           string
	Config       map[string]interface{}
	RootFS       *RootFS
}

// RootFS defines the root filesystem
type RootFS struct {
	Type    string
	DiffIDs []string
}

// ImageManifest represents the OCI image manifest
type ImageManifest struct {
	SchemaVersion int
	MediaType     string
	Config        *ManifestDescriptor
	Layers        []*ManifestDescriptor
}

// ManifestDescriptor describes a manifest entry
type ManifestDescriptor struct {
	MediaType string
	Size      int64
	Digest    string
}

// BuildSpec defines how to build an image
type BuildSpec struct {
	Version string
	Name    string
	Tag     string
	Base    *BaseImage
	Workdir string
	Env     map[string]string
	Expose  []int
	Steps   []BuildStep
	Health  *HealthCheck
	Labels  map[string]string
	Metadata map[string]interface{}
}

// BaseImage defines the base image
type BaseImage struct {
	Image string
	Tag   string
}

// BuildStep represents a build instruction
type BuildStep struct {
	Type    string // run, copy, volume, cmd, etc.
	Command string
	Source  string
	Dest    string
	Comment string
}

// HealthCheck defines container health check configuration
type HealthCheck struct {
	Type     string
	Path     string
	Port     int
	Interval int
	Timeout  int
	Retries  int
}

// GroupSpec defines a container group
type GroupSpec struct {
	Name      string
	Strategy  GroupStrategy
	Containers []GroupContainerSpec
	Networking *GroupNetworking
	Resources *GroupResourceLimits
}

// GroupStrategy defines how containers are grouped
type GroupStrategy string

const (
	StrategyMonolithic      GroupStrategy = "monolithic"
	StrategyFrontendBackend GroupStrategy = "frontend-backend"
	StrategyThreeTier       GroupStrategy = "three-tier"
	StrategyMicroservices   GroupStrategy = "microservices"
	StrategyCustom          GroupStrategy = "custom"
)

// GroupContainerSpec defines a container in a group
type GroupContainerSpec struct {
	Name    string
	Image   string
	Ports   []string
	Env     map[string]string
	Volumes map[string]string
	DependsOn []string // Container dependencies
}

// GroupNetworking defines networking for a group
type GroupNetworking struct {
	Internal bool
	Bridge   string
	Subnet   string
}

// GroupResourceLimits defines resource limits for a group
type GroupResourceLimits struct {
	MaxCPU    float64
	MaxMemory int64
	MaxDisk   int64
}

// Group represents a container group
type Group struct {
	ID          string
	Name        string
	Strategy    GroupStrategy
	Containers  []string // Container IDs
	State       GroupState
	Networking  *GroupNetworking
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GroupState represents the state of a group
type GroupState string

const (
	GroupStateCreated  GroupState = "created"
	GroupStateRunning  GroupState = "running"
	GroupStateStopped GroupState = "stopped"
	GroupStateRemoved GroupState = "removed"
)

// NetworkConfig represents container network configuration
type NetworkConfig struct {
	Mode       string
	Bridge     string
	IPAddress  string
	Gateway    string
	PortMappings map[string]string
}

// RuntimeInfo contains runtime information
type RuntimeInfo struct {
	Name        string
	Version     string
	OS          string
	Arch        string
	Kernel      string
	Containers  int
	Images      int
	Groups      int
	StoragePath string
	StatePath   string
}

// CreateContainer creates a new container from a spec
func (r *NebulaRuntime) CreateContainer(ctx context.Context, spec *ContainerSpec) (*Container, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if container already exists
	if _, exists := r.containers[spec.ID]; exists {
		return nil, fmt.Errorf("container %s already exists", spec.ID)
	}

	// Verify image exists
	image, exists := r.images[spec.Image]
	if !exists {
		return nil, fmt.Errorf("image %s not found", spec.Image)
	}

	// Create namespaces
	nsSet, err := r.namespaces.CreateNamespaceSet(spec.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create namespaces: %w", err)
	}

	// Create cgroup
	cgroupPath, err := r.cgroups.CreateCGroup(spec.ID, spec.Resources)
	if err != nil {
		return nil, fmt.Errorf("failed to create cgroup: %w", err)
	}

	// Create filesystem
	_, err = r.filesystem.CreateRootfs(spec.ID, image)
	if err != nil {
		return nil, fmt.Errorf("failed to create rootfs: %w", err)
	}

	// Setup network
	netConfig, err := r.network.SetupContainerNetwork(spec.ID, spec.Network, spec.Ports)
	if err != nil {
		return nil, fmt.Errorf("failed to setup network: %w", err)
	}

	// Create container object
	container := &Container{
		ID:          spec.ID,
		Name:        spec.Name,
		Image:       spec.Image,
		State:       StateCreated,
		CreatedAt:   time.Now(),
		Spec:        spec,
		Namespaces:  nsSet,
		CGroupPath:  cgroupPath,
		Network:     netConfig,
		Volumes:     spec.Volumes,
		Labels:      spec.Labels,
		GroupID:     spec.GroupID,
	}

	r.containers[spec.ID] = container
	return container, nil
}

// StartContainer starts a container
func (r *NebulaRuntime) StartContainer(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	container, exists := r.containers[id]
	if !exists {
		return fmt.Errorf("container %s not found", id)
	}

	if container.State == StateRunning {
		return fmt.Errorf("container %s is already running", id)
	}

	// Start the container process
	pid, err := r.process.StartContainer(ctx, container)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	now := time.Now()
	container.Pid = pid
	container.State = StateRunning
	container.StartedAt = &now
	container.StoppedAt = nil

	return nil
}

// StopContainer stops a running container
func (r *NebulaRuntime) StopContainer(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	container, exists := r.containers[id]
	if !exists {
		return fmt.Errorf("container %s not found", id)
	}

	if container.State != StateRunning {
		return fmt.Errorf("container %s is not running", id)
	}

	// Stop the container process
	exitCode, err := r.process.StopContainer(ctx, container)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	now := time.Now()
	container.State = StateStopped
	container.StoppedAt = &now
	container.ExitCode = exitCode

	return nil
}

// DeleteContainer removes a container
func (r *NebulaRuntime) DeleteContainer(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	container, exists := r.containers[id]
	if !exists {
		return fmt.Errorf("container %s not found", id)
	}

	if container.State == StateRunning {
		return fmt.Errorf("cannot delete running container %s", id)
	}

	// Cleanup namespaces
	if err := r.namespaces.CleanupNamespaceSet(container.Namespaces); err != nil {
		return fmt.Errorf("failed to cleanup namespaces: %w", err)
	}

	// Cleanup cgroup
	if err := r.cgroups.RemoveCGroup(container.CGroupPath); err != nil {
		return fmt.Errorf("failed to remove cgroup: %w", err)
	}

	// Cleanup filesystem
	if err := r.filesystem.RemoveRootfs(container.ID); err != nil {
		return fmt.Errorf("failed to remove rootfs: %w", err)
	}

	// Cleanup network
	if err := r.network.CleanupContainerNetwork(container.ID); err != nil {
		return fmt.Errorf("failed to cleanup network: %w", err)
	}

	delete(r.containers, id)
	return nil
}

// GetContainer retrieves a container by ID
func (r *NebulaRuntime) GetContainer(ctx context.Context, id string) (*Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	container, exists := r.containers[id]
	if !exists {
		return nil, fmt.Errorf("container %s not found", id)
	}

	return container, nil
}

// ListContainers returns all containers
func (r *NebulaRuntime) ListContainers(ctx context.Context) ([]*Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	containers := make([]*Container, 0, len(r.containers))
	for _, container := range r.containers {
		containers = append(containers, container)
	}

	return containers, nil
}

// BuildImage builds an image from a BuildSpec
func (r *NebulaRuntime) BuildImage(ctx context.Context, spec *BuildSpec) (*Image, error) {
	image, err := r.imageManager.BuildImage(ctx, spec, r)
	if err != nil {
		return nil, err
	}

	// Store image
	imageRef := fmt.Sprintf("%s:%s", spec.Name, spec.Tag)
	r.mu.Lock()
	r.images[imageRef] = image
	r.mu.Unlock()

	return image, nil
}

// PullImage pulls an image from a registry
func (r *NebulaRuntime) PullImage(ctx context.Context, image string) error {
	return r.imageManager.PullImage(ctx, image)
}

// DeleteImage removes an image
func (r *NebulaRuntime) DeleteImage(ctx context.Context, image string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.images[image]; !exists {
		return fmt.Errorf("image %s not found", image)
	}

	delete(r.images, image)
	return nil
}

// ListImages returns all images
func (r *NebulaRuntime) ListImages(ctx context.Context) ([]*Image, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	images := make([]*Image, 0, len(r.images))
	for _, image := range r.images {
		images = append(images, image)
	}

	return images, nil
}

// GetImage retrieves an image by name
func (r *NebulaRuntime) GetImage(ctx context.Context, image string) (*Image, error) {
	r.mu.RLock()
	img, exists := r.images[image]
	r.mu.RUnlock()

	if exists {
		return img, nil
	}

	// Try image manager
	return r.imageManager.GetImage(image)
}

// CreateGroup creates a new container group
func (r *NebulaRuntime) CreateGroup(ctx context.Context, spec *GroupSpec) (*Group, error) {
	groupMgr := NewGroupManager(r)
	return groupMgr.CreateGroup(ctx, spec)
}

// StartGroup starts all containers in a group
func (r *NebulaRuntime) StartGroup(ctx context.Context, groupID string) error {
	groupMgr := NewGroupManager(r)
	return groupMgr.StartGroup(ctx, groupID)
}

// StopGroup stops all containers in a group
func (r *NebulaRuntime) StopGroup(ctx context.Context, groupID string) error {
	groupMgr := NewGroupManager(r)
	return groupMgr.StopGroup(ctx, groupID)
}

// GetGroup retrieves a group by ID
func (r *NebulaRuntime) GetGroup(ctx context.Context, groupID string) (*Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	group, exists := r.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group %s not found", groupID)
	}

	return group, nil
}

// ListGroups returns all groups
func (r *NebulaRuntime) ListGroups(ctx context.Context) ([]*Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groups := make([]*Group, 0, len(r.groups))
	for _, group := range r.groups {
		groups = append(groups, group)
	}

	return groups, nil
}

// Info returns runtime information
func (r *NebulaRuntime) Info(ctx context.Context) (*RuntimeInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return &RuntimeInfo{
		Name:        "NebulaBox Runtime",
		Version:     "0.1.0",
		OS:          "linux",
		Arch:        "amd64",
		Kernel:      "unknown", // Will be populated from system
		Containers:  len(r.containers),
		Images:      len(r.images),
		Groups:      len(r.groups),
		StoragePath: r.storagePath,
		StatePath:   r.statePath,
	}, nil
}

// Version returns the runtime version
func (r *NebulaRuntime) Version(ctx context.Context) (string, error) {
	return "0.1.0", nil
}

