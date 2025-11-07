package groups

import (
	"fmt"
	"sync"
	"time"
)

// RelationshipManager manages container relationships
type RelationshipManager struct {
	relationships map[string][]*ContainerRelationship // parentID -> relationships
	childMap      map[string]*ContainerRelationship  // childID -> relationship
	mu            sync.RWMutex
}

// NewRelationshipManager creates a new relationship manager
func NewRelationshipManager() *RelationshipManager {
	return &RelationshipManager{
		relationships: make(map[string][]*ContainerRelationship),
		childMap:      make(map[string]*ContainerRelationship),
	}
}

// CreateRelationship creates a parent-child relationship between containers
func (rm *RelationshipManager) CreateRelationship(parentID, childID, relType string, properties map[string]interface{}) (*ContainerRelationship, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Check if child already has a parent
	if existing, exists := rm.childMap[childID]; exists {
		return nil, fmt.Errorf("container %s already has a parent relationship with %s", childID, existing.ParentID)
	}

	// Prevent circular relationships
	if rm.hasCircularRelationship(parentID, childID) {
		return nil, fmt.Errorf("circular relationship detected: cannot create relationship")
	}

	relationship := &ContainerRelationship{
		ParentID:   parentID,
		ChildID:    childID,
		Type:       relType,
		Properties: properties,
		CreatedAt:  time.Now(),
	}

	// Add to relationships map
	rm.relationships[parentID] = append(rm.relationships[parentID], relationship)
	rm.childMap[childID] = relationship

	return relationship, nil
}

// hasCircularRelationship checks if creating a relationship would create a cycle
func (rm *RelationshipManager) hasCircularRelationship(parentID, potentialChildID string) bool {
	// If potentialChildID is already a parent of parentID, it's circular
	visited := make(map[string]bool)
	return rm.isAncestor(potentialChildID, parentID, visited)
}

// isAncestor checks if ancestorID is an ancestor of descendantID
func (rm *RelationshipManager) isAncestor(ancestorID, descendantID string, visited map[string]bool) bool {
	if ancestorID == descendantID {
		return true
	}

	if visited[ancestorID] {
		return false // Already checked this path
	}
	visited[ancestorID] = true

	// Check all children of ancestorID
	for _, rel := range rm.relationships[ancestorID] {
		if rm.isAncestor(rel.ChildID, descendantID, visited) {
			return true
		}
	}

	return false
}

// GetChildren returns all children of a container
func (rm *RelationshipManager) GetChildren(containerID string) []*ContainerRelationship {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.relationships[containerID]
}

// GetParent returns the parent of a container
func (rm *RelationshipManager) GetParent(containerID string) (*ContainerRelationship, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	rel, exists := rm.childMap[containerID]
	return rel, exists
}

// DeleteRelationship removes a relationship
func (rm *RelationshipManager) DeleteRelationship(parentID, childID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Find and remove from relationships map
	rels := rm.relationships[parentID]
	newRels := []*ContainerRelationship{}
	found := false
	for _, rel := range rels {
		if rel.ChildID != childID {
			newRels = append(newRels, rel)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("relationship not found")
	}

	if len(newRels) == 0 {
		delete(rm.relationships, parentID)
	} else {
		rm.relationships[parentID] = newRels
	}

	// Remove from child map
	delete(rm.childMap, childID)

	return nil
}

// GetAncestry returns the full ancestry chain for a container
func (rm *RelationshipManager) GetAncestry(containerID string) []*ContainerRelationship {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var ancestry []*ContainerRelationship
	current := containerID

	for {
		rel, exists := rm.childMap[current]
		if !exists {
			break
		}
		ancestry = append(ancestry, rel)
		current = rel.ParentID
	}

	return ancestry
}

// GetDescendants returns all descendants of a container
func (rm *RelationshipManager) GetDescendants(containerID string) []*ContainerRelationship {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var descendants []*ContainerRelationship
	rm.collectDescendants(containerID, &descendants)
	return descendants
}

// collectDescendants recursively collects all descendants
func (rm *RelationshipManager) collectDescendants(containerID string, result *[]*ContainerRelationship) {
	for _, rel := range rm.relationships[containerID] {
		*result = append(*result, rel)
		rm.collectDescendants(rel.ChildID, result)
	}
}

// GetContainerTree returns the full tree structure for a container
type ContainerTree struct {
	ContainerID string          `json:"containerId"`
	Children    []*ContainerTree `json:"children"`
	ParentID    string          `json:"parentId,omitempty"`
	Type        string          `json:"type,omitempty"`
}

// GetContainerTree returns the tree structure starting from a container
func (rm *RelationshipManager) GetContainerTree(containerID string) *ContainerTree {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	tree := &ContainerTree{
		ContainerID: containerID,
		Children:    []*ContainerTree{},
	}

	// Get parent info if exists
	if parentRel, exists := rm.childMap[containerID]; exists {
		tree.ParentID = parentRel.ParentID
		tree.Type = parentRel.Type
	}

	// Build children trees
	for _, rel := range rm.relationships[containerID] {
		childTree := rm.buildTree(rel.ChildID)
		tree.Children = append(tree.Children, childTree)
	}

	return tree
}

// buildTree recursively builds the tree structure
func (rm *RelationshipManager) buildTree(containerID string) *ContainerTree {
	tree := &ContainerTree{
		ContainerID: containerID,
		Children:    []*ContainerTree{},
	}

	// Get parent info
	if parentRel, exists := rm.childMap[containerID]; exists {
		tree.ParentID = parentRel.ParentID
		tree.Type = parentRel.Type
	}

	// Build children trees
	for _, rel := range rm.relationships[containerID] {
		childTree := rm.buildTree(rel.ChildID)
		tree.Children = append(tree.Children, childTree)
	}

	return tree
}

