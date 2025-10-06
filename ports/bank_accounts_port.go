package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type BankAccountsService interface {
	CreateBankAccount(ctx context.Context, bankAccount dto.CreateBankAccountsDTO, claims *dto.JWTClaims) error
	UpdateBankAccountID(ctx context.Context, BankID string, update dto.UpdateBankAccountsDTO, claims *dto.JWTClaims) error
	DeleteBankAccountByID(ctx context.Context, BankID string, claims *dto.JWTClaims) error
	ListBankAccounts(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
	GetListBankAccountByBankID(ctx context.Context, BankID string, claims *dto.JWTClaims) (*dto.BankAccountsDTO, error)
}

type BankAccountsRepository interface {
	CreateBankAccount(ctx context.Context, bankAccount models.BankAccount) error
	UpdateBankAccountByID(ctx context.Context, id string, update models.BankAccount) (*models.BankAccount, error)
	SoftDeleteBankAccountByID(ctx context.Context, id string) error
	GetAllBankAccountsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.BankAccount, error)
	GetOneBankAccountByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.BankAccount, error)
	GetListBankAccountsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.BankAccount, int64, error)
}
