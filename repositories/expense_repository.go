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

type expenseRepo struct {
	coll *mongo.Collection
}

func NewExpenseRepository(db *mongo.Database) ports.ExpenseRepository {
	return &expenseRepo{coll: db.Collection(models.CollectionExpense)}
}

func (r *expenseRepo) CreateExpense(ctx context.Context, expense models.Expense) error {
	_, err := r.coll.InsertOne(ctx, expense)
	return err
}

func (r *expenseRepo) UpdateExpenseByID(ctx context.Context, expenseID string, update models.Expense) (*models.Expense, error) {
	filter := bson.M{"expense_id": expenseID}
	set := bson.M{
		"transaction_category_id": update.TransactionCategoryID,
		"bank_id":                 update.BankID,
		"description":             update.Description,
		"amount":                  update.Amount,
		"currency":                update.Currency,
		"txn_date":                update.TxnDate,
		"payment_method":          update.PaymentMethod,
		"reference_no":            update.ReferenceNo,
		"note":                    update.Note,
		"updated_at":              time.Now(),
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Expense
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *expenseRepo) SoftDeleteExpenseByID(ctx context.Context, expenseID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"expense_id": expenseID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *expenseRepo) GetAllExpenseByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Expense, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []*models.Expense
	for cursor.Next(ctx) {
		var expense models.Expense
		if err := cursor.Decode(&expense); err != nil {
			return nil, err
		}
		expenses = append(expenses, &expense)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *expenseRepo) GetOneExpenseByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Expense, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var expense models.Expense
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&expense); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepo) GetListExpensesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Expense, int64, error) {

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

	var results []models.Expense
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
