package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
)

type DropDownRepository interface {
	GetPositions(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error)
	GetDepartments(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error)
	GetProvinces(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Province, error)
	GetDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.District, error)
	GetSubDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SubDistrict, error)
}
