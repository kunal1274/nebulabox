package shareruntime

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ChangeType represents the type of change in a sync event
type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

// SyncEvent represents a change event in the shared runtime
type SyncEvent struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId"`
	ResourceType string           `json:"resourceType"` // workspace, member, session, tunnel, etc.
	ResourceID   string           `json:"resourceId"`
	ChangeType   ChangeType       `json:"changeType"`
	Data         json.RawMessage  `json:"data,omitempty"`
	Timestamp    time.Time        `json:"timestamp"`
	UserID       string           `json:"userId"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ReplicationAdapter defines the interface for replication adapters
type ReplicationAdapter interface {
	// RegisterChange registers a change event
	RegisterChange(event *SyncEvent) error
	
	// GetChanges retrieves changes since a given timestamp
	GetChanges(workspaceID string, since time.Time) ([]*SyncEvent, error)
	
	// GetLatestChangeID returns the latest change ID for a workspace
	GetLatestChangeID(workspaceID string) (string, error)
	
	// Subscribe subscribes to change events
	Subscribe(workspaceID string, callback func(*SyncEvent)) (string, error) // Returns subscription ID
	
	// Unsubscribe unsubscribes from change events
	Unsubscribe(subscriptionID string) error
}

// InMemoryReplicationAdapter is an in-memory implementation of ReplicationAdapter
type InMemoryReplicationAdapter struct {
	events          map[string]*SyncEvent // eventID -> event
	workspaceEvents map[string][]string   // workspaceID -> []eventID
	subscribers     map[string]*Subscription // subscriptionID -> Subscription
	subscriptionsByWorkspace map[string][]string // workspaceID -> []subscriptionID
	mu              sync.RWMutex
	nextEventID     int64
	nextSubID       int64
}

// Subscription represents an active subscription
type Subscription struct {
	ID          string
	WorkspaceID string
	Callback    func(*SyncEvent)
	CreatedAt   time.Time
}

// NewInMemoryReplicationAdapter creates a new in-memory replication adapter
func NewInMemoryReplicationAdapter() *InMemoryReplicationAdapter {
	return &InMemoryReplicationAdapter{
		events:                   make(map[string]*SyncEvent),
		workspaceEvents:          make(map[string][]string),
		subscribers:              make(map[string]*Subscription),
		subscriptionsByWorkspace: make(map[string][]string),
	}
}

// RegisterChange registers a change event
func (a *InMemoryReplicationAdapter) RegisterChange(event *SyncEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if event.ID == "" {
		event.ID = fmt.Sprintf("evt-%d", a.nextEventID)
		a.nextEventID++
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	a.events[event.ID] = event
	a.workspaceEvents[event.WorkspaceID] = append(a.workspaceEvents[event.WorkspaceID], event.ID)

	// Notify subscribers
	if subs, exists := a.subscriptionsByWorkspace[event.WorkspaceID]; exists {
		for _, subID := range subs {
			if sub, exists := a.subscribers[subID]; exists {
				go func(callback func(*SyncEvent), evt *SyncEvent) {
					defer func() {
						if r := recover(); r != nil {
							// Prevent panic from breaking the sync system
						}
					}()
					callback(evt)
				}(sub.Callback, event)
			}
		}
	}

	return nil
}

// GetChanges retrieves changes since a given timestamp
func (a *InMemoryReplicationAdapter) GetChanges(workspaceID string, since time.Time) ([]*SyncEvent, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	eventIDs, exists := a.workspaceEvents[workspaceID]
	if !exists {
		return []*SyncEvent{}, nil
	}

	var result []*SyncEvent
	for _, eventID := range eventIDs {
		if event, exists := a.events[eventID]; exists {
			if event.Timestamp.After(since) || event.Timestamp.Equal(since) {
				result = append(result, event)
			}
		}
	}

	return result, nil
}

// GetLatestChangeID returns the latest change ID for a workspace
func (a *InMemoryReplicationAdapter) GetLatestChangeID(workspaceID string) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	eventIDs, exists := a.workspaceEvents[workspaceID]
	if !exists || len(eventIDs) == 0 {
		return "", nil
	}

	// Find the latest event
	var latestEvent *SyncEvent
	var latestTime time.Time
	for _, eventID := range eventIDs {
		if event, exists := a.events[eventID]; exists {
			if event.Timestamp.After(latestTime) {
				latestTime = event.Timestamp
				latestEvent = event
			}
		}
	}

	if latestEvent != nil {
		return latestEvent.ID, nil
	}
	return "", nil
}

// Subscribe subscribes to change events
func (a *InMemoryReplicationAdapter) Subscribe(workspaceID string, callback func(*SyncEvent)) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	subID := fmt.Sprintf("sub-%d", a.nextSubID)
	a.nextSubID++

	sub := &Subscription{
		ID:          subID,
		WorkspaceID: workspaceID,
		Callback:    callback,
		CreatedAt:   time.Now(),
	}

	a.subscribers[subID] = sub
	a.subscriptionsByWorkspace[workspaceID] = append(a.subscriptionsByWorkspace[workspaceID], subID)

	return subID, nil
}

// Unsubscribe unsubscribes from change events
func (a *InMemoryReplicationAdapter) Unsubscribe(subscriptionID string) error {
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

// SyncManager manages database synchronization
type SyncManager struct {
	adapter ReplicationAdapter
	mu      sync.RWMutex
}

// NewSyncManager creates a new sync manager
func NewSyncManager(adapter ReplicationAdapter) *SyncManager {
	return &SyncManager{
		adapter: adapter,
	}
}

// RecordChange records a change in the sync system
func (sm *SyncManager) RecordChange(workspaceID, resourceType, resourceID string, changeType ChangeType, data interface{}, userID string) error {
	var dataJSON json.RawMessage
	if data != nil {
		var err error
		dataJSON, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	event := &SyncEvent{
		WorkspaceID: workspaceID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ChangeType:   changeType,
		Data:         dataJSON,
		UserID:       userID,
		Timestamp:    time.Now(),
	}

	return sm.adapter.RegisterChange(event)
}

// GetChangesSince retrieves changes since a timestamp
func (sm *SyncManager) GetChangesSince(workspaceID string, since time.Time) ([]*SyncEvent, error) {
	return sm.adapter.GetChanges(workspaceID, since)
}

// SubscribeToChanges subscribes to change events for a workspace
func (sm *SyncManager) SubscribeToChanges(workspaceID string, callback func(*SyncEvent)) (string, error) {
	return sm.adapter.Subscribe(workspaceID, callback)
}

// UnsubscribeFromChanges unsubscribes from change events
func (sm *SyncManager) UnsubscribeFromChanges(subscriptionID string) error {
	return sm.adapter.Unsubscribe(subscriptionID)
}

// GetLatestChangeID gets the latest change ID for a workspace
func (sm *SyncManager) GetLatestChangeID(workspaceID string) (string, error) {
	return sm.adapter.GetLatestChangeID(workspaceID)
}

