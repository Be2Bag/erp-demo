package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type AdminService interface {
	UpdateUserStatus(ctx context.Context, req dto.RequestUpdateUserStatus) error
	UpdateUserPosition(ctx context.Context, req dto.RequestUpdateUserPosition) error
}

type AdminRepository interface {
}
