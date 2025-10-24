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

type receivableRepo struct {
	coll *mongo.Collection
}

func NewReceivableRepository(db *mongo.Database) ports.ReceivableRepository {
	return &receivableRepo{coll: db.Collection(models.CollectionReceivable)}
}

func (r *receivableRepo) CreateReceivable(ctx context.Context, receivable models.Receivable) error {
	_, err := r.coll.InsertOne(ctx, receivable)
	return err
}

func (r *receivableRepo) UpdateReceivableByID(ctx context.Context, receivableID string, update models.Receivable) (*models.Receivable, error) {
	filter := bson.M{"id_receivable": receivableID}
	set := bson.M{
		"customer":   update.Customer,
		"invoice_no": update.InvoiceNo,
		"issue_date": update.IssueDate,
		"due_date":   update.DueDate,
		"amount":     update.Amount,
		"balance":    update.Balance,
		"status":     update.Status,
		"updated_at": time.Now(),
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Receivable
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *receivableRepo) SoftDeleteReceivableByID(ctx context.Context, receivableID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id_receivable": receivableID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *receivableRepo) GetAllReceivablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Receivable, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var receivables []*models.Receivable
	for cursor.Next(ctx) {
		var receivable models.Receivable
		if err := cursor.Decode(&receivable); err != nil {
			return nil, err
		}
		receivables = append(receivables, &receivable)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return receivables, nil
}

func (r *receivableRepo) GetOneReceivableByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Receivable, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var receivable models.Receivable
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&receivable); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &receivable, nil
}

func (r *receivableRepo) GetListReceivablesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Receivable, int64, error) {

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

	var results []models.Receivable
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
