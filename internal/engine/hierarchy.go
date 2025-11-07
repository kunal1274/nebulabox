package engine

import (
	"context"
	"fmt"
	"sync"
)

// HierarchicalContainer represents a container that can contain other containers
type HierarchicalContainer struct {
	*Container
	Children    []string // Container IDs of child containers
	ParentID    string   // Parent container ID (empty if root)
	Groups      []string // Group IDs this container belongs to
	NestedGroups []string // Group IDs that exist within this container
	Depth       int      // Nesting depth (0 = root level)
	mu          sync.RWMutex
}

// HierarchyManager manages hierarchical container relationships
type HierarchyManager struct {
	runtime     *NebulaRuntime
	hierarchies map[string]*HierarchicalContainer // containerID -> HierarchicalContainer
	mu          sync.RWMutex
}

// NewHierarchyManager creates a new hierarchy manager
func NewHierarchyManager(runtime *NebulaRuntime) *HierarchyManager {
	return &HierarchyManager{
		runtime:     runtime,
		hierarchies: make(map[string]*HierarchicalContainer),
	}
}

// CreateNestedContainer creates a container within another container
func (hm *HierarchyManager) CreateNestedContainer(ctx context.Context, parentID string, spec *ContainerSpec) (*Container, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Verify parent exists
	parent, err := hm.runtime.GetContainer(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("parent container %s not found: %w", parentID, err)
	}

	// Get or create hierarchical container for parent
	parentHier := hm.getOrCreateHierarchical(parent)

	// Create child container
	child, err := hm.runtime.CreateContainer(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create nested container: %w", err)
	}

	// Get or create hierarchical container for child
	childHier := hm.getOrCreateHierarchical(child)
	childHier.ParentID = parentID
	childHier.Depth = parentHier.Depth + 1

	// Add child to parent
	parentHier.mu.Lock()
	parentHier.Children = append(parentHier.Children, child.ID)
	parentHier.mu.Unlock()

	return child, nil
}

// CreateNestedGroup creates a group within a container
func (hm *HierarchyManager) CreateNestedGroup(ctx context.Context, containerID string, spec *GroupSpec) (*Group, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Verify container exists
	container, err := hm.runtime.GetContainer(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("container %s not found: %w", containerID, err)
	}

	// Get or create hierarchical container
	hier := hm.getOrCreateHierarchical(container)

	// Create group
	group, err := hm.runtime.CreateGroup(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create nested group: %w", err)
	}

	// Add group to container's nested groups
	hier.mu.Lock()
	hier.NestedGroups = append(hier.NestedGroups, group.ID)
	hier.mu.Unlock()

	return group, nil
}

// AddContainerToGroup adds a container to a group (can be from any level)
func (hm *HierarchyManager) AddContainerToGroup(ctx context.Context, containerID, groupID string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Verify container exists
	container, err := hm.runtime.GetContainer(ctx, containerID)
	if err != nil {
		return fmt.Errorf("container %s not found: %w", containerID, err)
	}

	// Verify group exists
	group, err := hm.runtime.GetGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("group %s not found: %w", groupID, err)
	}

	// Get or create hierarchical container
	hier := hm.getOrCreateHierarchical(container)

	// Add to group
	hier.mu.Lock()
	hier.Groups = append(hier.Groups, groupID)
	hier.mu.Unlock()

	// Add container to group
	group.Containers = append(group.Containers, containerID)

	return nil
}

// GetHierarchy returns the full hierarchy tree for a container
func (hm *HierarchyManager) GetHierarchy(ctx context.Context, containerID string) (*HierarchyTree, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	hier, exists := hm.hierarchies[containerID]
	if !exists {
		return nil, fmt.Errorf("container %s not found in hierarchy", containerID)
	}

	return hm.buildTree(ctx, hier), nil
}

// GetFullHierarchy returns the complete hierarchy from root
func (hm *HierarchyManager) GetFullHierarchy(ctx context.Context) ([]*HierarchyTree, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	// Find all root containers (no parent)
	roots := make([]*HierarchicalContainer, 0)
	for _, hier := range hm.hierarchies {
		if hier.ParentID == "" {
			roots = append(roots, hier)
		}
	}

	trees := make([]*HierarchyTree, 0, len(roots))
	for _, root := range roots {
		trees = append(trees, hm.buildTree(ctx, root))
	}

	return trees, nil
}

// HierarchyTree represents a container hierarchy tree
type HierarchyTree struct {
	Container   *Container
	Children    []*HierarchyTree
	Groups      []*Group
	NestedGroups []*Group
	Depth       int
}

// buildTree builds a hierarchy tree recursively
func (hm *HierarchyManager) buildTree(ctx context.Context, hier *HierarchicalContainer) *HierarchyTree {
	tree := &HierarchyTree{
		Container:    hier.Container,
		Children:     make([]*HierarchyTree, 0),
		Groups:       make([]*Group, 0),
		NestedGroups: make([]*Group, 0),
		Depth:        hier.Depth,
	}

	// Build children trees
	hier.mu.RLock()
	children := make([]string, len(hier.Children))
	copy(children, hier.Children)
	hier.mu.RUnlock()

	for _, childID := range children {
		if childHier, exists := hm.hierarchies[childID]; exists {
			tree.Children = append(tree.Children, hm.buildTree(ctx, childHier))
		}
	}

	// Get groups
	hier.mu.RLock()
	groups := make([]string, len(hier.Groups))
	copy(groups, hier.Groups)
	nestedGroups := make([]string, len(hier.NestedGroups))
	copy(nestedGroups, hier.NestedGroups)
	hier.mu.RUnlock()

	for _, groupID := range groups {
		if group, err := hm.runtime.GetGroup(ctx, groupID); err == nil {
			tree.Groups = append(tree.Groups, group)
		}
	}

	for _, groupID := range nestedGroups {
		if group, err := hm.runtime.GetGroup(ctx, groupID); err == nil {
			tree.NestedGroups = append(tree.NestedGroups, group)
		}
	}

	return tree
}

// getOrCreateHierarchical gets or creates a hierarchical container
func (hm *HierarchyManager) getOrCreateHierarchical(container *Container) *HierarchicalContainer {
	if hier, exists := hm.hierarchies[container.ID]; exists {
		return hier
	}

	hier := &HierarchicalContainer{
		Container:   container,
		Children:    make([]string, 0),
		Groups:      make([]string, 0),
		NestedGroups: make([]string, 0),
		Depth:       0,
	}

	hm.hierarchies[container.ID] = hier
	return hier
}

// RemoveContainer removes a container from hierarchy
func (hm *HierarchyManager) RemoveContainer(ctx context.Context, containerID string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hier, exists := hm.hierarchies[containerID]
	if !exists {
		return nil // Not in hierarchy, nothing to do
	}

	// Remove from parent's children
	if hier.ParentID != "" {
		if parentHier, exists := hm.hierarchies[hier.ParentID]; exists {
			parentHier.mu.Lock()
			for i, childID := range parentHier.Children {
				if childID == containerID {
					parentHier.Children = append(parentHier.Children[:i], parentHier.Children[i+1:]...)
					break
				}
			}
			parentHier.mu.Unlock()
		}
	}

	// Remove all children recursively
	for _, childID := range hier.Children {
		hm.RemoveContainer(ctx, childID)
	}

	delete(hm.hierarchies, containerID)
	return nil
}

// ListContainersInHierarchy lists all containers in a hierarchy
func (hm *HierarchyManager) ListContainersInHierarchy(ctx context.Context, rootID string) ([]*Container, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	hier, exists := hm.hierarchies[rootID]
	if !exists {
		return nil, fmt.Errorf("container %s not found in hierarchy", rootID)
	}

	containers := make([]*Container, 0)
	hm.collectContainers(ctx, hier, &containers)
	return containers, nil
}

// collectContainers collects all containers recursively
func (hm *HierarchyManager) collectContainers(ctx context.Context, hier *HierarchicalContainer, containers *[]*Container) {
	*containers = append(*containers, hier.Container)

	hier.mu.RLock()
	children := make([]string, len(hier.Children))
	copy(children, hier.Children)
	hier.mu.RUnlock()

	for _, childID := range children {
		if childHier, exists := hm.hierarchies[childID]; exists {
			hm.collectContainers(ctx, childHier, containers)
		}
	}
}

