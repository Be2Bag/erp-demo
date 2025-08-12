package repositories

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) ports.UserRepository {
	return &userRepo{coll: db.Collection(models.CollectionUsers)}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, map[string]interface{}{"user_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.coll.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) AggregateUser(ctx context.Context, pipeline mongo.Pipeline) ([]*models.User, error) {
	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) CountUsers(ctx context.Context, filter interface{}) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}

func (r *userRepo) UpdateUserByID(ctx context.Context, id string, user *models.User) (*models.User, error) {
	filter := map[string]interface{}{"user_id": id}
	update := map[string]interface{}{
		"$set": user,
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return user, nil
}

func (r *userRepo) GetUserByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.User, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) UpdateUserByFilter(ctx context.Context, filter interface{}, update interface{}) (*models.User, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var user models.User
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}
