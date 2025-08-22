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
	coll              *mongo.Collection
	collUserTaskStats *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) ports.TaskRepository {
	return &taskRepo{
		coll:              db.Collection(models.CollectionTasks),
		collUserTaskStats: db.Collection(models.CollectionUserTaskStats),
	}
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

// ดึง steps ทั้งชุด (เปลี่ยนชื่อจาก GetAllStepStatuses -> GetAllStepSteps)
func (r *taskRepo) GetAllStepSteps(ctx context.Context, taskID string) ([]models.TaskWorkflowStep, error) {
	filter := bson.M{"task_id": taskID, "deleted_at": nil}
	proj := bson.M{"applied_workflow.steps": 1, "_id": 0}

	var doc struct {
		AppliedWorkflow struct {
			Steps []models.TaskWorkflowStep `bson:"steps"`
		} `bson:"applied_workflow"`
	}
	if err := r.coll.FindOne(ctx, filter, options.FindOne().SetProjection(proj)).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return doc.AppliedWorkflow.Steps, nil
}

// อัปเดตฟิลด์ของ step เดียว (status optional, notes optional)
func (r *taskRepo) UpdateOneStepFields(ctx context.Context, taskID, stepID string, status *string, notes *string, now time.Time) error {
	filter := bson.M{
		"task_id":                        taskID,
		"deleted_at":                     nil,
		"applied_workflow.steps.step_id": stepID, // ensure exists
	}

	set := bson.M{
		"applied_workflow.steps.$[s].updated_at": now,
	}
	if status != nil {
		set["applied_workflow.steps.$[s].status"] = *status
	}
	if notes != nil {
		set["applied_workflow.steps.$[s].notes"] = *notes
	}

	arrayFilters := []interface{}{
		bson.M{"s.step_id": stepID},
	}

	// ตั้ง timestamp ตามสถานะใหม่ (เฉพาะตอนขออัปเดต status)
	if status != nil {
		switch *status {
		case "in_progress":
			// set started_at เมื่อยังไม่มี
			set["applied_workflow.steps.$[sNotStarted].started_at"] = now
			arrayFilters = append(arrayFilters, bson.M{
				"sNotStarted.step_id": stepID,
				"$or":                 []bson.M{{"sNotStarted.started_at": bson.M{"$exists": false}}, {"sNotStarted.started_at": nil}},
			})
		case "done", "skip":
			// set completed_at เมื่อยังไม่มี
			set["applied_workflow.steps.$[sNotCompleted].completed_at"] = now
			arrayFilters = append(arrayFilters, bson.M{
				"sNotCompleted.step_id": stepID,
				"$or":                   []bson.M{{"sNotCompleted.completed_at": bson.M{"$exists": false}}, {"sNotCompleted.completed_at": nil}},
			})
		}
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	res, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": set}, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *taskRepo) UpdateTaskStatus(ctx context.Context, taskID, status string, now time.Time) error {
	filter := bson.M{"task_id": taskID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": now}}
	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *taskRepo) GetOneUserTaskStatsByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.UserTaskStats, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var stats models.UserTaskStats
	if err := r.collUserTaskStats.FindOne(ctx, filter, opts).Decode(&stats); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &stats, nil
}

func (r *taskRepo) UpsertUserTaskStats(ctx context.Context, stats *models.UserTaskStats) error {
	filter := bson.M{"user_id": stats.UserID}
	update := bson.M{"$set": stats}
	opts := options.Update().SetUpsert(true)
	res, err := r.collUserTaskStats.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
