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
	config            config.Config
	taskRepo          ports.TaskRepository
	userRepo          ports.UserRepository
	workflowRepo      ports.WorkFlowRepository
	departmentRepo    ports.DepartmentRepository
	kpiEvaluationRepo ports.KPIEvaluationRepository
	kpiRepo           ports.KPIRepository
	signJobRepo       ports.SignJobRepository
}

func NewTaskService(cfg config.Config, taskRepo ports.TaskRepository, userRepo ports.UserRepository, workflowRepo ports.WorkFlowRepository, departmentRepo ports.DepartmentRepository, kpiEvaluationRepo ports.KPIEvaluationRepository, kpiRepo ports.KPIRepository, signJobRepo ports.SignJobRepository) ports.TaskService {
	return &taskService{config: cfg, taskRepo: taskRepo, userRepo: userRepo, workflowRepo: workflowRepo, departmentRepo: departmentRepo, kpiEvaluationRepo: kpiEvaluationRepo, kpiRepo: kpiRepo, signJobRepo: signJobRepo}
}

func (s *taskService) GetListTasks(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string, status string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department_id"] = department
	}

	status = strings.TrimSpace(status)
	if status != "" {
		filter["status"] = status
	} else {
		filter["status"] = bson.M{"$in": []string{"todo", "in_progress"}}
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"project_name": re},
			{"job_name": re},
			{"assignee_name": re},
			{"assignee_nickname": re},
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
		width := 0.0
		height := 0.0
		quantity := 0

		signJob, _ := s.signJobRepo.GetOneSignJobByFilter(ctx, bson.M{"job_id": m.JobID, "deleted_at": nil}, bson.M{"_id": 0, "width": 1, "height": 1, "quantity": 1})
		if signJob != nil {
			width = signJob.Width
			height = signJob.Height
			quantity = signJob.Quantity
		}

		createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)

		departments, _ := s.departmentRepo.GetOneDepartmentByFilter(ctx, bson.M{"department_id": m.Department, "deleted_at": nil}, bson.M{"_id": 0, "department_name": 1})

		if departments != nil {
			departmentsName = departments.DepartmentName
		}
		if createdBy != nil {
			createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
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

			Width:    width,
			Height:   height,
			Quantity: quantity,

			Department:       m.Department,
			DepartmentName:   departmentsName,
			Assignee:         m.Assignee,
			AssigneeName:     m.AssigneeName,
			AssigneeNickName: m.AssigneeNickName,
			Importance:       m.Importance,

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

	// ปรับความจุเริ่มต้นให้สอดคล้องกับชุดข้อมูลที่จะใช้จริง
	capacity := 0
	if !createTask.IsEdit {
		capacity = len(workflow.Steps)
	} else {
		capacity = len(createTask.ExtraSteps)
	}
	steps := make([]models.TaskWorkflowStep, 0, capacity)
	var total float64

	if !createTask.IsEdit {
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
	} else {
		for i, st := range createTask.ExtraSteps {
			steps = append(steps, models.TaskWorkflowStep{
				StepID:      uuid.NewString(),
				StepName:    st.StepName,
				Description: st.Description,
				Hours:       st.Hours,
				Order:       i + 1, // คง logic เดิมไว้
				Status:      "todo",
				Notes:       "",
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			total += st.Hours
		}
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

	// [ADD] เลือกชื่อสเต็ปปัจจุบัน (in_progress > todo > สุดท้าย) — วางไว้ก่อนสร้าง model
	curStepName := ""
	for _, s := range steps {
		if s.Status == "in_progress" {
			curStepName = s.StepName
			break
		}
	}
	if curStepName == "" {
		for _, s := range steps {
			if s.Status == "todo" {
				curStepName = s.StepName
				break
			}
		}
	}
	if curStepName == "" && len(steps) > 0 {
		curStepName = steps[len(steps)-1].StepName
	}

	assigneeName := "ไม่พบชื่อผู้รับผิดชอบ"
	assigneeNickName := "ไม่พบชื่อเล่น"

	getUsers, err := s.userRepo.GetByID(ctx, createTask.Assignee)
	if err != nil {
		return err
	}

	if getUsers != nil {
		assigneeName = fmt.Sprintf("%s %s %s", getUsers.TitleTH, getUsers.FirstNameTH, getUsers.LastNameTH)
		assigneeNickName = getUsers.NickName
	}

	model := models.Tasks{
		TaskID:      uuid.New().String(),
		ProjectID:   createTask.ProjectID,
		ProjectName: createTask.ProjectName,
		JobID:       createTask.JobID,
		JobName:     createTask.JobName,
		Description: createTask.Description,

		Department:       createTask.Department,
		Assignee:         createTask.Assignee,
		AssigneeName:     assigneeName,
		AssigneeNickName: assigneeNickName,
		Importance:       createTask.Importance,
		StartDate:        start,
		EndDate:          end,
		KPIID:            createTask.KPIID,
		WorkFlowID:       createTask.WorkflowID,

		AppliedWorkflow: AppliedWorkflow,

		Status:    "todo",
		StepName:  curStepName,
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
		return err
	}

	nowUTC := time.Now().UTC()

	var totals models.UserTaskTotals
	var kpi models.UserTaskKPI
	var createdAt time.Time

	if existingStats == nil {
		// สร้างเริ่มต้น
		totals = models.UserTaskTotals{
			Assigned:   1,
			Open:       1,
			InProgress: 0,
			Completed:  0,
			Skipped:    0,
		}
		createdAt = nowUTC
	} else {
		totals = existingStats.Totals
		totals.Assigned++
		totals.Open++
		// กันค่าติดลบ (เผื่อข้อมูลก่อนหน้าเพี้ยน)
		if totals.Assigned < 0 {
			totals.Assigned = 0
		}
		if totals.Open < 0 {
			totals.Open = 0
		}
		if totals.InProgress < 0 {
			totals.InProgress = 0
		}
		if totals.Completed < 0 {
			totals.Completed = 0
		}
		if totals.Skipped < 0 {
			totals.Skipped = 0
		}
		kpi = existingStats.KPI
		createdAt = existingStats.CreatedAt
	}

	statsDoc := &models.UserTaskStats{
		UserID:       createTask.Assignee,
		DepartmentID: createTask.Department,
		Totals:       totals,
		KPI:          kpi,
		CreatedAt:    createdAt,
		UpdatedAt:    nowUTC,
		DeletedAt:    nil,
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

	departmentsName := "ไม่พบแผนก"
	createdByName := "ไม่พบผู้สร้าง"
	width := 0.0
	height := 0.0
	quantity := 0

	createdBy, _ := s.userRepo.GetByID(ctx, m.CreatedBy)

	departments, _ := s.departmentRepo.GetOneDepartmentByFilter(ctx, bson.M{"department_id": m.Department, "deleted_at": nil}, bson.M{"_id": 0, "department_name": 1})

	signJob, _ := s.signJobRepo.GetOneSignJobByFilter(ctx, bson.M{"job_id": m.JobID, "deleted_at": nil}, bson.M{"_id": 0, "width": 1, "height": 1, "quantity": 1})
	if signJob != nil {
		width = signJob.Width
		height = signJob.Height
		quantity = signJob.Quantity
	}

	if departments != nil {
		departmentsName = departments.DepartmentName
	}
	if createdBy != nil {
		createdByName = fmt.Sprintf("%s %s %s", createdBy.TitleTH, createdBy.FirstNameTH, createdBy.LastNameTH)
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

		Width:    width,
		Height:   height,
		Quantity: quantity,

		DepartmentName:   departmentsName,
		Department:       m.Department,
		Assignee:         m.Assignee,
		AssigneeName:     m.AssigneeName,
		AssigneeNickName: m.AssigneeNickName,
		Importance:       m.Importance,

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
	}
	return dtoObj, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID string, claims *dto.JWTClaims) error {

	t, _ := s.taskRepo.GetOneTasksByFilter(ctx,
		bson.M{"task_id": taskID, "deleted_at": nil},
		bson.M{"assignee": 1, "status": 1, "department_id": 1},
	)

	err := s.taskRepo.SoftDeleteTaskByID(ctx, taskID)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	// <============================ Update Stats Inline ============================>
	if t != nil && strings.TrimSpace(t.Assignee) != "" {
		updateStats := func(userID string, assignedDelta, openDelta, inProgDelta, completedDelta int) error {
			existingStats, err := s.taskRepo.GetOneUserTaskStatsByFilter(ctx, bson.M{"user_id": userID}, bson.M{})
			if err != nil && err != mongo.ErrNoDocuments {
				return err
			}
			nowUTC := time.Now().UTC()

			totals := models.UserTaskTotals{}
			createdAt := nowUTC
			if existingStats != nil {
				totals = existingStats.Totals
				createdAt = existingStats.CreatedAt
			}
			totals.Assigned += assignedDelta
			totals.Open += openDelta
			totals.InProgress += inProgDelta
			totals.Completed += completedDelta
			if totals.Assigned < 0 {
				totals.Assigned = 0
			}
			if totals.Open < 0 {
				totals.Open = 0
			}
			if totals.InProgress < 0 {
				totals.InProgress = 0
			}
			if totals.Completed < 0 {
				totals.Completed = 0
			}

			return s.taskRepo.UpsertUserTaskStats(ctx, &models.UserTaskStats{
				UserID:       userID,
				DepartmentID: t.Department,
				Totals:       totals,
				KPI:          models.UserTaskKPI{Score: nil, LastCalculatedAt: nil},
				CreatedAt:    createdAt,
				UpdatedAt:    nowUTC,
			})
		}

		switch t.Status {
		case "done":
			_ = updateStats(t.Assignee, -1, 0, 0, -1)
		case "in_progress":
			_ = updateStats(t.Assignee, -1, -1, -1, 0)
		default: // "todo"
			_ = updateStats(t.Assignee, -1, -1, 0, 0)
		}
	}

	return nil

}

func (s *taskService) UpdateStepStatus(ctx context.Context, taskID, stepID string, req dto.UpdateStepStatusNoteRequest, claims *dto.JWTClaims) error {

	now := time.Now()

	// ดึงสถานะ/assignee ปัจจุบันไว้ก่อน (เพื่อเทียบหลังอัปเดต)
	prevTask, _ := s.taskRepo.GetOneTasksByFilter(ctx,
		bson.M{"task_id": taskID, "deleted_at": nil},
		bson.M{"assignee": 1, "status": 1, "department_id": 1},
	)
	var prevStatus, assignee, department string
	if prevTask != nil {
		if prevTask.Assignee != claims.UserID {
			return fmt.Errorf("user %s ไม่ใช่ผู้รับผิดชอบงานนี้ (%s) ไม่สามารถแก้ไขสถานะได้", claims.UserID, prevTask.Assignee)
		}

		prevStatus = prevTask.Status
		assignee = prevTask.Assignee
		department = prevTask.Department
	}

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

	// [ADD] คำนวณ step_name ปัจจุบันจาก steps ล่าสุด (in_progress > todo > สุดท้าย)
	curStepName := ""
	for _, s2 := range steps {
		if s2.Status == "in_progress" {
			curStepName = s2.StepName
			break
		}
	}
	if curStepName == "" {
		for _, s2 := range steps {
			if s2.Status == "todo" {
				curStepName = s2.StepName
				break
			}
		}
	}
	if curStepName == "" && len(steps) > 0 {
		curStepName = steps[len(steps)-1].StepName
	}

	if err := s.taskRepo.UpdateTaskStatus(ctx, taskID, newTaskStatus, curStepName, now); err != nil {
		return err
	}

	// [EVAL] ถ้าเพิ่งเปลี่ยนเป็น done ให้สร้างแบบประเมิน
	if prevStatus != "done" && newTaskStatus == "done" {
		_ = s.CreateEvaluationIfNeeded(ctx, taskID)
	}

	// <============================ Update Stats Inline ============================>
	// ปรับเฉพาะเมื่อ status งานเปลี่ยน (assignee เดิม)
	if assignee != "" && prevStatus != "" && prevStatus != newTaskStatus {
		updateStats := func(userID string, assignedDelta, openDelta, inProgDelta, completedDelta int) error {
			if strings.TrimSpace(userID) == "" {
				return nil
			}
			existingStats, err := s.taskRepo.GetOneUserTaskStatsByFilter(ctx, bson.M{"user_id": userID}, bson.M{})
			if err != nil && err != mongo.ErrNoDocuments {
				return err
			}
			nowUTC := time.Now()

			totals := models.UserTaskTotals{}
			createdAt := nowUTC
			if existingStats != nil {
				totals = existingStats.Totals
				createdAt = existingStats.CreatedAt
			}
			totals.Assigned += assignedDelta // ปกติ status change ไม่แตะ assigned (0)
			totals.Open += openDelta
			totals.InProgress += inProgDelta
			totals.Completed += completedDelta
			if totals.Assigned < 0 {
				totals.Assigned = 0
			}
			if totals.Open < 0 {
				totals.Open = 0
			}
			if totals.InProgress < 0 {
				totals.InProgress = 0
			}
			if totals.Completed < 0 {
				totals.Completed = 0
			}

			return s.taskRepo.UpsertUserTaskStats(ctx, &models.UserTaskStats{
				UserID:       userID,
				DepartmentID: department,
				Totals:       totals,
				KPI:          models.UserTaskKPI{Score: nil, LastCalculatedAt: nil},
				CreatedAt:    createdAt,
				UpdatedAt:    nowUTC,
			})
		}

		switch prevStatus {
		case "todo":
			switch newTaskStatus {
			case "in_progress":
				_ = updateStats(assignee, 0, 0, +1, 0)
			case "done":
				_ = updateStats(assignee, 0, -1, 0, +1)
			}
		case "in_progress":
			switch newTaskStatus {
			case "todo":
				_ = updateStats(assignee, 0, 0, -1, 0)
			case "done":
				_ = updateStats(assignee, 0, -1, -1, +1)
			}
		case "done":
			switch newTaskStatus {
			case "in_progress":
				_ = updateStats(assignee, 0, +1, +1, -1)
			case "todo":
				_ = updateStats(assignee, 0, +1, 0, -1)
			}
		}
	}

	return nil
}

// func (s *taskService) ReplaceTask(ctx context.Context, taskID string, req dto.UpdateTaskPutRequest, updatedBy string) error {
// 	now := time.Now()

// 	// 1) โหลดงานเดิม
// 	existing, err := s.taskRepo.GetOneTasksByFilter(ctx, bson.M{"task_id": taskID, "deleted_at": nil}, bson.M{})
// 	if err != nil {
// 		return err
// 	}
// 	if existing == nil {
// 		return mongo.ErrNoDocuments
// 	}
// 	oldAssignee := existing.Assignee
// 	oldStatus := existing.Status
// 	oldCreatedAt := existing.CreatedAt
// 	oldCreatedBy := existing.CreatedBy
// 	oldTaskID := existing.TaskID

// 	// 2) validate ขั้นพื้นฐาน
// 	imp := strings.ToLower(strings.TrimSpace(req.Importance))
// 	if imp != "low" && imp != "medium" && imp != "high" {
// 		return fmt.Errorf("importance must be one of: low|medium|high")
// 	}

// 	start, err := helpers.DateToISO(strings.TrimSpace(req.StartDate))
// 	if err != nil {
// 		return fmt.Errorf("invalid start_date: %w", err)
// 	}
// 	end, err := helpers.DateToISO(strings.TrimSpace(req.EndDate))
// 	if err != nil {
// 		return fmt.Errorf("invalid end_date: %w", err)
// 	}
// 	if !start.IsZero() && !end.IsZero() && end.Before(start) {
// 		return fmt.Errorf("end_date must be on or after start_date")
// 	}

// 	// 3) สร้าง snapshot steps ใหม่จาก req.AppliedWorkflow.Steps (ต้องส่งทั้งชุด)
// 	//    - validate hours > 0, status ถูกต้อง, reindex order, auto started/completed
// 	normalizeStatus := func(s string) (string, error) {
// 		ss := strings.ToLower(strings.TrimSpace(s))
// 		switch ss {
// 		case "todo", "in_progress", "skip", "done":
// 			return ss, nil
// 		default:
// 			return "", fmt.Errorf("invalid step status: %s", s)
// 		}
// 	}

// 	steps := make([]models.TaskWorkflowStep, 0, len(req.AppliedWorkflow.Steps))
// 	for i, st := range req.AppliedWorkflow.Steps {
// 		if st.Hours <= 0 {
// 			return fmt.Errorf("steps[%d].hours must be > 0", i)
// 		}
// 		ns, err := normalizeStatus(st.Status)
// 		if err != nil {
// 			return fmt.Errorf("steps[%d]: %w", i, err)
// 		}

// 		// auto time fields ตามลอจิกเดิม
// 		started := st.StartedAt
// 		completed := st.CompletedAt
// 		if ns == "in_progress" && started == nil {
// 			t := now
// 			started = &t
// 		}
// 		if (ns == "done" || ns == "skip") && completed == nil {
// 			t := now
// 			completed = &t
// 		}

// 		steps = append(steps, models.TaskWorkflowStep{
// 			StepID:      strings.TrimSpace(st.StepID),
// 			StepName:    strings.TrimSpace(st.StepName),
// 			Description: strings.TrimSpace(st.Description),
// 			Hours:       st.Hours,
// 			Order:       st.Order, // เดี๋ยว reindex อีกที
// 			Status:      ns,
// 			StartedAt:   started,
// 			CompletedAt: completed,
// 			Notes:       strings.TrimSpace(st.Notes),
// 			CreatedAt:   existing.CreatedAt, // หรือจะใช้เวลาตอนสร้างเดิมก็ได้ (ขึ้นกับนโยบาย)
// 			UpdatedAt:   now,
// 		})
// 	}

// 	// sort ตาม Order แล้ว reindex 1..N
// 	sort.SliceStable(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
// 	for i := range steps {
// 		steps[i].Order = i + 1
// 	}

// 	// คำนวณ total hours
// 	var total float64
// 	for _, s := range steps {
// 		total += s.Hours
// 	}

// 	// 4) คำนวณ step_name และ task.status จาก steps
// 	curStepName := ""
// 	for _, s2 := range steps {
// 		if s2.Status == "in_progress" {
// 			curStepName = s2.StepName
// 			break
// 		}
// 	}
// 	if curStepName == "" {
// 		for _, s2 := range steps {
// 			if s2.Status == "todo" {
// 				curStepName = s2.StepName
// 				break
// 			}
// 		}
// 	}
// 	if curStepName == "" && len(steps) > 0 {
// 		curStepName = steps[len(steps)-1].StepName
// 	}

// 	derived := helpers.DeriveTaskStatusFromSteps(steps)

// 	// 5) ประกอบเอกสารใหม่ทั้งก้อน (replace) โดย "คง" immutable เดิม
// 	newDoc := models.Tasks{
// 		TaskID:      oldTaskID,
// 		ProjectID:   strings.TrimSpace(req.ProjectID),
// 		ProjectName: strings.TrimSpace(req.ProjectName),
// 		JobID:       strings.TrimSpace(req.JobID),
// 		JobName:     strings.TrimSpace(req.JobName),
// 		Description: strings.TrimSpace(req.Description),

// 		Department: strings.TrimSpace(req.Department),
// 		Assignee:   strings.TrimSpace(req.Assignee),
// 		Importance: imp,

// 		StartDate: start,
// 		EndDate:   end,

// 		KPIID:      strings.TrimSpace(req.KPIID),
// 		WorkFlowID: strings.TrimSpace(req.WorkflowID),

// 		AppliedWorkflow: models.TaskAppliedWorkflow{
// 			WorkFlowID:   strings.TrimSpace(req.AppliedWorkflow.WorkFlowID),
// 			WorkFlowName: strings.TrimSpace(req.AppliedWorkflow.WorkFlowName),
// 			Department:   strings.TrimSpace(req.AppliedWorkflow.Department),
// 			Description:  strings.TrimSpace(req.AppliedWorkflow.Description),
// 			TotalHours:   total, // ทับ req
// 			Steps:        steps,
// 			Version:      existing.AppliedWorkflow.Version + 1,
// 		},

// 		Status:    derived,      // ทับ req.Status
// 		StepName:  curStepName,  // จากขั้นตอน
// 		CreatedBy: oldCreatedBy, // คงของเดิม
// 		CreatedAt: oldCreatedAt, // คงของเดิม
// 		UpdatedAt: now,
// 		DeletedAt: nil,
// 	}

// 	// 6) Replace ใน DB
// 	updated, err := s.taskRepo.ReplaceTaskByID(ctx, taskID, &newDoc)
// 	if err != nil {
// 		return err
// 	}
// 	if updated == nil {
// 		return mongo.ErrNoDocuments
// 	}

// 	// [EVAL] ถ้าเพิ่งเปลี่ยนเป็น done ให้สร้างแบบประเมิน
// 	if oldStatus != "done" && newDoc.Status == "done" {
// 		_ = s.CreateEvaluationIfNeeded(ctx, taskID)
// 	}

// 	// 7) อัปเดตสถิติ (เหมือนเดิม) — diff ผู้รับผิดชอบ/สถานะ
// 	newAssignee := newDoc.Assignee
// 	newStatus := newDoc.Status

// 	updateStats := func(userID string, assignedDelta, openDelta, inProgDelta, completedDelta int) error {
// 		if strings.TrimSpace(userID) == "" {
// 			return nil
// 		}
// 		existingStats, err := s.taskRepo.GetOneUserTaskStatsByFilter(ctx,
// 			bson.M{"user_id": userID},
// 			bson.M{},
// 		)
// 		if err != nil && err != mongo.ErrNoDocuments {
// 			return err
// 		}
// 		nowUTC := time.Now().UTC()
// 		totals := models.UserTaskTotals{}
// 		createdAt := nowUTC
// 		if existingStats != nil {
// 			totals = existingStats.Totals
// 			createdAt = existingStats.CreatedAt
// 		}
// 		totals.Assigned += assignedDelta
// 		totals.Open += openDelta
// 		totals.InProgress += inProgDelta
// 		totals.Completed += completedDelta
// 		if totals.Assigned < 0 {
// 			totals.Assigned = 0
// 		}
// 		if totals.Open < 0 {
// 			totals.Open = 0
// 		}
// 		if totals.InProgress < 0 {
// 			totals.InProgress = 0
// 		}
// 		if totals.Completed < 0 {
// 			totals.Completed = 0
// 		}

// 		statsDoc := &models.UserTaskStats{
// 			UserID:       userID,
// 			DepartmentID: newDoc.Department,
// 			Totals:       totals,
// 			KPI:          models.UserTaskKPI{Score: nil, LastCalculatedAt: nil},
// 			CreatedAt:    createdAt,
// 			UpdatedAt:    nowUTC,
// 		}
// 		return s.taskRepo.UpsertUserTaskStats(ctx, statsDoc)
// 	}

// 	if oldAssignee != newAssignee {
// 		// หักของเก่า
// 		switch oldStatus {
// 		case "done":
// 			_ = updateStats(oldAssignee, -1, 0, 0, -1)
// 		case "in_progress":
// 			_ = updateStats(oldAssignee, -1, -1, -1, 0)
// 		default:
// 			_ = updateStats(oldAssignee, -1, -1, 0, 0)
// 		}
// 		// บวกของใหม่
// 		switch newStatus {
// 		case "done":
// 			_ = updateStats(newAssignee, +1, 0, 0, +1)
// 		case "in_progress":
// 			_ = updateStats(newAssignee, +1, +1, +1, 0)
// 		default:
// 			_ = updateStats(newAssignee, +1, +1, 0, 0)
// 		}
// 		return nil
// 	}

// 	if oldStatus != newStatus {
// 		switch oldStatus {
// 		case "todo":
// 			switch newStatus {
// 			case "in_progress":
// 				_ = updateStats(newAssignee, 0, 0, +1, 0)
// 			case "done":
// 				_ = updateStats(newAssignee, 0, -1, 0, +1)
// 			}
// 		case "in_progress":
// 			switch newStatus {
// 			case "todo":
// 				_ = updateStats(newAssignee, 0, 0, -1, 0)
// 			case "done":
// 				_ = updateStats(newAssignee, 0, -1, -1, +1)
// 			}
// 		case "done":
// 			switch newStatus {
// 			case "in_progress":
// 				_ = updateStats(newAssignee, 0, +1, +1, -1)
// 			case "todo":
// 				_ = updateStats(newAssignee, 0, +1, 0, -1)
// 			}
// 		}
// 	}

// 	return nil
// }

func (s *taskService) CreateEvaluationIfNeeded(ctx context.Context, taskID string) error {
	// 1) โหลด Task ล่าสุด (ต้องมีข้อมูลพื้นฐานพอให้ตั้งแบบประเมิน)
	task, errOnGetOneTasksByFilter := s.taskRepo.GetOneTasksByFilter(ctx,
		bson.M{"task_id": taskID, "deleted_at": nil},
		bson.M{},
	)
	if errOnGetOneTasksByFilter != nil {
		log.Println("Error loading task for CreateEvaluationIfNeeded:", errOnGetOneTasksByFilter)
		return errOnGetOneTasksByFilter
	}
	if task == nil {
		return mongo.ErrNoDocuments
	}

	// 2) กันซ้ำด้วย unique key (เช่น unique index ที่ collection: {task_id:1} หรือ {task_id:1,kpi_id:1})
	exists, errOnGetOneKPIEvaluationByFilter := s.kpiEvaluationRepo.GetOneKPIEvaluationByFilter(ctx,
		bson.M{"task_id": task.TaskID, "deleted_at": nil},
		bson.M{"_id": 1},
	)
	if errOnGetOneKPIEvaluationByFilter != nil && errOnGetOneKPIEvaluationByFilter != mongo.ErrNoDocuments {
		log.Println("Error loading KPI evaluation for CreateEvaluationIfNeeded:", errOnGetOneKPIEvaluationByFilter)
		return errOnGetOneKPIEvaluationByFilter
	}
	if exists != nil {
		// มีแล้ว -> ไม่ต้องสร้างซ้ำ
		return nil
	}

	now := time.Now()

	// 3) เตรียมเอกสารแบบประเมินเริ่มต้น (ปรับ fields ให้ตรง model/eval DTO ของคุณ)
	doc := &models.KPIEvaluation{
		EvaluationID: uuid.NewString(),
		ProjectID:    task.ProjectID,
		JobID:        task.JobID,
		TaskID:       task.TaskID,
		KPIID:        task.KPIID,
		Version:      1, // default; update if KPI template provides version
		EvaluatorID:  "",
		EvaluateeID:  task.Assignee,
		Department:   task.Department,
		Scores:       []models.KPIScore{}, // will append after loading template
		TotalScore:   0,
		Feedback:     "",
		IsEvaluated:  false,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}

	// 4) (ทางเลือก) Auto-generate รายการ KPI items จาก KPI Template

	filter := bson.M{"kpi_id": task.KPIID, "deleted_at": nil}
	projection := bson.M{}

	tpl, errOnGetOneKPIByFilter := s.kpiRepo.GetOneKPIByFilter(ctx, filter, projection)
	if errOnGetOneKPIByFilter != nil {
		log.Println("Error loading KPI template for CreateEvaluationIfNeeded:", errOnGetOneKPIByFilter)
		return errOnGetOneKPIByFilter
	}

	if tpl != nil {
		for _, it := range tpl.Items {
			doc.Scores = append(doc.Scores, models.KPIScore{
				ItemID:   it.ItemID,
				Name:     it.Name,     // snapshot ชื่อ item
				Category: it.Category, // snapshot หมวดหมู่
				Weight:   it.Weight,   // weight ปัจจุบัน (int)
				MaxScore: it.MaxScore, // คะแนนเต็ม
				Score:    0,           // ยังไม่ประเมิน (zero value)
				Notes:    "",
			})
		}
	}

	// 5) Insert
	return s.kpiEvaluationRepo.CreateKPIEvaluations(ctx, *doc)
}

func (s *taskService) ReplaceTask(ctx context.Context, taskID string, req dto.UpdateTaskPutRequest, updatedBy string) error {
	now := time.Now()

	// 1) โหลดงานเดิม
	existing, err := s.taskRepo.GetOneTasksByFilter(ctx, bson.M{"task_id": taskID, "deleted_at": nil}, bson.M{})
	if err != nil {
		return err
	}
	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if updatedBy != existing.Assignee {
		return fmt.Errorf("user %s ไม่ใช่ผู้รับผิดชอบงานนี้ (%s) ไม่สามารถแก้ไขได้", updatedBy, existing.Assignee)
	}

	oldAssignee := existing.Assignee
	oldStatus := existing.Status
	oldCreatedAt := existing.CreatedAt
	oldCreatedBy := existing.CreatedBy
	oldTaskID := existing.TaskID

	// 2) validate ขั้นพื้นฐาน
	imp := strings.ToLower(strings.TrimSpace(req.Importance))
	if imp != "low" && imp != "medium" && imp != "high" {
		return fmt.Errorf("importance must be one of: low|medium|high")
	}

	start, err := helpers.DateToISO(strings.TrimSpace(req.StartDate))
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}
	end, err := helpers.DateToISO(strings.TrimSpace(req.EndDate))
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}
	if !start.IsZero() && !end.IsZero() && end.Before(start) {
		return fmt.Errorf("end_date must be on or after start_date")
	}

	// 3) สร้าง snapshot steps ใหม่ (รองรับกรณี FE ไม่ส่ง step_id/started_at/completed_at)
	normalizeStatus := func(s string) (string, error) {
		ss := strings.ToLower(strings.TrimSpace(s))
		switch ss {
		case "todo", "in_progress", "skip", "done":
			return ss, nil
		default:
			return "", fmt.Errorf("invalid step status: %s", s)
		}
	}

	// index ของสเต็ปเดิมไว้ให้แมป
	byID := make(map[string]*models.TaskWorkflowStep, len(existing.AppliedWorkflow.Steps))
	byName := make(map[string]*models.TaskWorkflowStep, len(existing.AppliedWorkflow.Steps)) // ถ้ากังวลชื่อซ้ำ ให้ใช้ name#order แทน
	for i := range existing.AppliedWorkflow.Steps {
		prev := &existing.AppliedWorkflow.Steps[i]
		if prev.StepID != "" {
			byID[prev.StepID] = prev
		}
		key := strings.ToLower(strings.TrimSpace(prev.StepName))
		// เก็บตัวแรกไว้พอ (ต้องการแค่คงค่าเดิม)
		if _, ok := byName[key]; !ok {
			byName[key] = prev
		}
	}

	steps := make([]models.TaskWorkflowStep, 0, len(req.AppliedWorkflow.Steps))

	for i, st := range req.AppliedWorkflow.Steps {
		ns, err := normalizeStatus(st.Status)
		if err != nil {
			return fmt.Errorf("steps[%d]: %w", i, err)
		}

		// หา prev โดย id ก่อน ถ้าไม่มี id ค่อยหาโดยชื่อ
		var prev *models.TaskWorkflowStep
		stepID := strings.TrimSpace(st.StepID)
		if stepID != "" {
			prev = byID[stepID]
		} else {
			key := strings.ToLower(strings.TrimSpace(st.StepName))
			prev = byName[key]
		}

		// hours: ถ้า FE ส่ง <= 0 และมี prev -> ใช้ของเดิม, ไม่งั้นต้อง > 0
		hours := st.Hours
		if hours <= 0 && prev != nil {
			hours = prev.Hours
		}
		if hours <= 0 {
			return fmt.Errorf("steps[%d].hours must be > 0", i)
		}

		// step_id: ถ้าไม่มีและแมปไม่เจอ -> gen ใหม่
		if stepID == "" {
			if prev != nil {
				stepID = prev.StepID
			} else {
				stepID = uuid.NewString()
			}
		}

		// เวลาต่างๆ: ใช้ของที่ส่งมา > ถ้าไม่ส่งและมี prev -> คงของเดิม > แล้วเติมออโต้ตามสถานะ
		started := st.StartedAt
		if started == nil && prev != nil {
			started = prev.StartedAt
		}
		completed := st.CompletedAt
		if completed == nil && prev != nil {
			completed = prev.CompletedAt
		}
		if ns == "in_progress" && started == nil {
			t := now
			started = &t
		}
		if (ns == "done" || ns == "skip") && completed == nil {
			t := now
			completed = &t
		}
		// (ออปชัน) ย้อนสถานะจาก done -> todo/in_progress ให้ล้าง completed_at
		if prev != nil && prev.Status == "done" && (ns == "todo" || ns == "in_progress") {
			completed = nil
		}

		// คง description/notes เดิม ถ้า FE ส่งว่าง
		desc := strings.TrimSpace(st.Description)
		if desc == "" && prev != nil {
			desc = prev.Description
		}
		notes := strings.TrimSpace(st.Notes)
		if notes == "" && prev != nil {
			notes = prev.Notes
		}

		// created_at: คงของเดิมถ้ามี prev, ไม่งั้น now
		createdAt := now
		if prev != nil {
			createdAt = prev.CreatedAt
		}

		steps = append(steps, models.TaskWorkflowStep{
			StepID:      stepID,
			StepName:    strings.TrimSpace(st.StepName),
			Description: desc,
			Hours:       hours,
			Order:       st.Order, // จะ reindex อีกที
			Status:      ns,
			StartedAt:   started,
			CompletedAt: completed,
			Notes:       notes,
			CreatedAt:   createdAt,
			UpdatedAt:   now,
		})
	}

	// sort + reindex
	sort.SliceStable(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	for i := range steps {
		steps[i].Order = i + 1
	}

	// คำนวณ total hours ใหม่
	var total float64
	for _, s2 := range steps {
		total += s2.Hours
	}

	// 4) คำนวณ step_name และ task.status จาก steps
	curStepName := ""
	for _, s2 := range steps {
		if s2.Status == "in_progress" {
			curStepName = s2.StepName
			break
		}
	}
	if curStepName == "" {
		for _, s2 := range steps {
			if s2.Status == "todo" {
				curStepName = s2.StepName
				break
			}
		}
	}
	if curStepName == "" && len(steps) > 0 {
		curStepName = steps[len(steps)-1].StepName
	}

	derived := helpers.DeriveTaskStatusFromSteps(steps)

	// 5) ประกอบเอกสารใหม่ทั้งก้อน (replace) โดยคง immutable เดิม
	newDoc := models.Tasks{
		TaskID:      oldTaskID,
		ProjectID:   strings.TrimSpace(req.ProjectID),
		ProjectName: strings.TrimSpace(req.ProjectName),
		JobID:       strings.TrimSpace(req.JobID),
		JobName:     strings.TrimSpace(req.JobName),
		Description: strings.TrimSpace(req.Description),

		Department: strings.TrimSpace(req.Department),
		Assignee:   strings.TrimSpace(req.Assignee),
		Importance: imp,

		StartDate: start,
		EndDate:   end,

		KPIID:      strings.TrimSpace(req.KPIID),
		WorkFlowID: strings.TrimSpace(req.WorkflowID),

		AppliedWorkflow: models.TaskAppliedWorkflow{
			WorkFlowID:   strings.TrimSpace(req.AppliedWorkflow.WorkFlowID),
			WorkFlowName: strings.TrimSpace(req.AppliedWorkflow.WorkFlowName),
			Department:   strings.TrimSpace(req.AppliedWorkflow.Department),
			Description:  strings.TrimSpace(req.AppliedWorkflow.Description),
			TotalHours:   total,
			Steps:        steps,
			Version:      existing.AppliedWorkflow.Version + 1, // bump version
		},

		Status:    derived,     // ทับ req.Status
		StepName:  curStepName, // จากขั้นตอน
		CreatedBy: oldCreatedBy,
		CreatedAt: oldCreatedAt,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	// 6) Replace ใน DB
	updated, err := s.taskRepo.ReplaceTaskByID(ctx, taskID, &newDoc)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}

	// (ออปชัน) ถ้าเพิ่งเปลี่ยนเป็น done ให้สร้างแบบประเมิน
	if oldStatus != "done" && newDoc.Status == "done" {
		_ = s.CreateEvaluationIfNeeded(ctx, taskID)
	}

	// 7) อัปเดตสถิติ — diff ผู้รับผิดชอบ/สถานะ
	newAssignee := newDoc.Assignee
	newStatus := newDoc.Status

	updateStats := func(userID string, assignedDelta, openDelta, inProgDelta, completedDelta int) error {
		if strings.TrimSpace(userID) == "" {
			return nil
		}
		existingStats, err := s.taskRepo.GetOneUserTaskStatsByFilter(ctx,
			bson.M{"user_id": userID},
			bson.M{},
		)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}
		nowUTC := time.Now().UTC()
		totals := models.UserTaskTotals{}
		createdAt := nowUTC
		if existingStats != nil {
			totals = existingStats.Totals
			createdAt = existingStats.CreatedAt
		}
		totals.Assigned += assignedDelta
		totals.Open += openDelta
		totals.InProgress += inProgDelta
		totals.Completed += completedDelta
		if totals.Assigned < 0 {
			totals.Assigned = 0
		}
		if totals.Open < 0 {
			totals.Open = 0
		}
		if totals.InProgress < 0 {
			totals.InProgress = 0
		}
		if totals.Completed < 0 {
			totals.Completed = 0
		}

		statsDoc := &models.UserTaskStats{
			UserID:       userID,
			DepartmentID: newDoc.Department,
			Totals:       totals,
			KPI:          models.UserTaskKPI{Score: nil, LastCalculatedAt: nil},
			CreatedAt:    createdAt,
			UpdatedAt:    nowUTC,
		}
		return s.taskRepo.UpsertUserTaskStats(ctx, statsDoc)
	}

	if oldAssignee != newAssignee {
		switch oldStatus {
		case "done":
			_ = updateStats(oldAssignee, -1, 0, 0, -1)
		case "in_progress":
			_ = updateStats(oldAssignee, -1, -1, -1, 0)
		default:
			_ = updateStats(oldAssignee, -1, -1, 0, 0)
		}
		switch newStatus {
		case "done":
			_ = updateStats(newAssignee, +1, 0, 0, +1)
		case "in_progress":
			_ = updateStats(newAssignee, +1, +1, +1, 0)
		default:
			_ = updateStats(newAssignee, +1, +1, 0, 0)
		}
		return nil
	}

	if oldStatus != newStatus {
		switch oldStatus {
		case "todo":
			switch newStatus {
			case "in_progress":
				_ = updateStats(newAssignee, 0, 0, +1, 0)
			case "done":
				_ = updateStats(newAssignee, 0, -1, 0, +1)
			}
		case "in_progress":
			switch newStatus {
			case "todo":
				_ = updateStats(newAssignee, 0, 0, -1, 0)
			case "done":
				_ = updateStats(newAssignee, 0, -1, -1, +1)
			}
		case "done":
			switch newStatus {
			case "in_progress":
				_ = updateStats(newAssignee, 0, +1, +1, -1)
			case "todo":
				_ = updateStats(newAssignee, 0, +1, 0, -1)
			}
		}
	}

	return nil
}
