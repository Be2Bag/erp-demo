package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type DropDownService interface {
	GetPositions(ctx context.Context) ([]dto.ResponseGetPositions, error)
	GetDepartments(ctx context.Context) ([]dto.ResponseGetDepartments, error)
	GetProvinces(ctx context.Context) ([]dto.ResponseGetProvinces, error)
	GetDistricts(ctx context.Context, provinceID string) ([]dto.ResponseGetDistricts, error)
	GetSubDistricts(ctx context.Context, districtID string) ([]dto.ResponseGetSubDistricts, error)
}
