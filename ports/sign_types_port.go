package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SignTypeRepository interface {
	CreateSignType(ctx context.Context, signType models.SignType) error
	UpdateSignTypeByTypeID(ctx context.Context, typeID string, update models.SignType) (*models.SignType, error)
	SoftDeleteSignTypeByTypeID(ctx context.Context, typeID string) error
	GetAllSignTypeByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignType, error)
	GetOneSignTypeByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.SignType, error)
	GetListSignTypesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.SignType, int64, error)
}

type SignTypeService interface {
	CreateSignType(ctx context.Context, signType dto.CreateSignTypeDTO, claims *dto.JWTClaims) error
	ListSignTypes(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetSignTypeByTypeID(ctx context.Context, TypeID string, claims *dto.JWTClaims) (*dto.SignTypeDTO, error)
	UpdateSignTypeByTypeID(ctx context.Context, typeID string, update dto.UpdateSignTypeDTO, claims *dto.JWTClaims) error
	DeleteSignTypeByTypeID(ctx context.Context, typeID string, claims *dto.JWTClaims) error
}
