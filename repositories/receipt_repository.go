package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type receiptRepo struct {
	coll *mongo.Collection
}

func NewReceiptRepository(db *mongo.Database) ports.ReceiptRepository {
	return &receiptRepo{coll: db.Collection(models.CollectionReceipts)}
}

func (r *receiptRepo) CreateReceipt(ctx context.Context, receipt models.Receipt) error {
	_, err := r.coll.InsertOne(ctx, receipt)
	return err
}

func (r *receiptRepo) SoftDeleteReceiptByID(ctx context.Context, receiptID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id_receipt": receiptID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *receiptRepo) GetListReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Receipt, int64, error) {

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

	var results []models.Receipt
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *receiptRepo) GetOneReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Receipt, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var receipt models.Receipt
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&receipt); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &receipt, nil
}

func (r *receiptRepo) GetAllReceiptsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Receipt, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var receipts []*models.Receipt
	for cursor.Next(ctx) {
		var receipt models.Receipt
		if err := cursor.Decode(&receipt); err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *receiptRepo) GetMaxReceiptNumber(ctx context.Context, prefix string) (string, error) {
	filter := bson.M{
		"receipt_number": bson.M{"$regex": fmt.Sprintf("^%s", prefix)},
		"deleted_at":     nil,
	}
	opts := options.FindOne().
		SetSort(bson.D{{Key: "receipt_number", Value: -1}}).
		SetProjection(bson.M{"receipt_number": 1})

	var result models.Receipt
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", fmt.Errorf("find max receipt number: %w", err)
	}
	return result.ReceiptNumber, nil
}

func (r *receiptRepo) UpdateReceiptByID(ctx context.Context, receiptID string, updateData interface{}) error {
	update := bson.M{
		"$set": updateData,
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"id_receipt": receiptID}, update)
	return err
}
