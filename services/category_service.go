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

type categoryService struct {
	categoryRepo ports.CategoryRepository
	config       config.Config
}

func NewCategoryService(cfg config.Config, categoryRepo ports.CategoryRepository) ports.CategoryService {
	return &categoryService{config: cfg, categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, createCategory dto.CreateCategoryDTO, claims *dto.JWTClaims) error {

	now := time.Now()

	model := models.Category{
		CategoryID:     uuid.NewString(),
		DepartmentID:   createCategory.DepartmentID,
		CategoryNameTH: createCategory.CategoryNameTH,
		CategoryNameEN: createCategory.CategoryNameEN,
		Description:    createCategory.Description,
		CreatedBy:      claims.UserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.categoryRepo.CreateCategory(ctx, model); err != nil {
		return err
	}
	return nil

}

func (s *categoryService) ListCategory(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error) {
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
			{"category_name_th": re},
			{"category_name_en": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":       "created_at",
		"updated_at":       "updated_at",
		"category_name_th": "category_name_th",
		"category_name_en": "category_name_en",
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

	items, total, err := s.categoryRepo.GetListCategoryByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list categories: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		list = append(list, dto.CategoryDTO{
			CategoryID:     m.CategoryID,
			DepartmentID:   m.DepartmentID,
			CategoryNameTH: m.CategoryNameTH,
			CategoryNameEN: m.CategoryNameEN,
			Description:    m.Description,
			CreatedBy:      m.CreatedBy,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
			DeletedAt:      m.DeletedAt,
			Note:           m.Note,
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

func (s *categoryService) GetCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) (*dto.CategoryDTO, error) {

	filter := bson.M{"category_id": categoryID}
	projection := bson.M{}

	m, err := s.categoryRepo.GetOneCategoryByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	dtoObj := &dto.CategoryDTO{
		CategoryID:     m.CategoryID,
		DepartmentID:   m.DepartmentID,
		CategoryNameTH: m.CategoryNameTH,
		CategoryNameEN: m.CategoryNameEN,
		Description:    m.Description,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
		Note:           m.Note,
	}
	return dtoObj, nil
}

func (s *categoryService) UpdateCategoryByID(ctx context.Context, categoryID string, update dto.UpdateCategoryDTO, claims *dto.JWTClaims) error {
	filter := bson.M{"category_id": categoryID, "deleted_at": nil}

	existing, err := s.categoryRepo.GetOneCategoryByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.DepartmentID != nil {
		existing.DepartmentID = *update.DepartmentID
	}
	if update.CategoryNameTH != nil {
		existing.CategoryNameTH = *update.CategoryNameTH
	}
	if update.CategoryNameEN != nil {
		existing.CategoryNameEN = *update.CategoryNameEN
	}
	if update.Description != nil {
		existing.Description = *update.Description
	}
	if update.Note != nil {
		existing.Note = update.Note
	}

	existing.UpdatedAt = time.Now()

	updated, err := s.categoryRepo.UpdateCategoryByID(ctx, categoryID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *categoryService) DeleteCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error {
	err := s.categoryRepo.SoftDeleteCategoryByID(ctx, categoryID, claims)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
