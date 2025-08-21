package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/pkg/helpers"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskService struct {
	config       config.Config
	taskRepo     ports.TaskRepository
	userRepo     ports.UserRepository
	workflowRepo ports.WorkFlowRepository
}

func NewTaskService(cfg config.Config, taskRepo ports.TaskRepository, userRepo ports.UserRepository, workflowRepo ports.WorkFlowRepository) ports.TaskService {
	return &taskService{config: cfg, taskRepo: taskRepo, userRepo: userRepo, workflowRepo: workflowRepo}
}

func (s *taskService) GetListTasks(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department"] = department
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

			Department: m.Department,
			Assignee:   m.Assignee,
			Importance: m.Importance,

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

			Status:    m.Status,
			CreatedBy: m.CreatedBy,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
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

func (s *taskService) CreateTask(ctx context.Context, createTask dto.CreateTaskRequest, claims *dto.JWTClaims) error {

	now := time.Now()
	var start time.Time
	var end time.Time

	filter := bson.M{"workflow_id": createTask.WorkflowID}
	projection := bson.M{}
	workflow, err := s.workflowRepo.GetOneWorkFlowTemplateByFilter(ctx, filter, projection)
	if err != nil {
		return err
	}

	if createTask.StartDate != "" {
		parsedDate, err := time.Parse("2006-01-02", createTask.StartDate)
		if err != nil {
			return err
		}
		start = parsedDate
	}

	if createTask.EndDate != "" {
		parsedDate, err := time.Parse("2006-01-02", createTask.EndDate)
		if err != nil {
			return err
		}
		end = parsedDate
	}

	steps := make([]models.TaskWorkflowStep, 0, len(workflow.Steps)+(len(createTask.ExtraSteps)))
	var total float64
	for _, st := range workflow.Steps {
		steps = append(steps, models.TaskWorkflowStep{
			StepID:      uuid.NewString(),
			StepName:    st.StepName,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       st.Order,
			Status:      "todo",
			Notes:       "",
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		total += st.Hours
	}

	for i, st := range createTask.ExtraSteps {
		steps = append(steps, models.TaskWorkflowStep{
			StepID:      uuid.NewString(),
			StepName:    st.StepName,
			Description: st.Description,
			Hours:       st.Hours,
			Order:       len(workflow.Steps) + i + 1,
			Status:      "todo",
			Notes:       "",
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		total += st.Hours
	}

	AppliedWorkflow := models.TaskAppliedWorkflow{
		WorkFlowID:   workflow.WorkFlowID,
		WorkFlowName: workflow.WorkFlowName,
		Department:   workflow.Department,
		Description:  workflow.Description,
		TotalHours:   total,
		Steps:        steps,
		Version:      1,
	}

	model := models.Tasks{
		TaskID:      uuid.New().String(),
		ProjectID:   createTask.ProjectID,
		ProjectName: createTask.ProjectName,
		JobID:       createTask.JobID,
		JobName:     createTask.JobName,
		Description: createTask.Description,

		Department: createTask.Department,
		Assignee:   createTask.Assignee,
		Importance: createTask.Importance,
		StartDate:  start,
		EndDate:    end,
		KPIID:      createTask.KPIID,
		WorkFlowID: createTask.WorkflowID,

		AppliedWorkflow: AppliedWorkflow,

		Status:    "todo",
		CreatedBy: claims.UserID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	if err := s.taskRepo.CreateTask(ctx, model); err != nil {
		return err
	}

	// <============================ UpsertUserTaskStats ============================>

	filterUserTaskStats := bson.M{"user_id": createTask.Assignee} // ถ้าใช้ period ให้ใส่เพิ่ม: filter["period"] = models.StatsPeriodAll
	projectionUserTaskStats := bson.M{}

	existingStats, err := s.taskRepo.GetOneUserTaskStatsByFilter(ctx, filterUserTaskStats, projectionUserTaskStats)
	if err != nil && err != mongo.ErrNoDocuments {
		return err // error จริง
	}

	totals := models.UserTaskTotals{}
	if existingStats == nil {
		// ยังไม่มีเอกสาร → เริ่มนับใหม่
		totals = models.UserTaskTotals{
			Assigned:   1,
			Open:       1, // งานที่เพิ่งสร้างสถานะ "todo" ถือว่า open
			InProgress: 0,
			Completed:  0,
			Skipped:    0,
		}
	} else {
		// มีเอกสารอยู่แล้ว → บวกเพิ่ม
		totals = existingStats.Totals
		totals.Assigned++
		totals.Open++
	}

	// เตรียม doc สำหรับ upsert
	statsDoc := &models.UserTaskStats{
		UserID:       createTask.Assignee,
		DepartmentID: createTask.Department,
		Totals:       totals,
		KPI:          models.UserTaskKPI{Score: nil, LastCalculatedAt: nil},
		UpdatedAt:    now,
	}
	if existingStats == nil {
		statsDoc.CreatedAt = now
	} else {
		statsDoc.CreatedAt = existingStats.CreatedAt
	}

	if err := s.taskRepo.UpsertUserTaskStats(ctx, statsDoc); err != nil {
		return err
	}

	return nil
}

func (s *taskService) GetTaskByID(ctx context.Context, taskID string) (*dto.TaskDTO, error) {

	filter := bson.M{"task_id": taskID, "deleted_at": nil}
	projection := bson.M{}

	m, err := s.taskRepo.GetOneTasksByFilter(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
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

	dtoObj := &dto.TaskDTO{
		TaskID:      m.TaskID,
		ProjectID:   m.ProjectID,
		ProjectName: m.ProjectName,
		JobID:       m.JobID,
		JobName:     m.JobName,
		Description: m.Description,

		Department: m.Department,
		Assignee:   m.Assignee,
		Importance: m.Importance,

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

		Status:    m.Status,
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
	return dtoObj, nil
}

func (s *taskService) UpdateTask(ctx context.Context, taskID string, req dto.UpdateTaskRequest, updatedBy string) error {
	now := time.Now()

	filter := bson.M{"task_id": taskID, "deleted_at": nil}
	existing, err := s.taskRepo.GetOneTasksByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	// --- 2) อัปเดตฟิลด์ระดับงาน (ตาม pointer) ---
	if v := req.ProjectID; v != nil {
		existing.ProjectID = strings.TrimSpace(*v) // trim ช่องว่างก่อน-หลัง
	}
	if v := req.ProjectName; v != nil {
		existing.ProjectName = strings.TrimSpace(*v)
	}
	if v := req.JobID; v != nil {
		existing.JobID = strings.TrimSpace(*v)
	}
	if v := req.JobName; v != nil {
		existing.JobName = strings.TrimSpace(*v)
	}
	if v := req.Description; v != nil {
		existing.Description = strings.TrimSpace(*v)
	}
	if v := req.Department; v != nil {
		existing.Department = strings.TrimSpace(*v)
	}
	if v := req.Assignee; v != nil {
		existing.Assignee = strings.TrimSpace(*v)
	}
	if v := req.Importance; v != nil {
		val := strings.ToLower(strings.TrimSpace(*v)) // normalize เป็น lower-case
		switch val {
		case "low", "medium", "high": // validate ชุดค่า
			existing.Importance = val
		default:
			return fmt.Errorf("importance must be one of: low|medium|high") // แจ้ง error ถ้าค่าไม่อยู่ในชุดที่ยอมรับ
		}
	}

	if v := req.StartDate; v != nil {
		t, err := helpers.DateToISO(*v) // แปลง start_date
		if err != nil {
			return fmt.Errorf("invalid start_date: %w", err) // แจ้ง error รูปแบบไม่ถูกต้อง
		}
		existing.StartDate = t
	}
	if v := req.EndDate; v != nil {
		t, err := helpers.DateToISO(*v) // แปลง end_date
		if err != nil {
			return fmt.Errorf("invalid end_date: %w", err)
		}
		existing.EndDate = t
	}
	if !existing.StartDate.IsZero() && !existing.EndDate.IsZero() && existing.EndDate.Before(existing.StartDate) {
		return fmt.Errorf("end_date must be on or after start_date") // ตรวจเงื่อนไข end >= start
	}

	if v := req.KPIID; v != nil {
		existing.KPIID = strings.TrimSpace(*v) // ปรับ kpi_id ถ้ามีส่งมา
	}

	// --- 3) จัดการกรณีเปลี่ยน Workflow ---
	if v := req.WorkflowID; v != nil { // ถ้าผู้ใช้ส่ง workflow_id มา
		newWfID := strings.TrimSpace(*v)
		if newWfID == "" {
			return fmt.Errorf("workflow_id cannot be empty when provided") // ป้องกันค่าเว้นว่าง
		}
		if newWfID != existing.WorkFlowID { // เฉพาะกรณีต่างจากของเดิม
			wfFilter := bson.M{"workflow_id": newWfID, "deleted_at": nil} // หา template workflow ใหม่
			wf, err := s.workflowRepo.GetOneWorkFlowTemplateByFilter(ctx, wfFilter, bson.M{})
			if err != nil {
				return err
			}
			if wf == nil {
				return mongo.ErrNoDocuments // ไม่พบ workflow ใหม่
			}

			// สร้าง snapshot ของ steps จาก template ใหม่ (เริ่มด้วย status=todo)
			steps := make([]models.TaskWorkflowStep, 0, len(wf.Steps)+len(req.NewSteps))
			var total float64
			for _, st := range wf.Steps {
				steps = append(steps, models.TaskWorkflowStep{
					StepID:      st.StepID,                      // ถ้าอยากให้ StepID แยกจาก template ให้ใช้ uuid.NewString()
					StepName:    strings.TrimSpace(st.StepName), // trim ข้อความ
					Description: strings.TrimSpace(st.Description),
					Hours:       st.Hours,
					Order:       st.Order,
					Status:      "todo",
					CreatedAt:   now,
					UpdatedAt:   now,
				})
				total += st.Hours // รวมชั่วโมง
			}

			// ถ้ามี new_steps และไม่ได้ตั้งใจ replace ทั้งชุด → append ต่อท้าย
			if req.ReplaceSteps == nil || !*req.ReplaceSteps {
				base := len(steps) // เริ่มลำดับต่อจากจำนวนเดิม
				for i, ns := range req.NewSteps {
					if ns.Hours <= 0 {
						return fmt.Errorf("new_steps[%d].hours must be > 0", i) // validate ชั่วโมง > 0
					}
					steps = append(steps, models.TaskWorkflowStep{
						StepID:      uuid.NewString(), // step ใหม่ → gen UUID
						StepName:    strings.TrimSpace(ns.StepName),
						Description: strings.TrimSpace(ns.Description),
						Hours:       ns.Hours,
						Order:       base + i + 1, // ลำดับต่อท้าย
						Status:      "todo",
						CreatedAt:   now,
						UpdatedAt:   now,
					})
					total += ns.Hours // บวกชั่วโมงของ step ใหม่เข้าไป
				}
			} else {
				// replace ทั้งชุดด้วย new_steps เท่านั้น (ไม่เอาของ template)
				steps = steps[:0] // เคลียร์ steps
				total = 0
				for i, ns := range req.NewSteps {
					if ns.Hours <= 0 {
						return fmt.Errorf("new_steps[%d].hours must be > 0", i)
					}
					steps = append(steps, models.TaskWorkflowStep{
						StepID:      uuid.NewString(),
						StepName:    strings.TrimSpace(ns.StepName),
						Description: strings.TrimSpace(ns.Description),
						Hours:       ns.Hours,
						Order:       i + 1, // เริ่มนับใหม่ตั้งแต่ 1
						Status:      "todo",
						CreatedAt:   now,
						UpdatedAt:   now,
					})
					total += ns.Hours
				}
			}

			// เขียน snapshot applied_workflow ใหม่ทับของเดิม
			existing.WorkFlowID = newWfID
			existing.AppliedWorkflow = models.TaskAppliedWorkflow{
				WorkFlowID:   wf.WorkFlowID,
				WorkFlowName: strings.TrimSpace(wf.WorkFlowName),
				Department:   strings.TrimSpace(wf.Department),
				Description:  strings.TrimSpace(wf.Description),
				TotalHours:   total,                         // ชั่วโมงรวมใหม่
				Steps:        steps,                         // steps snapshot ใหม่
				Version:      helpers.MaxInt(1, wf.Version), // ถ้ามีเวอร์ชันใน template ใช้ค่านั้น
			}
		}
	}

	// --- 4) ไม่ได้เปลี่ยน workflow แต่มีคำสั่งจัดการ steps ---
	if req.WorkflowID == nil || strings.TrimSpace(*req.WorkflowID) == existing.WorkFlowID {
		// (4.1) replace steps ทั้งชุดด้วย new_steps
		if req.ReplaceSteps != nil && *req.ReplaceSteps {
			steps := make([]models.TaskWorkflowStep, 0, len(req.NewSteps))
			var total float64
			for i, ns := range req.NewSteps {
				if ns.Hours <= 0 {
					return fmt.Errorf("new_steps[%d].hours must be > 0", i)
				}
				steps = append(steps, models.TaskWorkflowStep{
					StepID:      uuid.NewString(),
					StepName:    strings.TrimSpace(ns.StepName),
					Description: strings.TrimSpace(ns.Description),
					Hours:       ns.Hours,
					Order:       i + 1, // เริ่ม 1..N
					Status:      "todo",
					CreatedAt:   now,
					UpdatedAt:   now,
				})
				total += ns.Hours
			}
			existing.AppliedWorkflow.Steps = steps      // ทับ steps เดิมทั้งหมด
			existing.AppliedWorkflow.TotalHours = total // อัปเดต total_hours
		}

		// (4.2) แพตช์ step ที่มีอยู่ (เจาะจงด้วย step_id)
		if len(req.StepPatches) > 0 {
			idxByID := make(map[string]int, len(existing.AppliedWorkflow.Steps)) // map หา index ของแต่ละ step_id
			for i := range existing.AppliedWorkflow.Steps {
				idxByID[existing.AppliedWorkflow.Steps[i].StepID] = i
			}
			for j, p := range req.StepPatches {
				i, ok := idxByID[p.StepID] // ตรวจ step_id ว่ามีไหม
				if !ok {
					return fmt.Errorf("step_patches[%d]: step_id not found: %s", j, p.StepID)
				}
				st := existing.AppliedWorkflow.Steps[i]

				if p.StepName != nil {
					st.StepName = strings.TrimSpace(*p.StepName)
				}
				if p.Description != nil {
					st.Description = strings.TrimSpace(*p.Description)
				}
				if p.Hours != nil {
					if *p.Hours <= 0 {
						return fmt.Errorf("step_patches[%d].hours must be > 0", j) // hours ต้อง > 0
					}
					st.Hours = *p.Hours
				}
				if p.Order != nil {
					if *p.Order < 1 {
						return fmt.Errorf("step_patches[%d].order must be >= 1", j) // order ขั้นต่ำ 1
					}
					st.Order = *p.Order
				}
				if p.Notes != nil {
					st.Notes = strings.TrimSpace(*p.Notes)
				}
				if p.Status != nil {
					newStatus := strings.ToLower(strings.TrimSpace(*p.Status)) // normalize สถานะ
					if !helpers.InSet(newStatus, "todo", "in_progress", "blocked", "done") {
						return fmt.Errorf("step_patches[%d].status invalid (todo|in_progress|blocked|done)", j)
					}
					// ถ้าเปลี่ยนเป็น in_progress แต่ยังไม่มี started_at และ caller ไม่ได้ส่ง started_at มา → เซ็ตเดี๋ยวนี้
					if newStatus == "in_progress" && st.StartedAt == nil && p.StartedAt == nil {
						t := now
						st.StartedAt = &t
					}
					// ถ้าเปลี่ยนเป็น done แต่ยังไม่มี completed_at และ caller ไม่ได้ส่ง → เซ็ตเดี๋ยวนี้
					if newStatus == "done" && st.CompletedAt == nil && p.CompletedAt == nil {
						t := now
						st.CompletedAt = &t
					}
					st.Status = newStatus
				}
				if p.StartedAt != nil {
					st.StartedAt = p.StartedAt // ยอมรับเวลาจาก caller ถ้าต้องการควบคุมเอง
				}
				if p.CompletedAt != nil {
					st.CompletedAt = p.CompletedAt
				}

				st.UpdatedAt = now                     // อัปเดตเวลาแก้ไขของ step
				existing.AppliedWorkflow.Steps[i] = st // เขียนกลับเข้า slice
			}
		}

		// (4.3) ลบสเต็ปตามรายการ id
		if len(req.DeleteStepIDs) > 0 {
			toDel := make(map[string]struct{}, len(req.DeleteStepIDs))
			for _, id := range req.DeleteStepIDs {
				toDel[id] = struct{}{}
			}
			kept := make([]models.TaskWorkflowStep, 0, len(existing.AppliedWorkflow.Steps))
			for _, st := range existing.AppliedWorkflow.Steps {
				if _, del := toDel[st.StepID]; !del { // เก็บเฉพาะตัวที่ไม่ต้องลบ
					kept = append(kept, st)
				}
			}
			existing.AppliedWorkflow.Steps = kept // เขียนชุดใหม่ (ลบรายการที่ขอ)
		}

		// (4.4) เพิ่มสเต็ปใหม่ต่อท้าย (เฉพาะกรณีไม่ได้ replace ทั้งชุด)
		if len(req.NewSteps) > 0 && (req.ReplaceSteps == nil || !*req.ReplaceSteps) {
			base := len(existing.AppliedWorkflow.Steps) // เริ่มลำดับต่อท้าย
			for i, ns := range req.NewSteps {
				if ns.Hours <= 0 {
					return fmt.Errorf("new_steps[%d].hours must be > 0", i)
				}
				existing.AppliedWorkflow.Steps = append(existing.AppliedWorkflow.Steps, models.TaskWorkflowStep{
					StepID:      uuid.NewString(),
					StepName:    strings.TrimSpace(ns.StepName),
					Description: strings.TrimSpace(ns.Description),
					Hours:       ns.Hours,
					Order:       base + i + 1,
					Status:      "todo",
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			}
		}

		// (4.5) จัด order ให้เรียงและรีเซ็ตหมายเลข 1..N
		sort.SliceStable(existing.AppliedWorkflow.Steps, func(i, j int) bool {
			return existing.AppliedWorkflow.Steps[i].Order < existing.AppliedWorkflow.Steps[j].Order
		})
		for i := range existing.AppliedWorkflow.Steps {
			existing.AppliedWorkflow.Steps[i].Order = i + 1 // รีindex ให้ต่อเนื่อง
		}

		// (4.6) คำนวณ total_hours ใหม่จากชั่วโมงของทุก step
		var total float64
		for _, s := range existing.AppliedWorkflow.Steps {
			total += s.Hours
		}
		existing.AppliedWorkflow.TotalHours = total
	}

	// --- 5) สรุป task.status จากสถานะของ steps ทั้งหมด ---
	existing.Status = helpers.DeriveTaskStatusFromSteps(existing.AppliedWorkflow.Steps) // all done → done, มี in_progress → in_progress, มี ไม่งั้น → todo

	// --- 6) อัปเดตเวลาแก้งาน ---
	existing.UpdatedAt = now

	// --- 7) บันทึกลง DB ---
	updated, err := s.taskRepo.UpdateTaskByID(ctx, taskID, *existing) // เขียนกลับ
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments // ป้องกันกรณีไม่เจอระหว่างเขียน (edge case)
	}
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID string, claims *dto.JWTClaims) error {
	err := s.taskRepo.SoftDeleteTaskByJobID(ctx, taskID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (s *taskService) UpdateStepStatus(ctx context.Context, taskID, stepID string, req dto.UpdateStepStatusNoteRequest, claims *dto.JWTClaims) error {

	now := time.Now().UTC()

	var normalized *string
	if req.Status != nil {
		v := strings.ToLower(strings.TrimSpace(*req.Status))
		switch v {
		case "todo", "in_progress", "skip", "done":
			normalized = &v
		default:
			return fmt.Errorf("invalid status: %s (allow: todo|in_progress|skip|done)", v)
		}
	}

	// อัปเดตฟิลด์ในสเต็ป (status และ/หรือ notes)
	if err := s.taskRepo.UpdateOneStepFields(ctx, taskID, stepID, normalized, req.Notes, now); err != nil {
		return err
	}

	// โหลด steps ทั้งหมด เพื่อสรุปสถานะงาน + คืน step ที่อัปเดตแล้ว
	steps, err := s.taskRepo.GetAllStepSteps(ctx, taskID) // เปลี่ยนชื่อฟังก์ชันให้ตรงความจริง
	if err != nil {
		return err
	}
	if len(steps) == 0 {
		return mongo.ErrNoDocuments
	}

	// หา step ที่เพิ่งแก้
	var updatedStep *models.TaskWorkflowStep
	for i := range steps {
		if steps[i].StepID == stepID {
			updatedStep = &steps[i]
			break
		}
	}
	if updatedStep == nil {
		return mongo.ErrNoDocuments
	}

	// สรุปสถานะงานจาก steps (skip = ปิดสเต็ปเหมือน done)
	newTaskStatus := helpers.DeriveTaskStatusFromSteps(steps)

	if err := s.taskRepo.UpdateTaskStatus(ctx, taskID, newTaskStatus, now); err != nil {
		return err
	}

	// เผื่ออยากคำนวณ “สเต็ปถัดไปที่ยังเปิดอยู่ (todo)” เพื่อ UI เด้งไป
	var next *models.TaskWorkflowStep
	for i := range steps {
		if steps[i].Status == "todo" {
			next = &steps[i]
			break
		}
	}

	log.Println("Next step:", next)
	return nil
}
