package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ReceivableService interface {
	CreateReceivable(ctx context.Context, receivable dto.CreateReceivableDTO, claims *dto.JWTClaims) error
	ListReceivables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, Status string) (dto.Pagination, error)
	GetReceivableByID(ctx context.Context, receivableID string, claims *dto.JWTClaims) (*dto.ReceivableDTO, error)
	UpdateReceivableByID(ctx context.Context, receivableID string, update dto.UpdateReceivableDTO, claims *dto.JWTClaims) error
	DeleteReceivableByID(ctx context.Context, receivableID string, claims *dto.JWTClaims) error
	SummaryReceivableByFilter(ctx context.Context, claims *dto.JWTClaims) (dto.ReceivableSummaryDTO, error)
}

type ReceivableRepository interface {
	CreateReceivable(ctx context.Context, receivable models.Receivable) error
	UpdateReceivableByID(ctx context.Context, receivableID string, update models.Receivable) (*models.Receivable, error)
	SoftDeleteReceivableByID(ctx context.Context, receivableID string) error
	GetAllReceivablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Receivable, error)
	GetOneReceivableByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Receivable, error)
	GetListReceivablesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Receivable, int64, error)
}
