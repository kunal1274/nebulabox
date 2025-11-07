package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// Invite represents an invitation to join a shared workspace
type Invite struct {
	ID          string            `json:"id"`
	WorkspaceID string            `json:"workspaceId"`
	InviterID   string            `json:"inviterId"`
	InviterName string            `json:"inviterName"`
	Email       string            `json:"email,omitempty"`
	Role        string            `json:"role"` // admin, editor, viewer
	Token       string            `json:"token"`
	Status      string            `json:"status"` // pending, accepted, rejected, expired
	ExpiresAt   time.Time         `json:"expiresAt"`
	CreatedAt   time.Time         `json:"createdAt"`
	AcceptedAt  time.Time         `json:"acceptedAt,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// InviteManager manages workspace invitations
type InviteManager struct {
	invites     map[string]*Invite // inviteID -> Invite
	tokenMap    map[string]string // token -> inviteID
	workspaceInvites map[string][]string // workspaceID -> []inviteID
	mu          sync.RWMutex
}

// NewInviteManager creates a new invite manager
func NewInviteManager() *InviteManager {
	return &InviteManager{
		invites:          make(map[string]*Invite),
		tokenMap:         make(map[string]string),
		workspaceInvites: make(map[string][]string),
	}
}

// CreateInvite creates a new invitation
func (im *InviteManager) CreateInvite(workspaceID, inviterID, inviterName, email, role string, expiresInHours int) (*Invite, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"editor": true,
		"viewer": true,
	}
	if !validRoles[role] {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	token := fmt.Sprintf("inv-%d-%s", time.Now().UnixNano(), generateRandomString(16))
	expiresAt := time.Now().Add(time.Duration(expiresInHours) * time.Hour)

	invite := &Invite{
		ID:          fmt.Sprintf("inv-%d", time.Now().UnixNano()),
		WorkspaceID: workspaceID,
		InviterID:   inviterID,
		InviterName: inviterName,
		Email:       email,
		Role:        role,
		Token:       token,
		Status:      "pending",
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}

	im.invites[invite.ID] = invite
	im.tokenMap[token] = invite.ID
	im.workspaceInvites[workspaceID] = append(im.workspaceInvites[workspaceID], invite.ID)

	return invite, nil
}

// GetInvite retrieves an invite by ID
func (im *InviteManager) GetInvite(inviteID string) (*Invite, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	invite, exists := im.invites[inviteID]
	return invite, exists
}

// GetInviteByToken retrieves an invite by token
func (im *InviteManager) GetInviteByToken(token string) (*Invite, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	inviteID, exists := im.tokenMap[token]
	if !exists {
		return nil, false
	}

	invite, exists := im.invites[inviteID]
	return invite, exists
}

// ListWorkspaceInvites returns all invites for a workspace
func (im *InviteManager) ListWorkspaceInvites(workspaceID string) []*Invite {
	im.mu.RLock()
	defer im.mu.RUnlock()

	inviteIDs, exists := im.workspaceInvites[workspaceID]
	if !exists {
		return []*Invite{}
	}

	var result []*Invite
	for _, id := range inviteIDs {
		if invite, exists := im.invites[id]; exists {
			result = append(result, invite)
		}
	}

	return result
}

// AcceptInvite accepts an invitation
func (im *InviteManager) AcceptInvite(token, userID string) (*Invite, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	inviteID, exists := im.tokenMap[token]
	if !exists {
		return nil, fmt.Errorf("invalid invite token")
	}

	invite, exists := im.invites[inviteID]
	if !exists {
		return nil, fmt.Errorf("invite not found")
	}

	if invite.Status != "pending" {
		return nil, fmt.Errorf("invite is not pending")
	}

	if time.Now().After(invite.ExpiresAt) {
		invite.Status = "expired"
		return nil, fmt.Errorf("invite has expired")
	}

	invite.Status = "accepted"
	invite.AcceptedAt = time.Now()

	return invite, nil
}

// RejectInvite rejects an invitation
func (im *InviteManager) RejectInvite(inviteID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	invite, exists := im.invites[inviteID]
	if !exists {
		return fmt.Errorf("invite not found")
	}

	if invite.Status != "pending" {
		return fmt.Errorf("invite is not pending")
	}

	invite.Status = "rejected"
	return nil
}

// RevokeInvite revokes an invitation
func (im *InviteManager) RevokeInvite(inviteID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	invite, exists := im.invites[inviteID]
	if !exists {
		return fmt.Errorf("invite not found")
	}

	if invite.Status == "accepted" {
		return fmt.Errorf("cannot revoke accepted invite")
	}

	delete(im.invites, inviteID)
	delete(im.tokenMap, invite.Token)

	// Remove from workspace invites
	inviteIDs := im.workspaceInvites[invite.WorkspaceID]
	newIDs := []string{}
	for _, id := range inviteIDs {
		if id != inviteID {
			newIDs = append(newIDs, id)
		}
	}
	im.workspaceInvites[invite.WorkspaceID] = newIDs

	return nil
}

// CleanupExpiredInvites removes expired invites
func (im *InviteManager) CleanupExpiredInvites() {
	im.mu.Lock()
	defer im.mu.Unlock()

	now := time.Now()
	for _, invite := range im.invites {
		if invite.Status == "pending" && now.After(invite.ExpiresAt) {
			invite.Status = "expired"
		}
	}
}

// generateRandomString generates a random string for tokens
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

