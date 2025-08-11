package service

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
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
	uuid := uuid.New().String()
	var templates []models.KPITemplateList
	for _, kpi := range req.KPIs {
		templates = append(templates, models.KPITemplateList{
			KPIID:       uuid,
			Name:        kpi.Name,
			Description: kpi.Description,
			Category:    kpi.Category,
			MaxScore:    kpi.MaxScore,
			Value:       kpi.Weight,
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		})
	}

	kpiTemplate := models.KPITemplate{
		KPIID:      uuid,
		Name:       req.Name,
		Department: req.Department,
		Templates:  templates,
		CreatedBy:  claims.UserID,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  nil,
	}

	return s.kpiRepo.CreateKPITemplate(ctx, kpiTemplate)
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
