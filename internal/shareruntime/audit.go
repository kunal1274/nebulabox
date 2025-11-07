package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// AuditAction represents the type of action being audited
type AuditAction string

const (
	// Workspace actions
	ActionWorkspaceCreated   AuditAction = "workspace:created"
	ActionWorkspaceUpdated   AuditAction = "workspace:updated"
	ActionWorkspaceDeleted   AuditAction = "workspace:deleted"
	ActionWorkspaceStatusChanged AuditAction = "workspace:status_changed"

	// Member actions
	ActionMemberAdded    AuditAction = "member:added"
	ActionMemberRemoved  AuditAction = "member:removed"
	ActionMemberRoleChanged AuditAction = "member:role_changed"

	// Invite actions
	ActionInviteCreated  AuditAction = "invite:created"
	ActionInviteAccepted AuditAction = "invite:accepted"
	ActionInviteRejected AuditAction = "invite:rejected"
	ActionInviteRevoked  AuditAction = "invite:revoked"

	// Session actions
	ActionSessionCreated AuditAction = "session:created"
	ActionSessionClosed  AuditAction = "session:closed"
	ActionSessionActivity AuditAction = "session:activity"

	// Tunnel actions
	ActionTunnelCreated AuditAction = "tunnel:created"
	ActionTunnelClosed  AuditAction = "tunnel:closed"
	ActionTunnelConnected AuditAction = "tunnel:connected"

	// Permission actions
	ActionPermissionDenied AuditAction = "permission:denied"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId,omitempty"`
	UserID      string            `json:"userId"`
	Username    string            `json:"username"`
	Action      AuditAction       `json:"action"`
	ResourceType string           `json:"resourceType,omitempty"` // workspace, member, invite, session, tunnel
	ResourceID   string           `json:"resourceId,omitempty"`
	Details      map[string]string `json:"details,omitempty"`
	IPAddress    string           `json:"ipAddress,omitempty"`
	UserAgent    string           `json:"userAgent,omitempty"`
	Success      bool             `json:"success"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
	Timestamp    time.Time        `json:"timestamp"`
}

// AuditLogger manages audit logging for shared runtime
type AuditLogger struct {
	logs        []*AuditLog
	maxLogs     int
	mu          sync.RWMutex
	// MongoDB repository for persistence (optional)
	mongoRepo   AuditLogMongoRepo
}

// AuditLogMongoRepo interface for MongoDB persistence
type AuditLogMongoRepo interface {
	Insert(ctx interface{}, entry interface{}) error
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(maxLogs int) *AuditLogger {
	if maxLogs <= 0 {
		maxLogs = 10000 // Default 10k logs
	}
	return &AuditLogger{
		logs:    make([]*AuditLog, 0, maxLogs),
		maxLogs: maxLogs,
	}
}

// Log creates a new audit log entry
func (al *AuditLogger) Log(log *AuditLog) {
	al.mu.Lock()
	defer al.mu.Unlock()

	if log.ID == "" {
		log.ID = fmt.Sprintf("audit-%d", time.Now().UnixNano())
	}
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	al.logs = append(al.logs, log)

	// Maintain max size (FIFO)
	if len(al.logs) > al.maxLogs {
		al.logs = al.logs[1:]
	}
}

// LogAction creates an audit log entry with action details
func (al *AuditLogger) LogAction(userID, username string, action AuditAction, workspaceID string, details map[string]string) {
	log := &AuditLog{
		UserID:      userID,
		Username:    username,
		Action:      action,
		WorkspaceID: workspaceID,
		Details:     details,
		Success:     true,
	}
	al.Log(log)
}

// LogFailure creates an audit log entry for a failed action
func (al *AuditLogger) LogFailure(userID, username string, action AuditAction, workspaceID, errorMsg string, details map[string]string) {
	log := &AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		WorkspaceID:  workspaceID,
		Details:      details,
		Success:      false,
		ErrorMessage: errorMsg,
	}
	al.Log(log)
}

// GetLogs retrieves audit logs with optional filters
func (al *AuditLogger) GetLogs(filters AuditLogFilters) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var result []*AuditLog
	for _, log := range al.logs {
		if al.matchesFilters(log, filters) {
			result = append(result, log)
		}
	}

	// Reverse to get most recent first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// GetWorkspaceLogs retrieves audit logs for a specific workspace
func (al *AuditLogger) GetWorkspaceLogs(workspaceID string, limit int) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(al.logs) - 1; i >= 0 && count < limit; i-- {
		log := al.logs[i]
		if log.WorkspaceID == workspaceID {
			result = append(result, log)
			count++
		}
	}

	return result
}

// GetUserLogs retrieves audit logs for a specific user
func (al *AuditLogger) GetUserLogs(userID string, limit int) []*AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(al.logs) - 1; i >= 0 && count < limit; i-- {
		log := al.logs[i]
		if log.UserID == userID {
			result = append(result, log)
			count++
		}
	}

	return result
}

// AuditLogFilters defines filters for querying audit logs
type AuditLogFilters struct {
	WorkspaceID string
	UserID      string
	Action      AuditAction
	ResourceType string
	ResourceID   string
	Success     *bool
	StartTime   *time.Time
	EndTime     *time.Time
	Limit       int
}

// matchesFilters checks if a log matches the given filters
func (al *AuditLogger) matchesFilters(log *AuditLog, filters AuditLogFilters) bool {
	if filters.WorkspaceID != "" && log.WorkspaceID != filters.WorkspaceID {
		return false
	}
	if filters.UserID != "" && log.UserID != filters.UserID {
		return false
	}
	if filters.Action != "" && log.Action != filters.Action {
		return false
	}
	if filters.ResourceType != "" && log.ResourceType != filters.ResourceType {
		return false
	}
	if filters.ResourceID != "" && log.ResourceID != filters.ResourceID {
		return false
	}
	if filters.Success != nil && log.Success != *filters.Success {
		return false
	}
	if filters.StartTime != nil && log.Timestamp.Before(*filters.StartTime) {
		return false
	}
	if filters.EndTime != nil && log.Timestamp.After(*filters.EndTime) {
		return false
	}
	return true
}

// GetStats returns statistics about audit logs
func (al *AuditLogger) GetStats(workspaceID string, startTime, endTime time.Time) AuditStats {
	al.mu.RLock()
	defer al.mu.RUnlock()

	stats := AuditStats{
		TotalActions:      0,
		SuccessfulActions: 0,
		FailedActions:     0,
		ActionsByType:     make(map[string]int),
		UsersByActivity:   make(map[string]int),
	}

	for _, log := range al.logs {
		if workspaceID != "" && log.WorkspaceID != workspaceID {
			continue
		}
		if log.Timestamp.Before(startTime) || log.Timestamp.After(endTime) {
			continue
		}

		stats.TotalActions++
		if log.Success {
			stats.SuccessfulActions++
		} else {
			stats.FailedActions++
		}

		actionType := string(log.Action)
		stats.ActionsByType[actionType]++

		stats.UsersByActivity[log.Username]++
	}

	return stats
}

// AuditStats represents statistics about audit logs
type AuditStats struct {
	TotalActions      int            `json:"totalActions"`
	SuccessfulActions int            `json:"successfulActions"`
	FailedActions     int            `json:"failedActions"`
	ActionsByType     map[string]int `json:"actionsByType"`
	UsersByActivity   map[string]int `json:"usersByActivity"`
}

