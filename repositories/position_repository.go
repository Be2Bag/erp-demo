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

type positionRepo struct {
	coll *mongo.Collection
}

func NewPositionRepository(db *mongo.Database) ports.PositionRepository {
	return &positionRepo{coll: db.Collection(models.CollectionPositions)}
}

func (r *positionRepo) CreatePosition(ctx context.Context, position models.Position) error {
	_, err := r.coll.InsertOne(ctx, position)
	return err
}

func (r *positionRepo) GetListPositionByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Position, int64, error) {

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

	var results []models.Position
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *positionRepo) GetOnePositionByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Position, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var position models.Position
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&position); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &position, nil
}

func (r *positionRepo) UpdatePositionByID(ctx context.Context, positionID string, update models.Position) (*models.Position, error) {
	filter := bson.M{"position_id": positionID}
	set := bson.M{
		"position_name": update.PositionName,
		"level":         update.Level,
		"department_id": update.DepartmentID,
		"note":          update.Note,
		"updated_at":    update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Position
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *positionRepo) SoftDeletePositionByID(ctx context.Context, positionID string, claims *dto.JWTClaims) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"position_id": positionID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *positionRepo) GetAllPositionByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []*models.Position
	for cursor.Next(ctx) {
		var position models.Position
		if err := cursor.Decode(&position); err != nil {
			return nil, err
		}
		positions = append(positions, &position)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return positions, nil
}
