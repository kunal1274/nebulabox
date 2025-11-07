package mongodb_repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BuildLogEntry represents a build log entry in MongoDB
type BuildLogEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ImageID     string             `bson:"imageId" json:"imageId"`
	Tag         string             `bson:"tag" json:"tag"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	Line        string             `bson:"line" json:"line"`
	Level       string             `bson:"level,omitempty" json:"level,omitempty"` // info, error, warning
	Step        string             `bson:"step,omitempty" json:"step,omitempty"`
	Order       int                `bson:"order" json:"order"` // Order in the build log sequence
}

// BuildLogsRepository handles build logs in MongoDB
type BuildLogsRepository struct {
	collection *mongo.Collection
}

// NewBuildLogsRepository creates a new build logs repository
func NewBuildLogsRepository(db *mongo.Database) *BuildLogsRepository {
	return &BuildLogsRepository{
		collection: db.Collection("build_logs"),
	}
}

// Insert inserts a build log entry
func (r *BuildLogsRepository) Insert(ctx context.Context, entry *BuildLogEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// InsertBatch inserts multiple build log entries efficiently
func (r *BuildLogsRepository) InsertBatch(ctx context.Context, entries []*BuildLogEntry) error {
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

// GetBuildLogs retrieves build logs for an image
func (r *BuildLogsRepository) GetBuildLogs(ctx context.Context, imageID string, limit int) ([]*BuildLogEntry, error) {
	query := bson.M{"imageId": imageID}
	
	if limit <= 0 {
		limit = 1000
	}
	if limit > 10000 {
		limit = 10000
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "order", Value: 1}, {Key: "timestamp", Value: 1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*BuildLogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// GetBuildLogsByTag retrieves build logs by tag
func (r *BuildLogsRepository) GetBuildLogsByTag(ctx context.Context, tag string, limit int) ([]*BuildLogEntry, error) {
	query := bson.M{"tag": tag}
	
	if limit <= 0 {
		limit = 1000
	}
	if limit > 10000 {
		limit = 10000
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "order", Value: 1}, {Key: "timestamp", Value: 1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*BuildLogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// DeleteBuildLogs deletes build logs for an image
func (r *BuildLogsRepository) DeleteBuildLogs(ctx context.Context, imageID string) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{"imageId": imageID})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

