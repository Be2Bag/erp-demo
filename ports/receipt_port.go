package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ReceiptService interface {
	CreateReceipt(ctx context.Context, in dto.CreateReceiptDTO, claims *dto.JWTClaims) error
	ListReceipts(ctx context.Context, claims *dto.JWTClaims, page, size int, search, sortBy, sortOrder, status, startDate, endDate, billType string) (dto.Pagination, error)
	GetReceiptByID(ctx context.Context, receiptID string, claims *dto.JWTClaims) (*dto.ReceiptDTO, error)
	DeleteReceiptByID(ctx context.Context, receiptID string, claims *dto.JWTClaims) error
	SummaryReceiptByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestSummaryReceipt) (dto.ReceiptSummaryDTO, error)
}

type ReceiptRepository interface {
	CreateReceipt(ctx context.Context, receipt models.Receipt) error
	GetListReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Receipt, int64, error)
	GetOneReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Receipt, error)
	SoftDeleteReceiptByID(ctx context.Context, receiptID string) error
	GetAllReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Receipt, error)
	GetMaxReceiptNumber(ctx context.Context, prefix string) (string, error)
}
