package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type InComeService interface {
	CreateInCome(ctx context.Context, inCome dto.CreateIncomeDTO, claims *dto.JWTClaims) error
	ListInComes(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, transactionCategoryID string, startDate string, endDate string) (dto.Pagination, error)
	GetIncomeByID(ctx context.Context, incomeID string, claims *dto.JWTClaims) (*dto.IncomeDTO, error)
	UpdateInComeByID(ctx context.Context, incomeID string, update dto.UpdateIncomeDTO, claims *dto.JWTClaims) error
	DeleteInComeByInComeID(ctx context.Context, incomeID string, claims *dto.JWTClaims) error
	SummaryInComeByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestIncomeSummary) (dto.IncomeSummaryDTO, error)
}

type InComeRepository interface {
	CreateInCome(ctx context.Context, inCome models.Income) error
	UpdateInComeByID(ctx context.Context, incomeID string, update models.Income) (*models.Income, error)
	SoftDeleteInComeByincomeID(ctx context.Context, incomeID string) error
	GetAllInComeByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Income, error)
	GetOneInComeByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Income, error)
	GetListInComesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Income, int64, error)
}
