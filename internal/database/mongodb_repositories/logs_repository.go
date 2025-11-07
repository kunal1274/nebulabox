package mongodb_repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ContainerLogEntry represents a container log entry in MongoDB
type ContainerLogEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContainerID string             `bson:"containerId" json:"containerId"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	Level       string             `bson:"level" json:"level"`
	Message     string             `bson:"message" json:"message"`
	Source      string             `bson:"source,omitempty" json:"source,omitempty"` // stdout, stderr
}

// ContainerLogsRepository handles container logs in MongoDB
type ContainerLogsRepository struct {
	collection *mongo.Collection
}

// NewContainerLogsRepository creates a new container logs repository
func NewContainerLogsRepository(db *mongo.Database) *ContainerLogsRepository {
	return &ContainerLogsRepository{
		collection: db.Collection("container_logs"),
	}
}

// Insert inserts a log entry
func (r *ContainerLogsRepository) Insert(ctx context.Context, entry *ContainerLogEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// InsertBatch inserts multiple log entries efficiently
func (r *ContainerLogsRepository) InsertBatch(ctx context.Context, entries []*ContainerLogEntry) error {
	if len(entries) == 0 {
		return nil
	}
	
	docs := make([]interface{}, len(entries))
	for i, entry := range entries {
		if entry.Timestamp.IsZero() {
			entry.Timestamp = time.Now()
		}
		docs[i] = entry
	}
	
	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

// Search searches logs with filters
func (r *ContainerLogsRepository) Search(ctx context.Context, filters ContainerLogFilters) ([]*ContainerLogEntry, error) {
	query := bson.M{}
	
	if filters.ContainerID != "" {
		query["containerId"] = filters.ContainerID
	}
	
	if filters.Level != "" {
		query["level"] = filters.Level
	}
	
	if filters.Query != "" {
		query["message"] = bson.M{"$regex": filters.Query, "$options": "i"}
	}
	
	if !filters.Since.IsZero() {
		query["timestamp"] = bson.M{"$gte": filters.Since}
	}
	
	if !filters.Until.IsZero() {
		if query["timestamp"] == nil {
			query["timestamp"] = bson.M{}
		}
		query["timestamp"].(bson.M)["$lte"] = filters.Until
	}
	
	limit := filters.Limit
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*ContainerLogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// ContainerLogFilters represents filters for log search
type ContainerLogFilters struct {
	ContainerID string
	Level       string
	Query       string
	Since       time.Time
	Until       time.Time
	Limit       int
}

// GetLatest gets the latest logs for a container
func (r *ContainerLogsRepository) GetLatest(ctx context.Context, containerID string, limit int) ([]*ContainerLogEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{"containerId": containerID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*ContainerLogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// Count returns the count of logs matching filters
func (r *ContainerLogsRepository) Count(ctx context.Context, filters ContainerLogFilters) (int64, error) {
	query := bson.M{}
	
	if filters.ContainerID != "" {
		query["containerId"] = filters.ContainerID
	}
	
	if filters.Level != "" {
		query["level"] = filters.Level
	}
	
	if !filters.Since.IsZero() {
		query["timestamp"] = bson.M{"$gte": filters.Since}
	}
	
	if !filters.Until.IsZero() {
		if query["timestamp"] == nil {
			query["timestamp"] = bson.M{}
		}
		query["timestamp"].(bson.M)["$lte"] = filters.Until
	}
	
	return r.collection.CountDocuments(ctx, query)
}

// DeleteOld deletes logs older than the specified duration
func (r *ContainerLogsRepository) DeleteOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result, err := r.collection.DeleteMany(ctx, bson.M{"timestamp": bson.M{"$lt": cutoff}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

