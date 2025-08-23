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
)

type kpiEvaluationRepoService struct {
	config            config.Config
	kpiRepo           ports.KPIRepository
	userRepo          ports.UserRepository
	kpiEvaluationRepo ports.KPIEvaluationRepository
	taskRepo          ports.TaskRepository
	departmentRepo    ports.DepartmentRepository
}

func NewKPIEvaluationService(cfg config.Config, kpiRepo ports.KPIRepository, userRepo ports.UserRepository, kpiEvaluationRepo ports.KPIEvaluationRepository, taskRepo ports.TaskRepository, departmentRepo ports.DepartmentRepository) ports.KPIEvaluationService {
	return &kpiEvaluationRepoService{config: cfg, kpiRepo: kpiRepo, userRepo: userRepo, kpiEvaluationRepo: kpiEvaluationRepo, taskRepo: taskRepo, departmentRepo: departmentRepo}
}

func (s *kpiEvaluationRepoService) CreateKPIEvaluation(ctx context.Context, req dto.CreateKPIEvaluationRequest, claims *dto.JWTClaims) error {
	now := time.Now()

	// 1) ดึง KPI Template ตาม KPIID ที่ส่งมา (ต้องยังไม่ถูกลบ)
	tplFilter := bson.M{"kpi_id": req.KPIID, "deleted_at": nil}
	tpl, err := s.kpiRepo.GetOneKPIByFilter(ctx, tplFilter, bson.M{})
	if err != nil {
		return fmt.Errorf("get kpi template: %w", err)
	}
	if tpl == nil {
		return fmt.Errorf("kpi template not found")
	}

	// ต้องมีคะแนนส่งเข้ามา
	if len(req.Scores) == 0 {
		return fmt.Errorf("scores required")
	}

	// 2) สร้าง map ของ items ใน template เพื่อ lookup เร็ว O(1)
	itemMap := make(map[string]models.KPITemplateItem, len(tpl.Items))
	for _, it := range tpl.Items {
		itemMap[it.ItemID] = it
	}

	// 3) วน validate + สร้างรายการคะแนน (KPIScore) พร้อมคำนวณคะแนนถ่วงน้ำหนักรวม
	scores := make([]models.KPIScore, 0, len(req.Scores))
	seen := make(map[string]struct{})
	totalWeighted := 0 // ผลรวม (score * weight) (สมมติ weight รวม 100)

	for i, sc := range req.Scores {
		// ตรวจ ItemID
		if sc.ItemID == "" {
			return fmt.Errorf("scores[%d].item_id empty", i)
		}
		if _, ok := seen[sc.ItemID]; ok {
			return fmt.Errorf("scores[%d].item_id duplicated: %s", i, sc.ItemID)
		}
		seen[sc.ItemID] = struct{}{}

		// ตรวจว่าอยู่ใน template
		tplItem, ok := itemMap[sc.ItemID]
		if !ok {
			return fmt.Errorf("scores[%d].item_id not in template: %s", i, sc.ItemID)
		}
		// ช่วงคะแนนต้องอยู่ใน 0..MaxScore
		if sc.Score < 0 || sc.Score > tplItem.MaxScore {
			return fmt.Errorf("scores[%d].score out of range (0-%d)", i, tplItem.MaxScore)
		}

		// เพิ่มเข้า slice
		scores = append(scores, models.KPIScore{
			ItemID:   tplItem.ItemID,
			Name:     tplItem.Name,
			Category: tplItem.Category,
			Weight:   tplItem.Weight,
			MaxScore: tplItem.MaxScore,
			Score:    sc.Score,
			Notes:    strings.TrimSpace(sc.Notes),
		})

		// คำนวณคะแนนรวมแบบถ่วงน้ำหนัก (ยังไม่ได้ normalize เป็นเปอร์เซ็นต์)
		if tplItem.MaxScore > 0 {
			totalWeighted += sc.Score * tplItem.Weight
		}
	}

	// หมายเหตุ: ตอนนี้ TotalScore เก็บเป็นผลรวม (score * weight)
	// ถ้าต้องการเปอร์เซ็นต์ อาจต้อง / (MaxScore * 100) หรือสูตรอื่นภายหลัง

	// 4) สร้างเอกสาร KPIEvaluation
	evaluation := models.KPIEvaluation{
		EvaluationID: uuid.NewString(),
		JobID:        req.JobID,
		TaskID:       req.TaskID,
		KPIID:        tpl.KPIID,
		Version:      1,
		EvaluatorID:  claims.UserID,
		EvaluateeID:  req.EvaluateeID,
		Department:   tpl.Department,
		Scores:       scores,
		TotalScore:   totalWeighted,
		Feedback:     strings.TrimSpace(req.Feedback),
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}

	// 5) บันทึกลงฐานข้อมูล
	if err := s.kpiEvaluationRepo.CreateKPIEvaluations(ctx, evaluation); err != nil {
		return fmt.Errorf("create kpi evaluation: %w", err)
	}

	return nil
}

// func (s *kpiService) GetKPITemplateByID(ctx context.Context, kpiID string) (*dto.KPITemplateDTO, error) {

// 	filter := bson.M{"kpi_id": kpiID, "deleted_at": nil}
// 	projection := bson.M{}

// 	m, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, projection)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if m == nil {
// 		return nil, nil
// 	}

// 	ItemsDTO := make([]dto.KPITemplateItemDTO, 0, len(m.Items))
// 	for _, st := range m.Items {
// 		ItemsDTO = append(ItemsDTO, dto.KPITemplateItemDTO{
// 			ItemID:      st.ItemID,
// 			Name:        st.Name,
// 			Description: st.Description,
// 			Category:    st.Category,
// 			MaxScore:    st.MaxScore,
// 			Weight:      st.Weight,
// 			CreatedAt:   st.CreatedAt,
// 			UpdatedAt:   st.UpdatedAt,
// 		})
// 	}

// 	dtoObj := &dto.KPITemplateDTO{
// 		KPIID:       m.KPIID,
// 		KPIName:     m.KPIName,
// 		Department:  m.Department,
// 		TotalWeight: m.TotalWeight,
// 		Items:       ItemsDTO,
// 		Version:     m.Version,
// 		CreatedBy:   m.CreatedBy,
// 		CreatedAt:   m.CreatedAt,
// 		UpdatedAt:   m.UpdatedAt,
// 	}
// 	return dtoObj, nil
// }

// func (s *kpiService) UpdateKPITemplate(ctx context.Context, kpiID string, req dto.UpdateKPITemplateDTO, claims *dto.JWTClaims) error {

// 	now := time.Now()
// 	// ดึงข้อมูลเดิม
// 	filter := bson.M{"kpi_id": kpiID, "deleted_at": nil}
// 	existing, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, bson.M{})
// 	if err != nil {
// 		return err
// 	}
// 	if existing == nil {
// 		return mongo.ErrNoDocuments
// 	}

// 	// อัปเดตฟิลด์ข้อความ (trim)
// 	if strings.TrimSpace(req.KPIName) != "" {
// 		existing.KPIName = strings.TrimSpace(req.KPIName)
// 	}
// 	if strings.TrimSpace(req.Department) != "" {
// 		existing.Department = strings.TrimSpace(req.Department)
// 	}

// 	// ถ้า client ส่ง items มา (pointer != nil) → validate + แทนที่ทั้งชุด
// 	if req.Items != nil {
// 		if len(*req.Items) == 0 {
// 			return fmt.Errorf("items must not be empty")
// 		}
// 		sumWeight := 0
// 		seen := make(map[string]struct{})
// 		newItems := make([]models.KPITemplateItem, 0, len(*req.Items))

// 		for i, it := range *req.Items {
// 			name := strings.TrimSpace(it.Name)
// 			if name == "" {
// 				return fmt.Errorf("items[%d].name is required", i)
// 			}
// 			key := strings.ToLower(name)
// 			if _, ok := seen[key]; ok {
// 				return fmt.Errorf("items[%d].name duplicated: %s", i, name)
// 			}
// 			seen[key] = struct{}{}

// 			if it.MaxScore <= 0 {
// 				return fmt.Errorf("items[%d].max_score must be > 0", i)
// 			}
// 			if it.Weight <= 0 {
// 				return fmt.Errorf("items[%d].weight must be > 0", i)
// 			}
// 			sumWeight += it.Weight

// 			newItems = append(newItems, models.KPITemplateItem{
// 				ItemID:      uuid.NewString(),
// 				Name:        name,
// 				Description: strings.TrimSpace(it.Description),
// 				Category:    strings.TrimSpace(it.Category),
// 				MaxScore:    it.MaxScore,
// 				Weight:      it.Weight,
// 				CreatedAt:   now,
// 				UpdatedAt:   now,
// 				DeletedAt:   nil,
// 			})
// 		}
// 		// ให้ข้อความตรงกับ handler ตอนนี้ (เช็คเท่ากับ 100 ตรง ๆ)
// 		if sumWeight != 100 {
// 			return fmt.Errorf("sum of weights must be 100")
// 		}

// 		existing.Items = newItems
// 		existing.TotalWeight = 100 // คงเป็น 100 เสมอ เพื่อความชัดเจน
// 	}

// 	existing.Version += 1    // เพิ่มเวอร์ชัน
// 	existing.UpdatedAt = now // อัปเดตเวลา

// 	updated, err := s.kpiRepo.UpdateKPIByID(ctx, kpiID, *existing)
// 	if err != nil {
// 		return err
// 	}
// 	if updated == nil {
// 		return mongo.ErrNoDocuments
// 	}
// 	return nil
// }

// func (s *kpiService) DeleteKPITemplate(ctx context.Context, kpiID string) error {
// 	err := s.kpiRepo.SoftDeleteKPIByID(ctx, kpiID)
// 	if err == mongo.ErrNoDocuments {
// 		return nil
// 	}
// 	return err
// }

func (s *kpiEvaluationRepoService) ListKPIEvaluation(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
		"status":     "done",
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
			{"project_name": re},
			{"job_name": re},
		}
	}

	projection := bson.M{}

	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"project_name": "project_name",
		"job_name":     "job_name",
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

	items, total, err := s.taskRepo.GetListTasksByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list tasks: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		departmentsName := "ไม่พบแผนก"
		createdByName := "ไม่พบผู้สร้าง"
		assigneeName := "ไม่พบผู้รับผิดชอบ"

		createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)
		assignee, _ := s.userRepo.GetByID(ctx, m.Assignee)
		departments, _ := s.departmentRepo.GetOneDepartmentByFilter(ctx, bson.M{"department_id": m.Department, "deleted_at": nil}, bson.M{"_id": 0, "department_name": 1})

		if departments != nil {
			departmentsName = departments.DepartmentName
		}
		if createdBy != nil {
			createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
		}
		if assignee != nil {
			assigneeName = fmt.Sprintf("%s %s %s", assignee.TitleTH, assignee.FirstNameTH, assignee.LastNameTH)
		}

		steps := make([]dto.TaskWorkflowStep, 0, len(m.AppliedWorkflow.Steps))
		for _, st := range m.AppliedWorkflow.Steps {
			steps = append(steps, dto.TaskWorkflowStep{
				StepID:      st.StepID,
				StepName:    st.StepName,
				Description: st.Description,
				Hours:       st.Hours,
				Order:       st.Order,
				Status:      st.Status,
				StartedAt:   st.StartedAt,
				CompletedAt: st.CompletedAt,
				Notes:       st.Notes,
				CreatedAt:   st.CreatedAt,
				UpdatedAt:   st.UpdatedAt,
			})
		}
		list = append(list, dto.TaskDTO{
			TaskID:      m.TaskID,
			ProjectID:   m.ProjectID,
			ProjectName: m.ProjectName,
			JobID:       m.JobID,
			JobName:     m.JobName,
			Description: m.Description,

			Department:     m.Department,
			DepartmentName: departmentsName,
			Assignee:       m.Assignee,
			AssigneeName:   assigneeName,
			Importance:     m.Importance,

			StartDate: m.StartDate,
			EndDate:   m.EndDate,

			KPIID:      m.KPIID,
			WorkFlowID: m.WorkFlowID,

			AppliedWorkflow: dto.TaskAppliedWorkflow{
				WorkFlowID:   m.AppliedWorkflow.WorkFlowID,
				WorkFlowName: m.AppliedWorkflow.WorkFlowName,
				Department:   m.AppliedWorkflow.Department,
				Description:  m.AppliedWorkflow.Description,
				TotalHours:   m.AppliedWorkflow.TotalHours,
				Steps:        steps,
				Version:      m.AppliedWorkflow.Version,
			},

			Status:        m.Status,
			StepName:      m.StepName,
			CreatedBy:     m.CreatedBy,
			CreatedByName: createdByName,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeletedAt:     m.DeletedAt,
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
