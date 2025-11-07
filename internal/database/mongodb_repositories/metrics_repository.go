package mongodb_repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIMetricEntry represents an API metric entry in MongoDB
type APIMetricEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	Endpoint    string             `bson:"endpoint" json:"endpoint"`
	Method      string             `bson:"method" json:"method"`
	StatusCode  int                `bson:"statusCode" json:"statusCode"`
	DurationMs  int64              `bson:"durationMs" json:"durationMs"`
	RequestSize int64              `bson:"requestSize,omitempty" json:"requestSize,omitempty"`
	ResponseSize int64             `bson:"responseSize,omitempty" json:"responseSize,omitempty"`
	Error       string             `bson:"error,omitempty" json:"error,omitempty"`
	UserID      string             `bson:"userId,omitempty" json:"userId,omitempty"`
}

// APIMetricsRepository handles API metrics in MongoDB
type APIMetricsRepository struct {
	collection *mongo.Collection
}

// NewAPIMetricsRepository creates a new API metrics repository
func NewAPIMetricsRepository(db *mongo.Database) *APIMetricsRepository {
	return &APIMetricsRepository{
		collection: db.Collection("api_metrics"),
	}
}

// Insert inserts a metric entry
func (r *APIMetricsRepository) Insert(ctx context.Context, entry *APIMetricEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// InsertBatch inserts multiple metric entries efficiently
func (r *APIMetricsRepository) InsertBatch(ctx context.Context, entries []*APIMetricEntry) error {
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

// GetMetrics retrieves metrics with filters
func (r *APIMetricsRepository) GetMetrics(ctx context.Context, filters APIMetricFilters) ([]*APIMetricEntry, error) {
	query := bson.M{}
	
	if filters.Endpoint != "" {
		query["endpoint"] = filters.Endpoint
	}
	
	if filters.Method != "" {
		query["method"] = filters.Method
	}
	
	if filters.StatusCode > 0 {
		query["statusCode"] = filters.StatusCode
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
		limit = 1000
	}
	if limit > 10000 {
		limit = 10000
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*APIMetricEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// APIMetricFilters represents filters for metric queries
type APIMetricFilters struct {
	Endpoint   string
	Method     string
	StatusCode int
	Since      time.Time
	Until      time.Time
	Limit      int
}

// GetAggregatedMetrics returns aggregated metrics (e.g., per endpoint)
func (r *APIMetricsRepository) GetAggregatedMetrics(ctx context.Context, since time.Time, until time.Time) (map[string]interface{}, error) {
	matchStage := bson.M{
		"timestamp": bson.M{
			"$gte": since,
			"$lte": until,
		},
	}
	
	groupStage := bson.M{
		"_id": "$endpoint",
		"count": bson.M{"$sum": 1},
		"avgDuration": bson.M{"$avg": "$durationMs"},
		"minDuration": bson.M{"$min": "$durationMs"},
		"maxDuration": bson.M{"$max": "$durationMs"},
		"errorCount": bson.M{
			"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$ne": []interface{}{"$error", ""}},
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
	
	results := make(map[string]interface{})
	results["endpoints"] = []bson.M{}
	
	var endpointStats []bson.M
	if err := cursor.All(ctx, &endpointStats); err != nil {
		return nil, err
	}
	
	results["endpoints"] = endpointStats
	return results, nil
}

// SystemMetricEntry represents a system metric entry in MongoDB
type SystemMetricEntry struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp         time.Time         `bson:"timestamp" json:"timestamp"`
	CPUUsage          float64           `bson:"cpuUsage" json:"cpuUsage"`
	MemoryUsage       float64           `bson:"memoryUsage" json:"memoryUsage"`
	DiskUsage         float64           `bson:"diskUsage" json:"diskUsage"`
	ContainersRunning int               `bson:"containersRunning" json:"containersRunning"`
	ContainersTotal   int               `bson:"containersTotal" json:"containersTotal"`
	NetworkInBytes     int64             `bson:"networkInBytes,omitempty" json:"networkInBytes,omitempty"`
	NetworkOutBytes    int64             `bson:"networkOutBytes,omitempty" json:"networkOutBytes,omitempty"`
}

// SystemMetricsRepository handles system metrics in MongoDB
type SystemMetricsRepository struct {
	collection *mongo.Collection
}

// NewSystemMetricsRepository creates a new system metrics repository
func NewSystemMetricsRepository(db *mongo.Database) *SystemMetricsRepository {
	return &SystemMetricsRepository{
		collection: db.Collection("system_metrics"),
	}
}

// Insert inserts a system metric entry
func (r *SystemMetricsRepository) Insert(ctx context.Context, entry *SystemMetricEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// GetMetrics retrieves system metrics within a time range
func (r *SystemMetricsRepository) GetMetrics(ctx context.Context, since time.Time, until time.Time, limit int) ([]*SystemMetricEntry, error) {
	query := bson.M{
		"timestamp": bson.M{
			"$gte": since,
			"$lte": until,
		},
	}
	
	if limit <= 0 {
		limit = 1000
	}
	if limit > 10000 {
		limit = 10000
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: 1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entries []*SystemMetricEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// GetLatest gets the latest system metric
func (r *SystemMetricsRepository) GetLatest(ctx context.Context) (*SystemMetricEntry, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	
	var entry SystemMetricEntry
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	
	return &entry, nil
}

