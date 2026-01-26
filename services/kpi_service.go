package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type kpiService struct {
	kpiRepo  ports.KPIRepository
	userRepo ports.UserRepository
	config   config.Config
}

func NewKPIService(cfg config.Config, kpiRepo ports.KPIRepository, userRepo ports.UserRepository) ports.KPIService {
	return &kpiService{config: cfg, kpiRepo: kpiRepo, userRepo: userRepo}
}

func (s *kpiService) CreateKPITemplate(ctx context.Context, req dto.CreateKPITemplateDTO, claims *dto.JWTClaims) error {

	now := time.Now()
	// --- validate basic fields ---
	if strings.TrimSpace(req.KPIName) == "" {
		return fmt.Errorf("kpi_name is required")
	}
	if strings.TrimSpace(req.Department) == "" {
		return fmt.Errorf("department is required")
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("items must not be empty")
	}

	// --- validate items ---
	sumWeight := 0
	seen := make(map[string]struct{}) // กันชื่อซ้ำใน template เดียวกัน (normalize เป็น lower+trim)
	for i, it := range req.Items {
		itName := strings.TrimSpace(it.Name)
		if itName == "" {
			return fmt.Errorf("items[%d].name is required", i)
		}
		key := strings.ToLower(itName)
		if _, ok := seen[key]; ok {
			return fmt.Errorf("items[%d].name duplicated: %s", i, itName)
		}
		seen[key] = struct{}{}

		if it.MaxScore <= 0 {
			return fmt.Errorf("items[%d].max_score must be > 0", i)
		}
		if it.Weight <= 0 {
			return fmt.Errorf("items[%d].weight must be > 0", i)
		}
		sumWeight += it.Weight
	}
	if sumWeight != 100 {
		return fmt.Errorf("sum of weights must be 100, got %d", sumWeight)
	}

	filter := bson.M{
		"kpi_name":      req.KPIName,
		"department_id": req.Department,
		"deleted_at":    nil,
	}
	opts := options.Find().SetProjection(bson.M{"_id": 1})

	existing, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, opts)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("template with the same name already exists in this department")
	}

	// ---- build model ----
	var items []models.KPITemplateItem
	items = make([]models.KPITemplateItem, 0, len(req.Items))
	for _, kpi := range req.Items {
		items = append(items, models.KPITemplateItem{
			ItemID:      uuid.NewString(),
			Name:        strings.TrimSpace(kpi.Name),
			Description: strings.TrimSpace(kpi.Description),
			Category:    strings.TrimSpace(kpi.Category),
			MaxScore:    kpi.MaxScore,
			Weight:      kpi.Weight,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	tplID := uuid.NewString()
	doc := models.KPITemplate{
		KPIID:       tplID,
		KPIName:     req.KPIName,
		Department:  req.Department,
		TotalWeight: 100,
		Items:       items,
		IsActive:    true, // ตั้งค่าเอง ไม่เชื่อ client
		Version:     1,
		CreatedBy:   claims.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.kpiRepo.CreateKPI(ctx, doc); err != nil {
		return err
	}

	return nil
}

func (s *kpiService) GetKPITemplateByID(ctx context.Context, kpiID string) (*dto.KPITemplateDTO, error) {

	filter := bson.M{"kpi_id": kpiID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	ItemsDTO := make([]dto.KPITemplateItemDTO, 0, len(m.Items))
	for _, st := range m.Items {
		ItemsDTO = append(ItemsDTO, dto.KPITemplateItemDTO{
			ItemID:      st.ItemID,
			Name:        st.Name,
			Description: st.Description,
			Category:    st.Category,
			MaxScore:    st.MaxScore,
			Weight:      st.Weight,
			CreatedAt:   st.CreatedAt,
			UpdatedAt:   st.UpdatedAt,
		})
	}

	dtoObj := &dto.KPITemplateDTO{
		KPIID:       m.KPIID,
		KPIName:     m.KPIName,
		Department:  m.Department,
		TotalWeight: m.TotalWeight,
		Items:       ItemsDTO,
		Version:     m.Version,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	return dtoObj, nil
}

func (s *kpiService) UpdateKPITemplate(ctx context.Context, kpiID string, req dto.UpdateKPITemplateDTO, claims *dto.JWTClaims) error {

	now := time.Now()
	// ดึงข้อมูลเดิม
	filter := bson.M{"kpi_id": kpiID, "deleted_at": nil}
	existing, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	// อัปเดตฟิลด์ข้อความ (trim)
	if strings.TrimSpace(req.KPIName) != "" {
		existing.KPIName = strings.TrimSpace(req.KPIName)
	}
	if strings.TrimSpace(req.Department) != "" {
		existing.Department = strings.TrimSpace(req.Department)
	}

	// ถ้า client ส่ง items มา (pointer != nil) → validate + แทนที่ทั้งชุด
	if req.Items != nil {
		if len(*req.Items) == 0 {
			return fmt.Errorf("items must not be empty")
		}
		sumWeight := 0
		seen := make(map[string]struct{})
		newItems := make([]models.KPITemplateItem, 0, len(*req.Items))

		for i, it := range *req.Items {
			name := strings.TrimSpace(it.Name)
			if name == "" {
				return fmt.Errorf("items[%d].name is required", i)
			}
			key := strings.ToLower(name)
			if _, ok := seen[key]; ok {
				return fmt.Errorf("items[%d].name duplicated: %s", i, name)
			}
			seen[key] = struct{}{}

			if it.MaxScore <= 0 {
				return fmt.Errorf("items[%d].max_score must be > 0", i)
			}
			if it.Weight <= 0 {
				return fmt.Errorf("items[%d].weight must be > 0", i)
			}
			sumWeight += it.Weight

			newItems = append(newItems, models.KPITemplateItem{
				ItemID:      uuid.NewString(),
				Name:        name,
				Description: strings.TrimSpace(it.Description),
				Category:    strings.TrimSpace(it.Category),
				MaxScore:    it.MaxScore,
				Weight:      it.Weight,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			})
		}
		// ให้ข้อความตรงกับ handler ตอนนี้ (เช็คเท่ากับ 100 ตรง ๆ)
		if sumWeight != 100 {
			return fmt.Errorf("sum of weights must be 100")
		}

		existing.Items = newItems
		existing.TotalWeight = 100 // คงเป็น 100 เสมอ เพื่อความชัดเจน
	}

	existing.Version += 1    // เพิ่มเวอร์ชัน
	existing.UpdatedAt = now // อัปเดตเวลา

	updated, err := s.kpiRepo.UpdateKPIByID(ctx, kpiID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *kpiService) DeleteKPITemplate(ctx context.Context, kpiID string) error {
	err := s.kpiRepo.SoftDeleteKPIByID(ctx, kpiID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *kpiService) ListKPITemplates(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department_id"] = department
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"workflow_name": re},
		}
	}

	projection := bson.M{}

	// sort
	allowedSortFields := map[string]string{
		"created_at":    "created_at",
		"updated_at":    "updated_at",
		"workflow_name": "workflow_name",
	}
	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = "created_at"
	}
	order := int32(-1)
	if strings.EqualFold(sortOrder, "asc") {
		order = 1
	}

	sort := bson.D{
		{Key: field, Value: order},
		{Key: "_id", Value: -1},
	}

	items, total, err := s.kpiRepo.GetListKPIByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list kpi templates: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {
		ItemsDTO := make([]dto.KPITemplateItemDTO, 0, len(m.Items))
		for _, st := range m.Items {
			ItemsDTO = append(ItemsDTO, dto.KPITemplateItemDTO{
				ItemID:      st.ItemID,
				Name:        st.Name,
				Description: st.Description,
				Category:    st.Category,
				MaxScore:    st.MaxScore,
				Weight:      st.Weight,
				CreatedAt:   st.CreatedAt,
				UpdatedAt:   st.UpdatedAt,
			})
		}
		list = append(list, dto.KPITemplateDTO{
			KPIID:       m.KPIID,
			KPIName:     m.KPIName,
			Department:  m.Department,
			TotalWeight: m.TotalWeight,
			Items:       ItemsDTO,
			Version:     m.Version,
			CreatedBy:   m.CreatedBy,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}

	totalPages := 0
	if total > 0 && size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	return dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}, nil
}
