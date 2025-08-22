package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type workflowService struct {
	config       config.Config
	workflowRepo ports.WorkFlowRepository
}

func NewWorkflowService(cfg config.Config, workflowRepo ports.WorkFlowRepository) ports.WorkFlowService {
	return &workflowService{config: cfg, workflowRepo: workflowRepo}
}

func (s *workflowService) CreateWorkflowTemplate(ctx context.Context, req dto.CreateWorkflowTemplateDTO, claims *dto.JWTClaims) error {
	if strings.TrimSpace(req.WorkFlowName) == "" {
		return errors.New("workflow_name is required")
	}
	if len(req.Steps) == 0 {
		return errors.New("steps are required")
	}

	now := time.Now()
	steps := make([]models.WorkFlowStep, 0, len(req.Steps))
	var total float64
	for _, st := range req.Steps {
		if st.Hours < 0 {
			return errors.New("step hours must be >= 0")
		}
		steps = append(steps, models.WorkFlowStep{
			StepID:      uuid.NewString(),
			StepName:    st.StepName,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       st.Order,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		total += st.Hours
	}

	tmpl := models.WorkFlowTemplate{
		WorkFlowID:   uuid.NewString(),
		WorkFlowName: req.WorkFlowName,
		Department:   req.Department,
		Description:  req.Description,
		TotalHours:   total,
		Steps:        steps,
		IsActive:     true,
		Version:      1,
		CreatedBy:    claims.UserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.workflowRepo.CreateWorkFlowTemplate(ctx, &tmpl); err != nil {
		return err
	}

	return nil
}

func (s *workflowService) GetWorkflowTemplateByID(ctx context.Context, workflowID string) (*dto.WorkflowTemplateDTO, error) {

	if strings.TrimSpace(workflowID) == "" {
		return nil, errors.New("workflowID is required")
	}

	filter := bson.M{
		"workflow_id": workflowID,
		"deleted_at":  nil,
	}
	var projection bson.M // ใช้ nil ก็ได้ ถ้าไม่ต้องการเลือก field

	m, err := s.workflowRepo.GetOneWorkFlowTemplateByFilter(ctx, filter, projection)
	if err != nil {
		return nil, fmt.Errorf("get workflow template by id: %w", err)
	}
	if m == nil {
		return nil, nil
	}

	stepsDTO := make([]dto.WorkflowStepDTO, 0, len(m.Steps))
	for _, st := range m.Steps {
		stepsDTO = append(stepsDTO, dto.WorkflowStepDTO{
			StepID:      st.StepID,
			StepName:    st.StepName,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       st.Order,
			CreatedAt:   st.CreatedAt,
			UpdatedAt:   st.UpdatedAt,
		})
	}

	dtoObj := &dto.WorkflowTemplateDTO{
		WorkFlowID:   m.WorkFlowID,
		WorkFlowName: m.WorkFlowName,
		Department:   m.Department,
		Description:  m.Description,
		TotalHours:   m.TotalHours,
		Steps:        stepsDTO,
		IsActive:     m.IsActive,
		Version:      m.Version,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	return dtoObj, nil
}

func (s *workflowService) ListWorkflowTemplates(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department_id"] = department
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"workflow_name": re},
		}
	}

	projection := bson.M{}

	allowedSortFields := map[string]string{
		"created_at":    "created_at",
		"updated_at":    "updated_at",
		"workflow_name": "workflow_name",
	}
	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = "created_at"
	}
	order := int32(-1)
	if strings.EqualFold(sortOrder, "asc") {
		order = 1
	}

	sort := bson.D{
		{Key: field, Value: order},
		{Key: "_id", Value: -1},
	}

	items, total, err := s.workflowRepo.GetListWorkFlowTemplatesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list sign jobs: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		stepsDTO := make([]dto.WorkflowStepDTO, 0, len(m.Steps))
		for _, st := range m.Steps {
			stepsDTO = append(stepsDTO, dto.WorkflowStepDTO{
				StepID:      st.StepID,
				StepName:    st.StepName,
				Description: st.Description,
				Hours:       st.Hours,
				Order:       st.Order,
				CreatedAt:   st.CreatedAt,
				UpdatedAt:   st.UpdatedAt,
			})
		}

		list = append(list, dto.WorkflowTemplateDTO{
			WorkFlowID:   m.WorkFlowID,
			WorkFlowName: m.WorkFlowName,
			Department:   m.Department,
			Description:  m.Description,
			TotalHours:   m.TotalHours,
			Steps:        stepsDTO,
			IsActive:     m.IsActive,
			Version:      m.Version,
			CreatedBy:    m.CreatedBy,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		})
	}

	totalPages := 0
	if total > 0 && size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	return dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}, nil
}

func (s *workflowService) UpdateWorkflowTemplate(ctx context.Context, workflowID string, req dto.UpdateWorkflowTemplateDTO, updatedBy string) error {

	filter := bson.M{"workflow_id": workflowID, "deleted_at": nil}
	existing, err := s.workflowRepo.GetOneWorkFlowTemplateByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if req.WorkFlowName != "" {
		existing.WorkFlowName = req.WorkFlowName
	}
	if req.Department != "" {
		existing.Department = req.Department
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Steps != nil {
		nowSteps := time.Now()
		var total float64
		newSteps := make([]models.WorkFlowStep, 0, len(*req.Steps))
		for _, st := range *req.Steps {
			if st.Hours < 0 {
				return errors.New("step hours must be >= 0")
			}
			newSteps = append(newSteps, models.WorkFlowStep{
				StepID:      uuid.NewString(),
				StepName:    st.StepName,
				Description: st.Description,
				Hours:       st.Hours,
				Order:       st.Order,
				CreatedAt:   nowSteps,
				UpdatedAt:   nowSteps,
			})
			total += st.Hours
		}
		existing.Steps = newSteps
		existing.TotalHours = total
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.workflowRepo.UpdateWorkFlowTemplateByID(ctx, workflowID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *workflowService) DeleteWorkflowTemplate(ctx context.Context, workflowID string) error {
	err := s.workflowRepo.SoftDeleteWorkFlowTemplateByID(ctx, workflowID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
