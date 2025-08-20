package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type AdminService interface {
	UpdateUserStatus(ctx context.Context, req dto.RequestUpdateUserStatus) error
	UpdateUserRole(ctx context.Context, req dto.RequestUpdateUserRole) error
	UpdateUserPosition(ctx context.Context, req dto.RequestUpdateUserPosition) error
}

type AdminRepository interface {
}
