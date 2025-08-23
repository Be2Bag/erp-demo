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

type kpiEvaluationRepo struct {
	coll *mongo.Collection
}

func NewKPIEvaluationRepository(db *mongo.Database) ports.KPIEvaluationRepository {
	return &kpiEvaluationRepo{
		coll: db.Collection(models.CollectionKPIEvaluations),
	}
}

func (r *kpiEvaluationRepo) CreateKPIEvaluations(ctx context.Context, kpi models.KPIEvaluation) error {
	_, err := r.coll.InsertOne(ctx, kpi)
	return err
}

func (r *kpiEvaluationRepo) UpdateKPIEvaluationByID(ctx context.Context, evaluationID string, update models.KPIEvaluation) (*models.KPIEvaluation, error) {
	filter := bson.M{"evaluation_id": evaluationID, "deleted_at": bson.M{"$exists": false}}

	set := bson.M{
		"job_id":        update.JobID,
		"task_id":       update.TaskID,
		"kpi_id":        update.KPIID,
		"version":       update.Version,
		"evaluator_id":  update.EvaluatorID,
		"evaluatee_id":  update.EvaluateeID,
		"department_id": update.Department,
		"scores":        update.Scores,
		"total_score":   update.TotalScore,
		"feedback":      update.Feedback,
		"updated_at":    update.UpdatedAt,
	}

	if update.TaskID == "" {
		delete(set, "task_id")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.KPIEvaluation
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *kpiEvaluationRepo) SoftDeleteKPIEvaluationByID(ctx context.Context, evaluationID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"evaluation_id": evaluationID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *kpiEvaluationRepo) GetAllKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.KPIEvaluation, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var kpiEvaluations []*models.KPIEvaluation
	for cursor.Next(ctx) {
		var kpiEvaluation models.KPIEvaluation
		if err := cursor.Decode(&kpiEvaluation); err != nil {
			return nil, err
		}
		kpiEvaluations = append(kpiEvaluations, &kpiEvaluation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return kpiEvaluations, nil
}

func (r *kpiEvaluationRepo) GetOneKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.KPIEvaluation, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var kpiEvaluation models.KPIEvaluation
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&kpiEvaluation); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &kpiEvaluation, nil
}

func (r *kpiEvaluationRepo) GetListKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.KPIEvaluation, int64, error) {

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

	var results []models.KPIEvaluation
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
