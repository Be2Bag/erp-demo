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

type payableRepo struct {
	collPayables   *mongo.Collection
	collPaymentsTx *mongo.Collection
}

func NewPayableRepository(db *mongo.Database) ports.PayableRepository {
	return &payableRepo{
		collPayables:   db.Collection(models.CollectionPayable),
		collPaymentsTx: db.Collection(models.CollectionPaymentTransaction),
	}
}

func (r *payableRepo) CreatePayable(ctx context.Context, payable models.Payable) error {
	_, err := r.collPayables.InsertOne(ctx, payable)
	return err
}

func (r *payableRepo) UpdatePayableByID(ctx context.Context, payableID string, update models.Payable) (*models.Payable, error) {
	filter := bson.M{"id_payable": payableID}
	set := bson.M{
		"supplier":    update.Supplier,
		"purchase_no": update.PurchaseNo,
		"invoice_no":  update.InvoiceNo,
		"issue_date":  update.IssueDate,
		"due_date":    update.DueDate,
		"amount":      update.Amount,
		"balance":     update.Balance,
		"status":      update.Status,
		"payment_ref": update.PaymentRef,
		"note":        update.Note,
		"updated_at":  time.Now(),
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Payable
	if err := r.collPayables.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *payableRepo) SoftDeletePayableByID(ctx context.Context, payableID string) error {
	_, err := r.collPayables.UpdateOne(ctx, bson.M{"id_payable": payableID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *payableRepo) GetAllPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Payable, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.collPayables.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payables []*models.Payable
	for cursor.Next(ctx) {
		var payable models.Payable
		if err := cursor.Decode(&payable); err != nil {
			return nil, err
		}
		payables = append(payables, &payable)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return payables, nil
}

func (r *payableRepo) GetOnePayableByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Payable, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var payable models.Payable
	if err := r.collPayables.FindOne(ctx, filter, opts).Decode(&payable); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &payable, nil
}

func (r *payableRepo) GetListPayablesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Payable, int64, error) {

	findOpts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(limit)

	if projection != nil {
		findOpts.SetProjection(projection)
	}

	cur, err := r.collPayables.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find: %w", err)
	}
	defer cur.Close(ctx)

	var results []models.Payable
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.collPayables.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *payableRepo) CreatePaymentTransaction(ctx context.Context, tx models.PaymentTransaction) error {
	_, err := r.collPaymentsTx.InsertOne(ctx, tx)
	return err
}

func (r *payableRepo) GetAllPaymentTransactionByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.PaymentTransaction, error) {
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
