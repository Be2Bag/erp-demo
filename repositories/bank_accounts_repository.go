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

type BankAccountsRepo struct {
	coll *mongo.Collection
}

func NewBankAccountsRepository(db *mongo.Database) ports.BankAccountsRepository {
	return &BankAccountsRepo{coll: db.Collection(models.CollectionSBankAccounts)}
}

func (r *BankAccountsRepo) CreateBankAccount(ctx context.Context, bankAccount models.BankAccount) error {
	_, err := r.coll.InsertOne(ctx, bankAccount)
	return err
}

func (r *BankAccountsRepo) UpdateBankAccountByID(ctx context.Context, id string, update models.BankAccount) (*models.BankAccount, error) {
	filter := bson.M{"bank_id": id}
	set := bson.M{
		"bank_name":    update.BankName,
		"account_no":   update.AccountNo,
		"account_name": update.AccountName,
		"updated_at":   update.UpdatedAt,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.BankAccount
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *BankAccountsRepo) SoftDeleteBankAccountByID(ctx context.Context, id string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"bank_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *BankAccountsRepo) GetAllBankAccountsByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.BankAccount, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bankAccounts []*models.BankAccount
	for cursor.Next(ctx) {
		var bankAccount models.BankAccount
		if err := cursor.Decode(&bankAccount); err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, &bankAccount)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bankAccounts, nil
}

func (r *BankAccountsRepo) GetOneBankAccountByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.BankAccount, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var bankAccount models.BankAccount
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&bankAccount); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bankAccount, nil
}

func (r *BankAccountsRepo) GetListBankAccountsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.BankAccount, int64, error) {

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

	var results []models.BankAccount
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
