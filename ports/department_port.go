package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, createDepartment dto.CreateDepartmentDTO, claims *dto.JWTClaims) error
	UpdateDepartment(ctx context.Context, departmentID string, updateDepartment dto.UpdateDepartmentDTO, claims *dto.JWTClaims) error
	DeleteDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) error
	GetDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) (*dto.DepartmentDTO, error)
	GetDepartmentList(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, sortBy string, sortOrder string) (dto.Pagination, error)
}
type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, department models.Department) error
	GetListDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Department, int64, error)
	GetOneDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Department, error)
	UpdateDepartmentByID(ctx context.Context, departmentID string, update models.Department) (*models.Department, error)
	SoftDeleteDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) error
	GetAllDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error)
}
