package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoopt "go.mongodb.org/mongo-driver/mongo/options"
)

type kpiRepo struct {
	coll *mongo.Collection
}

func NewKPIRepository(db *mongo.Database) ports.KPIRepository {
	return &kpiRepo{coll: db.Collection(models.CollectionKPITemplates)}
}

func (r *kpiRepo) GetKPITemplates(ctx context.Context, filter interface{}, options interface{}) ([]models.KPITemplate, error) {
	var templates []models.KPITemplate
	var findOptions *mongoopt.FindOptions
	if options != nil {
		var ok bool
		findOptions, ok = options.(*mongoopt.FindOptions)
		if !ok {
			return nil, errors.New("options must be of type *options.FindOptions")
		}
	}

	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(cctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(cctx)

	for cursor.Next(cctx) {
		var template models.KPITemplate
		if err := cursor.Decode(&template); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

// สร้างแม่แบบ KPI ใหม่
func (r *kpiRepo) CreateKPITemplate(ctx context.Context, template models.KPITemplate) error {
	_, err := r.coll.InsertOne(ctx, template)
	return err
}

// ดึงข้อมูลแม่แบบ KPI ตามรหัส (template_id)
func (r *kpiRepo) GetKPITemplateByID(ctx context.Context, id string) (*models.KPITemplate, error) {
	var tpl models.KPITemplate
	err := r.coll.FindOne(ctx, bson.M{"template_id": id}).Decode(&tpl)
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

// อัปเดตแม่แบบ KPI (แทนที่ฟิลด์หลัก + items)
func (r *kpiRepo) UpdateKPITemplate(ctx context.Context, id string, updated models.KPITemplate) (*models.KPITemplate, error) {
	update := bson.M{
		"$set": bson.M{
			"name":         updated.Name,
			"department":   updated.Department,
			"items":        updated.Items,
			"total_weight": updated.TotalWeight,
			"version":      updated.Version,
			"is_active":    updated.IsActive,
			"updated_at":   updated.UpdatedAt,
		},
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"template_id": id}, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return r.GetKPITemplateByID(ctx, id)
}

// ลบแม่แบบ KPI
func (r *kpiRepo) DeleteKPITemplate(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"template_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// added: count documents for pagination
func (r *kpiRepo) CountKPITemplates(ctx context.Context, filter interface{}) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}
