package mongodb_repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repositories holds all MongoDB repositories
type Repositories struct {
	ContainerLogs *ContainerLogsRepository
	APIMetrics    *APIMetricsRepository
	SystemMetrics *SystemMetricsRepository
	AuditLogs     *AuditLogsRepository
	BuildLogs     *BuildLogsRepository
}

var (
	mongoReposInstance *Repositories
)

// InitMongoRepositories initializes all MongoDB repositories
func InitMongoRepositories(db *mongo.Database) *Repositories {
	if mongoReposInstance != nil {
		return mongoReposInstance
	}
	
	if db == nil {
		return nil
	}
	
	mongoReposInstance = &Repositories{
		ContainerLogs: NewContainerLogsRepository(db),
		APIMetrics:    NewAPIMetricsRepository(db),
		SystemMetrics: NewSystemMetricsRepository(db),
		AuditLogs:     NewAuditLogsRepository(db),
		BuildLogs:     NewBuildLogsRepository(db),
	}
	
	return mongoReposInstance
}

// GetMongoRepositories returns the MongoDB repositories instance
func GetMongoRepositories() *Repositories {
	return mongoReposInstance
}

// NewMongoRepositories creates new repositories with a specific database
// This is useful for testing
func NewMongoRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		ContainerLogs: NewContainerLogsRepository(db),
		APIMetrics:    NewAPIMetricsRepository(db),
		SystemMetrics: NewSystemMetricsRepository(db),
		AuditLogs:     NewAuditLogsRepository(db),
		BuildLogs:     NewBuildLogsRepository(db),
	}
}

