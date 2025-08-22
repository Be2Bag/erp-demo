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

type workFlowRepo struct {
	coll *mongo.Collection
}

func NewWorkFlowRepository(db *mongo.Database) ports.WorkFlowRepository {
	return &workFlowRepo{coll: db.Collection(models.CollectionWorkflowTemplates)}
}

func (r *workFlowRepo) CreateWorkFlowTemplate(ctx context.Context, tmpl *models.WorkFlowTemplate) error {
	_, err := r.coll.InsertOne(ctx, tmpl)
	return err
}

func (r *workFlowRepo) UpdateWorkFlowTemplateByID(ctx context.Context, workflowID string, update models.WorkFlowTemplate) (*models.WorkFlowTemplate, error) {
	filter := bson.M{"workflow_id": workflowID}
	set := bson.M{
		"workflow_name": update.WorkFlowName,
		"workflow_id":   update.WorkFlowID,
		"department_id": update.Department,
		"description":   update.Description,
		"total_hours":   update.TotalHours,
		"steps":         update.Steps,
		"is_active":     update.IsActive,
		"version":       update.Version,
		"created_by":    update.CreatedBy,
		"created_at":    update.CreatedAt,
		"updated_at":    update.UpdatedAt,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.WorkFlowTemplate
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *workFlowRepo) SoftDeleteWorkFlowTemplateByID(ctx context.Context, workflowID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"workflow_id": workflowID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *workFlowRepo) GetAllWorkFlowTemplatesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.WorkFlowTemplate, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*models.WorkFlowTemplate
	for cursor.Next(ctx) {
		var template models.WorkFlowTemplate
		if err := cursor.Decode(&template); err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

func (r *workFlowRepo) GetOneWorkFlowTemplateByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.WorkFlowTemplate, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var template models.WorkFlowTemplate
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&template); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func (r *workFlowRepo) GetListWorkFlowTemplatesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.WorkFlowTemplate, int64, error) {

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

	var results []models.WorkFlowTemplate
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
