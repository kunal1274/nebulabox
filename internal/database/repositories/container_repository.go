package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nebulabox/nebulabox/internal/containerd"
	"github.com/nebulabox/nebulabox/internal/database"
	"gorm.io/gorm"
)

// ContainerStorageOptions contains additional container data to store
type ContainerStorageOptions struct {
	Ports       []string          `json:"ports,omitempty"`
	Env         []string          `json:"env,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Network     string            `json:"network,omitempty"`
	Command     string            `json:"command,omitempty"`
	WorkspaceID string            `json:"workspaceId,omitempty"`
}

// ContainerRepository handles container database operations
type ContainerRepository struct {
	db *gorm.DB
}

// NewContainerRepository creates a new container repository
func NewContainerRepository(db *gorm.DB) *ContainerRepository {
	return &ContainerRepository{db: db}
}

// CreateOrUpdate creates or updates a container in the database
// Additional fields (Ports, Env, Volumes, Labels, Network) can be passed via options
func (r *ContainerRepository) CreateOrUpdate(container *containerd.Container, opts *ContainerStorageOptions) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Convert containerd.Container to database.Container
	dbContainer := &database.Container{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Status:  container.Status,
		Created: container.Created,
	}

	// Apply additional options if provided
	if opts != nil {
		dbContainer.Network = opts.Network
		if opts.WorkspaceID != "" {
			dbContainer.WorkspaceID = &opts.WorkspaceID
		}

		// Parse ports if available
		if len(opts.Ports) > 0 {
			portsJSON, err := json.Marshal(opts.Ports)
			if err == nil {
				dbContainer.Ports = string(portsJSON)
			}
		}

		// Parse environment variables if available
		if len(opts.Env) > 0 {
			envJSON, err := json.Marshal(opts.Env)
			if err == nil {
				dbContainer.Env = string(envJSON)
			}
		}

		// Parse volumes if available
		if len(opts.Volumes) > 0 {
			volumesJSON, err := json.Marshal(opts.Volumes)
			if err == nil {
				dbContainer.Volumes = string(volumesJSON)
			}
		}

		// Parse labels if available
		if len(opts.Labels) > 0 {
			labelsJSON, err := json.Marshal(opts.Labels)
			if err == nil {
				dbContainer.Labels = string(labelsJSON)
			}
		}

		// Parse command if available
		if opts.Command != "" {
			dbContainer.Command = opts.Command
		}
	}

	// Set timestamps
	now := time.Now()
	if dbContainer.CreatedAt.IsZero() {
		dbContainer.CreatedAt = now
	}
	dbContainer.UpdatedAt = now

	// Update status timestamps
	if container.Status == "running" && dbContainer.Started == nil {
		dbContainer.Started = &now
	} else if (container.Status == "stopped" || container.Status == "exited") && dbContainer.Stopped == nil {
		dbContainer.Stopped = &now
	}

	// Use Save to create or update
	return r.db.Save(dbContainer).Error
}

// Get retrieves a container by ID
func (r *ContainerRepository) Get(id string) (*containerd.Container, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbContainer database.Container
	if err := r.db.Where("id = ?", id).First(&dbContainer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toContainerdContainer(&dbContainer), nil
}

// List retrieves all containers matching the criteria
func (r *ContainerRepository) List(all bool, workspaceID string) ([]*containerd.Container, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbContainers []database.Container
	query := r.db

	// Filter by status if not all
	if !all {
		query = query.Where("status = ?", "running")
	}

	// Filter by workspace if specified
	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	}

	if err := query.Find(&dbContainers).Error; err != nil {
		return nil, err
	}

	containers := make([]*containerd.Container, len(dbContainers))
	for i, dbContainer := range dbContainers {
		containers[i] = r.toContainerdContainer(&dbContainer)
	}

	return containers, nil
}

// UpdateStatus updates container status
func (r *ContainerRepository) UpdateStatus(id, status string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	now := time.Now()
	if status == "running" {
		updates["started"] = &now
	} else if status == "stopped" || status == "exited" {
		updates["stopped"] = &now
	}

	return r.db.Model(&database.Container{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a container
func (r *ContainerRepository) Delete(id string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Delete(&database.Container{}, "id = ?", id).Error
}

// AssociateWorkspace associates a container with a workspace
func (r *ContainerRepository) AssociateWorkspace(containerID, workspaceID string) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.Model(&database.Container{}).
		Where("id = ?", containerID).
		Update("workspace_id", workspaceID).Error
}

// toContainerdContainer converts database.Container to containerd.Container
// Note: containerd.Container only has basic fields, additional data is stored separately
func (r *ContainerRepository) toContainerdContainer(dbContainer *database.Container) *containerd.Container {
	return &containerd.Container{
		ID:      dbContainer.ID,
		Name:    dbContainer.Name,
		Image:   dbContainer.Image,
		Status:  dbContainer.Status,
		Created: dbContainer.Created,
	}
}

// GetStorageOptions retrieves additional container storage options
func (r *ContainerRepository) GetStorageOptions(id string) (*ContainerStorageOptions, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var dbContainer database.Container
	if err := r.db.Where("id = ?", id).First(&dbContainer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	opts := &ContainerStorageOptions{
		Network:     dbContainer.Network,
		WorkspaceID: "",
		Command:     dbContainer.Command,
	}

	if dbContainer.WorkspaceID != nil {
		opts.WorkspaceID = *dbContainer.WorkspaceID
	}

	// Parse ports
	if dbContainer.Ports != "" {
		if err := json.Unmarshal([]byte(dbContainer.Ports), &opts.Ports); err != nil {
			// Log error but continue
		}
	}

	// Parse environment variables
	if dbContainer.Env != "" {
		if err := json.Unmarshal([]byte(dbContainer.Env), &opts.Env); err != nil {
			// Log error but continue
		}
	}

	// Parse volumes
	if dbContainer.Volumes != "" {
		if err := json.Unmarshal([]byte(dbContainer.Volumes), &opts.Volumes); err != nil {
			// Log error but continue
		}
	}

	// Parse labels
	if dbContainer.Labels != "" {
		if err := json.Unmarshal([]byte(dbContainer.Labels), &opts.Labels); err != nil {
			// Log error but continue
		}
	}

	return opts, nil
}

