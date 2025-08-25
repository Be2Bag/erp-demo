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

type categoryRepo struct {
	coll *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) ports.CategoryRepository {
	return &categoryRepo{coll: db.Collection(models.CollectionCategory)}
}

func (r *categoryRepo) CreateCategory(ctx context.Context, category models.Category) error {
	_, err := r.coll.InsertOne(ctx, category)
	return err
}

func (r *categoryRepo) GetListCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Category, int64, error) {

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

	var results []models.Category
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *categoryRepo) GetOneCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Category, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var category models.Category
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) UpdateCategoryByID(ctx context.Context, categoryID string, update models.Category) (*models.Category, error) {
	filter := bson.M{"category_id": categoryID}
	set := bson.M{
		"department_id":    update.DepartmentID,
		"category_name_th": update.CategoryNameTH,
		"category_name_en": update.CategoryNameEN,
		"description":      update.Description,
		"note":             update.Note,
		"updated_at":       update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Category
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *categoryRepo) SoftDeleteCategoryByID(ctx context.Context, categoryID string, claims *dto.JWTClaims) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"category_id": categoryID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *categoryRepo) GetAllCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Category, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
