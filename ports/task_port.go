package ports

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TaskService interface {
	GetListTasks(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department_id string, sortBy string, sortOrder string) (dto.Pagination, error)
	CreateTask(ctx context.Context, createTask dto.CreateTaskRequest, claims *dto.JWTClaims) error
	GetTaskByID(ctx context.Context, taskID string) (*dto.TaskDTO, error)
	// UpdateTask(ctx context.Context, taskID string, req dto.UpdateTaskRequest, updatedBy string) error
	DeleteTask(ctx context.Context, taskID string, claims *dto.JWTClaims) error
	UpdateStepStatus(ctx context.Context, taskID, stepID string, req dto.UpdateStepStatusNoteRequest, claims *dto.JWTClaims) error

	ReplaceTask(ctx context.Context, taskID string, req dto.UpdateTaskPutRequest, updatedBy string) error
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Tasks) error
	UpdateTaskByID(ctx context.Context, taskID string, update models.Tasks) (*models.Tasks, error)
	SoftDeleteTaskByID(ctx context.Context, taskID string) error
	GetAllTaskByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Tasks, error)
	GetOneTasksByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Tasks, error)
	GetListTasksByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Tasks, int64, error)
	//<------------------------------------------------------------------------------->
	GetAllStepSteps(ctx context.Context, taskID string) ([]models.TaskWorkflowStep, error)
	UpdateOneStepFields(ctx context.Context, taskID, stepID string, status *string, notes *string, now time.Time) error
	UpdateTaskStatus(ctx context.Context, taskID, status, stepName string, now time.Time) error

	GetOneUserTaskStatsByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.UserTaskStats, error)
	GetAllUserTaskStatsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.UserTaskStats, error)
	UpsertUserTaskStats(ctx context.Context, stats *models.UserTaskStats) error

	ReplaceTaskByID(ctx context.Context, taskID string, doc *models.Tasks) (*models.Tasks, error)
}
