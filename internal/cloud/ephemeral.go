package cloud

import (
	"fmt"
	"sync"
	"time"
)

// EphemeralRuntimeStatus represents the status of an ephemeral runtime
type EphemeralRuntimeStatus string

const (
	EphemeralStatusProvisioning EphemeralRuntimeStatus = "provisioning"
	EphemeralStatusActive        EphemeralRuntimeStatus = "active"
	EphemeralStatusIdle         EphemeralRuntimeStatus = "idle"
	EphemeralStatusSleeping     EphemeralRuntimeStatus = "sleeping"
	EphemeralStatusTerminating   EphemeralRuntimeStatus = "terminating"
	EphemeralStatusTerminated   EphemeralRuntimeStatus = "terminated"
)

// EphemeralRuntime represents a temporary remote runtime for team testing
type EphemeralRuntime struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	WorkspaceID     string                 `json:"workspaceId"`
	Region          string                 `json:"region"`
	Status          EphemeralRuntimeStatus `json:"status"`
	InstanceType    string                 `json:"instanceType"`
	Image           string                 `json:"image"`
	CreatedAt       time.Time              `json:"createdAt"`
	ExpiresAt       time.Time              `json:"expiresAt"`
	LastActivityAt  time.Time              `json:"lastActivityAt"`
	CreatedBy       string                 `json:"createdBy"`
	Members         []string               `json:"members"` // User IDs with access
	AccessURL       string                 `json:"accessUrl"`
	SSHKey          string                 `json:"sshKey,omitempty"`
	Resources       *EphemeralResources    `json:"resources,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	SnapshotID      string                 `json:"snapshotId,omitempty"` // Snapshot before sleep
}

// EphemeralResources represents resource allocation for ephemeral runtime
type EphemeralResources struct {
	CPU        int    `json:"cpu"`
	Memory     string `json:"memory"`
	Disk       string `json:"disk"`
	Network    string `json:"network,omitempty"`
	PublicIP   string `json:"publicIp,omitempty"`
	PrivateIP  string `json:"privateIp,omitempty"`
}

// EphemeralRuntimeManager manages ephemeral runtime lifecycle
type EphemeralRuntimeManager struct {
	runtimes map[string]*EphemeralRuntime // runtimeID -> EphemeralRuntime
	byWorkspace map[string][]string        // workspaceID -> []runtimeID
	mu        sync.RWMutex
	nextID    int64
	idleTimeout time.Duration
}

// NewEphemeralRuntimeManager creates a new ephemeral runtime manager
func NewEphemeralRuntimeManager(idleTimeout time.Duration) *EphemeralRuntimeManager {
	if idleTimeout == 0 {
		idleTimeout = 30 * time.Minute // Default: 30 minutes
	}
	return &EphemeralRuntimeManager{
		runtimes:    make(map[string]*EphemeralRuntime),
		byWorkspace: make(map[string][]string),
		idleTimeout: idleTimeout,
	}
}

// ProvisionRuntime provisions a new ephemeral runtime
func (erm *EphemeralRuntimeManager) ProvisionRuntime(name, workspaceID, region, instanceType, image, createdBy string, members []string, duration time.Duration) (*EphemeralRuntime, error) {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtimeID := fmt.Sprintf("ephemeral-%d", erm.nextID)
	erm.nextID++

	now := time.Now()
	expiresAt := now.Add(duration)
	if duration == 0 {
		expiresAt = now.Add(24 * time.Hour) // Default: 24 hours
	}

	runtime := &EphemeralRuntime{
		ID:            runtimeID,
		Name:          name,
		WorkspaceID:   workspaceID,
		Region:        region,
		Status:        EphemeralStatusProvisioning,
		InstanceType:  instanceType,
		Image:         image,
		CreatedAt:     now,
		ExpiresAt:     expiresAt,
		LastActivityAt: now,
		CreatedBy:     createdBy,
		Members:       members,
		Metadata:      make(map[string]interface{}),
	}

	// Initialize resources based on instance type
	runtime.Resources = erm.getResourcesForInstanceType(instanceType)
	runtime.AccessURL = fmt.Sprintf("https://%s.nebula.cloud", runtimeID)

	erm.runtimes[runtimeID] = runtime
	erm.byWorkspace[workspaceID] = append(erm.byWorkspace[workspaceID], runtimeID)

	// Simulate provisioning (in real implementation, would provision actual infrastructure)
	go erm.completeProvisioning(runtimeID)

	return runtime, nil
}

// completeProvisioning completes the provisioning process
func (erm *EphemeralRuntimeManager) completeProvisioning(runtimeID string) {
	time.Sleep(2 * time.Second) // Simulate provisioning time

	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if exists && runtime.Status == EphemeralStatusProvisioning {
		runtime.Status = EphemeralStatusActive
		runtime.LastActivityAt = time.Now()
	}
}

// GetRuntime retrieves an ephemeral runtime by ID
func (erm *EphemeralRuntimeManager) GetRuntime(runtimeID string) (*EphemeralRuntime, error) {
	erm.mu.RLock()
	defer erm.mu.RUnlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return nil, fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	return copyEphemeralRuntime(runtime), nil
}

// ListRuntimes lists all ephemeral runtimes, optionally filtered by workspace
func (erm *EphemeralRuntimeManager) ListRuntimes(workspaceID string) ([]*EphemeralRuntime, error) {
	erm.mu.RLock()
	defer erm.mu.RUnlock()

	var result []*EphemeralRuntime
	for _, runtime := range erm.runtimes {
		if workspaceID == "" || runtime.WorkspaceID == workspaceID {
			result = append(result, copyEphemeralRuntime(runtime))
		}
	}

	return result, nil
}

// UpdateActivity updates the last activity timestamp for a runtime
func (erm *EphemeralRuntimeManager) UpdateActivity(runtimeID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	runtime.LastActivityAt = time.Now()

	// Wake from sleep if idle
	if runtime.Status == EphemeralStatusSleeping {
		runtime.Status = EphemeralStatusActive
		runtime.SnapshotID = "" // Clear snapshot reference
	}

	return nil
}

// SleepRuntime puts a runtime to sleep and creates a snapshot
func (erm *EphemeralRuntimeManager) SleepRuntime(runtimeID string, snapshotID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	if runtime.Status != EphemeralStatusActive && runtime.Status != EphemeralStatusIdle {
		return fmt.Errorf("runtime is not in a sleepable state: %s", runtime.Status)
	}

	runtime.Status = EphemeralStatusSleeping
	runtime.SnapshotID = snapshotID
	runtime.Metadata["sleptAt"] = time.Now().Unix()

	return nil
}

// WakeRuntime wakes a runtime from sleep
func (erm *EphemeralRuntimeManager) WakeRuntime(runtimeID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	if runtime.Status != EphemeralStatusSleeping {
		return fmt.Errorf("runtime is not sleeping: %s", runtime.Status)
	}

	runtime.Status = EphemeralStatusActive
	runtime.LastActivityAt = time.Now()
	runtime.SnapshotID = "" // Clear snapshot reference

	return nil
}

// TerminateRuntime terminates an ephemeral runtime
func (erm *EphemeralRuntimeManager) TerminateRuntime(runtimeID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	runtime.Status = EphemeralStatusTerminating

	// Remove from workspace index
	runtimeIDs := erm.byWorkspace[runtime.WorkspaceID]
	newIDs := []string{}
	for _, id := range runtimeIDs {
		if id != runtimeID {
			newIDs = append(newIDs, id)
		}
	}
	erm.byWorkspace[runtime.WorkspaceID] = newIDs

	// Simulate termination (in real implementation, would deprovision infrastructure)
	go erm.completeTermination(runtimeID)

	return nil
}

// completeTermination completes the termination process
func (erm *EphemeralRuntimeManager) completeTermination(runtimeID string) {
	time.Sleep(1 * time.Second) // Simulate termination time

	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if exists {
		runtime.Status = EphemeralStatusTerminated
	}
}

// CheckIdleRuntimes checks for idle runtimes and automatically sleeps them
func (erm *EphemeralRuntimeManager) CheckIdleRuntimes() []string {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	var sleptRuntimes []string
	now := time.Now()

	for _, runtime := range erm.runtimes {
		if runtime.Status == EphemeralStatusActive {
			idleDuration := now.Sub(runtime.LastActivityAt)
			if idleDuration > erm.idleTimeout {
				runtime.Status = EphemeralStatusIdle
				runtime.Metadata["idleSince"] = now.Unix()
			}
		}

		// Auto-sleep idle runtimes after extended idle time
		if runtime.Status == EphemeralStatusIdle {
			idleSince, ok := runtime.Metadata["idleSince"].(int64)
			if ok && now.Sub(time.Unix(idleSince, 0)) > 15*time.Minute {
				runtime.Status = EphemeralStatusSleeping
				runtime.SnapshotID = fmt.Sprintf("auto-snapshot-%s", runtime.ID)
				sleptRuntimes = append(sleptRuntimes, runtime.ID)
			}
		}
	}

	return sleptRuntimes
}

// CheckExpiredRuntimes checks for expired runtimes and terminates them
func (erm *EphemeralRuntimeManager) CheckExpiredRuntimes() []string {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	var expiredRuntimes []string
	now := time.Now()

	for id, runtime := range erm.runtimes {
		if runtime.ExpiresAt.Before(now) && runtime.Status != EphemeralStatusTerminated && runtime.Status != EphemeralStatusTerminating {
			expiredRuntimes = append(expiredRuntimes, id)
			// Mark for termination (actual termination happens asynchronously)
			runtime.Status = EphemeralStatusTerminating
		}
	}

	return expiredRuntimes
}

// AddMember adds a member to an ephemeral runtime
func (erm *EphemeralRuntimeManager) AddMember(runtimeID, userID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	// Check if already a member
	for _, member := range runtime.Members {
		if member == userID {
			return nil // Already a member
		}
	}

	runtime.Members = append(runtime.Members, userID)
	return nil
}

// RemoveMember removes a member from an ephemeral runtime
func (erm *EphemeralRuntimeManager) RemoveMember(runtimeID, userID string) error {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	runtime, exists := erm.runtimes[runtimeID]
	if !exists {
		return fmt.Errorf("ephemeral runtime not found: %s", runtimeID)
	}

	newMembers := []string{}
	for _, member := range runtime.Members {
		if member != userID {
			newMembers = append(newMembers, member)
		}
	}
	runtime.Members = newMembers

	return nil
}

// getResourcesForInstanceType returns resource allocation for an instance type
func (erm *EphemeralRuntimeManager) getResourcesForInstanceType(instanceType string) *EphemeralResources {
	switch instanceType {
	case "small":
		return &EphemeralResources{CPU: 1, Memory: "1GB", Disk: "10GB"}
	case "medium":
		return &EphemeralResources{CPU: 2, Memory: "4GB", Disk: "20GB"}
	case "large":
		return &EphemeralResources{CPU: 4, Memory: "8GB", Disk: "40GB"}
	default:
		return &EphemeralResources{CPU: 1, Memory: "1GB", Disk: "10GB"}
	}
}

// copyEphemeralRuntime creates a deep copy of an ephemeral runtime
func copyEphemeralRuntime(er *EphemeralRuntime) *EphemeralRuntime {
	copy := &EphemeralRuntime{
		ID:            er.ID,
		Name:          er.Name,
		WorkspaceID:   er.WorkspaceID,
		Region:        er.Region,
		Status:        er.Status,
		InstanceType:  er.InstanceType,
		Image:         er.Image,
		CreatedAt:     er.CreatedAt,
		ExpiresAt:     er.ExpiresAt,
		LastActivityAt: er.LastActivityAt,
		CreatedBy:     er.CreatedBy,
		AccessURL:     er.AccessURL,
		SSHKey:        er.SSHKey,
		SnapshotID:    er.SnapshotID,
	}

	if er.Members != nil {
		copy.Members = make([]string, len(er.Members))
		copySlice(copy.Members, er.Members)
	}

	if er.Resources != nil {
		copy.Resources = &EphemeralResources{
			CPU:       er.Resources.CPU,
			Memory:    er.Resources.Memory,
			Disk:      er.Resources.Disk,
			Network:   er.Resources.Network,
			PublicIP:  er.Resources.PublicIP,
			PrivateIP: er.Resources.PrivateIP,
		}
	}

	if er.Metadata != nil {
		copy.Metadata = make(map[string]interface{})
		for k, v := range er.Metadata {
			copy.Metadata[k] = v
		}
	}

	return copy
}

// copySlice copies elements from src to dst (helper to avoid conflict with built-in copy)
func copySlice(dst, src []string) {
	for i := range src {
		dst[i] = src[i]
	}
}

