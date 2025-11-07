package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nebulabox/nebulabox/internal/database"
	"github.com/nebulabox/nebulabox/internal/shareruntime"
	"gorm.io/gorm"
)

// WorkspaceRepository handles workspace database operations
type WorkspaceRepository struct {
	db *gorm.DB
}

// NewWorkspaceRepository creates a new workspace repository
func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

// Create creates a new workspace
func (r *WorkspaceRepository) Create(workspace *shareruntime.Workspace) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	dbWorkspace := &database.Workspace{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Description: workspace.Description,
		Status:      workspace.Status, // Status is already a string
		OwnerID:     workspace.OwnerID,
	}

	// Serialize settings (workspace.Settings is not a pointer, it's a value)
	var settingsJSON []byte
	var err error
	if workspace.Settings.AllowGuestAccess || workspace.Settings.MaxMembers > 0 || len(workspace.Settings.Permissions) > 0 {
		settingsJSON, err = json.Marshal(workspace.Settings)
		if err == nil {
			dbWorkspace.Settings = string(settingsJSON)
		}
	}

	// Serialize metadata
	if len(workspace.Metadata) > 0 {
		metadataJSON, err := json.Marshal(workspace.Metadata)
		if err == nil {
			dbWorkspace.Metadata = string(metadataJSON)
		}
	}

	now := time.Now()
	dbWorkspace.CreatedAt = now
	dbWorkspace.UpdatedAt = now

	return r.db.Create(dbWorkspace).Error
}

// Get retrieves a workspace by ID
func (r *WorkspaceRepository) Get(id string) (*shareruntime.Workspace, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbWorkspace database.Workspace
	if err := r.db.Preload("Members").Preload("Invites").Preload("Sessions").
		Where("id = ?", id).First(&dbWorkspace).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toWorkspace(&dbWorkspace), nil
}

// List retrieves all workspaces for a user
func (r *WorkspaceRepository) List(userID string) ([]*shareruntime.Workspace, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbWorkspaces []database.Workspace
	query := r.db

	// If userID provided, filter by membership
	if userID != "" {
		query = query.Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
			Where("workspace_members.user_id = ? OR owner_id = ?", userID, userID)
	}

	if err := query.Find(&dbWorkspaces).Error; err != nil {
		return nil, err
	}

	workspaces := make([]*shareruntime.Workspace, len(dbWorkspaces))
	for i, dbWorkspace := range dbWorkspaces {
		workspaces[i] = r.toWorkspace(&dbWorkspace)
	}

	return workspaces, nil
}

// UpdateStatus updates workspace status
func (r *WorkspaceRepository) UpdateStatus(id string, status string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Model(&database.Workspace{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// Delete soft deletes a workspace
func (r *WorkspaceRepository) Delete(id string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Delete(&database.Workspace{}, "id = ?", id).Error
}

// AddMember adds a member to a workspace
func (r *WorkspaceRepository) AddMember(workspaceID, userID, role string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	member := &database.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return r.db.Create(member).Error
}

// RemoveMember removes a member from a workspace
func (r *WorkspaceRepository) RemoveMember(workspaceID, userID string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Delete(&database.WorkspaceMember{}).Error
}

// toWorkspace converts database.Workspace to shareruntime.Workspace
func (r *WorkspaceRepository) toWorkspace(dbWorkspace *database.Workspace) *shareruntime.Workspace {
	workspace := &shareruntime.Workspace{
		ID:          dbWorkspace.ID,
		Name:        dbWorkspace.Name,
		Description: dbWorkspace.Description,
		Status:      dbWorkspace.Status, // Status is already a string in shareruntime.Workspace
		OwnerID:     dbWorkspace.OwnerID,
		CreatedAt:   dbWorkspace.CreatedAt,
		UpdatedAt:   dbWorkspace.UpdatedAt,
	}

	// Parse settings
	if dbWorkspace.Settings != "" {
		var settings shareruntime.WorkspaceSettings
		if err := json.Unmarshal([]byte(dbWorkspace.Settings), &settings); err == nil {
			workspace.Settings = settings
		}
	}

	// Parse metadata
	if dbWorkspace.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(dbWorkspace.Metadata), &metadata); err == nil {
			workspace.Metadata = metadata
		}
	}

	// Convert members
	if len(dbWorkspace.Members) > 0 {
		workspace.Members = make([]shareruntime.WorkspaceMember, len(dbWorkspace.Members))
		for i, dbMember := range dbWorkspace.Members {
			workspace.Members[i] = shareruntime.WorkspaceMember{
				UserID:   dbMember.UserID,
				Username: dbMember.UserID, // Use UserID as username if not available
				Role:     dbMember.Role,
				JoinedAt: dbMember.JoinedAt,
			}
		}
	}

	return workspace
}

