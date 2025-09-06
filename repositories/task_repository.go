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
		"department_id":    update.Department,
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

func (r *taskRepo) SoftDeleteTaskByID(ctx context.Context, taskID string) error {
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

func (r *taskRepo) UpdateTaskStatus(ctx context.Context, taskID, status, stepName string, now time.Time) error {
	filter := bson.M{"task_id": taskID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"status": status, "step_name": stepName, "updated_at": now}}
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

func (r *taskRepo) GetAllUserTaskStatsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.UserTaskStats, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.collUserTaskStats.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userTaskStats []*models.UserTaskStats
	for cursor.Next(ctx) {
		var userTaskStat models.UserTaskStats
		if err := cursor.Decode(&userTaskStat); err != nil {
			return nil, err
		}
		userTaskStats = append(userTaskStats, &userTaskStat)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userTaskStats, nil
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

	// --- เพิ่มส่วนสำหรับ unset เมื่อ revert ---
	unset := bson.M{}
	needAnyFilter := false

	arrayFilters := []interface{}{
		bson.M{"s.step_id": stepID}, // ตัวหลัก
	}

	// ตั้ง/เคลียร์ timestamp ตามสถานะใหม่ (เฉพาะตอนขออัปเดต status)
	if status != nil {
		switch *status {
		case "in_progress":
			// ถ้ายังไม่เคย start ให้ set started_at
			set["applied_workflow.steps.$[sNotStarted].started_at"] = now
			arrayFilters = append(arrayFilters, bson.M{
				"sNotStarted.step_id": stepID,
				"$or": []bson.M{
					{"sNotStarted.started_at": bson.M{"$exists": false}},
					{"sNotStarted.started_at": nil},
				},
			})
			// ถ้าเคย done/skip มาก่อน ให้เคลียร์ completed_at
			unset["applied_workflow.steps.$[sAny].completed_at"] = ""
			needAnyFilter = true

		case "done", "skip":
			// ถ้ายังไม่เคย complete ให้ set completed_at
			set["applied_workflow.steps.$[sNotCompleted].completed_at"] = now
			arrayFilters = append(arrayFilters, bson.M{
				"sNotCompleted.step_id": stepID,
				"$or": []bson.M{
					{"sNotCompleted.completed_at": bson.M{"$exists": false}},
					{"sNotCompleted.completed_at": nil},
				},
			})
			// (ออปชัน) จะบังคับ set started_at ถ้ายังไม่มี ก็ทำเพิ่มได้:
			// set["applied_workflow.steps.$[sNotStarted].started_at"] = now
			// arrayFilters = append(arrayFilters, bson.M{
			// 	"sNotStarted.step_id": stepID,
			// 	"$or": []bson.M{
			// 		{"sNotStarted.started_at": bson.M{"$exists": false}},
			// 		{"sNotStarted.started_at": nil},
			// 	},
			// })

		case "todo":
			// กลับไปยังไม่เริ่ม → เคลียร์ทั้ง started_at และ completed_at
			unset["applied_workflow.steps.$[sAny].started_at"] = ""
			unset["applied_workflow.steps.$[sAny].completed_at"] = ""
			needAnyFilter = true
		}
	}

	if needAnyFilter {
		arrayFilters = append(arrayFilters, bson.M{"sAny.step_id": stepID})
	}

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
	res, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

//<===============================================================================================>

// package repositories
func (r *taskRepo) ReplaceTaskByID(ctx context.Context, taskID string, doc *models.Tasks) (*models.Tasks, error) {
	filter := bson.M{"task_id": taskID}
	opts := options.FindOneAndReplace().SetReturnDocument(options.After)
	var out models.Tasks
	if err := r.coll.FindOneAndReplace(ctx, filter, doc, opts).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *taskRepo) UpdateManyTaskByFilter(ctx context.Context, filter interface{}, update models.Tasks) (int64, error) {
	result, err := r.coll.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
