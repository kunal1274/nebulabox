package snapshot

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SnapshotState represents the state of a snapshot
type SnapshotState string

const (
	SnapshotStateCreating SnapshotState = "creating"
	SnapshotStateReady    SnapshotState = "ready"
	SnapshotStateFailed   SnapshotState = "failed"
	SnapshotStateRestoring SnapshotState = "restoring"
)

// SnapshotType represents what is being snapshotted
type SnapshotType string

const (
	SnapshotTypeContainer SnapshotType = "container"
	SnapshotTypeWorkspace SnapshotType = "workspace"
	SnapshotTypeVolume    SnapshotType = "volume"
)

// Snapshot represents a snapshot of a container or workspace
type Snapshot struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        SnapshotType           `json:"type"`
	ResourceID  string                 `json:"resourceId"` // Container ID, Workspace ID, etc.
	State       SnapshotState          `json:"state"`
	Size        int64                  `json:"size"` // Size in bytes
	CreatedAt   time.Time              `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	
	// Container-specific fields
	Image         string            `json:"image,omitempty"`
	Command       []string          `json:"command,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Ports         map[string]string `json:"ports,omitempty"`
	Volumes       []string          `json:"volumes,omitempty"`
	Network       string            `json:"network,omitempty"`
	Resources     *ResourceLimits  `json:"resources,omitempty"`
	
	// Filesystem snapshot reference
	FilesystemHash string `json:"filesystemHash,omitempty"`
	
	// Workspace-specific fields
	WorkspaceSettings map[string]interface{} `json:"workspaceSettings,omitempty"`
	Members           []string                `json:"members,omitempty"`
}

// ResourceLimits represents container resource limits
type ResourceLimits struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	IOPS   int64  `json:"iops,omitempty"`
	PIDs   int64  `json:"pids,omitempty"`
}

// SnapshotManager manages snapshots
type SnapshotManager struct {
	snapshots map[string]*Snapshot // snapshotID -> Snapshot
	byResource map[string][]string  // resourceID -> []snapshotID
	mu        sync.RWMutex
	nextID    int64
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager() *SnapshotManager {
	return &SnapshotManager{
		snapshots:  make(map[string]*Snapshot),
		byResource: make(map[string][]string),
	}
}

// CreateSnapshot creates a new snapshot
func (sm *SnapshotManager) CreateSnapshot(name, description, resourceID string, snapshotType SnapshotType, createdBy string, metadata map[string]interface{}) (*Snapshot, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	snapshotID := fmt.Sprintf("snap-%d", sm.nextID)
	sm.nextID++

	snapshot := &Snapshot{
		ID:          snapshotID,
		Name:        name,
		Description: description,
		Type:        snapshotType,
		ResourceID:  resourceID,
		State:       SnapshotStateCreating,
		Size:        0, // Will be calculated later
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		Metadata:    metadata,
	}

	sm.snapshots[snapshotID] = snapshot
	sm.byResource[resourceID] = append(sm.byResource[resourceID], snapshotID)

	return snapshot, nil
}

// UpdateSnapshot updates snapshot data
func (sm *SnapshotManager) UpdateSnapshot(snapshotID string, updates func(*Snapshot) error) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	snapshot, exists := sm.snapshots[snapshotID]
	if !exists {
		return fmt.Errorf("snapshot not found: %s", snapshotID)
	}

	return updates(snapshot)
}

// GetSnapshot retrieves a snapshot by ID
func (sm *SnapshotManager) GetSnapshot(snapshotID string) (*Snapshot, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshot, exists := sm.snapshots[snapshotID]
	if !exists {
		return nil, fmt.Errorf("snapshot not found: %s", snapshotID)
	}

	// Return a copy to prevent external modifications
	return copySnapshot(snapshot), nil
}

// ListSnapshots lists all snapshots, optionally filtered by resource ID or type
func (sm *SnapshotManager) ListSnapshots(resourceID string, snapshotType SnapshotType) ([]*Snapshot, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var result []*Snapshot
	for _, snapshot := range sm.snapshots {
		if resourceID != "" && snapshot.ResourceID != resourceID {
			continue
		}
		if snapshotType != "" && snapshot.Type != snapshotType {
			continue
		}
		result = append(result, copySnapshot(snapshot))
	}

	return result, nil
}

// ListResourceSnapshots lists all snapshots for a specific resource
func (sm *SnapshotManager) ListResourceSnapshots(resourceID string) ([]*Snapshot, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshotIDs, exists := sm.byResource[resourceID]
	if !exists {
		return []*Snapshot{}, nil
	}

	var result []*Snapshot
	for _, snapshotID := range snapshotIDs {
		if snapshot, exists := sm.snapshots[snapshotID]; exists {
			result = append(result, copySnapshot(snapshot))
		}
	}

	return result, nil
}

// DeleteSnapshot deletes a snapshot
func (sm *SnapshotManager) DeleteSnapshot(snapshotID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	snapshot, exists := sm.snapshots[snapshotID]
	if !exists {
		return fmt.Errorf("snapshot not found: %s", snapshotID)
	}

	delete(sm.snapshots, snapshotID)

	// Remove from resource index
	snapshotIDs := sm.byResource[snapshot.ResourceID]
	newIDs := []string{}
	for _, id := range snapshotIDs {
		if id != snapshotID {
			newIDs = append(newIDs, id)
		}
	}
	sm.byResource[snapshot.ResourceID] = newIDs

	return nil
}

// SetSnapshotState updates the state of a snapshot
func (sm *SnapshotManager) SetSnapshotState(snapshotID string, state SnapshotState) error {
	return sm.UpdateSnapshot(snapshotID, func(s *Snapshot) error {
		s.State = state
		return nil
	})
}

// SetSnapshotSize updates the size of a snapshot
func (sm *SnapshotManager) SetSnapshotSize(snapshotID string, size int64) error {
	return sm.UpdateSnapshot(snapshotID, func(s *Snapshot) error {
		s.Size = size
		return nil
	})
}

// SetSnapshotContainerData sets container-specific snapshot data
func (sm *SnapshotManager) SetSnapshotContainerData(snapshotID string, image string, command []string, env map[string]string, ports map[string]string, volumes []string, network string, resources *ResourceLimits) error {
	return sm.UpdateSnapshot(snapshotID, func(s *Snapshot) error {
		s.Image = image
		s.Command = command
		s.Env = env
		s.Ports = ports
		s.Volumes = volumes
		s.Network = network
		s.Resources = resources
		return nil
	})
}

// SetSnapshotFilesystemHash sets the filesystem hash for the snapshot
func (sm *SnapshotManager) SetSnapshotFilesystemHash(snapshotID string, hash string) error {
	return sm.UpdateSnapshot(snapshotID, func(s *Snapshot) error {
		s.FilesystemHash = hash
		return nil
	})
}

// copySnapshot creates a deep copy of a snapshot
func copySnapshot(s *Snapshot) *Snapshot {
	snap := &Snapshot{
		ID:              s.ID,
		Name:            s.Name,
		Description:     s.Description,
		Type:            s.Type,
		ResourceID:      s.ResourceID,
		State:           s.State,
		Size:            s.Size,
		CreatedAt:       s.CreatedAt,
		CreatedBy:       s.CreatedBy,
		Image:           s.Image,
		FilesystemHash:  s.FilesystemHash,
	}

	if s.Command != nil {
		snap.Command = make([]string, len(s.Command))
		copy(snap.Command, s.Command)
	}

	if s.Env != nil {
		snap.Env = make(map[string]string)
		for k, v := range s.Env {
			snap.Env[k] = v
		}
	}

	if s.Ports != nil {
		snap.Ports = make(map[string]string)
		for k, v := range s.Ports {
			snap.Ports[k] = v
		}
	}

	if s.Volumes != nil {
		snap.Volumes = make([]string, len(s.Volumes))
		copy(snap.Volumes, s.Volumes)
	}

	if s.Network != "" {
		snap.Network = s.Network
	}

	if s.Resources != nil {
		snap.Resources = &ResourceLimits{
			CPU:    s.Resources.CPU,
			Memory: s.Resources.Memory,
			IOPS:   s.Resources.IOPS,
			PIDs:   s.Resources.PIDs,
		}
	}

	if s.Metadata != nil {
		// Deep copy metadata
		metadataJSON, _ := json.Marshal(s.Metadata)
		json.Unmarshal(metadataJSON, &snap.Metadata)
	}

	if s.WorkspaceSettings != nil {
		settingsJSON, _ := json.Marshal(s.WorkspaceSettings)
		json.Unmarshal(settingsJSON, &snap.WorkspaceSettings)
	}

	if s.Members != nil {
		snap.Members = make([]string, len(s.Members))
		copy(snap.Members, s.Members)
	}

	return snap
}

