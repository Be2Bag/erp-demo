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

func (r *signJobRepo) ListSignJobs(ctx context.Context, page, size int, search string) ([]models.SignJob, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"project_name": bson.M{"$regex": search, "$options": "i"}},
			{"job_name": bson.M{"$regex": search, "$options": "i"}},
			{"company_name": bson.M{"$regex": search, "$options": "i"}},
			{"contact_person": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	findOpts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var results []models.SignJob
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
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
