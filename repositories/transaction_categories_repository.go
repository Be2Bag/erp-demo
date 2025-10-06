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

type transactionCategoryRepo struct {
	coll *mongo.Collection
}

func NewTransactionCategoryRepository(db *mongo.Database) ports.TransactionCategoryRepository {
	return &transactionCategoryRepo{coll: db.Collection(models.CollectionTransactionCategory)}
}

func (r *transactionCategoryRepo) CreateTransactionCategory(ctx context.Context, transactionCategory models.TransactionCategory) error {
	_, err := r.coll.InsertOne(ctx, transactionCategory)
	return err
}

func (r *transactionCategoryRepo) GetListTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.TransactionCategory, int64, error) {

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

	var results []models.TransactionCategory
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *transactionCategoryRepo) GetOneTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.TransactionCategory, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var category models.TransactionCategory
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *transactionCategoryRepo) UpdateTransactionCategoryByID(ctx context.Context, transactionCategoryID string, update models.TransactionCategory) (*models.TransactionCategory, error) {
	filter := bson.M{"transaction_category_id": transactionCategoryID}
	set := bson.M{
		"type":                         update.Type,
		"transaction_category_name_th": update.TransactionCategoryNameTH,
		"description":                  update.Description,
		"created_by":                   update.CreatedBy,
		"note":                         update.Note,
		"updated_at":                   update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.TransactionCategory
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *transactionCategoryRepo) SoftDeleteTransactionCategoryByID(ctx context.Context, transactionCategoryID string, claims *dto.JWTClaims) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"transaction_category_id": transactionCategoryID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *transactionCategoryRepo) GetAllTransactionCategoryByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.TransactionCategory, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.TransactionCategory
	for cursor.Next(ctx) {
		var category models.TransactionCategory
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
