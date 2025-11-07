package groups

import (
	"fmt"
	"sync"
	"time"
)

// ContainerGroup represents a logical grouping of containers
type ContainerGroup struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	ParentGroupID   string            `json:"parentGroupId,omitempty"` // For hierarchy
	ContainerIDs    []string          `json:"containerIds"`
	Labels          map[string]string `json:"labels,omitempty"`
	SharedResources *SharedResources  `json:"sharedResources,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

// SharedResources defines resources shared across a group
type SharedResources struct {
	Network   string            `json:"network,omitempty"`   // Shared network
	Volumes   []string          `json:"volumes,omitempty"`   // Shared volume names
	EnvVars   map[string]string `json:"envVars,omitempty"`   // Shared environment variables
	Ports     []string          `json:"ports,omitempty"`      // Shared port ranges
	Labels    map[string]string `json:"labels,omitempty"`    // Shared labels
}

// ContainerRelationship represents a parent-child relationship between containers
type ContainerRelationship struct {
	ParentID   string    `json:"parentId"`
	ChildID    string    `json:"childId"`
	Type       string    `json:"type"`       // dependency, composition, nested
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// GroupManager manages container groups and hierarchies
type GroupManager struct {
	groups        map[string]*ContainerGroup
	relationships map[string][]*ContainerRelationship // parentID -> relationships
	childMap      map[string]*ContainerRelationship  // childID -> relationship
	mu            sync.RWMutex
}

// NewGroupManager creates a new group manager
func NewGroupManager() *GroupManager {
	return &GroupManager{
		groups:        make(map[string]*ContainerGroup),
		relationships: make(map[string][]*ContainerRelationship),
		childMap:      make(map[string]*ContainerRelationship),
	}
}

// CreateGroup creates a new container group
func (gm *GroupManager) CreateGroup(name, description string, parentGroupID string, sharedResources *SharedResources) (*ContainerGroup, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group := &ContainerGroup{
		ID:              fmt.Sprintf("group-%d", time.Now().UnixNano()),
		Name:            name,
		Description:     description,
		ParentGroupID:   parentGroupID,
		ContainerIDs:   []string{},
		SharedResources: sharedResources,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validate parent exists if specified
	if parentGroupID != "" {
		if _, exists := gm.groups[parentGroupID]; !exists {
			return nil, fmt.Errorf("parent group not found: %s", parentGroupID)
		}
	}

	gm.groups[group.ID] = group
	return group, nil
}

// GetGroup retrieves a group by ID
func (gm *GroupManager) GetGroup(groupID string) (*ContainerGroup, bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	group, exists := gm.groups[groupID]
	return group, exists
}

// ListGroups returns all groups, optionally filtered by parent
func (gm *GroupManager) ListGroups(parentGroupID *string) []*ContainerGroup {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	var result []*ContainerGroup
	for _, group := range gm.groups {
		if parentGroupID == nil {
			result = append(result, group)
		} else if group.ParentGroupID == *parentGroupID {
			result = append(result, group)
		}
	}
	return result
}

// UpdateGroup updates a group
func (gm *GroupManager) UpdateGroup(groupID string, name, description *string, sharedResources *SharedResources) (*ContainerGroup, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group, exists := gm.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}

	if name != nil {
		group.Name = *name
	}
	if description != nil {
		group.Description = *description
	}
	if sharedResources != nil {
		group.SharedResources = sharedResources
	}

	group.UpdatedAt = time.Now()
	return group, nil
}

// DeleteGroup removes a group (only if empty or force=true)
func (gm *GroupManager) DeleteGroup(groupID string, force bool) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group, exists := gm.groups[groupID]
	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Check if group has containers
	if len(group.ContainerIDs) > 0 && !force {
		return fmt.Errorf("group has containers, cannot delete without force")
	}

	// Check if group has child groups
	for _, g := range gm.groups {
		if g.ParentGroupID == groupID {
			if !force {
				return fmt.Errorf("group has child groups, cannot delete without force")
			}
			// Force delete: orphan child groups
			g.ParentGroupID = ""
		}
	}

	delete(gm.groups, groupID)
	return nil
}

// AddContainerToGroup adds a container to a group
func (gm *GroupManager) AddContainerToGroup(groupID, containerID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group, exists := gm.groups[groupID]
	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Check if container already in group
	for _, id := range group.ContainerIDs {
		if id == containerID {
			return nil // Already in group
		}
	}

	group.ContainerIDs = append(group.ContainerIDs, containerID)
	group.UpdatedAt = time.Now()
	return nil
}

// RemoveContainerFromGroup removes a container from a group
func (gm *GroupManager) RemoveContainerFromGroup(groupID, containerID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group, exists := gm.groups[groupID]
	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Find and remove container
	newIDs := []string{}
	for _, id := range group.ContainerIDs {
		if id != containerID {
			newIDs = append(newIDs, id)
		}
	}

	group.ContainerIDs = newIDs
	group.UpdatedAt = time.Now()
	return nil
}

// GetGroupHierarchy returns the full hierarchy tree for a group
func (gm *GroupManager) GetGroupHierarchy(groupID string) (*GroupHierarchy, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	group, exists := gm.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}

	hierarchy := &GroupHierarchy{
		Group:      group,
		Children:   []*GroupHierarchy{},
		Containers: group.ContainerIDs,
	}

	// Find child groups
	for _, g := range gm.groups {
		if g.ParentGroupID == groupID {
			childHierarchy, err := gm.buildHierarchy(g.ID)
			if err == nil {
				hierarchy.Children = append(hierarchy.Children, childHierarchy)
			}
		}
	}

	return hierarchy, nil
}

// buildHierarchy recursively builds hierarchy tree
func (gm *GroupManager) buildHierarchy(groupID string) (*GroupHierarchy, error) {
	group, exists := gm.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}

	hierarchy := &GroupHierarchy{
		Group:      group,
		Children:   []*GroupHierarchy{},
		Containers: group.ContainerIDs,
	}

	// Find child groups
	for _, g := range gm.groups {
		if g.ParentGroupID == groupID {
			childHierarchy, err := gm.buildHierarchy(g.ID)
			if err == nil {
				hierarchy.Children = append(hierarchy.Children, childHierarchy)
			}
		}
	}

	return hierarchy, nil
}

// GroupHierarchy represents the hierarchical structure of groups
type GroupHierarchy struct {
	Group      *ContainerGroup `json:"group"`
	Children   []*GroupHierarchy `json:"children"`
	Containers []string        `json:"containers"`
}

