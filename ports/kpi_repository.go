package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
)

type KPIRepository interface {
	GetKPITemplates(ctx context.Context, filter interface{}, options interface{}) ([]models.KPITemplate, error)
	CountKPITemplates(ctx context.Context, filter interface{}) (int64, error) // added
	CreateKPITemplate(ctx context.Context, template models.KPITemplate) error
	GetKPITemplateByID(ctx context.Context, id string) (*models.KPITemplate, error)
	UpdateKPITemplate(ctx context.Context, id string, updated models.KPITemplate) (*models.KPITemplate, error)
	DeleteKPITemplate(ctx context.Context, id string) error
}
