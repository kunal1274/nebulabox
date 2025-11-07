package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Collections defines MongoDB collection names and accessors
type Collections struct {
	AuditLogs     *mongo.Collection
	ContainerLogs *mongo.Collection
	APIMetrics    *mongo.Collection
	SystemMetrics *mongo.Collection
	BuildLogs     *mongo.Collection
	TestRuns      *mongo.Collection
}

// GetCollections returns all MongoDB collections
func GetCollections() *Collections {
	if mongoInstance == nil || mongoInstance.DB == nil {
		return nil
	}

	return &Collections{
		AuditLogs:     mongoInstance.DB.Collection("audit_logs"),
		ContainerLogs: mongoInstance.DB.Collection("container_logs"),
		APIMetrics:    mongoInstance.DB.Collection("api_metrics"),
		SystemMetrics: mongoInstance.DB.Collection("system_metrics"),
		BuildLogs:     mongoInstance.DB.Collection("build_logs"),
		TestRuns:      mongoInstance.DB.Collection("test_runs"),
	}
}

