package repositories

import (
	"context"
	"errors"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auditLogRepo struct {
	coll *mongo.Collection
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *mongo.Database) ports.AuditLogRepository {
	coll := db.Collection(models.CollectionAuditLogs)

	// Create indexes for better query performance
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "resource", Value: 1}},
			Options: options.Index().SetName("idx_resource"),
		},
		{
			Keys:    bson.D{{Key: "action", Value: 1}},
			Options: options.Index().SetName("idx_action"),
		},
		{
			Keys:    bson.D{{Key: "log_id", Value: 1}},
			Options: options.Index().SetName("idx_log_id").SetUnique(true),
		},
	}

	// Create indexes in the background (ignore errors if already exist)
	_, _ = coll.Indexes().CreateMany(context.Background(), indexModels)

	return &auditLogRepo{coll: coll}
}

// Create inserts a new audit log record
func (r *auditLogRepo) Create(ctx context.Context, log *models.AuditLog) error {
	_, err := r.coll.InsertOne(ctx, log)
	return err
}

// GetByFilter retrieves audit logs with pagination and sorting
func (r *auditLogRepo) GetByFilter(ctx context.Context, filter interface{}, page, limit int, sortBy, sortOrder string) ([]*models.AuditLog, int64, error) {
	// Count total documents
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Set default sort
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDirection := -1 // desc
	if sortOrder == "asc" {
		sortDirection = 1
	}

	// Calculate skip
	skip := (page - 1) * limit

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortDirection}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	for cursor.Next(ctx) {
		var log models.AuditLog
		if err := cursor.Decode(&log); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByID retrieves a single audit log by its log_id
func (r *auditLogRepo) GetByID(ctx context.Context, logID string) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.coll.FindOne(ctx, bson.M{"log_id": logID}).Decode(&log)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}
