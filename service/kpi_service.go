package service

import (
	"context"
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

func (s *kpiService) GetKPITemplates(ctx context.Context, filter interface{}) ([]interface{}, error) {
	// ฟังก์ชันสำหรับดึงข้อมูล KPI templates
	return nil, nil
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
	// ฟังก์ชันสำหรับดึง KPI template ตาม ID
	return nil, nil
}

func (s *kpiService) UpdateKPITemplate(ctx context.Context, id string, updatedTemplate interface{}) error {
	// ฟังก์ชันสำหรับอัปเดต KPI template
	return nil
}

func (s *kpiService) DeleteKPITemplate(ctx context.Context, id string) error {
	// ฟังก์ชันสำหรับลบ KPI template
	return nil
}
