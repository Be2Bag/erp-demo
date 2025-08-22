package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type kpiRepo struct {
	coll *mongo.Collection
}

func NewKPIRepository(db *mongo.Database) ports.KPIRepository {
	return &kpiRepo{coll: db.Collection(models.CollectionKPITemplates)}
}

func (r *kpiRepo) CreateKPI(ctx context.Context, kpi models.KPITemplate) error {
	_, err := r.coll.InsertOne(ctx, kpi)
	return err
}

func (r *kpiRepo) UpdateKPIByID(ctx context.Context, kpiID string, update models.KPITemplate) (*models.KPITemplate, error) {
	filter := bson.M{"kpi_id": kpiID}
	set := bson.M{
		"kpi_name":      update.KPIName,
		"department_id": update.Department,
		"total_weight":  update.TotalWeight,
		"items":         update.Items,
		"is_active":     update.IsActive,
		"version":       update.Version,
		"updated_at":    update.UpdatedAt,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.KPITemplate
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *kpiRepo) SoftDeleteKPIByID(ctx context.Context, kpiID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"kpi_id": kpiID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *kpiRepo) GetAllKPIByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.KPITemplate, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var kpiTemplates []*models.KPITemplate
	for cursor.Next(ctx) {
		var kpiTemplate models.KPITemplate
		if err := cursor.Decode(&kpiTemplate); err != nil {
			return nil, err
		}
		kpiTemplates = append(kpiTemplates, &kpiTemplate)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return kpiTemplates, nil
}

func (r *kpiRepo) GetOneKPIByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.KPITemplate, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var kpiTemplate models.KPITemplate
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&kpiTemplate); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &kpiTemplate, nil
}

func (r *kpiRepo) GetListKPIByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.KPITemplate, int64, error) {

	findOpts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(limit)

	if projection != nil {
		findOpts.SetProjection(projection)
	}

	cur, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find: %w", err)
	}
	defer cur.Close(ctx)

	var results []models.KPITemplate
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
