package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type kpiService struct {
	config   config.Config
	kpiRepo  ports.KPIRepository
	userRepo ports.UserRepository
}

func NewKPIService(cfg config.Config, kpiRepo ports.KPIRepository, userRepo ports.UserRepository) ports.KPIService {
	// สร้างและคืนค่าอินสแตนซ์ใหม่ของ kpiService
	return &kpiService{config: cfg, kpiRepo: kpiRepo, userRepo: userRepo}
}

func (s *kpiService) CreateKPITemplate(ctx context.Context, req dto.KPITemplateDTO, claims *dto.JWTClaims) error {

	now := time.Now()

	if len(req.Items) == 0 {
		return fmt.Errorf("items must not be empty")
	}
	sumWeight := 0
	for i, it := range req.Items {
		if it.MaxScore <= 0 {
			return fmt.Errorf("items[%d].max_score must be > 0", i)
		}
		if it.Weight <= 0 {
			return fmt.Errorf("items[%d].weight must be > 0", i)
		}
		sumWeight += it.Weight
	}
	// บังคับ = 100 (หรือจะตรวจเท่ากับ req.TotalWeight ก็ได้)
	if sumWeight != 100 {
		return fmt.Errorf("sum of weights must be 100, got %d", sumWeight)
	}

	// (optional) ป้องกันชื่อซ้ำในแผนก
	filter := bson.M{
		"name":       req.Name,
		"department": req.Department,
	}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	existing, err := s.kpiRepo.GetKPITemplates(ctx, filter, opts)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("template with the same name already exists in this department")
	}

	// ---- build model ----
	var items []models.KPITemplateItem
	items = make([]models.KPITemplateItem, 0, len(req.Items))
	for _, kpi := range req.Items {
		items = append(items, models.KPITemplateItem{
			ItemID:      uuid.NewString(),
			Name:        kpi.Name,
			Description: kpi.Description,
			Category:    kpi.Category,
			MaxScore:    kpi.MaxScore,
			Weight:      kpi.Weight,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	tplID := uuid.NewString()
	doc := models.KPITemplate{
		TemplateID:  tplID,
		Name:        req.Name,
		Department:  req.Department,
		TotalWeight: 100, // เก็บค่าคงที่ 100 ชัดเจน
		Items:       items,
		IsActive:    true, // ตั้งค่าเอง ไม่เชื่อ client
		Version:     1,
		CreatedBy:   claims.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.kpiRepo.CreateKPITemplate(ctx, doc); err != nil {
		return err
	}

	return nil
}

func (s *kpiService) GetKPITemplateByID(ctx context.Context, id string) (interface{}, error) {
	tpl, err := s.kpiRepo.GetKPITemplateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *kpiService) UpdateKPITemplate(ctx context.Context, id string, updated dto.KPITemplateDTO, claims *dto.JWTClaims) (interface{}, error) {
	// fetch existing
	existing, err := s.kpiRepo.GetKPITemplateByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// validate items
	if len(updated.Items) == 0 {
		return nil, errors.New("items must not be empty")
	}
	sum := 0
	now := time.Now()
	items := make([]models.KPITemplateItem, 0, len(updated.Items))
	for i, it := range updated.Items {
		if it.MaxScore <= 0 {
			return nil, fmt.Errorf("items[%d].max_score must be > 0", i)
		}
		if it.Weight <= 0 {
			return nil, fmt.Errorf("items[%d].weight must be > 0", i)
		}
		sum += it.Weight
		items = append(items, models.KPITemplateItem{
			ItemID:      uuid.NewString(),
			Name:        it.Name,
			Description: it.Description,
			Category:    it.Category,
			MaxScore:    it.MaxScore,
			Weight:      it.Weight,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	if sum != 100 {
		return nil, errors.New("sum of weights must be 100")
	}

	// duplicate name in same department (exclude self)
	filter := bson.M{
		"name":       updated.Name,
		"department": updated.Department,
		"template_id": bson.M{
			"$ne": id,
		},
	}
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	exist, err := s.kpiRepo.GetKPITemplates(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if len(exist) > 0 {
		return nil, errors.New("template with the same name already exists in this department")
	}

	existing.Name = updated.Name
	existing.Department = updated.Department
	existing.Items = items
	existing.TotalWeight = 100
	existing.Version = existing.Version + 1
	existing.UpdatedAt = now
	// keep IsActive and CreatedBy as-is (optionally could audit claims.UserID)

	updatedTpl, err := s.kpiRepo.UpdateKPITemplate(ctx, id, *existing)
	if err != nil {
		return nil, err
	}

	return updatedTpl, nil
}

func (s *kpiService) DeleteKPITemplate(ctx context.Context, id string) error {
	// ensure exists first for clearer error (optional)
	_, err := s.kpiRepo.GetKPITemplateByID(ctx, id)
	if err != nil {
		return err
	}
	return s.kpiRepo.DeleteKPITemplate(ctx, id)
}

// added: list with search + pagination
func (s *kpiService) ListKPITemplates(ctx context.Context, q dto.KPITemplateListQuery) ([]interface{}, int64, error) {
	filter := bson.M{}
	if q.Search != "" {
		filter["name"] = bson.M{"$regex": q.Search, "$options": "i"}
	}
	if q.Department != "" {
		filter["department"] = q.Department
	}
	if q.IsActive != nil {
		filter["is_active"] = *q.IsActive
	}

	total, err := s.kpiRepo.CountKPITemplates(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []interface{}{}, 0, nil
	}

	skip := int64((q.Page - 1) * q.Limit)
	limit := int64(q.Limit)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1})

	list, err := s.kpiRepo.GetKPITemplates(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	// convert to []interface{} to keep interface signature minimal change
	out := make([]interface{}, 0, len(list))
	for i := range list {
		out = append(out, list[i])
	}

	return out, total, nil
}
