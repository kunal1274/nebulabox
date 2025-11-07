package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// getWorkspaceInviteLink handles GET /api/shareruntime/workspaces/:id/invites/:token/link
func (s *Server) getWorkspaceInviteLink(c *gin.Context) {
	workspaceID := c.Param("id")
	token := c.Param("token")
	
	invite, exists := s.sharedRuntimeInviteManager.GetInviteByToken(token)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invite not found",
		})
		return
	}

	// Verify invite belongs to workspace
	if invite.WorkspaceID != workspaceID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invite does not belong to this workspace",
		})
		return
	}

	// Generate shareable link
	baseURL := os.Getenv("NEBULABOX_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	link := fmt.Sprintf("%s/shareruntime/accept/%s", baseURL, token)

	c.JSON(http.StatusOK, gin.H{
		"link": link,
	})
}

// getInviteLink handles GET /api/shareruntime/invites/:token
func (s *Server) getInviteLink(c *gin.Context) {
	token := c.Param("token")
	invite, exists := s.sharedRuntimeInviteManager.GetInviteByToken(token)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invite not found",
		})
		return
	}

	// Generate shareable link
	inviteLink := fmt.Sprintf("/shareruntime/accept/%s", token)

	c.JSON(http.StatusOK, gin.H{
		"invite": invite,
		"link":   inviteLink,
		"token":  token,
	})
}

// getInviteInfo handles GET /api/shareruntime/invites/:token/info
func (s *Server) getInviteInfo(c *gin.Context) {
	token := c.Param("token")
	invite, exists := s.sharedRuntimeInviteManager.GetInviteByToken(token)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invite not found",
		})
		return
	}

	// Return public info only
	c.JSON(http.StatusOK, gin.H{
		"workspaceId": invite.WorkspaceID,
		"role":        invite.Role,
		"inviterName": invite.InviterName,
		"expiresAt":   invite.ExpiresAt,
		"status":      invite.Status,
	})
}

// getUserPermissions handles GET /api/shareruntime/workspaces/:id/permissions
func (s *Server) getUserPermissions(c *gin.Context) {
	workspaceID := c.Param("id")
	username, _ := c.Get("username")
	userID := username.(string)

	role, err := s.sharedRuntimeManager.GetUserRole(workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "User is not a member",
		})
		return
	}

	// Get permissions for role
	permissions := make(map[string]bool)
	if perms, exists := map[string][]string{
		"owner": {
			"workspace:view", "workspace:edit", "workspace:delete", "workspace:manage",
			"member:view", "member:invite", "member:remove", "member:manage",
			"session:create", "session:view", "session:manage",
			"tunnel:create", "tunnel:view", "tunnel:manage",
			"container:view", "container:exec", "container:modify",
		},
		"admin": {
			"workspace:view", "workspace:edit", "workspace:manage",
			"member:view", "member:invite", "member:remove", "member:manage",
			"session:create", "session:view", "session:manage",
			"tunnel:create", "tunnel:view", "tunnel:manage",
			"container:view", "container:exec", "container:modify",
		},
		"editor": {
			"workspace:view",
			"member:view",
			"session:create", "session:view",
			"tunnel:create", "tunnel:view",
			"container:view", "container:exec", "container:modify",
		},
		"viewer": {
			"workspace:view",
			"member:view",
			"session:view",
			"tunnel:view",
			"container:view",
		},
	}[role]; exists {
		for _, perm := range perms {
			permissions[perm] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"role":        role,
		"permissions": permissions,
	})
}

