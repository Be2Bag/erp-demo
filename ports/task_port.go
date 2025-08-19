package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TaskService interface {
	GetTasks(ctx context.Context, filter interface{}) ([]interface{}, error)
	CreateTask(ctx context.Context, createTask dto.CreateTaskRequest, claims *dto.JWTClaims) error
	GetTaskByID(ctx context.Context, id string) (interface{}, error)
	UpdateTask(ctx context.Context, id string, updatedTask interface{}) error
	DeleteTask(ctx context.Context, id string) error
	UpdateTaskWorkflow(ctx context.Context, id string, workflowStep interface{}) error
	GetTaskStatistics(ctx context.Context, filter interface{}) (map[string]interface{}, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Tasks) error
	UpdateTaskByID(ctx context.Context, taskID string, update models.Tasks) (*models.Tasks, error)
	SoftDeleteTaskByJobID(ctx context.Context, taskID string) error
	GetAllTaskByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Tasks, error)
	GetOneTasksByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Tasks, error)
	GetListTasksByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Tasks, int64, error)
}
