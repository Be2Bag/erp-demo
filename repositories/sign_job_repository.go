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

func (r *signJobRepo) UpdateSignJobByJobID(ctx context.Context, jobID string, update models.SignJob) (*models.SignJob, error) {
	filter := bson.M{"job_id": jobID}
	set := bson.M{
		"company_name":     update.CompanyName,
		"contact_person":   update.ContactPerson,
		"phone":            update.Phone,
		"email":            update.Email,
		"customer_type_id": update.CustomerTypeID,
		"address":          update.Address,
		"project_name":     update.ProjectName,
		"job_name":         update.JobName,
		"sign_type_id":     update.SignTypeID,
		"width":            update.Width,
		"height":           update.Height,
		"quantity":         update.Quantity,
		"price_thb":        update.PriceTHB,
		"content":          update.Content,
		"main_color":       update.MainColor,
		"payment_method":   update.PaymentMethod,
		"production_time":  update.ProductionTime,
		"design_option":    update.DesignOption,
		"install_option":   update.InstallOption,
		"notes":            update.Notes,
		"status":           update.Status,
		"updated_at":       update.UpdatedAt,
	}
	if !update.DueDate.IsZero() {
		set["due_date"] = update.DueDate
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

func (r *signJobRepo) SoftDeleteSignJobByJobID(ctx context.Context, jobID string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"job_id": jobID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *signJobRepo) GetAllSignJobByFilter(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignJob, error) {
	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var signJobs []*models.SignJob
	for cursor.Next(ctx) {
		var signJob models.SignJob
		if err := cursor.Decode(&signJob); err != nil {
			return nil, err
		}
		signJobs = append(signJobs, &signJob)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return signJobs, nil
}

func (r *signJobRepo) GetOneSignJobByFilter(ctx context.Context, filter interface{}, projection interface{}) (*models.SignJob, error) {
	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}
	var signJob models.SignJob
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&signJob); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &signJob, nil
}

func (r *signJobRepo) GetListSignJobsByFilter(ctx context.Context, filter interface{}, projection interface{}, sort bson.D, skip, limit int64) ([]models.SignJob, int64, error) {

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

	var results []models.SignJob
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return results, total, nil
}
