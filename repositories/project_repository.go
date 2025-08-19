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

type projectRepo struct {
	coll *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) ports.ProjectRepository {
	return &projectRepo{coll: db.Collection(models.CollectionSProject)}
}

func (r *projectRepo) CreateProject(ctx context.Context, project models.Project) error {
	_, err := r.coll.InsertOne(ctx, project)
	return err
}

func (r *projectRepo) GetListProjectByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.Project, int64, error) {

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

	var results []models.Project
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}

func (r *projectRepo) GetOneProjectByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.Project, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var project models.Project
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&project); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) UpdateProjectByID(ctx context.Context, projectID string, update models.Project) (*models.Project, error) {
	filter := bson.M{"project_id": projectID}
	set := bson.M{
		"project_name": update.ProjectName,
		"updated_at":   update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Project
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (s *projectRepo) SoftDeleteProjectByID(ctx context.Context, projectID string, claims *dto.JWTClaims) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"project_id": projectID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}
