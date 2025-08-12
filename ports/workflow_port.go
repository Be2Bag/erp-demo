package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

type WorkFlowService interface {
	CreateWorkflowTemplate(ctx context.Context, req dto.CreateWorkflowTemplateDTO, createdBy string) (*dto.WorkflowTemplateDTO, error)
	GetWorkflowTemplateByID(ctx context.Context, templateID string) (*dto.WorkflowTemplateDTO, error)
	ListWorkflowTemplates(ctx context.Context, search string, department string, page, limit int64, sort string) ([]dto.WorkflowTemplateDTO, int64, error)
	UpdateWorkflowTemplate(ctx context.Context, templateID string, req dto.UpdateWorkflowTemplateDTO, updatedBy string) (*dto.WorkflowTemplateDTO, error)
	DeleteWorkflowTemplate(ctx context.Context, templateID string) error
}

type WorkFlowRepository interface {
	CreateWorkFlowTemplate(ctx context.Context, tmpl *models.WorkFlowTemplate) error
	GetWorkFlowTemplateByTemplateID(ctx context.Context, templateID string) (*models.WorkFlowTemplate, error)
	GetWorkFlowTemplates(ctx context.Context, filter interface{}, options interface{}) ([]models.WorkFlowTemplate, error)
	CountWorkFlowTemplates(ctx context.Context, filter interface{}) (int64, error)
	UpdateWorkFlowTemplateByTemplateID(ctx context.Context, templateID string, update interface{}) error
	DeleteWorkFlowTemplateByTemplateID(ctx context.Context, templateID string) error
}
