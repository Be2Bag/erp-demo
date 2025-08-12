package repositories

import (
	"context"
	"errors"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoopt "go.mongodb.org/mongo-driver/mongo/options"
)

type workFlowRepo struct {
	coll *mongo.Collection
}

func NewWorkFlowRepository(db *mongo.Database) ports.WorkFlowRepository {
	return &workFlowRepo{coll: db.Collection(models.CollectionWorkflowTemplates)}
}

func (r *workFlowRepo) GetWorkFlowTemplates(ctx context.Context, filter interface{}, options interface{}) ([]models.WorkFlowTemplate, error) {
	var templates []models.WorkFlowTemplate
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
		var template models.WorkFlowTemplate
		if err := cursor.Decode(&template); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (r *workFlowRepo) CreateWorkFlowTemplate(ctx context.Context, tmpl *models.WorkFlowTemplate) error {
	_, err := r.coll.InsertOne(ctx, tmpl)
	return err
}

func (r *workFlowRepo) GetWorkFlowTemplateByTemplateID(ctx context.Context, templateID string) (*models.WorkFlowTemplate, error) {
	var tmpl models.WorkFlowTemplate
	if err := r.coll.FindOne(ctx, bson.M{"template_id": templateID}).Decode(&tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *workFlowRepo) CountWorkFlowTemplates(ctx context.Context, filter interface{}) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}

func (r *workFlowRepo) UpdateWorkFlowTemplateByTemplateID(ctx context.Context, templateID string, update interface{}) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"template_id": templateID}, update)
	return err
}

func (r *workFlowRepo) DeleteWorkFlowTemplateByTemplateID(ctx context.Context, templateID string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"template_id": templateID})
	return err
}
