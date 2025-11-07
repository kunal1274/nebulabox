package shareruntime

import "fmt"

// Permission represents a specific permission action
type Permission string

const (
	// Workspace permissions
	PermissionWorkspaceView   Permission = "workspace:view"
	PermissionWorkspaceEdit   Permission = "workspace:edit"
	PermissionWorkspaceDelete Permission = "workspace:delete"
	PermissionWorkspaceManage Permission = "workspace:manage"

	// Member permissions
	PermissionMemberView   Permission = "member:view"
	PermissionMemberInvite Permission = "member:invite"
	PermissionMemberRemove Permission = "member:remove"
	PermissionMemberManage Permission = "member:manage"

	// Session permissions
	PermissionSessionCreate Permission = "session:create"
	PermissionSessionView   Permission = "session:view"
	PermissionSessionManage Permission = "session:manage"

	// Tunnel permissions
	PermissionTunnelCreate Permission = "tunnel:create"
	PermissionTunnelView   Permission = "tunnel:view"
	PermissionTunnelManage Permission = "tunnel:manage"

	// Container permissions
	PermissionContainerView Permission = "container:view"
	PermissionContainerExec Permission = "container:exec"
	PermissionContainerModify Permission = "container:modify"
)

// RolePermissions defines permissions for each role
var RolePermissions = map[string][]Permission{
	"owner": {
		// Owners have all permissions
		PermissionWorkspaceView,
		PermissionWorkspaceEdit,
		PermissionWorkspaceDelete,
		PermissionWorkspaceManage,
		PermissionMemberView,
		PermissionMemberInvite,
		PermissionMemberRemove,
		PermissionMemberManage,
		PermissionSessionCreate,
		PermissionSessionView,
		PermissionSessionManage,
		PermissionTunnelCreate,
		PermissionTunnelView,
		PermissionTunnelManage,
		PermissionContainerView,
		PermissionContainerExec,
		PermissionContainerModify,
	},
	"admin": {
		// Admins have most permissions except workspace deletion
		PermissionWorkspaceView,
		PermissionWorkspaceEdit,
		PermissionWorkspaceManage,
		PermissionMemberView,
		PermissionMemberInvite,
		PermissionMemberRemove,
		PermissionMemberManage,
		PermissionSessionCreate,
		PermissionSessionView,
		PermissionSessionManage,
		PermissionTunnelCreate,
		PermissionTunnelView,
		PermissionTunnelManage,
		PermissionContainerView,
		PermissionContainerExec,
		PermissionContainerModify,
	},
	"editor": {
		// Editors can work with sessions and tunnels but not manage members
		PermissionWorkspaceView,
		PermissionMemberView,
		PermissionSessionCreate,
		PermissionSessionView,
		PermissionTunnelCreate,
		PermissionTunnelView,
		PermissionContainerView,
		PermissionContainerExec,
		PermissionContainerModify,
	},
	"viewer": {
		// Viewers have read-only access
		PermissionWorkspaceView,
		PermissionMemberView,
		PermissionSessionView,
		PermissionTunnelView,
		PermissionContainerView,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role string, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// CheckWorkspacePermission checks if a user has permission in a workspace
func (wm *WorkspaceManager) CheckWorkspacePermission(workspaceID, userID string, permission Permission) (bool, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return false, fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Find user role
	var userRole string
	for _, member := range workspace.Members {
		if member.UserID == userID {
			userRole = member.Role
			break
		}
	}

	if userRole == "" {
		return false, nil
	}

	return HasPermission(userRole, permission), nil
}

// RequirePermission checks permission and returns error if denied
func (wm *WorkspaceManager) RequirePermission(workspaceID, userID string, permission Permission) error {
	hasPermission, err := wm.CheckWorkspacePermission(workspaceID, userID, permission)
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("permission denied: %s", permission)
	}
	return nil
}

// CanInvite checks if a user can invite members to a workspace
func (wm *WorkspaceManager) CanInvite(workspaceID, userID string) (bool, error) {
	return wm.CheckWorkspacePermission(workspaceID, userID, PermissionMemberInvite)
}

// CanRemoveMember checks if a user can remove members from a workspace
func (wm *WorkspaceManager) CanRemoveMember(workspaceID, userID, targetUserID string) (bool, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return false, fmt.Errorf("workspace not found: %s", workspaceID)
	}

	// Cannot remove owner
	if workspace.OwnerID == targetUserID {
		return false, nil
	}

	// Find user role
	var userRole string
	for _, member := range workspace.Members {
		if member.UserID == userID {
			userRole = member.Role
			break
		}
	}

	if userRole == "" {
		return false, nil
	}

	// Only owner and admin can remove members
	return HasPermission(userRole, PermissionMemberRemove), nil
}

// CanCreateSession checks if a user can create sessions
func (wm *WorkspaceManager) CanCreateSession(workspaceID, userID string) (bool, error) {
	return wm.CheckWorkspacePermission(workspaceID, userID, PermissionSessionCreate)
}

// CanCreateTunnel checks if a user can create tunnels
func (wm *WorkspaceManager) CanCreateTunnel(workspaceID, userID string) (bool, error) {
	return wm.CheckWorkspacePermission(workspaceID, userID, PermissionTunnelCreate)
}

// GetUserRole returns the role of a user in a workspace
func (wm *WorkspaceManager) GetUserRole(workspaceID, userID string) (string, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return "", fmt.Errorf("workspace not found: %s", workspaceID)
	}

	for _, member := range workspace.Members {
		if member.UserID == userID {
			return member.Role, nil
		}
	}

	return "", fmt.Errorf("user is not a member")
}

