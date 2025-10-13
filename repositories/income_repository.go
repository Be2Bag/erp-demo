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

type inComeRepo struct {
	coll *mongo.Collection
}

func NewInComeRepository(db *mongo.Database) ports.InComeRepository {
	return &inComeRepo{coll: db.Collection(models.CollectionIncome)}
}

func (r *inComeRepo) CreateInCome(ctx context.Context, inCome models.Income) error {
	_, err := r.coll.InsertOne(ctx, inCome)
	return err
}

func (r *inComeRepo) UpdateInComeByID(ctx context.Context, incomeID string, update models.Income) (*models.Income, error) {
	filter := bson.M{"income_id": incomeID}
	set := bson.M{
		"transaction_category_id": update.TransactionCategoryID,
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
	var updated models.Income
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *inComeRepo) SoftDeleteInComeByincomeID(ctx context.Context, incomeID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"income_id": incomeID}, bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}})
	return err
}

func (r *inComeRepo) GetAllInComeByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Income, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var inComes []*models.Income
	for cursor.Next(ctx) {
		var inCome models.Income
		if err := cursor.Decode(&inCome); err != nil {
			return nil, err
		}
		inComes = append(inComes, &inCome)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return inComes, nil
}

func (r *inComeRepo) GetOneInComeByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Income, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var inCome models.Income
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&inCome); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &inCome, nil
}

func (r *inComeRepo) GetListInComesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Income, int64, error) {

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

	var results []models.Income
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
