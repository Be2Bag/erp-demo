package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

type DropDownService interface {
	GetPositions(ctx context.Context, departmentID string) ([]dto.ResponseGetPositions, error)
	GetDepartments(ctx context.Context) ([]dto.ResponseGetDepartments, error)
	GetProvinces(ctx context.Context) ([]dto.ResponseGetProvinces, error)
	GetDistricts(ctx context.Context, provinceID string) ([]dto.ResponseGetDistricts, error)
	GetSubDistricts(ctx context.Context, districtID string) ([]dto.ResponseGetSubDistricts, error)
	GetSignTypes(ctx context.Context) ([]dto.ResponseGetSignTypes, error)
	GetCustomerTypes(ctx context.Context) ([]dto.ResponseGetCustomerTypes, error)
}

type DropDownRepository interface {
	GetPositions(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error)
	GetDepartments(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error)
	GetProvinces(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Province, error)
	GetDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.District, error)
	GetSubDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SubDistrict, error)
	GetSignTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignType, error)
	GetCustomerTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.CustomerType, error)
}
