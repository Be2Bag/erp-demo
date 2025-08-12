package repositories

import (
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepo struct {
	coll *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) ports.AuthRepository {
	return &authRepo{coll: db.Collection(models.CollectionUsers)}
}
