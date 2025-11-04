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
	collReceivables *mongo.Collection
	collPaymentsTx  *mongo.Collection
}

func NewReceivableRepository(db *mongo.Database) ports.ReceivableRepository {
	return &receivableRepo{
		collReceivables: db.Collection(models.CollectionReceivable),
		collPaymentsTx:  db.Collection(models.CollectionPaymentTransaction),
	}
}

func (r *receivableRepo) CreateReceivable(ctx context.Context, receivable models.Receivable) error {
	_, err := r.collReceivables.InsertOne(ctx, receivable)
	return err
}

func (r *receivableRepo) UpdateReceivableByID(ctx context.Context, receivableID string, update models.Receivable) (*models.Receivable, error) {
	filter := bson.M{"id_receivable": receivableID}
	set := bson.M{
		"customer":   update.Customer,
		"bank_id":    update.BankID,
		"invoice_no": update.InvoiceNo,
		"issue_date": update.IssueDate,
		"due_date":   update.DueDate,
		"amount":     update.Amount,
		"balance":    update.Balance,
		"status":     update.Status,
		"note":       update.Note,
		"phone":      update.Phone,
		"address":    update.Address,
		"updated_at": time.Now(),
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Receivable
	if err := r.collReceivables.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *receivableRepo) SoftDeleteReceivableByID(ctx context.Context, receivableID string) error {
	_, err := r.collReceivables.UpdateOne(ctx, bson.M{"id_receivable": receivableID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *receivableRepo) GetAllReceivablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Receivable, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.collReceivables.Find(ctx, filter, opts)
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
	if err := r.collReceivables.FindOne(ctx, filter, opts).Decode(&receivable); err != nil {
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

	cur, err := r.collReceivables.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find: %w", err)
	}
	defer cur.Close(ctx)

	var results []models.Receivable
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.collReceivables.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *receivableRepo) CreatePaymentTransaction(ctx context.Context, tx models.PaymentTransaction) error {
	_, err := r.collPaymentsTx.InsertOne(ctx, tx)
	return err
}

func (r *receivableRepo) GetAllPaymentTransactionsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.PaymentTransaction, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.collPaymentsTx.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*models.PaymentTransaction
	for cursor.Next(ctx) {
		var transaction models.PaymentTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
