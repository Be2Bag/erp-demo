package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bankAccountService struct {
	config          config.Config
	bankAccountRepo ports.BankAccountsRepository
}

func NewBankAccountService(cfg config.Config, bankAccountRepo ports.BankAccountsRepository) ports.BankAccountsService {
	return &bankAccountService{config: cfg, bankAccountRepo: bankAccountRepo}
}

func (s *bankAccountService) CreateBankAccount(ctx context.Context, bankAccount dto.CreateBankAccountsDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	model := models.BankAccount{
		BankID:      uuid.New().String(),
		BankName:    bankAccount.BankName,
		AccountNo:   bankAccount.AccountNo,
		AccountName: bankAccount.AccountName,
		CreatedBy:   claims.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.bankAccountRepo.CreateBankAccount(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *bankAccountService) UpdateBankAccountID(ctx context.Context, BankID string, update dto.UpdateBankAccountsDTO, claims *dto.JWTClaims) error {

	filter := bson.M{"bank_id": BankID, "deleted_at": nil}
	existing, err := s.bankAccountRepo.GetOneBankAccountByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.BankName != "" {
		existing.BankName = update.BankName
	}
	if update.AccountNo != "" {
		existing.AccountNo = update.AccountNo
	}
	if update.AccountName != "" {
		existing.AccountName = update.AccountName
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.bankAccountRepo.UpdateBankAccountByID(ctx, BankID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *bankAccountService) DeleteBankAccountByID(ctx context.Context, BankID string, claims *dto.JWTClaims) error {
	filter := bson.M{"bank_id": BankID, "deleted_at": nil}
	existing, err := s.bankAccountRepo.GetOneBankAccountByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if err := s.bankAccountRepo.SoftDeleteBankAccountByID(ctx, BankID); err != nil {
		return err
	}

	return nil
}

func (s *bankAccountService) ListBankAccounts(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"bank_name": re},
			{"account_no": re},
			{"account_name": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"bank_name":  "bank_name",
	}
	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = "created_at"
	}
	order := int32(-1)
	if strings.EqualFold(sortOrder, "asc") {
		order = 1
	}

	sort := bson.D{
		{Key: field, Value: order},
		{Key: "_id", Value: -1},
	}

	items, total, err := s.bankAccountRepo.GetListBankAccountsByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list bank accounts: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		list = append(list, dto.BankAccountsDTO{
			BankID:      m.BankID,
			BankName:    m.BankName,
			AccountNo:   m.AccountNo,
			AccountName: m.AccountName,
			CreatedBy:   m.CreatedBy,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}

	totalPages := 0
	if total > 0 && size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	return dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}, nil
}

func (s *bankAccountService) GetListBankAccountByBankID(ctx context.Context, BankID string, claims *dto.JWTClaims) (*dto.BankAccountsDTO, error) {

	filter := bson.M{"bank_id": BankID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.bankAccountRepo.GetOneBankAccountByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	dtoObj := &dto.BankAccountsDTO{
		BankID:      m.BankID,
		BankName:    m.BankName,
		AccountNo:   m.AccountNo,
		AccountName: m.AccountName,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	return dtoObj, nil
}
