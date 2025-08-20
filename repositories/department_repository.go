package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type departmentRepo struct {
	coll *mongo.Collection
}

func NewDepartmentRepository(db *mongo.Database) ports.DepartmentRepository {
	return &departmentRepo{coll: db.Collection(models.CollectionDepartments)}
}

func (r *departmentRepo) CreateDepartment(ctx context.Context, department models.Department) error {
	_, err := r.coll.InsertOne(ctx, department)
	return err
}

func (r *departmentRepo) GetListDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Department, int64, error) {

	findOpts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(limit)

	if projection != nil {
		findOpts.SetProjection(projection)
	}

	cur, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find: %w", err)
	}
	defer cur.Close(ctx)

	var results []models.Department
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *departmentRepo) GetOneDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Department, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var department models.Department
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&department); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepo) UpdateDepartmentByID(ctx context.Context, departmentID string, update models.Department) (*models.Department, error) {
	filter := bson.M{"department_id": departmentID}
	set := bson.M{
		"department_name": update.DepartmentName,
		"manager_id":      update.ManagerID,
		"updated_at":      update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Department
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *departmentRepo) SoftDeleteDepartmentByID(ctx context.Context, departmentID string, claims *dto.JWTClaims) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"department_id": departmentID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *departmentRepo) GetAllDepartmentByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var departments []*models.Department
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
