package ports

import "github.com/Be2Bag/erp-demo/dto"

type DropDownService interface {
	GetPositions() ([]dto.ResponseGetPositions, error)
	GetDepartments() ([]dto.ResponseGetDepartments, error)
}
