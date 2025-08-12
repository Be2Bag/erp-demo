package repositories

import (
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type adminRepo struct {
	coll *mongo.Collection
}

func NewAdminRepository(db *mongo.Database) ports.AdminRepository {
	return &adminRepo{coll: db.Collection(models.CollectionUsers)}
}
