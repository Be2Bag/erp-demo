package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionCategoryService interface {
	CreateTransactionCategory(ctx context.Context, createCategoryTransaction dto.CreateTransactionCategoryDTO, claims *dto.JWTClaims) error
	ListTransactionCategory(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, types string) (dto.Pagination, error)
	GetTransactionCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) (*dto.TransactionCategoryDTO, error)
	UpdateTransactionCategoryByID(ctx context.Context, transactionCategoryID string, update dto.UpdateTransactionCategoryDTO, claims *dto.JWTClaims) error
	DeleteTransactionCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error
}

type TransactionCategoryRepository interface {
	CreateTransactionCategory(ctx context.Context, transactionCategory models.TransactionCategory) error
	GetListTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.TransactionCategory, int64, error)
	GetOneTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.TransactionCategory, error)
	UpdateTransactionCategoryByID(ctx context.Context, transactionCategoryID string, update models.TransactionCategory) (*models.TransactionCategory, error)
	SoftDeleteTransactionCategoryByID(ctx context.Context, transactionCategoryID string, claims *dto.JWTClaims) error
	GetAllTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.TransactionCategory, error)
}
