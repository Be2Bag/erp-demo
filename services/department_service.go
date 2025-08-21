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

type departmentService struct {
	config         config.Config
	departmentRepo ports.DepartmentRepository
	userRepo       ports.UserRepository
}

func NewDepartmentService(cfg config.Config, departmentRepo ports.DepartmentRepository, userRepo ports.UserRepository) ports.DepartmentService {
	return &departmentService{config: cfg, departmentRepo: departmentRepo, userRepo: userRepo}
}

func (s *departmentService) CreateDepartment(ctx context.Context, createDepartment dto.CreateDepartmentDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	model := models.Department{
		DepartmentID:   uuid.New().String(),
		DepartmentName: createDepartment.DepartmentName,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.departmentRepo.CreateDepartment(ctx, model); err != nil {
		return err
	}
	return nil
}

func (s *departmentService) UpdateDepartment(ctx context.Context, departmentID string, updateDepartment dto.UpdateDepartmentDTO, claims *dto.JWTClaims) error {
	now := time.Now()

	update := models.Department{
		DepartmentName: updateDepartment.DepartmentName,
		UpdatedAt:      now,
	}

	if _, err := s.departmentRepo.UpdateDepartmentByID(ctx, departmentID, update); err != nil {
		return err
	}
	return nil
}

func (s *departmentService) DeleteDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) error {
	err := s.departmentRepo.SoftDeleteDepartmentByID(ctx, departmentID, claims)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *departmentService) GetDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) (*dto.DepartmentDTO, error) {
	filter := bson.M{"department_id": departmentID, "deleted_at": nil}
	projection := bson.M{}
	department, err := s.departmentRepo.GetOneDepartmentByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	dtoObj := &dto.DepartmentDTO{
		DepartmentID:   department.DepartmentID,
		DepartmentName: department.DepartmentName,
		ManagerID:      department.ManagerID,
		CreatedAt:      department.CreatedAt,
		UpdatedAt:      department.UpdatedAt,
	}
	return dtoObj, nil
}

func (s *departmentService) GetDepartmentList(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error) {
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
			{"department_name": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":      "created_at",
		"updated_at":      "updated_at",
		"department_name": "department_name",
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

	items, total, err := s.departmentRepo.GetListDepartmentByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list departments: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		managerName := "ไม่มีชื่อผู้จัดการแผนก"

		if m.ManagerID != "" {
			manager, err := s.userRepo.GetUserByFilter(ctx, bson.M{"user_id": m.ManagerID}, bson.M{"title_th": 1, "first_name_th": 1, "last_name_th": 1})
			if err != nil {
				return dto.Pagination{}, fmt.Errorf("get manager: %w", err)
			}
			if manager != nil {
				managerName = fmt.Sprintf("%s %s %s", manager[0].TitleTH, manager[0].FirstNameTH, manager[0].LastNameTH)
			}
		}

		list = append(list, dto.DepartmentDTO{
			DepartmentID:   m.DepartmentID,
			DepartmentName: m.DepartmentName,
			ManagerID:      m.ManagerID,
			ManagerName:    managerName,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
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
