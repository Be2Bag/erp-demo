package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type KPIService interface {
	GetKPITemplates(ctx context.Context, filter interface{}) ([]interface{}, error)
	CreateKPITemplate(ctx context.Context, template dto.KPITemplateDTO) error
	GetKPITemplateByID(ctx context.Context, id string) (interface{}, error)
	UpdateKPITemplate(ctx context.Context, id string, updatedTemplate interface{}) error
	DeleteKPITemplate(ctx context.Context, id string) error
}
