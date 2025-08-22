package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type PositionService interface {
	CreatePosition(ctx context.Context, createPosition dto.CreatePositionDTO, claims *dto.JWTClaims) error
	UpdatePosition(ctx context.Context, positionID string, updatePosition dto.UpdatePositionDTO, claims *dto.JWTClaims) error
	DeletePositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) error
	GetPositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) (*dto.PositionDTO, error)
	GetPositionList(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department_id string, sortBy string, sortOrder string) (dto.Pagination, error)
}
type PositionRepository interface {
	CreatePosition(ctx context.Context, position models.Position) error
	GetListPositionByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Position, int64, error)
	GetOnePositionByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Position, error)
	UpdatePositionByID(ctx context.Context, positionID string, update models.Position) (*models.Position, error)
	SoftDeletePositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) error
	GetAllPositionByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error)
}
