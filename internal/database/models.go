package database

import (
	"time"

	"gorm.io/gorm"
)

// Container represents a container in the database
type Container struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Image       string    `gorm:"index" json:"image"`
	Status      string    `gorm:"index" json:"status"` // running, stopped, exited, etc.
	Created     time.Time `json:"created"`
	Started     *time.Time `json:"started,omitempty"`
	Stopped     *time.Time `json:"stopped,omitempty"`
	
	// Configuration
	Command     string    `gorm:"type:text" json:"command,omitempty"`
	Env         string    `gorm:"type:text" json:"env,omitempty"` // JSON array of env vars
	Ports       string    `gorm:"type:text" json:"ports,omitempty"` // JSON array of port mappings
	Volumes     string    `gorm:"type:text" json:"volumes,omitempty"` // JSON array of volume mounts
	Network     string    `json:"network,omitempty"`
	
	// Workspace association
	WorkspaceID *string   `gorm:"index" json:"workspaceId,omitempty"`
	
	// Metadata
	Labels      string    `gorm:"type:text" json:"labels,omitempty"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Image represents an image in the database
type Image struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Tag         string    `gorm:"index" json:"tag"`
	Digest      string    `gorm:"index" json:"digest,omitempty"`
	Size        int64     `json:"size"` // Size in bytes
	Created     time.Time `json:"created"`
	
	// Source information
	Registry    string    `json:"registry,omitempty"`
	Repository  string    `gorm:"index" json:"repository,omitempty"`
	
	// Metadata
	Labels      string    `gorm:"type:text" json:"labels,omitempty"` // JSON object
	Description string    `gorm:"type:text" json:"description,omitempty"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Workspace represents a shared runtime workspace
type Workspace struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Status      string    `gorm:"index" json:"status"` // active, sleeping, archived
	OwnerID     string    `gorm:"index" json:"ownerId"`
	
	// Settings
	Settings    string    `gorm:"type:text" json:"settings,omitempty"` // JSON object
	
	// Metadata
	Metadata    string    `gorm:"type:text" json:"metadata,omitempty"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	
	// Relationships
	Members     []WorkspaceMember `gorm:"foreignKey:WorkspaceID" json:"members,omitempty"`
	Invites     []Invite          `gorm:"foreignKey:WorkspaceID" json:"invites,omitempty"`
	Sessions    []Session         `gorm:"foreignKey:WorkspaceID" json:"sessions,omitempty"`
}

// WorkspaceMember represents a member of a workspace
type WorkspaceMember struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID string    `gorm:"index;uniqueIndex:idx_workspace_user" json:"workspaceId"`
	UserID      string    `gorm:"index;uniqueIndex:idx_workspace_user" json:"userId"`
	Role        string    `gorm:"index" json:"role"` // owner, editor, viewer
	JoinedAt    time.Time `json:"joinedAt"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Invite represents a workspace invitation
type Invite struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID string    `gorm:"index" json:"workspaceId"`
	Token       string    `gorm:"uniqueIndex" json:"token"`
	CreatedBy   string    `gorm:"index" json:"createdBy"`
	Role        string    `json:"role"` // editor, viewer
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	AcceptedAt  *time.Time `json:"acceptedAt,omitempty"`
	AcceptedBy  *string   `json:"acceptedBy,omitempty"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Session represents an active workspace session
type Session struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	WorkspaceID string    `gorm:"index" json:"workspaceId"`
	UserID      string    `gorm:"index" json:"userId"`
	ContainerID string    `gorm:"index" json:"containerId,omitempty"`
	Status      string    `gorm:"index" json:"status"` // active, closed
	State       string    `gorm:"type:text" json:"state,omitempty"` // JSON object
	LastActivity time.Time `json:"lastActivity"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ClosedAt    *time.Time `json:"closedAt,omitempty"`
}

// Snapshot represents a snapshot of a container, workspace, or volume
type Snapshot struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	ResourceType string    `gorm:"index" json:"resourceType"` // container, workspace, volume
	ResourceID  string    `gorm:"index" json:"resourceId"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Status      string    `gorm:"index" json:"status"` // creating, ready, failed, deleting
	Size        int64     `json:"size,omitempty"` // Size in bytes
	Data        string    `gorm:"type:text" json:"data,omitempty"` // JSON snapshot data
	CreatedBy   string    `gorm:"index" json:"createdBy"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Deployment represents an orchestrator deployment
type Deployment struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Image       string    `json:"image"`
	Replicas    int       `json:"replicas"`
	Status      string    `gorm:"index" json:"status"` // active, paused, failed
	Strategy    string    `json:"strategy"` // rolling, recreate, canary
	ServiceName string    `gorm:"index" json:"serviceName,omitempty"`
	NetworkName string    `json:"networkName,omitempty"`
	Ports       string    `gorm:"type:text" json:"ports,omitempty"` // JSON array
	Config      string    `gorm:"type:text" json:"config,omitempty"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Node represents a cluster node
type Node struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Address     string    `gorm:"index" json:"address"`
	Status      string    `gorm:"index" json:"status"` // active, inactive, draining
	Resources   string    `gorm:"type:text" json:"resources,omitempty"` // JSON object
	ContainerCount int    `json:"containerCount"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// ContainerGroup represents a container group
type ContainerGroup struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	ParentID    *string   `gorm:"index" json:"parentId,omitempty"`
	Resources   string    `gorm:"type:text" json:"resources,omitempty"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Template represents a stack template
type Template struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Category    string    `gorm:"index" json:"category"` // web, database, microservices, etc.
	Tags        string    `gorm:"type:text" json:"tags,omitempty"` // JSON array
	Config      string    `gorm:"type:text" json:"config"` // JSON template configuration
	IsDefault   bool      `json:"isDefault"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// User represents a user account
type User struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"uniqueIndex" json:"username"`
	Email       string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string   `json:"-"` // Never expose in JSON
	Role        string    `gorm:"index" json:"role"` // admin, user
	Active      bool      `gorm:"index" json:"active"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Team represents a team
type Team struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	OwnerID     string    `gorm:"index" json:"ownerId"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Tenant represents a tenant for multi-tenancy
type Tenant struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Quota       string    `gorm:"type:text" json:"quota,omitempty"` // JSON object
	Active      bool      `gorm:"index" json:"active"`
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Network represents a custom network
type Network struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex" json:"name"`
	Driver      string    `json:"driver"` // bridge, overlay, etc.
	Subnet      string    `json:"subnet,omitempty"`
	Gateway     string    `json:"gateway,omitempty"`
	Config      string    `gorm:"type:text" json:"config,omitempty"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Service represents a service discovery entry
type Service struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Type        string    `json:"type"` // dns, load-balancer, etc.
	Config      string    `gorm:"type:text" json:"config"` // JSON object
	
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

