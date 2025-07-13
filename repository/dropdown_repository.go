package repository

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dropDownRepo struct {
	departmentsColl *mongo.Collection
	positionsColl   *mongo.Collection
}

func NewDropDownRepository(db *mongo.Database) ports.DropDownRepository {
	return &dropDownRepo{
		departmentsColl: db.Collection("departments"),
		positionsColl:   db.Collection("positions"),
	}
}
func (r *dropDownRepo) GetPositions(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error) {
	var positions []*models.Position
	cursor, err := r.positionsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var position models.Position
		if err := cursor.Decode(&position); err != nil {
			return nil, err
		}
		positions = append(positions, &position)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *dropDownRepo) GetDepartments(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error) {
	var departments []*models.Department
	cursor, err := r.departmentsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var department models.Department
		if err := cursor.Decode(&department); err != nil {
			return nil, err
		}
		departments = append(departments, &department)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return departments, nil
}
