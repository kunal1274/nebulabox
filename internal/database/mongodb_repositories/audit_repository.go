package mongodb_repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditLogEntry represents an audit log entry in MongoDB
type AuditLogEntry struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	WorkspaceID  string             `bson:"workspaceId,omitempty" json:"workspaceId,omitempty"`
	UserID       string             `bson:"userId" json:"userId"`
	Username     string             `bson:"username,omitempty" json:"username,omitempty"`
	Action       string             `bson:"action" json:"action"`
	ResourceType string             `bson:"resourceType,omitempty" json:"resourceType,omitempty"`
	ResourceID   string             `bson:"resourceId,omitempty" json:"resourceId,omitempty"`
	Success      bool               `bson:"success" json:"success"`
	Message      string             `bson:"message,omitempty" json:"message,omitempty"`
	Details      map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
	IPAddress    string             `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	UserAgent    string             `bson:"userAgent,omitempty" json:"userAgent,omitempty"`
}

// AuditLogsRepository handles audit logs in MongoDB
type AuditLogsRepository struct {
	collection *mongo.Collection
}

// NewAuditLogsRepository creates a new audit logs repository
func NewAuditLogsRepository(db *mongo.Database) *AuditLogsRepository {
	return &AuditLogsRepository{
		collection: db.Collection("audit_logs"),
	}
}

// Insert inserts an audit log entry
func (r *AuditLogsRepository) Insert(ctx context.Context, entry *AuditLogEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// InsertBatch inserts multiple audit log entries efficiently
func (r *AuditLogsRepository) InsertBatch(ctx context.Context, entries []*AuditLogEntry) error {
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

// GetLogs retrieves audit logs with filters
func (r *AuditLogsRepository) GetLogs(ctx context.Context, filters AuditLogFilters) ([]*AuditLogEntry, error) {
	query := bson.M{}
	
	if filters.WorkspaceID != "" {
		query["workspaceId"] = filters.WorkspaceID
	}
	
	if filters.UserID != "" {
		query["userId"] = filters.UserID
	}
	
	if filters.Action != "" {
		query["action"] = filters.Action
	}
	
	if filters.ResourceType != "" {
		query["resourceType"] = filters.ResourceType
	}
	
	if filters.ResourceID != "" {
		query["resourceId"] = filters.ResourceID
	}
	
	if filters.Success != nil {
		query["success"] = *filters.Success
	}
	
	if !filters.StartTime.IsZero() || !filters.EndTime.IsZero() {
		timeRange := bson.M{}
		if !filters.StartTime.IsZero() {
			timeRange["$gte"] = filters.StartTime
		}
		if !filters.EndTime.IsZero() {
			timeRange["$lte"] = filters.EndTime
		}
		query["timestamp"] = timeRange
	}
	
	limit := filters.Limit
	if limit <= 0 {
		limit = 100
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
	
	var entries []*AuditLogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// AuditLogFilters represents filters for audit log queries
type AuditLogFilters struct {
	WorkspaceID  string
	UserID       string
	Action       string
	ResourceType string
	ResourceID   string
	Success      *bool
	StartTime    time.Time
	EndTime      time.Time
	Limit        int
}

// GetUserLogs retrieves audit logs for a specific user
func (r *AuditLogsRepository) GetUserLogs(ctx context.Context, userID string, limit int) ([]*AuditLogEntry, error) {
	filters := AuditLogFilters{
		UserID: userID,
		Limit:  limit,
	}
	return r.GetLogs(ctx, filters)
}

// GetStats returns aggregated statistics
func (r *AuditLogsRepository) GetStats(ctx context.Context, workspaceID string, startTime time.Time, endTime time.Time) (map[string]interface{}, error) {
	matchStage := bson.M{
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	
	if workspaceID != "" {
		matchStage["workspaceId"] = workspaceID
	}
	
	groupStage := bson.M{
		"_id": "$action",
		"count": bson.M{"$sum": 1},
		"successCount": bson.M{
			"$sum": bson.M{
				"$cond": []interface{}{
					"$success",
					1,
					0,
				},
			},
		},
		"failureCount": bson.M{
			"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$success", false}},
					1,
					0,
				},
			},
		},
	}
	
	pipeline := []bson.M{
		{"$match": matchStage},
		{"$group": groupStage},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var stats []bson.M
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, err
	}
	
	result := map[string]interface{}{
		"workspaceId": workspaceID,
		"startTime":   startTime,
		"endTime":     endTime,
		"actions":     stats,
	}
	
	return result, nil
}

