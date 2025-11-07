package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB connection and configuration
type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var (
	mongoInstance *MongoDB
)

// InitMongoDB initializes MongoDB connection
func InitMongoDB() (*MongoDB, error) {
	if mongoInstance != nil {
		return mongoInstance, nil
	}

	// Get MongoDB connection string from environment
	uri := os.Getenv("NEBULABOX_MONGODB_URI")
	if uri == "" {
		host := os.Getenv("NEBULABOX_MONGODB_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("NEBULABOX_MONGODB_PORT")
		if port == "" {
			port = "27017"
		}

		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}

	dbname := os.Getenv("NEBULABOX_MONGODB_DB")
	if dbname == "" {
		dbname = "nebulabox"
	}

	// Set client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbname)

	mongoInstance = &MongoDB{
		Client: client,
		DB:     db,
	}

	log.Printf("✅ Connected to MongoDB database: %s", dbname)

	// Set up indexes and TTL
	if err := mongoInstance.setupIndexes(); err != nil {
		log.Printf("⚠️  Warning: failed to set up indexes: %v", err)
	}

	return mongoInstance, nil
}

// GetMongoDB returns the MongoDB instance
func GetMongoDB() *MongoDB {
	return mongoInstance
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.Client.Disconnect(ctx)
	}
	return nil
}

// setupIndexes sets up TTL indexes and other indexes for log collections
func (m *MongoDB) setupIndexes() error {
	ctx := context.Background()

	// Audit logs collection
	auditLogs := m.DB.Collection("audit_logs")
	auditIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(90 * 24 * 3600), // 90 days TTL
		},
		{
			Keys: map[string]interface{}{"workspaceId": 1, "timestamp": -1},
		},
		{
			Keys: map[string]interface{}{"userId": 1, "timestamp": -1},
		},
	}
	if _, err := auditLogs.Indexes().CreateMany(ctx, auditIndexes); err != nil {
		return fmt.Errorf("failed to create audit_logs indexes: %w", err)
	}

	// Container logs collection
	containerLogs := m.DB.Collection("container_logs")
	containerLogIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(30 * 24 * 3600), // 30 days TTL
		},
		{
			Keys: map[string]interface{}{"containerId": 1, "timestamp": -1},
		},
	}
	if _, err := containerLogs.Indexes().CreateMany(ctx, containerLogIndexes); err != nil {
		return fmt.Errorf("failed to create container_logs indexes: %w", err)
	}

	// API metrics collection
	apiMetrics := m.DB.Collection("api_metrics")
	apiMetricsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(7 * 24 * 3600), // 7 days TTL
		},
		{
			Keys: map[string]interface{}{"endpoint": 1, "timestamp": -1},
		},
	}
	if _, err := apiMetrics.Indexes().CreateMany(ctx, apiMetricsIndexes); err != nil {
		return fmt.Errorf("failed to create api_metrics indexes: %w", err)
	}

	// System metrics collection
	systemMetrics := m.DB.Collection("system_metrics")
	systemMetricsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(30 * 24 * 3600), // 30 days TTL
		},
	}
	if _, err := systemMetrics.Indexes().CreateMany(ctx, systemMetricsIndexes); err != nil {
		return fmt.Errorf("failed to create system_metrics indexes: %w", err)
	}

	// Build logs collection
	buildLogs := m.DB.Collection("build_logs")
	buildLogsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(14 * 24 * 3600), // 14 days TTL
		},
		{
			Keys: map[string]interface{}{"imageId": 1, "timestamp": -1},
		},
	}
	if _, err := buildLogs.Indexes().CreateMany(ctx, buildLogsIndexes); err != nil {
		return fmt.Errorf("failed to create build_logs indexes: %w", err)
	}

	// Test runs collection
	testRuns := m.DB.Collection("test_runs")
	testRunsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"timestamp": -1},
			Options: options.Index().SetExpireAfterSeconds(90 * 24 * 3600), // 90 days TTL
		},
		{
			Keys: map[string]interface{}{"suiteId": 1, "timestamp": -1},
		},
	}
	if _, err := testRuns.Indexes().CreateMany(ctx, testRunsIndexes); err != nil {
		return fmt.Errorf("failed to create test_runs indexes: %w", err)
	}

	log.Println("✅ MongoDB indexes set up successfully")
	return nil
}

// HealthCheck checks if the MongoDB connection is healthy
func (m *MongoDB) HealthCheck() error {
	if m.Client == nil {
		return fmt.Errorf("MongoDB client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.Client.Ping(ctx, readpref.Primary())
}

// Collections returns the database collections
func (m *MongoDB) Collections() map[string]*mongo.Collection {
	return map[string]*mongo.Collection{
		"audit_logs":     m.DB.Collection("audit_logs"),
		"container_logs": m.DB.Collection("container_logs"),
		"api_metrics":    m.DB.Collection("api_metrics"),
		"system_metrics": m.DB.Collection("system_metrics"),
		"build_logs":     m.DB.Collection("build_logs"),
		"test_runs":      m.DB.Collection("test_runs"),
	}
}

