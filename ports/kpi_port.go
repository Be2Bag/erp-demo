package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

type KPIService interface {
	CreateKPITemplate(ctx context.Context, req dto.KPITemplateDTO, claims *dto.JWTClaims) error
	GetKPITemplateByID(ctx context.Context, id string) (interface{}, error)
	UpdateKPITemplate(ctx context.Context, id string, updated dto.KPITemplateDTO, claims *dto.JWTClaims) (interface{}, error)
	DeleteKPITemplate(ctx context.Context, id string) error
	ListKPITemplates(ctx context.Context, q dto.KPITemplateListQuery) ([]interface{}, int64, error)
}
type KPIRepository interface {
	GetKPITemplates(ctx context.Context, filter interface{}, options interface{}) ([]models.KPITemplate, error)
	CountKPITemplates(ctx context.Context, filter interface{}) (int64, error) // added
	CreateKPITemplate(ctx context.Context, template models.KPITemplate) error
	GetKPITemplateByID(ctx context.Context, id string) (*models.KPITemplate, error)
	UpdateKPITemplate(ctx context.Context, id string, updated models.KPITemplate) (*models.KPITemplate, error)
	DeleteKPITemplate(ctx context.Context, id string) error
}
