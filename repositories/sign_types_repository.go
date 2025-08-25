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

type signTypeRepo struct {
	coll *mongo.Collection
}

func NewSignTypeRepository(db *mongo.Database) ports.SignTypeRepository {
	return &signTypeRepo{coll: db.Collection(models.CollectionSignTypes)}
}

func (r *signTypeRepo) CreateSignType(ctx context.Context, signType models.SignType) error {
	_, err := r.coll.InsertOne(ctx, signType)
	return err
}

func (r *signTypeRepo) UpdateSignTypeByTypeID(ctx context.Context, typeID string, update models.SignType) (*models.SignType, error) {
	filter := bson.M{"type_id": typeID}
	set := bson.M{
		"name_th":    update.NameTH,
		"name_en":    update.NameEN,
		"updated_at": update.UpdatedAt,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.SignType
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *signTypeRepo) SoftDeleteSignTypeByTypeID(ctx context.Context, typeID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"type_id": typeID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *signTypeRepo) GetAllSignTypeByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignType, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var signTypes []*models.SignType
	for cursor.Next(ctx) {
		var signType models.SignType
		if err := cursor.Decode(&signType); err != nil {
			return nil, err
		}
		signTypes = append(signTypes, &signType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return signTypes, nil
}

func (r *signTypeRepo) GetOneSignTypeByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.SignType, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var signType models.SignType
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&signType); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &signType, nil
}

func (r *signTypeRepo) GetListSignTypesByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.SignType, int64, error) {

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

	var results []models.SignType
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
