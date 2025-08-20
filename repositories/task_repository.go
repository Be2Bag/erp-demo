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

type taskRepo struct {
	coll *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) ports.TaskRepository {
	return &taskRepo{coll: db.Collection(models.CollectionTasks)}
}

func (r *taskRepo) CreateTask(ctx context.Context, task models.Tasks) error {
	_, err := r.coll.InsertOne(ctx, task)
	return err
}

func (r *taskRepo) UpdateTaskByID(ctx context.Context, taskID string, update models.Tasks) (*models.Tasks, error) {
	filter := bson.M{"task_id": taskID}
	set := bson.M{
		"project_id":       update.ProjectID,
		"project_name":     update.ProjectName,
		"job_id":           update.JobID,
		"job_name":         update.JobName,
		"description":      update.Description,
		"department":       update.Department,
		"assignee":         update.Assignee,
		"importance":       update.Importance,
		"start_date":       update.StartDate,
		"end_date":         update.EndDate,
		"kpi_id":           update.KPIID,
		"workflow_id":      update.WorkFlowID,
		"applied_workflow": update.AppliedWorkflow,
		"status":           update.Status,
		"updated_at":       time.Now(),
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Tasks
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *taskRepo) SoftDeleteTaskByJobID(ctx context.Context, taskID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"task_id": taskID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *taskRepo) GetAllTaskByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Tasks, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Tasks
	for cursor.Next(ctx) {
		var task models.Tasks
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepo) GetOneTasksByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Tasks, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var task models.Tasks
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) GetListTasksByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Tasks, int64, error) {

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

	var results []models.Tasks
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
