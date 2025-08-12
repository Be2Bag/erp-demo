package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongoopt "go.mongodb.org/mongo-driver/mongo/options"
)

type workFlowService struct {
	config       config.Config
	workFlowRepo ports.WorkFlowRepository
}

func NewWorkFlowService(cfg config.Config, workFlowRepo ports.WorkFlowRepository) ports.WorkFlowService {
	return &workFlowService{config: cfg, workFlowRepo: workFlowRepo}
}

func (s *workFlowService) CreateWorkflowTemplate(ctx context.Context, req dto.CreateWorkflowTemplateDTO, createdBy string) (*dto.WorkflowTemplateDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}
	if len(req.Steps) == 0 {
		return nil, errors.New("steps are required")
	}

	now := time.Now().UTC()
	steps := make([]models.WorkFlowStep, 0, len(req.Steps))
	var total float64
	for _, st := range req.Steps {
		if st.Hours < 0 {
			return nil, errors.New("step hours must be >= 0")
		}
		steps = append(steps, models.WorkFlowStep{
			StepID:      uuid.NewString(),
			Name:        st.Name,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       st.Order,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		total += st.Hours
	}

	tmpl := models.WorkFlowTemplate{
		TemplateID:  uuid.NewString(),
		Name:        req.Name,
		Department:  req.Department,
		Description: req.Description,
		TotalHours:  total,
		Steps:       steps,
		IsActive:    true,
		Version:     1,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.workFlowRepo.CreateWorkFlowTemplate(ctx, &tmpl); err != nil {
		return nil, err
	}
	out := toTemplateDTO(&tmpl)
	return &out, nil
}

func (s *workFlowService) GetWorkflowTemplateByID(ctx context.Context, templateID string) (*dto.WorkflowTemplateDTO, error) {
	tmpl, err := s.workFlowRepo.GetWorkFlowTemplateByTemplateID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	out := toTemplateDTO(tmpl)
	return &out, nil
}

func (s *workFlowService) ListWorkflowTemplates(ctx context.Context, search string, department string, page, limit int64, sort string) ([]dto.WorkflowTemplateDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	filter := bson.M{}
	if strings.TrimSpace(search) != "" {
		reg := bson.M{"$regex": search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": reg},
			{"department": reg},
			{"description": reg},
		}
	}
	if strings.TrimSpace(department) != "" {
		filter["department"] = department
	}

	findOpts := mongoopt.Find()
	findOpts.SetLimit(limit)
	findOpts.SetSkip((page - 1) * limit)

	// sort format: "field:asc|desc", default updated_at desc
	sortField := "updated_at"
	sortDir := -1
	if strings.TrimSpace(sort) != "" {
		parts := strings.Split(sort, ":")
		if len(parts) >= 1 && parts[0] != "" {
			sortField = parts[0]
		}
		if len(parts) == 2 && strings.EqualFold(parts[1], "asc") {
			sortDir = 1
		}
	}
	findOpts.SetSort(bson.D{{Key: sortField, Value: sortDir}})

	total, err := s.workFlowRepo.CountWorkFlowTemplates(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	items, err := s.workFlowRepo.GetWorkFlowTemplates(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}

	out := make([]dto.WorkflowTemplateDTO, 0, len(items))
	for i := range items {
		out = append(out, toTemplateDTO(&items[i]))
	}
	return out, total, nil
}

func (s *workFlowService) UpdateWorkflowTemplate(ctx context.Context, templateID string, req dto.UpdateWorkflowTemplateDTO, updatedBy string) (*dto.WorkflowTemplateDTO, error) {
	set := bson.M{}
	now := time.Now().UTC()
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, errors.New("name cannot be empty")
		}
		set["name"] = *req.Name
	}
	if req.Department != nil {
		set["department"] = *req.Department
	}
	if req.Description != nil {
		set["description"] = *req.Description
	}
	if req.Steps != nil {
		steps := make([]models.WorkFlowStep, 0, len(*req.Steps))
		var total float64
		for _, st := range *req.Steps {
			if st.Hours < 0 {
				return nil, errors.New("step hours must be >= 0")
			}
			steps = append(steps, models.WorkFlowStep{
				StepID:      uuid.NewString(),
				Name:        st.Name,
				Description: st.Description,
				Hours:       st.Hours,
				Order:       st.Order,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			total += st.Hours
		}
		set["steps"] = steps
		set["total_hours"] = total
	}
	set["updated_at"] = now

	updateDoc := bson.M{
		"$set": set,
		"$inc": bson.M{"version": 1},
	}

	if err := s.workFlowRepo.UpdateWorkFlowTemplateByTemplateID(ctx, templateID, updateDoc); err != nil {
		return nil, err
	}
	updated, err := s.workFlowRepo.GetWorkFlowTemplateByTemplateID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	out := toTemplateDTO(updated)
	return &out, nil
}

func (s *workFlowService) DeleteWorkflowTemplate(ctx context.Context, templateID string) error {
	return s.workFlowRepo.DeleteWorkFlowTemplateByTemplateID(ctx, templateID)
}

func toTemplateDTO(m *models.WorkFlowTemplate) dto.WorkflowTemplateDTO {
	steps := make([]dto.WorkflowStepDTO, 0, len(m.Steps))
	for _, st := range m.Steps {
		steps = append(steps, dto.WorkflowStepDTO{
			StepID:      st.StepID,
			Name:        st.Name,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       st.Order,
			CreatedAt:   st.CreatedAt,
			UpdatedAt:   st.UpdatedAt,
		})
	}
	return dto.WorkflowTemplateDTO{
		TemplateID:  m.TemplateID,
		Name:        m.Name,
		Department:  m.Department,
		Description: m.Description,
		TotalHours:  m.TotalHours,
		Steps:       steps,
		IsActive:    m.IsActive,
		Version:     m.Version,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
