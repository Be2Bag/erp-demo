package repositories

import (
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpLoadRepo struct {
	coll *mongo.Collection
}

func NewUpLoadRepository(db *mongo.Database) ports.UpLoadRepository {
	return &UpLoadRepo{coll: db.Collection("uploads")}
}
