package ports

import "context"

type TaskService interface {
	GetTasks(ctx context.Context, filter interface{}) ([]interface{}, error)
	CreateTask(ctx context.Context, task interface{}) error
	GetTaskByID(ctx context.Context, id string) (interface{}, error)
	UpdateTask(ctx context.Context, id string, updatedTask interface{}) error
	DeleteTask(ctx context.Context, id string) error
	UpdateTaskWorkflow(ctx context.Context, id string, workflowStep interface{}) error
	GetTaskStatistics(ctx context.Context, filter interface{}) (map[string]interface{}, error)
}
