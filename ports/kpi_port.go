package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type KPIService interface {
	CreateKPITemplate(ctx context.Context, req dto.CreateKPITemplateDTO, claims *dto.JWTClaims) error
	GetKPITemplateByID(ctx context.Context, id string) (*dto.KPITemplateDTO, error)
	UpdateKPITemplate(ctx context.Context, kpiID string, req dto.UpdateKPITemplateDTO, claims *dto.JWTClaims) error
	DeleteKPITemplate(ctx context.Context, kpiID string) error
	ListKPITemplates(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department_id string, sortBy string, sortOrder string) (dto.Pagination, error)
}
type KPIRepository interface {
	CreateKPI(ctx context.Context, kpi models.KPITemplate) error
	UpdateKPIByID(ctx context.Context, kpiID string, update models.KPITemplate) (*models.KPITemplate, error)
	SoftDeleteKPIByID(ctx context.Context, kpiID string) error
	GetAllKPIByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.KPITemplate, error)
	GetOneKPIByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.KPITemplate, error)
	GetListKPIByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.KPITemplate, int64, error)
}
