package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type PayableService interface {
}

type PayableRepository interface {
	CreatePayable(ctx context.Context, payable models.Payable) error
	UpdatePayableByID(ctx context.Context, payableID string, update models.Payable) (*models.Payable, error)
	SoftDeletePayableByID(ctx context.Context, payableID string) error
	GetAllPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Payable, error)
	GetOnePayableByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Payable, error)
	GetListPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Payable, int64, error)
}
