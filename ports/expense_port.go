package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ExpenseService interface {
	CreateExpense(ctx context.Context, expense dto.CreateExpenseDTO, claims *dto.JWTClaims) error
	ListExpenses(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetExpenseByID(ctx context.Context, expenseID string, claims *dto.JWTClaims) (*dto.ExpenseDTO, error)
	UpdateExpenseByID(ctx context.Context, expenseID string, update dto.UpdateExpenseDTO, claims *dto.JWTClaims) error
	DeleteExpenseByID(ctx context.Context, expenseID string, claims *dto.JWTClaims) error
	SummaryExpenseByFilter(ctx context.Context, claims *dto.JWTClaims) (dto.ExpenseSummaryDTO, error)
}

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, expense models.Expense) error
	UpdateExpenseByID(ctx context.Context, expenseID string, update models.Expense) (*models.Expense, error)
	SoftDeleteExpenseByID(ctx context.Context, expenseID string) error
	GetAllExpenseByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Expense, error)
	GetOneExpenseByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Expense, error)
	GetListExpensesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Expense, int64, error)
}
