package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	AggregateUser(ctx context.Context, pipeline mongo.Pipeline) ([]*models.User, error)
	CountUsers(ctx context.Context, filter interface{}) (int64, error)
	UpdateUserByID(ctx context.Context, id string, user *models.User) (*models.User, error)
	GetUserByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.User, error)
}
