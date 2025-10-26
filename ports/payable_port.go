package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type PayableService interface {
	CreatePayable(ctx context.Context, payable dto.CreatePayableDTO, claims *dto.JWTClaims) error
	ListPayables(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, status string) (dto.Pagination, error)
	GetPayableByID(ctx context.Context, payableID string, claims *dto.JWTClaims) (*dto.PayableDTO, error)
	UpdatePayableByID(ctx context.Context, payableID string, update dto.UpdatePayableDTO, claims *dto.JWTClaims) error
	DeletePayableByID(ctx context.Context, payableID string, claims *dto.JWTClaims) error
	SummaryPayableByFilter(ctx context.Context, claims *dto.JWTClaims, report dto.RequestSummaryPayable) (dto.PayableSummaryDTO, error)
	// RecordPayment inserts a payment transaction and updates the payable balance/status accordingly
	RecordPayment(ctx context.Context, input dto.RecordPaymentDTO, claims *dto.JWTClaims) error
}

type PayableRepository interface {
	CreatePayable(ctx context.Context, payable models.Payable) error
	UpdatePayableByID(ctx context.Context, payableID string, update models.Payable) (*models.Payable, error)
	SoftDeletePayableByID(ctx context.Context, payableID string) error
	GetAllPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Payable, error)
	GetOnePayableByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Payable, error)
	GetListPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Payable, int64, error)
	// CreatePaymentTransaction inserts a payment transaction document
	CreatePaymentTransaction(ctx context.Context, tx models.PaymentTransaction) error
}
