package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// SessionState represents the current state of a session
type SessionState struct {
	SessionID   string            `json:"sessionId"`
	UserID      string            `json:"userId"`
	Username    string            `json:"username"`
	Type        string            `json:"type"`
	State       string            `json:"state"` // active, idle, paused, disconnected
	LastActivity time.Time        `json:"lastActivity"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SessionMultiplexer handles multiplexing multiple sessions to a shared container
type SessionMultiplexer struct {
	workspaceID string
	sessions    map[string]*SessionState // sessionID -> state
	mu          sync.RWMutex
	activitySubscribers map[string]chan *SessionState // userID -> channel
}

// NewSessionMultiplexer creates a new session multiplexer for a workspace
func NewSessionMultiplexer(workspaceID string) *SessionMultiplexer {
	return &SessionMultiplexer{
		workspaceID:        workspaceID,
		sessions:           make(map[string]*SessionState),
		activitySubscribers: make(map[string]chan *SessionState),
	}
}

// RegisterSession registers a new session
func (sm *SessionMultiplexer) RegisterSession(sessionID, userID, username, sessionType string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sessions[sessionID] = &SessionState{
		SessionID:   sessionID,
		UserID:      userID,
		Username:    username,
		Type:        sessionType,
		State:       "active",
		LastActivity: time.Now(),
		Metadata:    make(map[string]string),
	}

	return nil
}

// UpdateActivity updates the last activity time for a session
func (sm *SessionMultiplexer) UpdateActivity(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.LastActivity = time.Now()
	if session.State == "idle" || session.State == "paused" {
		session.State = "active"
	}

	// Notify subscribers
	sm.notifySubscribers(session)

	return nil
}

// UpdateSessionState updates the state of a session
func (sm *SessionMultiplexer) UpdateSessionState(sessionID, state string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	validStates := map[string]bool{
		"active":       true,
		"idle":         true,
		"paused":       true,
		"disconnected": true,
	}
	if !validStates[state] {
		return fmt.Errorf("invalid state: %s", state)
	}

	session.State = state
	sm.notifySubscribers(session)

	return nil
}

// GetSessionState retrieves the state of a session
func (sm *SessionMultiplexer) GetSessionState(sessionID string) (*SessionState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	state, exists := sm.sessions[sessionID]
	return state, exists
}

// ListActiveSessions returns all active sessions
func (sm *SessionMultiplexer) ListActiveSessions() []*SessionState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var result []*SessionState
	for _, state := range sm.sessions {
		if state.State == "active" || state.State == "idle" {
			result = append(result, state)
		}
	}
	return result
}

// UnregisterSession removes a session
func (sm *SessionMultiplexer) UnregisterSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
}

// SubscribeActivity subscribes to session activity updates
func (sm *SessionMultiplexer) SubscribeActivity(userID string) <-chan *SessionState {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	ch := make(chan *SessionState, 10)
	sm.activitySubscribers[userID] = ch
	return ch
}

// UnsubscribeActivity unsubscribes from activity updates
func (sm *SessionMultiplexer) UnsubscribeActivity(userID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if ch, exists := sm.activitySubscribers[userID]; exists {
		close(ch)
		delete(sm.activitySubscribers, userID)
	}
}

// notifySubscribers notifies all subscribers of session activity
func (sm *SessionMultiplexer) notifySubscribers(session *SessionState) {
	for _, ch := range sm.activitySubscribers {
		select {
		case ch <- session:
		default:
			// Channel full, skip notification
		}
	}
}

// SessionManager manages all session multiplexers
type SessionManager struct {
	multiplexers map[string]*SessionMultiplexer // workspaceID -> multiplexer
	mu           sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		multiplexers: make(map[string]*SessionMultiplexer),
	}
}

// GetMultiplexer gets or creates a multiplexer for a workspace
func (sm *SessionManager) GetMultiplexer(workspaceID string) *SessionMultiplexer {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if mux, exists := sm.multiplexers[workspaceID]; exists {
		return mux
	}

	mux := NewSessionMultiplexer(workspaceID)
	sm.multiplexers[workspaceID] = mux
	return mux
}

// RemoveMultiplexer removes a multiplexer for a workspace
func (sm *SessionManager) RemoveMultiplexer(workspaceID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.multiplexers, workspaceID)
}

