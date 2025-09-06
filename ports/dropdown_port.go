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
	GetSignJobList(ctx context.Context, projectID string) ([]dto.ResponseGetSignList, error)
	GetProjectList(ctx context.Context) ([]dto.ResponseGetProjects, error)
	GetUserList(ctx context.Context, usersID string) ([]dto.ResponseGetUsers, error)
	GetKPI(ctx context.Context, departmentID string) ([]dto.ResponseGetKPI, error)
	GetWorkflows(ctx context.Context, departmentID string) ([]dto.ResponseGetWorkflows, error)
	GetCategorys(ctx context.Context) ([]dto.ResponseGetCategorys, error)
}

type DropDownRepository interface {
	GetPositions(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error)
	GetDepartments(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error)
	GetProvinces(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Province, error)
	GetDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.District, error)
	GetSubDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SubDistrict, error)
	GetSignTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignType, error)
	GetCustomerTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.CustomerType, error)
	GetSignJobsList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignJob, error)
	GetProjectsList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Project, error)
	GetUsersList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.User, error)
	GetKPIList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.KPITemplate, error)
	GetWorkflowsList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.WorkFlowTemplate, error)
	GetCategorysList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Category, error)
}
