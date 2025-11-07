package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
)

// CreateWorkspaceRequest represents a request to create a shared workspace
type CreateWorkspaceRequest struct {
	Name        string                     `json:"name" binding:"required"`
	Description string                     `json:"description,omitempty"`
	ContainerID string                     `json:"containerId" binding:"required"`
	Settings    shareruntime.WorkspaceSettings `json:"settings,omitempty"`
}

// AddMemberRequest represents a request to add a member to a workspace
type AddMemberRequest struct {
	UserID   string `json:"userId" binding:"required"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// CreateInviteRequest represents a request to create an invitation
type CreateInviteRequest struct {
	Email         string `json:"email,omitempty"`
	Role          string `json:"role" binding:"required"`
	ExpiresInHours int   `json:"expiresInHours,omitempty"`
}

// CreateSessionRequest represents a request to create a session
type CreateSessionRequest struct {
	Type       string                            `json:"type" binding:"required"`
	Connection shareruntime.SessionConnection    `json:"connection,omitempty"`
}

// createWorkspace handles POST /api/shareruntime/workspaces
func (s *Server) createWorkspace(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)

	var req CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionWorkspaceCreated, "", "Invalid request body", map[string]string{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Settings.MaxMembers == 0 {
		req.Settings.MaxMembers = 10 // Default
	}
	if req.Settings.SessionTimeout == 0 {
		req.Settings.SessionTimeout = 60 // Default 60 minutes
	}

	workspace, err := s.sharedRuntimeManager.CreateWorkspace(
		req.Name,
		req.Description,
		userID,
		req.ContainerID,
		req.Settings,
	)
	if err != nil {
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionWorkspaceCreated, "", "Failed to create workspace", map[string]string{"error": err.Error(), "name": req.Name})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create workspace",
			"details": err.Error(),
		})
		return
	}

	s.auditLogger.LogAction(userID, userID, shareruntime.ActionWorkspaceCreated, workspace.ID, map[string]string{
		"workspaceName": workspace.Name,
		"containerId":   workspace.ContainerID,
	})

	// Record sync event
	s.syncManager.RecordChange(workspace.ID, "workspace", workspace.ID, shareruntime.ChangeTypeCreate, workspace, userID)

	c.JSON(http.StatusOK, workspace)
}

// listWorkspaces handles GET /api/shareruntime/workspaces
func (s *Server) listWorkspaces(c *gin.Context) {
	username, _ := c.Get("username")
	userID := username.(string)

	workspaces := s.sharedRuntimeManager.ListWorkspaces(userID)
	c.JSON(http.StatusOK, gin.H{
		"workspaces": workspaces,
		"count":      len(workspaces),
	})
}

// getWorkspace handles GET /api/shareruntime/workspaces/:id
func (s *Server) getWorkspace(c *gin.Context) {
	workspaceID := c.Param("id")
	workspace, exists := s.sharedRuntimeManager.GetWorkspace(workspaceID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Workspace not found",
		})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

// addWorkspaceMember handles POST /api/shareruntime/workspaces/:id/members
func (s *Server) addWorkspaceMember(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	// Check permission
	if canInvite, err := s.sharedRuntimeManager.CanInvite(workspaceID, userID); err != nil || !canInvite {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: cannot invite members",
		})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.sharedRuntimeManager.AddMember(workspaceID, req.UserID, req.Username, req.Role); err != nil {
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionMemberAdded, workspaceID, "Failed to add member", map[string]string{"targetUserId": req.UserID, "role": req.Role})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add member",
			"details": err.Error(),
		})
		return
	}

	s.auditLogger.LogAction(userID, userID, shareruntime.ActionMemberAdded, workspaceID, map[string]string{
		"targetUserId":  req.UserID,
		"targetUsername": req.Username,
		"role":          req.Role,
	})

	// Record sync event
	memberData := map[string]interface{}{
		"userId":   req.UserID,
		"username": req.Username,
		"role":     req.Role,
	}
	s.syncManager.RecordChange(workspaceID, "member", req.UserID, shareruntime.ChangeTypeCreate, memberData, userID)

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

// removeWorkspaceMember handles DELETE /api/shareruntime/workspaces/:id/members/:userId
func (s *Server) removeWorkspaceMember(c *gin.Context) {
	workspaceID := c.Param("id")
	targetUserID := c.Param("userId")
	username, _ := c.Get("username")
	userID := username.(string)

	// Check permission
	if canRemove, err := s.sharedRuntimeManager.CanRemoveMember(workspaceID, userID, targetUserID); err != nil || !canRemove {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: cannot remove members",
		})
		return
	}

	if err := s.sharedRuntimeManager.RemoveMember(workspaceID, targetUserID); err != nil {
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionMemberRemoved, workspaceID, "Failed to remove member", map[string]string{"targetUserId": targetUserID})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to remove member",
			"details": err.Error(),
		})
		return
	}

	s.auditLogger.LogAction(userID, userID, shareruntime.ActionMemberRemoved, workspaceID, map[string]string{
		"targetUserId": targetUserID,
	})

	// Record sync event
	s.syncManager.RecordChange(workspaceID, "member", targetUserID, shareruntime.ChangeTypeDelete, nil, userID)

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// updateWorkspaceMemberRole handles PATCH /api/shareruntime/workspaces/:id/members/:userId/role
func (s *Server) updateWorkspaceMemberRole(c *gin.Context) {
	workspaceID := c.Param("id")
	userID := c.Param("userId")
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.sharedRuntimeManager.UpdateMemberRole(workspaceID, userID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

// createInvite handles POST /api/shareruntime/workspaces/:id/invites
func (s *Server) createInvite(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	// Check permission
	if canInvite, err := s.sharedRuntimeManager.CanInvite(workspaceID, userID); err != nil || !canInvite {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: cannot create invites",
		})
		return
	}

	var req CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if req.ExpiresInHours == 0 {
		req.ExpiresInHours = 24 // Default 24 hours
	}

	// Validate role
	validRoles := map[string]bool{
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Must be one of: admin, editor, viewer",
		})
		return
	}

	invite, err := s.sharedRuntimeInviteManager.CreateInvite(
		workspaceID,
		userID,
		userID, // Use userID as name for now
		req.Email,
		req.Role,
		req.ExpiresInHours,
	)
	if err != nil {
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionInviteCreated, workspaceID, "Failed to create invite", map[string]string{"role": req.Role, "error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create invite",
			"details": err.Error(),
		})
		return
	}

	s.auditLogger.LogAction(userID, userID, shareruntime.ActionInviteCreated, workspaceID, map[string]string{
		"inviteId": invite.ID,
		"role":     invite.Role,
		"email":    req.Email,
	})

	c.JSON(http.StatusOK, invite)
}

// listInvites handles GET /api/shareruntime/workspaces/:id/invites
func (s *Server) listInvites(c *gin.Context) {
	workspaceID := c.Param("id")
	invites := s.sharedRuntimeInviteManager.ListWorkspaceInvites(workspaceID)
	c.JSON(http.StatusOK, gin.H{
		"invites": invites,
		"count":   len(invites),
	})
}

// acceptInvite handles POST /api/shareruntime/invites/:token/accept
func (s *Server) acceptInvite(c *gin.Context) {
	token := c.Param("token")
	username, _ := c.Get("username")
	userID := username.(string)

	invite, err := s.sharedRuntimeInviteManager.AcceptInvite(token, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to accept invite",
			"details": err.Error(),
		})
		return
	}

	// Automatically add user as member to workspace
	if err := s.sharedRuntimeManager.AddMember(
		invite.WorkspaceID,
		userID,
		userID,
		invite.Role,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add user to workspace",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invite": invite,
		"message": "Invite accepted and user added to workspace",
	})
}

// createSession handles POST /api/shareruntime/workspaces/:id/sessions
func (s *Server) createSession(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	// Check permission
	if canCreate, err := s.sharedRuntimeManager.CanCreateSession(workspaceID, userID); err != nil || !canCreate {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: cannot create sessions",
		})
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Generate connection details
	connection := req.Connection
	if connection.Protocol == "" {
		connection.Protocol = "websocket"
	}
	if connection.Endpoint == "" {
		connection.Endpoint = fmt.Sprintf("/ws/workspace/%s", workspaceID)
	}

	session, err := s.sharedRuntimeManager.CreateSession(
		workspaceID,
		userID,
		userID,
		req.Type,
		connection,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create session",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, session)
}

// listWorkspaceSessions handles GET /api/shareruntime/workspaces/:id/sessions
func (s *Server) listWorkspaceSessions(c *gin.Context) {
	workspaceID := c.Param("id")
	sessions := s.sharedRuntimeManager.ListWorkspaceSessions(workspaceID)
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// closeSession handles DELETE /api/shareruntime/sessions/:id
func (s *Server) closeSession(c *gin.Context) {
	sessionID := c.Param("id")
	if err := s.sharedRuntimeManager.CloseSession(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to close session",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session closed"})
}

// updateWorkspaceStatus handles PATCH /api/shareruntime/workspaces/:id/status
func (s *Server) updateWorkspaceStatus(c *gin.Context) {
	workspaceID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := s.sharedRuntimeManager.UpdateWorkspaceStatus(workspaceID, req.Status); err != nil {
		username, _ := c.Get("username")
		userID := username.(string)
		s.auditLogger.LogFailure(userID, userID, shareruntime.ActionWorkspaceStatusChanged, workspaceID, "Failed to update status", map[string]string{"status": req.Status})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update status",
			"details": err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	userID := username.(string)
	s.auditLogger.LogAction(userID, userID, shareruntime.ActionWorkspaceStatusChanged, workspaceID, map[string]string{
		"status": req.Status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// deleteWorkspace handles DELETE /api/shareruntime/workspaces/:id
func (s *Server) deleteWorkspace(c *gin.Context) {
	workspaceID := c.Param("id")
	if err := s.sharedRuntimeManager.DeleteWorkspace(workspaceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete workspace",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted"})
}

