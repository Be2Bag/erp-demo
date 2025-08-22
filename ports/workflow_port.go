package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type WorkFlowService interface {
	CreateWorkflowTemplate(ctx context.Context, req dto.CreateWorkflowTemplateDTO, claims *dto.JWTClaims) error
	GetWorkflowTemplateByID(ctx context.Context, workflowID string) (*dto.WorkflowTemplateDTO, error)
	ListWorkflowTemplates(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department_id string, sortBy string, sortOrder string) (dto.Pagination, error)
	UpdateWorkflowTemplate(ctx context.Context, workflowID string, req dto.UpdateWorkflowTemplateDTO, updatedBy string) error
	DeleteWorkflowTemplate(ctx context.Context, workflowID string) error
}

type WorkFlowRepository interface {
	CreateWorkFlowTemplate(ctx context.Context, tmpl *models.WorkFlowTemplate) error
	UpdateWorkFlowTemplateByID(ctx context.Context, workflowID string, update models.WorkFlowTemplate) (*models.WorkFlowTemplate, error)
	SoftDeleteWorkFlowTemplateByID(ctx context.Context, workflowID string) error
	GetAllWorkFlowTemplatesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.WorkFlowTemplate, error)
	GetOneWorkFlowTemplateByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.WorkFlowTemplate, error)
	GetListWorkFlowTemplatesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.WorkFlowTemplate, int64, error)
}
