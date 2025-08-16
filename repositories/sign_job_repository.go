package repositories

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signJobRepo struct {
	coll *mongo.Collection
}

func NewSignJobRepository(db *mongo.Database) ports.SignJobRepository {
	return &signJobRepo{coll: db.Collection(models.CollectionSignJobs)}
}

func (r *signJobRepo) CreateSignJob(ctx context.Context, signJob models.SignJob) error {
	_, err := r.coll.InsertOne(ctx, signJob)
	return err
}

func (r *signJobRepo) ListSignJobs(ctx context.Context, createdBy string) ([]models.SignJob, error) {
	filter := bson.M{"created_by": createdBy}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []models.SignJob
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *signJobRepo) GetSignJobByJobID(ctx context.Context, jobID string, createdBy string) (*models.SignJob, error) {
	filter := bson.M{"job_id": jobID, "created_by": createdBy}
	var m models.SignJob
	if err := r.coll.FindOne(ctx, filter).Decode(&m); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *signJobRepo) UpdateSignJobByJobID(ctx context.Context, jobID string, createdBy string, update models.SignJob) (*models.SignJob, error) {
	filter := bson.M{"job_id": jobID, "created_by": createdBy}
	set := bson.M{
		"project_name":     update.ProjectName,
		"job_name":         update.JobName,
		"customer_name":    update.CustomerName,
		"contact_person":   update.ContactPerson,
		"phone":            update.Phone,
		"email":            update.Email,
		"customer_type_id": update.CustomerTypeID,
		"address":          update.Address,
		"sign_type_id":     update.SignTypeID,
		"size":             update.Size,
		"quantity":         update.Quantity,
		"content":          update.Content,
		"main_color":       update.MainColor,
		"design_option":    update.DesignOption,
		"production_time":  update.ProductionTime,
		"due_date":         update.DueDate,
		"install_option":   update.InstallOption,
		"notes":            update.Notes,
		"status":           update.Status,
		"updated_at":       update.UpdatedAt,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.SignJob
	if err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": set}, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &updated, nil
}

func (r *signJobRepo) DeleteSignJobByJobID(ctx context.Context, jobID string, createdBy string) error {
	filter := bson.M{"job_id": jobID, "created_by": createdBy}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
