package api

import (
	"github.com/nebulabox/nebulabox/internal/shareruntime"
	"github.com/nebulabox/nebulabox/internal/snapshot"
)

// SnapshotManagerAdapter adapts snapshot.SnapshotManager to shareruntime.SnapshotManagerInterface
type SnapshotManagerAdapter struct {
	snapshotManager *snapshot.SnapshotManager
}

// CreateSnapshot creates a snapshot
func (sma *SnapshotManagerAdapter) CreateSnapshot(resourceType string, resourceID string, name string, description string) (string, error) {
	var snapshotType snapshot.SnapshotType
	switch resourceType {
	case "container":
		snapshotType = snapshot.SnapshotTypeContainer
	case "workspace":
		snapshotType = snapshot.SnapshotTypeWorkspace
	case "volume":
		snapshotType = snapshot.SnapshotTypeVolume
	default:
		snapshotType = snapshot.SnapshotTypeContainer
	}
	
	snap, err := sma.snapshotManager.CreateSnapshot(name, description, resourceID, snapshotType, "system", nil)
	if err != nil {
		return "", err
	}
	return snap.ID, nil
}

// GetSnapshot gets a snapshot
func (sma *SnapshotManagerAdapter) GetSnapshot(id string) (*shareruntime.SnapshotInfo, error) {
	snap, err := sma.snapshotManager.GetSnapshot(id)
	if err != nil {
		return nil, err
	}
	return &shareruntime.SnapshotInfo{
		ID:          snap.ID,
		Name:        snap.Name,
		ResourceType: string(snap.Type),
		ResourceID:  snap.ResourceID,
	}, nil
}

// RestoreSnapshot restores a snapshot
func (sma *SnapshotManagerAdapter) RestoreSnapshot(id string) error {
	// The SnapshotManager doesn't have RestoreSnapshot, it's handled by the API
	// For now, we'll just verify the snapshot exists
	_, err := sma.snapshotManager.GetSnapshot(id)
	return err
}

