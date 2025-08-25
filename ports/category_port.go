package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, createCategory dto.CreateCategoryDTO, claims *dto.JWTClaims) error
	ListCategory(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) (*dto.CategoryDTO, error)
	UpdateCategoryByID(ctx context.Context, categoryID string, update dto.UpdateCategoryDTO, claims *dto.JWTClaims) error
	DeleteCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error
}
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category models.Category) error
	GetListCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Category, int64, error)
	GetOneCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Category, error)
	UpdateCategoryByID(ctx context.Context, categoryID string, update models.Category) (*models.Category, error)
	SoftDeleteCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error
	GetAllCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Category, error)
}
