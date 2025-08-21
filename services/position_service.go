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

type positionService struct {
	config       config.Config
	positionRepo ports.PositionRepository
}

func NewPositionService(cfg config.Config, positionRepo ports.PositionRepository) ports.PositionService {
	return &positionService{config: cfg, positionRepo: positionRepo}
}

func (s *positionService) CreatePosition(ctx context.Context, createPosition dto.CreatePositionDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	model := models.Position{
		PositionID:   uuid.New().String(),
		DepartmentID: createPosition.DepartmentID,
		PositionName: createPosition.PositionName,
		Level:        createPosition.Level,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.positionRepo.CreatePosition(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *positionService) UpdatePosition(ctx context.Context, positionID string, updatePosition dto.UpdatePositionDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	filter := bson.M{"position_id": positionID, "deleted_at": nil}
	existing, err := s.positionRepo.GetOnePositionByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if updatePosition.DepartmentID != "" {
		existing.DepartmentID = updatePosition.DepartmentID
	}

	if updatePosition.PositionName != "" {
		existing.PositionName = updatePosition.PositionName
	}

	if updatePosition.Level != "" {
		existing.Level = updatePosition.Level
	}

	if updatePosition.Note != "" {
		existing.Note = &updatePosition.Note
	}

	existing.UpdatedAt = now

	if _, err := s.positionRepo.UpdatePositionByID(ctx, positionID, *existing); err != nil {
		return err
	}
	return nil
}

func (s *positionService) DeletePositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) error {
	err := s.positionRepo.SoftDeletePositionByID(ctx, positionID, claims)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *positionService) GetPositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) (*dto.PositionDTO, error) {
	filter := bson.M{"position_id": positionID, "deleted_at": nil}
	projection := bson.M{}
	position, err := s.positionRepo.GetOnePositionByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	dtoObj := &dto.PositionDTO{
		PositionID:   position.PositionID,
		DepartmentID: position.DepartmentID,
		PositionName: position.PositionName,
		Level:        position.Level,
		CreatedAt:    position.CreatedAt,
		UpdatedAt:    position.UpdatedAt,
	}
	return dtoObj, nil
}

func (s *positionService) GetPositionList(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department_id"] = department
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"department_name": re},
			{"position_name": re},
			{"level": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":      "created_at",
		"updated_at":      "updated_at",
		"department_name": "department_name",
		"position_name":   "position_name",
		"level":           "level",
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

	items, total, err := s.positionRepo.GetListPositionByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list positions: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		list = append(list, dto.PositionDTO{
			PositionID:   m.PositionID,
			DepartmentID: m.DepartmentID,
			PositionName: m.PositionName,
			Level:        m.Level,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
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
