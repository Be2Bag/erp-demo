package repository

import (
	"context"
	"errors"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
	mongoopt "go.mongodb.org/mongo-driver/mongo/options"
)

type kpiRepo struct {
	coll *mongo.Collection
}

func NewKPIRepository(db *mongo.Database) ports.KPIRepository {
	return &kpiRepo{coll: db.Collection("kpi")}
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

	cursor, err := r.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var template models.KPITemplate
		if err := cursor.Decode(&template); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// สร้างแม่แบบ KPI ใหม่
func (r *kpiRepo) CreateKPITemplate(ctx context.Context, template models.KPITemplate) error {
	_, err := r.coll.InsertOne(ctx, template)
	return err
}

// ดึงข้อมูลแม่แบบ KPI ตามรหัส
func (r *kpiRepo) GetKPITemplateByID(ctx context.Context, id string) (interface{}, error) {
	// การนำไปใช้สำหรับดึงข้อมูลแม่แบบ KPI ตามรหัส
	return nil, nil
}

// อัปเดตแม่แบบ KPI
func (r *kpiRepo) UpdateKPITemplate(ctx context.Context, id string, updatedTemplate interface{}) error {
	// การนำไปใช้สำหรับอัปเดตแม่แบบ KPI
	return nil
}

// ลบแม่แบบ KPI
func (r *kpiRepo) DeleteKPITemplate(ctx context.Context, id string) error {
	// การนำไปใช้สำหรับลบแม่แบบ KPI
	return nil
}
