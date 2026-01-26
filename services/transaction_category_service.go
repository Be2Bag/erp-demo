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

type transactionCategoryService struct {
	transactionCategoryRepo ports.TransactionCategoryRepository
	config                  config.Config
}

func NewTransactionCategoryService(cfg config.Config, transactionCategoryRepo ports.TransactionCategoryRepository) ports.TransactionCategoryService {
	return &transactionCategoryService{config: cfg, transactionCategoryRepo: transactionCategoryRepo}
}

func (s *transactionCategoryService) CreateTransactionCategory(ctx context.Context, createCategoryTransaction dto.CreateTransactionCategoryDTO, claims *dto.JWTClaims) error {

	now := time.Now()

	model := models.TransactionCategory{
		TransactionCategoryID:     uuid.NewString(),
		Type:                      createCategoryTransaction.Type,
		TransactionCategoryNameTH: createCategoryTransaction.TransactionCategoryNameTH,
		Description:               createCategoryTransaction.Description,
		CreatedBy:                 claims.UserID,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	if err := s.transactionCategoryRepo.CreateTransactionCategory(ctx, model); err != nil {
		return err
	}
	return nil

}

func (s *transactionCategoryService) ListTransactionCategory(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string, types string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"type":       types,
		"deleted_at": nil,
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"category_name_th": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":       "created_at",
		"updated_at":       "updated_at",
		"category_name_th": "category_name_th",
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

	items, total, err := s.transactionCategoryRepo.GetListTransactionCategoryByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list categories: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		list = append(list, dto.TransactionCategoryDTO{
			TransactionCategoryID:     m.TransactionCategoryID,
			Type:                      m.Type,
			TransactionCategoryNameTH: m.TransactionCategoryNameTH,
			Description:               m.Description,
			CreatedBy:                 m.CreatedBy,
			CreatedAt:                 m.CreatedAt,
			UpdatedAt:                 m.UpdatedAt,
			DeletedAt:                 m.DeletedAt,
			Note:                      m.Note,
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

func (s *transactionCategoryService) GetTransactionCategoryByID(ctx context.Context, TransactionCategoryID string, claims *dto.JWTClaims) (*dto.TransactionCategoryDTO, error) {

	filter := bson.M{"transaction_category_id": TransactionCategoryID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.transactionCategoryRepo.GetOneTransactionCategoryByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, nil
	}

	dtoObj := &dto.TransactionCategoryDTO{
		TransactionCategoryID:     m.TransactionCategoryID,
		Type:                      m.Type,
		TransactionCategoryNameTH: m.TransactionCategoryNameTH,
		Description:               m.Description,
		CreatedBy:                 m.CreatedBy,
		CreatedAt:                 m.CreatedAt,
		UpdatedAt:                 m.UpdatedAt,
		DeletedAt:                 m.DeletedAt,
		Note:                      m.Note,
	}
	return dtoObj, nil
}

func (s *transactionCategoryService) UpdateTransactionCategoryByID(ctx context.Context, transactionCategoryID string, update dto.UpdateTransactionCategoryDTO, claims *dto.JWTClaims) error {
	filter := bson.M{"transaction_category_id": transactionCategoryID, "deleted_at": nil}

	existing, err := s.transactionCategoryRepo.GetOneTransactionCategoryByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.Type != nil {
		existing.Type = *update.Type
	}
	if update.TransactionCategoryNameTH != nil {
		existing.TransactionCategoryNameTH = *update.TransactionCategoryNameTH
	}
	if update.Description != nil {
		existing.Description = *update.Description
	}

	if update.Note != nil {
		existing.Note = update.Note
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.transactionCategoryRepo.UpdateTransactionCategoryByID(ctx, transactionCategoryID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *transactionCategoryService) DeleteTransactionCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error {
	err := s.transactionCategoryRepo.SoftDeleteTransactionCategoryByID(ctx, categoryID, claims)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
