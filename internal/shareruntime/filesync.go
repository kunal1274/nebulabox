package shareruntime

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// FileChangeType represents the type of file change
type FileChangeType string

const (
	FileChangeCreated FileChangeType = "created"
	FileChangeModified FileChangeType = "modified"
	FileChangeDeleted FileChangeType = "deleted"
	FileChangeRenamed  FileChangeType = "renamed"
)

// FileChange represents a change to a file in the shared filesystem
type FileChange struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId"`
	ContainerID string            `json:"containerId"`
	Path        string            `json:"path"`
	ChangeType  FileChangeType    `json:"changeType"`
	OldPath     string            `json:"oldPath,omitempty"` // For rename operations
	Hash        string            `json:"hash,omitempty"`     // SHA256 hash of file content
	Size        int64             `json:"size"`
	IsDirectory bool              `json:"isDirectory"`
	UserID      string            `json:"userId"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// FileSyncAdapter defines the interface for filesystem synchronization
type FileSyncAdapter interface {
	// RecordFileChange records a file change event
	RecordFileChange(change *FileChange) error
	
	// GetFileChanges retrieves file changes since a timestamp
	GetFileChanges(workspaceID, containerID string, since time.Time) ([]*FileChange, error)
	
	// GetFileHash calculates and returns the hash of a file
	GetFileHash(filePath string) (string, error)
	
	// SubscribeToFileChanges subscribes to file change events
	SubscribeToFileChanges(workspaceID, containerID string, callback func(*FileChange)) (string, error)
	
	// UnsubscribeFromFileChanges unsubscribes from file change events
	UnsubscribeFromFileChanges(subscriptionID string) error
	
	// SyncFile syncs a file (for rsync-like operations)
	SyncFile(workspaceID, containerID, filePath string, targetPath string) error
}

// InMemoryFileSyncAdapter is an in-memory implementation of FileSyncAdapter
type InMemoryFileSyncAdapter struct {
	changes          map[string]*FileChange // changeID -> FileChange
	changesByWorkspace map[string][]string   // workspaceID -> []changeID
	changesByContainer map[string][]string   // containerID -> []changeID
	subscribers      map[string]*FileSubscription // subscriptionID -> Subscription
	subscriptionsByWorkspace map[string][]string // workspaceID -> []subscriptionID
	mu               sync.RWMutex
	nextChangeID     int64
	nextSubID        int64
}

// FileSubscription represents an active file change subscription
type FileSubscription struct {
	ID          string
	WorkspaceID string
	ContainerID string
	Callback    func(*FileChange)
	CreatedAt   time.Time
}

// NewInMemoryFileSyncAdapter creates a new in-memory file sync adapter
func NewInMemoryFileSyncAdapter() *InMemoryFileSyncAdapter {
	return &InMemoryFileSyncAdapter{
		changes:                  make(map[string]*FileChange),
		changesByWorkspace:        make(map[string][]string),
		changesByContainer:        make(map[string][]string),
		subscribers:               make(map[string]*FileSubscription),
		subscriptionsByWorkspace: make(map[string][]string),
	}
}

// RecordFileChange records a file change event
func (a *InMemoryFileSyncAdapter) RecordFileChange(change *FileChange) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if change.ID == "" {
		change.ID = fmt.Sprintf("file-%d", a.nextChangeID)
		a.nextChangeID++
	}
	if change.Timestamp.IsZero() {
		change.Timestamp = time.Now()
	}

	a.changes[change.ID] = change
	a.changesByWorkspace[change.WorkspaceID] = append(a.changesByWorkspace[change.WorkspaceID], change.ID)
	a.changesByContainer[change.ContainerID] = append(a.changesByContainer[change.ContainerID], change.ID)

	// Notify subscribers
	if subs, exists := a.subscriptionsByWorkspace[change.WorkspaceID]; exists {
		for _, subID := range subs {
			if sub, exists := a.subscribers[subID]; exists {
				// Check if subscription is for this container or all containers
				if sub.ContainerID == "" || sub.ContainerID == change.ContainerID {
					go func(callback func(*FileChange), ch *FileChange) {
						defer func() {
							if r := recover(); r != nil {
								// Prevent panic from breaking the sync system
							}
						}()
						callback(ch)
					}(sub.Callback, change)
				}
			}
		}
	}

	return nil
}

// GetFileChanges retrieves file changes since a timestamp
func (a *InMemoryFileSyncAdapter) GetFileChanges(workspaceID, containerID string, since time.Time) ([]*FileChange, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var changeIDs []string
	if containerID != "" {
		changeIDs = a.changesByContainer[containerID]
	} else {
		changeIDs = a.changesByWorkspace[workspaceID]
	}

	var result []*FileChange
	for _, changeID := range changeIDs {
		if change, exists := a.changes[changeID]; exists {
			if change.WorkspaceID == workspaceID && (containerID == "" || change.ContainerID == containerID) {
				if change.Timestamp.After(since) || change.Timestamp.Equal(since) {
					result = append(result, change)
				}
			}
		}
	}

	return result, nil
}

// GetFileHash calculates and returns the SHA256 hash of a file
func (a *InMemoryFileSyncAdapter) GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SubscribeToFileChanges subscribes to file change events
func (a *InMemoryFileSyncAdapter) SubscribeToFileChanges(workspaceID, containerID string, callback func(*FileChange)) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	subID := fmt.Sprintf("filesub-%d", a.nextSubID)
	a.nextSubID++

	sub := &FileSubscription{
		ID:          subID,
		WorkspaceID: workspaceID,
		ContainerID: containerID,
		Callback:    callback,
		CreatedAt:   time.Now(),
	}

	a.subscribers[subID] = sub
	a.subscriptionsByWorkspace[workspaceID] = append(a.subscriptionsByWorkspace[workspaceID], subID)

	return subID, nil
}

// UnsubscribeFromFileChanges unsubscribes from file change events
func (a *InMemoryFileSyncAdapter) UnsubscribeFromFileChanges(subscriptionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	sub, exists := a.subscribers[subscriptionID]
	if !exists {
		return fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	delete(a.subscribers, subscriptionID)

	// Remove from workspace subscriptions
	subs := a.subscriptionsByWorkspace[sub.WorkspaceID]
	newSubs := []string{}
	for _, sid := range subs {
		if sid != subscriptionID {
			newSubs = append(newSubs, sid)
		}
	}
	a.subscriptionsByWorkspace[sub.WorkspaceID] = newSubs

	return nil
}

// SyncFile syncs a file (mock implementation for rsync-like operations)
func (a *InMemoryFileSyncAdapter) SyncFile(workspaceID, containerID, filePath string, targetPath string) error {
	// In a real implementation, this would:
	// 1. Read file from source
	// 2. Calculate hash
	// 3. Compare with target hash
	// 4. Transfer differences (rsync delta)
	// 5. Update target file
	
	// For now, just validate paths
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("source file not found: %w", err)
	}
	
	return nil
}

// FileSyncManager manages filesystem synchronization
type FileSyncManager struct {
	adapter FileSyncAdapter
	mu      sync.RWMutex
}

// NewFileSyncManager creates a new file sync manager
func NewFileSyncManager(adapter FileSyncAdapter) *FileSyncManager {
	return &FileSyncManager{
		adapter: adapter,
	}
}

// RecordFileChange records a file change
func (fsm *FileSyncManager) RecordFileChange(workspaceID, containerID, filePath string, changeType FileChangeType, userID string, isDirectory bool, size int64) error {
	// Calculate file hash if it's a file (not directory)
	var hash string
	var err error
	if !isDirectory {
		hash, err = fsm.adapter.GetFileHash(filePath)
		if err != nil {
			// File might not exist if it's a delete operation
			if changeType != FileChangeDeleted {
				return err
			}
		}
	}

	change := &FileChange{
		WorkspaceID: workspaceID,
		ContainerID: containerID,
		Path:        filePath,
		ChangeType:  changeType,
		Hash:         hash,
		Size:         size,
		IsDirectory: isDirectory,
		UserID:       userID,
		Timestamp:    time.Now(),
	}

	return fsm.adapter.RecordFileChange(change)
}

// GetFileChangesSince retrieves file changes since a timestamp
func (fsm *FileSyncManager) GetFileChangesSince(workspaceID, containerID string, since time.Time) ([]*FileChange, error) {
	return fsm.adapter.GetFileChanges(workspaceID, containerID, since)
}

// SubscribeToFileChanges subscribes to file change events
func (fsm *FileSyncManager) SubscribeToFileChanges(workspaceID, containerID string, callback func(*FileChange)) (string, error) {
	return fsm.adapter.SubscribeToFileChanges(workspaceID, containerID, callback)
}

// UnsubscribeFromFileChanges unsubscribes from file change events
func (fsm *FileSyncManager) UnsubscribeFromFileChanges(subscriptionID string) error {
	return fsm.adapter.UnsubscribeFromFileChanges(subscriptionID)
}

// SyncFile syncs a file using the adapter
func (fsm *FileSyncManager) SyncFile(workspaceID, containerID, filePath, targetPath string) error {
	return fsm.adapter.SyncFile(workspaceID, containerID, filePath, targetPath)
}

// GetFileHash gets the hash of a file
func (fsm *FileSyncManager) GetFileHash(filePath string) (string, error) {
	return fsm.adapter.GetFileHash(filePath)
}

