package repositories

import (
	"github.com/nebulabox/nebulabox/internal/database"
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	Container *ContainerRepository
	Image     *ImageRepository
	Workspace *WorkspaceRepository
}

var (
	repositoriesInstance *Repositories
)

// InitRepositories initializes all repositories
func InitRepositories() (*Repositories, error) {
	if repositoriesInstance != nil {
		return repositoriesInstance, nil
	}

	postgres := database.GetPostgreSQL()
	if postgres == nil || postgres.DB == nil {
		// Return nil repositories if database not available
		// This allows graceful fallback to in-memory storage
		return nil, nil
	}

	repositoriesInstance = &Repositories{
		Container: NewContainerRepository(postgres.DB),
		Image:     NewImageRepository(postgres.DB),
		Workspace: NewWorkspaceRepository(postgres.DB),
	}

	return repositoriesInstance, nil
}

// GetRepositories returns the repositories instance
func GetRepositories() *Repositories {
	return repositoriesInstance
}

// NewRepositories creates new repositories with a specific database connection
// This is useful for testing
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Container: NewContainerRepository(db),
		Image:     NewImageRepository(db),
		Workspace: NewWorkspaceRepository(db),
	}
}

