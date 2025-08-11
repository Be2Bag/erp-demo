package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type KPIService interface {
	CreateKPITemplate(ctx context.Context, req dto.KPITemplateDTO, claims *dto.JWTClaims) error
	GetKPITemplateByID(ctx context.Context, id string) (interface{}, error)
	UpdateKPITemplate(ctx context.Context, id string, updated dto.KPITemplateDTO, claims *dto.JWTClaims) (interface{}, error)
	DeleteKPITemplate(ctx context.Context, id string) error
	ListKPITemplates(ctx context.Context, q dto.KPITemplateListQuery) ([]interface{}, int64, error)
}
