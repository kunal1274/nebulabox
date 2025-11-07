package orchestrator

import (
	"fmt"
	"sync"
	"time"
)

// Node represents a NebulaBox node in the cluster
type Node struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Address       string            `json:"address"` // IP or hostname
	Port          int               `json:"port"`
	Status        string            `json:"status"` // online, offline, maintenance
	LastSeen      time.Time         `json:"lastSeen"`
	Labels        map[string]string `json:"labels,omitempty"` // For node selection
	Resources     NodeResources     `json:"resources"`
	Containers    []string          `json:"containers"` // Container IDs on this node
	ContainerCount int              `json:"containerCount"`
	Region        string            `json:"region,omitempty"`
	Zone          string            `json:"zone,omitempty"`
}

// NodeResources represents node resource capacity and usage
type NodeResources struct {
	CPUCores     int     `json:"cpuCores"`
	MemoryMB     int     `json:"memoryMB"`
	DiskGB       int     `json:"diskGB"`
	CPUUsed      float64 `json:"cpuUsed"`      // Percentage
	MemoryUsed   float64 `json:"memoryUsed"`   // Percentage
	DiskUsed     float64 `json:"diskUsed"`     // Percentage
	ContainersRunning int `json:"containersRunning"`
}

// NodeManager manages cluster nodes
type NodeManager struct {
	nodes map[string]*Node
	mu    sync.RWMutex
}

// NewNodeManager creates a new node manager
func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: make(map[string]*Node),
	}
}

// RegisterNode registers a new node or updates an existing one
func (nm *NodeManager) RegisterNode(node *Node) error {
	if node.ID == "" {
		return fmt.Errorf("node ID is required")
	}
	if node.Address == "" {
		return fmt.Errorf("node address is required")
	}

	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Set defaults
	if node.Name == "" {
		node.Name = node.ID
	}
	if node.Port == 0 {
		node.Port = 8081 // Default API port
	}
	if node.Status == "" {
		node.Status = "online"
	}
	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}

	node.LastSeen = time.Now()
	nm.nodes[node.ID] = node

	return nil
}

// GetNode retrieves a node by ID
func (nm *NodeManager) GetNode(id string) (*Node, bool) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	node, exists := nm.nodes[id]
	return node, exists
}

// ListNodes returns all registered nodes
func (nm *NodeManager) ListNodes() []*Node {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	nodes := make([]*Node, 0, len(nm.nodes))
	for _, node := range nm.nodes {
		// Create a copy to avoid race conditions
		nodeCopy := *node
		nodes = append(nodes, &nodeCopy)
	}
	return nodes
}

// ListOnlineNodes returns only online nodes
func (nm *NodeManager) ListOnlineNodes() []*Node {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	nodes := make([]*Node, 0)
	for _, node := range nm.nodes {
		if node.Status == "online" {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}
	return nodes
}

// UpdateNodeStatus updates node status
func (nm *NodeManager) UpdateNodeStatus(id, status string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.nodes[id]
	if !exists {
		return fmt.Errorf("node not found: %s", id)
	}

	validStatuses := map[string]bool{
		"online":      true,
		"offline":     true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	node.Status = status
	node.LastSeen = time.Now()
	return nil
}

// UpdateNodeResources updates node resource usage
func (nm *NodeManager) UpdateNodeResources(id string, resources NodeResources) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.nodes[id]
	if !exists {
		return fmt.Errorf("node not found: %s", id)
	}

	node.Resources = resources
	node.Resources.ContainersRunning = resources.ContainersRunning
	node.ContainerCount = resources.ContainersRunning
	node.LastSeen = time.Now()
	return nil
}

// AddContainerToNode adds a container to a node
func (nm *NodeManager) AddContainerToNode(nodeID, containerID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	// Check if container already exists
	for _, cid := range node.Containers {
		if cid == containerID {
			return nil // Already added
		}
	}

	node.Containers = append(node.Containers, containerID)
	node.ContainerCount = len(node.Containers)
	node.Resources.ContainersRunning = node.ContainerCount
	return nil
}

// RemoveContainerFromNode removes a container from a node
func (nm *NodeManager) RemoveContainerFromNode(nodeID, containerID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	// Find and remove container
	for i, cid := range node.Containers {
		if cid == containerID {
			node.Containers = append(node.Containers[:i], node.Containers[i+1:]...)
			node.ContainerCount = len(node.Containers)
			node.Resources.ContainersRunning = node.ContainerCount
			return nil
		}
	}

	return fmt.Errorf("container not found on node: %s", containerID)
}

// SelectNode selects the best node for a container based on criteria
func (nm *NodeManager) SelectNode(criteria NodeSelectionCriteria) (*Node, error) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	onlineNodes := make([]*Node, 0)
	for _, node := range nm.nodes {
		if node.Status == "online" {
			onlineNodes = append(onlineNodes, node)
		}
	}

	if len(onlineNodes) == 0 {
		return nil, fmt.Errorf("no online nodes available")
	}

	// Simple selection: return first online node
	// TODO: Implement more sophisticated selection (resource availability, labels, etc.)
	selected := onlineNodes[0]

	// Filter by labels if specified
	if len(criteria.RequiredLabels) > 0 {
		filtered := make([]*Node, 0)
		for _, node := range onlineNodes {
			match := true
			for key, value := range criteria.RequiredLabels {
				if node.Labels[key] != value {
					match = false
					break
				}
			}
			if match {
				filtered = append(filtered, node)
			}
		}
		if len(filtered) > 0 {
			selected = filtered[0]
		} else {
			return nil, fmt.Errorf("no nodes match required labels")
		}
	}

	// Filter by resource requirements if specified
	if criteria.MinMemoryMB > 0 || criteria.MinCPU > 0 {
		filtered := make([]*Node, 0)
		for _, node := range onlineNodes {
			availableCPU := float64(node.Resources.CPUCores) * (1.0 - node.Resources.CPUUsed/100.0)
			availableMemory := float64(node.Resources.MemoryMB) * (1.0 - node.Resources.MemoryUsed/100.0)

			if criteria.MinCPU > 0 && availableCPU < float64(criteria.MinCPU) {
				continue
			}
			if criteria.MinMemoryMB > 0 && availableMemory < float64(criteria.MinMemoryMB) {
				continue
			}
			filtered = append(filtered, node)
		}
		if len(filtered) > 0 {
			selected = filtered[0]
		} else {
			return nil, fmt.Errorf("no nodes have sufficient resources")
		}
	}

	return selected, nil
}

// NodeSelectionCriteria defines criteria for node selection
type NodeSelectionCriteria struct {
	RequiredLabels map[string]string
	MinCPU         int
	MinMemoryMB    int
	PreferRegion   string
	PreferZone     string
}

// DeleteNode removes a node from the cluster
func (nm *NodeManager) DeleteNode(id string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if _, exists := nm.nodes[id]; !exists {
		return fmt.Errorf("node not found: %s", id)
	}

	if len(nm.nodes[id].Containers) > 0 {
		return fmt.Errorf("cannot delete node with containers: %s", id)
	}

	delete(nm.nodes, id)
	return nil
}

// CleanupStaleNodes removes nodes that haven't been seen in a while
func (nm *NodeManager) CleanupStaleNodes(timeout time.Duration) []string {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	stale := make([]string, 0)
	now := time.Now()

	for id, node := range nm.nodes {
		if node.Status == "offline" || now.Sub(node.LastSeen) > timeout {
			if len(node.Containers) == 0 {
				delete(nm.nodes, id)
				stale = append(stale, id)
			} else {
				// Mark as offline but keep in registry
				node.Status = "offline"
			}
		}
	}

	return stale
}

