package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type KPIEvaluationService interface {
	ListKPIEvaluation(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error)
	CreateKPIEvaluation(ctx context.Context, req dto.CreateKPIEvaluationRequest, claims *dto.JWTClaims) error
}
type KPIEvaluationRepository interface {
	CreateKPIEvaluations(ctx context.Context, kpi models.KPIEvaluation) error
	UpdateKPIEvaluationByID(ctx context.Context, evaluationID string, update models.KPIEvaluation) (*models.KPIEvaluation, error)
	SoftDeleteKPIEvaluationByID(ctx context.Context, evaluationID string) error
	GetAllKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.KPIEvaluation, error)
	GetOneKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.KPIEvaluation, error)
	GetListKPIEvaluationByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.KPIEvaluation, int64, error)
}
