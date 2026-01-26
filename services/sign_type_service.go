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

type signTypeService struct {
	signTypeRepo ports.SignTypeRepository
	userRepo     ports.UserRepository
	config       config.Config
}

func NewSignTypeService(cfg config.Config, signTypeRepo ports.SignTypeRepository, userRepo ports.UserRepository) ports.SignTypeService {
	return &signTypeService{config: cfg, signTypeRepo: signTypeRepo, userRepo: userRepo}
}

func (s *signTypeService) CreateSignType(ctx context.Context, signType dto.CreateSignTypeDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	model := models.SignType{
		TypeID:    uuid.NewString(),
		NameTH:    signType.NameTH,
		NameEN:    signType.NameEN,
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.signTypeRepo.CreateSignType(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *signTypeService) ListSignTypes(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error) {
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
			{"name_th": re},
			{"name_en": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"name_th":    "name_th",
		"name_en":    "name_en",
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

	items, total, err := s.signTypeRepo.GetListSignTypesByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list sign types: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		createdByName := "ไม่พบผู้สร้าง"
		createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)
		if createdBy != nil {
			createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
		}

		list = append(list, dto.SignTypeDTO{
			TypeID:    m.TypeID,
			NameTH:    m.NameTH,
			NameEN:    m.NameEN,
			CreatedBy: createdByName,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
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

func (s *signTypeService) GetSignTypeByTypeID(ctx context.Context, TypeID string, claims *dto.JWTClaims) (*dto.SignTypeDTO, error) {

	filter := bson.M{"type_id": TypeID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.signTypeRepo.GetOneSignTypeByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	createdByName := "ไม่พบผู้สร้าง"
	createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)
	if createdBy != nil {
		createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
	}

	dtoObj := &dto.SignTypeDTO{
		TypeID:    m.TypeID,
		NameTH:    m.NameTH,
		NameEN:    m.NameEN,
		CreatedBy: createdByName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *signTypeService) UpdateSignTypeByTypeID(ctx context.Context, typeID string, update dto.UpdateSignTypeDTO, claims *dto.JWTClaims) error {
	// ดึงข้อมูลเดิม
	filter := bson.M{"type_id": typeID, "deleted_at": nil}
	existing, err := s.signTypeRepo.GetOneSignTypeByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if update.NameTH != "" {
		existing.NameTH = update.NameTH
	}
	if update.NameEN != "" {
		existing.NameEN = update.NameEN
	}

	existing.CreatedBy = claims.UserID

	existing.UpdatedAt = time.Now()

	updated, err := s.signTypeRepo.UpdateSignTypeByTypeID(ctx, typeID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *signTypeService) DeleteSignTypeByTypeID(ctx context.Context, typeID string, claims *dto.JWTClaims) error {
	err := s.signTypeRepo.SoftDeleteSignTypeByTypeID(ctx, typeID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
