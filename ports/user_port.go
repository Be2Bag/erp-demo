package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	Create(ctx context.Context, u dto.RequestCreateUser) error
	GetByID(ctx context.Context, id string) (*dto.ResponseGetUserByID, error)
	GetAll(ctx context.Context, req dto.RequestGetUserAll) (dto.Pagination, error)
	UpdateUserByID(ctx context.Context, id string, req dto.RequestUpdateUser) (*models.User, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateDocuments(ctx context.Context, req dto.RequestUpdateDocuments) (*models.User, error)
	CountUsers(ctx context.Context) (dto.ResponseGetCountUsers, error)
}
type UserRepository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	AggregateUser(ctx context.Context, pipeline mongo.Pipeline) ([]*models.User, error)
	CountUsers(ctx context.Context, filter interface{}) (int64, error)
	UpdateUserByID(ctx context.Context, id string, user *models.User) (*models.User, error)
	GetUserByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.User, error)
	UpdateUserByFilter(ctx context.Context, filter interface{}, update interface{}) (*models.User, error)
}
