package services

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
)

type taskService struct {
	config   config.Config
	taskRepo ports.TaskRepository
	userRepo ports.UserRepository
}

func NewTaskService(cfg config.Config, taskRepo ports.TaskRepository, userRepo ports.UserRepository) ports.TaskService {
	return &taskService{config: cfg, taskRepo: taskRepo, userRepo: userRepo}
}

func (s *taskService) GetTasks(ctx context.Context, filter interface{}) ([]interface{}, error) {
	// Implementation for fetching tasks
	return nil, nil
}

func (s *taskService) CreateTask(ctx context.Context, createTask dto.CreateTaskRequest, claims *dto.JWTClaims) error {
	now := time.Now()
	var start time.Time
	var end time.Time

	if createTask.StartDate != "" {
		parsedDate, err := time.Parse("2006-01-02", createTask.StartDate)
		if err != nil {
			return err
		}
		start = parsedDate
	}

	if createTask.EndDate != "" {
		parsedDate, err := time.Parse("2006-01-02", createTask.EndDate)
		if err != nil {
			return err
		}
		end = parsedDate
	}

	model := models.Tasks{
		TaskID:      uuid.New().String(),
		ProjectID:   createTask.ProjectID,
		ProjectName: createTask.ProjectName,
		JobName:     createTask.JobName,
		Description: createTask.Description,

		Department: createTask.Department,
		Assignee:   createTask.Assignee,
		Importance: createTask.Importance,
		StartDate:  start,
		EndDate:    end,
		KPIID:      createTask.KPIID,
		WorkFlowID: createTask.WorkflowID,

		// AppliedWorkflow: TaskAppliedWorkflow{
		// 	WorkFlowID:   createTask.WorkflowID,
		// 	WorkFlowName: createTask.WorkflowName,
		// 	Department:   createTask.Department,
		// 	Description:  createTask.Description,
		// 	TotalHours:   createTask.TotalHours,
		// 	Steps:       createTask.Steps,
		// 	Version:     1,
		// },

		Status:    "todo",
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	if err := s.taskRepo.CreateTask(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) (interface{}, error) {
	// Implementation for fetching a task by ID
	return nil, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id string, updatedTask interface{}) error {
	// Implementation for updating a task
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	// Implementation for deleting a task
	return nil
}

func (s *taskService) UpdateTaskWorkflow(ctx context.Context, id string, workflowStep interface{}) error {
	// Implementation for updating task workflow
	return nil
}

func (s *taskService) GetTaskStatistics(ctx context.Context, filter interface{}) (map[string]interface{}, error) {
	// Implementation for fetching task statistics
	return nil, nil
}
