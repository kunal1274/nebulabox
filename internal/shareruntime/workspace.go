package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// Workspace represents a shared runtime workspace
type Workspace struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	OwnerID     string                 `json:"ownerId"`
	ContainerID string                 `json:"containerId"` // The shared container
	Status      string                 `json:"status"`       // active, paused, stopped, sleeping
	Members     []WorkspaceMember      `json:"members"`
	Settings    WorkspaceSettings      `json:"settings"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// WorkspaceMember represents a member in a shared workspace
type WorkspaceMember struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Role      string    `json:"role"` // owner, admin, editor, viewer
	JoinedAt  time.Time `json:"joinedAt"`
	LastSeen  time.Time `json:"lastSeen,omitempty"`
	IsActive  bool      `json:"isActive"`
}

// WorkspaceSettings defines workspace configuration
type WorkspaceSettings struct {
	AllowGuestAccess   bool              `json:"allowGuestAccess,omitempty"`
	MaxMembers         int               `json:"maxMembers,omitempty"`
	SessionTimeout     int               `json:"sessionTimeout,omitempty"` // minutes
	AutoPauseOnIdle    bool              `json:"autoPauseOnIdle,omitempty"`
	IdleTimeout        int               `json:"idleTimeout,omitempty"` // minutes
	AllowedIPs         []string          `json:"allowedIPs,omitempty"`
	Permissions        map[string]string `json:"permissions,omitempty"` // action -> role
	AuditLogging       bool              `json:"auditLogging,omitempty"`
	ResourceLimits     *ResourceLimits   `json:"resourceLimits,omitempty"`
}

// ResourceLimits defines resource constraints for workspace
type ResourceLimits struct {
	MaxCPU    float64 `json:"maxCpu,omitempty"`
	MaxMemory int64   `json:"maxMemory,omitempty"` // bytes
	MaxDisk   int64   `json:"maxDisk,omitempty"`   // bytes
}

// Session represents an active session in a shared workspace
type Session struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId"`
	UserID      string            `json:"userId"`
	Username    string            `json:"username"`
	Type        string            `json:"type"` // terminal, file, api, exec
	Connection  SessionConnection `json:"connection"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	StartedAt   time.Time         `json:"startedAt"`
	LastActivity time.Time        `json:"lastActivity"`
}

// SessionConnection represents connection details for a session
type SessionConnection struct {
	Protocol string `json:"protocol"` // websocket, ssh, http
	Endpoint string `json:"endpoint"`
	Token    string `json:"token,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// WorkspaceManager manages shared runtime workspaces
type WorkspaceManager struct {
	workspaces map[string]*Workspace
	sessions   map[string]*Session // sessionID -> Session
	userSessions map[string][]string // userID -> []sessionID
	sessionManager *SessionManager
	mu         sync.RWMutex
}

// NewWorkspaceManager creates a new workspace manager
func NewWorkspaceManager() *WorkspaceManager {
	return &WorkspaceManager{
		workspaces:    make(map[string]*Workspace),
		sessions:      make(map[string]*Session),
		userSessions:  make(map[string][]string),
		sessionManager: NewSessionManager(),
	}
}

// CreateWorkspace creates a new shared workspace
func (wm *WorkspaceManager) CreateWorkspace(name, description, ownerID, containerID string, settings WorkspaceSettings) (*Workspace, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace := &Workspace{
		ID:          fmt.Sprintf("ws-%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		ContainerID: containerID,
		Status:      "active",
		Members: []WorkspaceMember{
			{
				UserID:   ownerID,
				Role:     "owner",
				JoinedAt: time.Now(),
				IsActive: true,
			},
		},
		Settings:  settings,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	wm.workspaces[workspace.ID] = workspace
	return workspace, nil
}

// GetWorkspace retrieves a workspace by ID
func (wm *WorkspaceManager) GetWorkspace(workspaceID string) (*Workspace, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	workspace, exists := wm.workspaces[workspaceID]
	return workspace, exists
}

// ListWorkspaces returns workspaces for a user
func (wm *WorkspaceManager) ListWorkspaces(userID string) []*Workspace {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	var result []*Workspace
	for _, workspace := range wm.workspaces {
		// Check if user is a member
		for _, member := range workspace.Members {
			if member.UserID == userID {
				result = append(result, workspace)
				break
			}
		}
	}
	return result
}

// AddMember adds a member to a workspace
func (wm *WorkspaceManager) AddMember(workspaceID, userID, username, role string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Check max members limit
	if workspace.Settings.MaxMembers > 0 && len(workspace.Members) >= workspace.Settings.MaxMembers {
		return fmt.Errorf("workspace has reached maximum member limit")
	}

	// Check if already a member
	for _, member := range workspace.Members {
		if member.UserID == userID {
			return fmt.Errorf("user is already a member")
		}
	}

	workspace.Members = append(workspace.Members, WorkspaceMember{
		UserID:   userID,
		Username: username,
		Role:     role,
		JoinedAt: time.Now(),
		IsActive: false,
	})

	workspace.UpdatedAt = time.Now()
	return nil
}

// RemoveMember removes a member from a workspace
func (wm *WorkspaceManager) RemoveMember(workspaceID, userID string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Cannot remove owner
	if workspace.OwnerID == userID {
		return fmt.Errorf("cannot remove workspace owner")
	}

	// Find and remove member
	newMembers := []WorkspaceMember{}
	for _, member := range workspace.Members {
		if member.UserID != userID {
			newMembers = append(newMembers, member)
		}
	}

	if len(newMembers) == len(workspace.Members) {
		return fmt.Errorf("user is not a member")
	}

	workspace.Members = newMembers
	workspace.UpdatedAt = time.Now()

	// Close all sessions for this user in this workspace
	wm.closeUserSessions(workspaceID, userID)

	return nil
}

// UpdateMemberRole updates a member's role in a workspace
func (wm *WorkspaceManager) UpdateMemberRole(workspaceID, userID, newRole string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	for i := range workspace.Members {
		if workspace.Members[i].UserID == userID {
			// Cannot change owner role
			if workspace.Members[i].Role == "owner" && newRole != "owner" {
				return fmt.Errorf("cannot change owner role")
			}
			workspace.Members[i].Role = newRole
			workspace.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("user is not a member")
}

// CreateSession creates a new session in a workspace
func (wm *WorkspaceManager) CreateSession(workspaceID, userID, username, sessionType string, connection SessionConnection) (*Session, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return nil, fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Check if user is a member
	var isMember bool
	for _, member := range workspace.Members {
		if member.UserID == userID {
			isMember = true
			// Update last seen
			member.LastSeen = time.Now()
			member.IsActive = true
			break
		}
	}

	if !isMember {
		return nil, fmt.Errorf("user is not a member of this workspace")
	}

	// Check permissions based on session type
	if !wm.checkPermission(workspace, userID, sessionType) {
		return nil, fmt.Errorf("permission denied for session type: %s", sessionType)
	}

	session := &Session{
		ID:          fmt.Sprintf("sess-%d", time.Now().UnixNano()),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Username:    username,
		Type:        sessionType,
		Connection:  connection,
		StartedAt:   time.Now(),
		LastActivity: time.Now(),
	}

	wm.sessions[session.ID] = session
	wm.userSessions[userID] = append(wm.userSessions[userID], session.ID)

	// Register with session multiplexer
	mux := wm.sessionManager.GetMultiplexer(workspaceID)
	mux.RegisterSession(session.ID, userID, username, sessionType)

	return session, nil
}

// GetSession retrieves a session by ID
func (wm *WorkspaceManager) GetSession(sessionID string) (*Session, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	session, exists := wm.sessions[sessionID]
	return session, exists
}

// ListWorkspaceSessions returns all sessions for a workspace
func (wm *WorkspaceManager) ListWorkspaceSessions(workspaceID string) []*Session {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	var result []*Session
	for _, session := range wm.sessions {
		if session.WorkspaceID == workspaceID {
			result = append(result, session)
		}
	}
	return result
}

// CloseSession closes a session
func (wm *WorkspaceManager) CloseSession(sessionID string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	session, exists := wm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	delete(wm.sessions, sessionID)

	// Unregister from session multiplexer
	mux := wm.sessionManager.GetMultiplexer(session.WorkspaceID)
	mux.UnregisterSession(sessionID)

	// Remove from user sessions
	sessions := wm.userSessions[session.UserID]
	newSessions := []string{}
	for _, sid := range sessions {
		if sid != sessionID {
			newSessions = append(newSessions, sid)
		}
	}
	wm.userSessions[session.UserID] = newSessions

	return nil
}

// closeUserSessions closes all sessions for a user in a workspace
func (wm *WorkspaceManager) closeUserSessions(workspaceID, userID string) {
	sessions := wm.ListWorkspaceSessions(workspaceID)
	for _, session := range sessions {
		if session.UserID == userID {
			wm.CloseSession(session.ID)
		}
	}
}

// checkPermission checks if a user has permission for an action
func (wm *WorkspaceManager) checkPermission(workspace *Workspace, userID, action string) bool {
	// Find user role
	var userRole string
	for _, member := range workspace.Members {
		if member.UserID == userID {
			userRole = member.Role
			break
		}
	}

	if userRole == "" {
		return false
	}

	// Owner and admin have all permissions
	if userRole == "owner" || userRole == "admin" {
		return true
	}

	// Check workspace-specific permissions
	if requiredRole, exists := workspace.Settings.Permissions[action]; exists {
		roleHierarchy := map[string]int{
			"owner":  4,
			"admin":  3,
			"editor": 2,
			"viewer": 1,
		}
		userLevel := roleHierarchy[userRole]
		requiredLevel := roleHierarchy[requiredRole]
		return userLevel >= requiredLevel
	}

	// Default permissions based on role
	switch userRole {
	case "editor":
		return action == "terminal" || action == "file" || action == "exec"
	case "viewer":
		return action == "api" || action == "file"
	default:
		return false
	}
}

// UpdateWorkspaceStatus updates workspace status
func (wm *WorkspaceManager) UpdateWorkspaceStatus(workspaceID, status string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	validStatuses := map[string]bool{
		"active": true,
		"paused": true,
		"stopped": true,
		"sleeping": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	workspace.Status = status
	workspace.UpdatedAt = time.Now()
	return nil
}

// DeleteWorkspace deletes a workspace
func (wm *WorkspaceManager) DeleteWorkspace(workspaceID string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	_, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Close all sessions
	sessions := wm.ListWorkspaceSessions(workspaceID)
	for _, session := range sessions {
		delete(wm.sessions, session.ID)
	}

	// Remove from user sessions
	for userID := range wm.userSessions {
		sessions := wm.userSessions[userID]
		newSessions := []string{}
		for _, sid := range sessions {
			if _, exists := wm.sessions[sid]; exists {
				newSessions = append(newSessions, sid)
			}
		}
		wm.userSessions[userID] = newSessions
	}

	delete(wm.workspaces, workspaceID)
	
	// Remove session multiplexer
	wm.sessionManager.RemoveMultiplexer(workspaceID)
	
	return nil
}

// GetSessionManager returns the session manager
func (wm *WorkspaceManager) GetSessionManager() *SessionManager {
	return wm.sessionManager
}
