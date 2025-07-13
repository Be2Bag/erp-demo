package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
)

type UserService interface {
	Create(ctx context.Context, u dto.RequestCreateUser) error
	GetByID(ctx context.Context, id string) (*dto.ResponseGetUserByID, error)
	GetAll(ctx context.Context, req dto.RequestGetUserAll) (dto.Pagination, error)
	UpdateUserByID(ctx context.Context, id string, req dto.RequestUpdateUser) (*models.User, error)
	DeleteUserByID(ctx context.Context, id string) error
}
