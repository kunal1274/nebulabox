package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// AutoSleepConfig represents configuration for workspace auto-sleep
type AutoSleepConfig struct {
	Enabled          bool          `json:"enabled"`
	IdleTimeout      time.Duration `json:"idleTimeout"`      // Time before workspace is considered idle
	SleepTimeout     time.Duration `json:"sleepTimeout"`     // Time after idle before sleep
	CreateSnapshot   bool          `json:"createSnapshot"`   // Whether to create snapshot before sleep
	AutoWakeOnAccess bool          `json:"autoWakeOnAccess"` // Auto-wake when user accesses
}

// DefaultAutoSleepConfig returns default auto-sleep configuration
func DefaultAutoSleepConfig() AutoSleepConfig {
	return AutoSleepConfig{
		Enabled:          true,
		IdleTimeout:      30 * time.Minute,
		SleepTimeout:     15 * time.Minute, // Sleep after 15 minutes of idle
		CreateSnapshot:   true,
		AutoWakeOnAccess: true,
	}
}

// WorkspaceActivityTracker tracks activity for workspaces
type WorkspaceActivityTracker struct {
	activities map[string]time.Time // workspaceID -> last activity time
	mu         sync.RWMutex
}

// NewWorkspaceActivityTracker creates a new activity tracker
func NewWorkspaceActivityTracker() *WorkspaceActivityTracker {
	return &WorkspaceActivityTracker{
		activities: make(map[string]time.Time),
	}
}

// UpdateActivity updates the last activity time for a workspace
func (wat *WorkspaceActivityTracker) UpdateActivity(workspaceID string) {
	wat.mu.Lock()
	defer wat.mu.Unlock()
	wat.activities[workspaceID] = time.Now()
}

// GetLastActivity gets the last activity time for a workspace
func (wat *WorkspaceActivityTracker) GetLastActivity(workspaceID string) (time.Time, bool) {
	wat.mu.RLock()
	defer wat.mu.RUnlock()
	activity, exists := wat.activities[workspaceID]
	return activity, exists
}

// GetIdleDuration returns how long a workspace has been idle
func (wat *WorkspaceActivityTracker) GetIdleDuration(workspaceID string) time.Duration {
	wat.mu.RLock()
	defer wat.mu.RUnlock()
	lastActivity, exists := wat.activities[workspaceID]
	if !exists {
		return 0
	}
	return time.Since(lastActivity)
}

// ClearActivity clears activity for a workspace
func (wat *WorkspaceActivityTracker) ClearActivity(workspaceID string) {
	wat.mu.Lock()
	defer wat.mu.Unlock()
	delete(wat.activities, workspaceID)
}

// SnapshotManagerInterface defines the interface for snapshot operations
type SnapshotManagerInterface interface {
	CreateSnapshot(resourceType string, resourceID string, name string, description string) (string, error)
	GetSnapshot(id string) (*SnapshotInfo, error)
	RestoreSnapshot(id string) error
}

// SnapshotInfo represents snapshot information
type SnapshotInfo struct {
	ID          string
	Name        string
	ResourceType string
	ResourceID  string
}

// AutoSleepManager manages automatic sleep and wake for workspaces
type AutoSleepManager struct {
	workspaceManager    *WorkspaceManager
	snapshotManager     SnapshotManagerInterface
	activityTracker     *WorkspaceActivityTracker
	configs             map[string]AutoSleepConfig // workspaceID -> config
	mu                  sync.RWMutex
	ticker              *time.Ticker
	stop                chan struct{}
}

// NewAutoSleepManager creates a new auto-sleep manager
func NewAutoSleepManager(workspaceManager *WorkspaceManager, snapshotManager SnapshotManagerInterface) *AutoSleepManager {
	return &AutoSleepManager{
		workspaceManager: workspaceManager,
		snapshotManager:  snapshotManager,
		activityTracker:  NewWorkspaceActivityTracker(),
		configs:          make(map[string]AutoSleepConfig),
		stop:             make(chan struct{}),
	}
}

// SetConfig sets auto-sleep configuration for a workspace
func (asm *AutoSleepManager) SetConfig(workspaceID string, config AutoSleepConfig) {
	asm.mu.Lock()
	defer asm.mu.Unlock()
	asm.configs[workspaceID] = config
}

// GetConfig gets auto-sleep configuration for a workspace
func (asm *AutoSleepManager) GetConfig(workspaceID string) AutoSleepConfig {
	asm.mu.RLock()
	defer asm.mu.RUnlock()
	
	config, exists := asm.configs[workspaceID]
	if !exists {
		return DefaultAutoSleepConfig()
	}
	return config
}

// Start starts the auto-sleep monitoring loop
func (asm *AutoSleepManager) Start() {
	asm.ticker = time.NewTicker(5 * time.Minute) // Check every 5 minutes
	
	go func() {
		for {
			select {
			case <-asm.ticker.C:
				asm.checkWorkspaces()
			case <-asm.stop:
				return
			}
		}
	}()
}

// Stop stops the auto-sleep monitoring loop
func (asm *AutoSleepManager) Stop() {
	if asm.ticker != nil {
		asm.ticker.Stop()
	}
	close(asm.stop)
}

// RecordActivity records activity for a workspace
func (asm *AutoSleepManager) RecordActivity(workspaceID string) {
	asm.activityTracker.UpdateActivity(workspaceID)
	
	// Check if workspace is sleeping and should be woken
	config := asm.GetConfig(workspaceID)
	if config.AutoWakeOnAccess {
		workspace, exists := asm.workspaceManager.GetWorkspace(workspaceID)
		if exists && workspace.Status == "sleeping" {
			go asm.WakeWorkspace(workspaceID)
		}
	}
}

// checkWorkspaces checks all workspaces for idle/sleep conditions
func (asm *AutoSleepManager) checkWorkspaces() {
	workspaces := asm.workspaceManager.ListWorkspaces("")

	for _, workspace := range workspaces {
		config := asm.GetConfig(workspace.ID)
		if !config.Enabled {
			continue
		}

		// Skip workspaces that are already sleeping or terminating
		if workspace.Status == "sleeping" || workspace.Status == "terminating" {
			continue
		}

		idleDuration := asm.activityTracker.GetIdleDuration(workspace.ID)
		
		// Check if workspace should be put to sleep
		if idleDuration >= (config.IdleTimeout + config.SleepTimeout) {
			go asm.sleepWorkspace(workspace.ID, config)
		}
	}
}

// sleepWorkspace puts a workspace to sleep with optional snapshot
func (asm *AutoSleepManager) sleepWorkspace(workspaceID string, config AutoSleepConfig) {
	asm.mu.Lock()
	defer asm.mu.Unlock()

	workspace, exists := asm.workspaceManager.GetWorkspace(workspaceID)
	if !exists {
		return
	}

	// Double-check status (may have changed)
	if workspace.Status == "sleeping" || workspace.Status == "terminating" {
		return
	}

	var snapshotID string
	if config.CreateSnapshot {
		// Create snapshot before sleep
		name := fmt.Sprintf("auto-sleep-%s-%d", workspaceID, time.Now().Unix())
		description := fmt.Sprintf("Auto-generated snapshot before sleep at %s", time.Now().Format(time.RFC3339))
		
		id, err := asm.snapshotManager.CreateSnapshot("workspace", workspaceID, name, description)
		if err == nil {
			snapshotID = id
		}
	}

	// Store snapshot info in metadata and update status
	if workspace.Metadata == nil {
		workspace.Metadata = make(map[string]interface{})
	}
	workspace.Metadata["snapshotId"] = snapshotID
	workspace.Metadata["sleptAt"] = time.Now().Unix()
	
	// Update workspace status to sleeping
	asm.workspaceManager.UpdateWorkspaceStatus(workspaceID, "sleeping")
}

// WakeWorkspace wakes a workspace from sleep, optionally restoring from snapshot
func (asm *AutoSleepManager) WakeWorkspace(workspaceID string) error {
	asm.mu.Lock()
	defer asm.mu.Unlock()

	workspace, exists := asm.workspaceManager.GetWorkspace(workspaceID)
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	if workspace.Status != "sleeping" {
		return nil // Already awake or in different state
	}

	// Get snapshot ID if available
	var snapshotID string
	if workspace.Metadata != nil {
		if id, ok := workspace.Metadata["snapshotId"].(string); ok {
			snapshotID = id
		}
	}

	// Restore from snapshot if available
	if snapshotID != "" {
		err := asm.snapshotManager.RestoreSnapshot(snapshotID)
		if err != nil {
			// Log error but continue with wake
			fmt.Printf("Warning: failed to restore snapshot %s for workspace %s: %v\n", snapshotID, workspaceID, err)
		}
	}

	// Update workspace status to active and clear sleep metadata
	workspace.Status = "active"
	if workspace.Metadata != nil {
		delete(workspace.Metadata, "snapshotId")
		delete(workspace.Metadata, "sleptAt")
	}

	// Update activity
	asm.activityTracker.UpdateActivity(workspaceID)

	return nil
}

// GetIdleWorkspaces returns list of workspaces that are idle
func (asm *AutoSleepManager) GetIdleWorkspaces() []WorkspaceActivity {
	asm.mu.RLock()
	defer asm.mu.RUnlock()

	var result []WorkspaceActivity
	workspaces := asm.workspaceManager.ListWorkspaces("")

	for _, workspace := range workspaces {
		config := asm.GetConfig(workspace.ID)
		if !config.Enabled {
			continue
		}

		if workspace.Status == "sleeping" || workspace.Status == "terminating" {
			continue
		}

		idleDuration := asm.activityTracker.GetIdleDuration(workspace.ID)
		if idleDuration >= config.IdleTimeout {
			lastActivity, _ := asm.activityTracker.GetLastActivity(workspace.ID)
			result = append(result, WorkspaceActivity{
				WorkspaceID: workspace.ID,
				WorkspaceName: workspace.Name,
				LastActivity: lastActivity,
				IdleDuration: idleDuration,
				Status: workspace.Status,
			})
		}
	}

	return result
}

// WorkspaceActivity represents activity information for a workspace
type WorkspaceActivity struct {
	WorkspaceID   string        `json:"workspaceId"`
	WorkspaceName string        `json:"workspaceName"`
	LastActivity  time.Time     `json:"lastActivity"`
	IdleDuration  time.Duration `json:"idleDuration"`
	Status        string        `json:"status"`
}

