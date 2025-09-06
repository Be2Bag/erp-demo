package dto

import "time"

// ===== Task Create Request =====
type CreateTaskRequest struct {
	ProjectID        string             `json:"project_id"`
	ProjectName      string             `json:"project_name"`
	JobID            string             `json:"job_id"`
	JobName          string             `json:"job_name"`
	Description      string             `json:"description,omitempty"`
	Department       string             `json:"department_id"`
	Assignee         string             `json:"assignee"`
	AssigneeName     string             `json:"assignee_name"`     // ชื่อผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	AssigneeNickName string             `json:"assignee_nickname"` // ชื่อเล่นผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	Importance       string             `json:"importance"`
	StartDate        string             `json:"start_date"`
	EndDate          string             `json:"end_date"`
	KPIID            string             `json:"kpi_id"`
	WorkflowID       string             `json:"workflow_id"`
	IsEdit           bool               `json:"is_edit"`
	ExtraSteps       []ExtraStepRequest `json:"extra_steps,omitempty"`
}

type ExtraStepRequest struct {
	StepName    string  `json:"step_name"`
	Description string  `json:"description,omitempty"`
	Hours       float64 `json:"hours"`
}

type UpdateTaskRequest struct {
	// ฟิลด์ระดับงาน (อัปเดตเฉพาะที่ส่งมา)
	ProjectID        *string `json:"project_id,omitempty"`
	ProjectName      *string `json:"project_name,omitempty"`
	JobID            *string `json:"job_id,omitempty"`
	JobName          *string `json:"job_name,omitempty"`
	Description      *string `json:"description,omitempty"`
	Department       *string `json:"department_id,omitempty"`
	Assignee         *string `json:"assignee,omitempty"`
	AssigneeName     string  `json:"assignee_name"`        // ชื่อผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	AssigneeNickName string  `json:"assignee_nickname"`    // ชื่อเล่นผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	Importance       *string `json:"importance,omitempty"` // low|medium|high
	StartDate        *string `json:"start_date,omitempty"` // "YYYY-MM-DD"
	EndDate          *string `json:"end_date,omitempty"`   // "YYYY-MM-DD"
	KPIID            *string `json:"kpi_id,omitempty"`
	WorkflowID       *string `json:"workflow_id,omitempty"` // ถ้าเปลี่ยน workflow ควรรีเพลส steps

	Status *string `json:"status,omitempty"` // todo|in_progress|skip|done

	// การจัดการ Steps
	// 1) เพิ่มสเต็ปใหม่ (append ต่อท้าย หรือถ้าจะ replace ทั้งชุดให้ดู ReplaceSteps)
	NewSteps []ExtraStepRequest `json:"new_steps,omitempty"`

	// 2) แก้ไขสเต็ปเดิมเป็นรายตัว (ระบุ step_id)
	StepPatches []TaskStepPatch `json:"step_patches,omitempty"`

	// 3) ลบสเต็ปตาม id
	DeleteStepIDs []string `json:"delete_step_ids,omitempty"`

	// 4) รีเพลสทั้งชุด (ถ้า true จะทับ steps เดิมด้วย NewSteps)
	ReplaceSteps *bool `json:"replace_steps,omitempty"`
}

type TaskStepPatch struct {
	StepID      string     `json:"step_id"` // target
	StepName    *string    `json:"step_name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Hours       *float64   `json:"hours,omitempty"`
	Order       *int       `json:"order,omitempty"`
	Status      *string    `json:"status,omitempty"` // todo|in_progress|skip|done
	Notes       *string    `json:"notes,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"` // หรือใช้สตริง ISO8601 ก็ได้
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// อัปเดตสเต็ปเดียว
type UpdateStepStatusNoteRequest struct {
	Status *string `json:"status,omitempty"` // todo|in_progress|skip|done (optional)
	Notes  *string `json:"notes,omitempty"`  // optional
}

// ===== Response =====
type TaskDTO struct {
	TaskID      string `json:"task_id"`      // รหัสงาน (UUID/unique)
	ProjectID   string `json:"project_id"`   // รหัสโปรเจกต์
	ProjectName string `json:"project_name"` // ชื่อโปรเจกต์
	JobID       string `json:"job_id"`
	JobName     string `json:"job_name"`    // ชื่องาน
	Description string `json:"description"` // รายละเอียดงาน

	Department       string `json:"department_id"`     // แผนกที่เกี่ยวข้อง
	DepartmentName   string `json:"department_name"`   // ชื่อแผนก (อาจจะไม่ต้องมีถ้าไม่ใช้)
	Assignee         string `json:"assignee"`          // ผู้รับผิดชอบหลัก
	AssigneeName     string `json:"assignee_name"`     // ชื่อผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	AssigneeNickName string `json:"assignee_nickname"` // ชื่อเล่นผู้รับผิดชอบ (อาจจะไม่ต้องมีถ้าไม่ใช้)
	Importance       string `json:"importance"`        // ความสำคัญ (low|medium|high)

	StartDate time.Time `json:"start_date"` // วันที่เริ่มงาน
	EndDate   time.Time `json:"end_date"`   // วันที่สิ้นสุดงาน

	KPIID      string `json:"kpi_id"`      // รหัส KPI ที่เกี่ยวข้อง
	WorkFlowID string `json:"workflow_id"` // รหัส Workflow (อ้างอิง template/ค้นสถิติ)

	AppliedWorkflow TaskAppliedWorkflow `json:"applied_workflow"` // Snapshot workflow ที่ใช้ในงานนี้

	Status        string     `json:"status"`          // สถานะปัจจุบันของงาน (todos|in_progress|skip|done)
	StepName      string     `json:"step_name"`       // ชื่อขั้นตอนปัจจุบัน
	CreatedBy     string     `json:"created_by"`      // ผู้สร้างงาน
	CreatedByName string     `json:"created_by_name"` // ชื่อผู้สร้างงาน
	CreatedAt     time.Time  `json:"created_at"`      // วันที่สร้าง
	UpdatedAt     time.Time  `json:"updated_at"`      // วันที่อัปเดตล่าสุด
	DeletedAt     *time.Time `json:"deleted_at"`      // วันที่ลบ (soft delete)
}

type TaskAppliedWorkflow struct {
	WorkFlowID   string             `json:"workflow_id"`   // รหัส Workflow (UUID)
	WorkFlowName string             `json:"workflow_name"` // ชื่อ Workflow
	Department   string             `json:"department_id"` // แผนกที่เกี่ยวข้อง
	Description  string             `json:"description"`   // รายละเอียดเพิ่มเติม
	TotalHours   float64            `json:"total_hours"`   // ชั่วโมงรวม (แคชจากผลรวม step)
	Steps        []TaskWorkflowStep `json:"steps"`         // ลำดับขั้นตอนทั้งหมด
	Version      int                `json:"version"`       // เวอร์ชันของ template
}

type TaskWorkflowStep struct {
	StepID      string     `json:"step_id"`                // รหัส Step (UUID)
	StepName    string     `json:"step_name"`              // ชื่อ Step
	Description string     `json:"description,omitempty"`  // รายละเอียด (ไม่บังคับ)
	Hours       float64    `json:"hours"`                  // ชั่วโมงที่ใช้ (รองรับทศนิยม เช่น 0.5)
	Order       int        `json:"order"`                  // ลำดับขั้นตอน (1..N)
	Status      string     `json:"status"`                 // สถานะ (todo|in_progress|skip|done)
	StartedAt   *time.Time `json:"started_at,omitempty"`   // เวลาที่เริ่ม (optional)
	CompletedAt *time.Time `json:"completed_at,omitempty"` // เวลาที่เสร็จ (optional)
	Notes       string     `json:"notes,omitempty"`        // บันทึก/หมายเหตุ
	CreatedAt   time.Time  `json:"created_at"`             // วันที่สร้างขั้นตอนนี้
	UpdatedAt   time.Time  `json:"updated_at"`             // วันที่อัปเดตขั้นตอนล่าสุด
}

// NEW

// ===== Task Replace (PUT) Request =====
// ข้อกำหนด:
// - ต้องส่งครบทุกคีย์ที่ระบุ (required)
// - applied_workflow.steps ต้องส่งเป็น "ทั้งชุด" ที่ต้องการให้เก็บหลังอัปเดต
type UpdateTaskPutRequest struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	JobID       string `json:"job_id"`
	JobName     string `json:"job_name"`
	Description string `json:"description"`

	Department string `json:"department_id"`
	Assignee   string `json:"assignee"`
	Importance string `json:"importance"` // low|medium|high

	// แนะนำให้ client ส่งรูปแบบ YYYY-MM-DD (ฝั่งเซิร์ฟเวอร์แปลงเป็น time)
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	KPIID      string `json:"kpi_id"`
	WorkflowID string `json:"workflow_id"`

	// Snapshot ใหม่ทั้งก้อน (จำเป็น)
	AppliedWorkflow PutTaskAppliedWorkflow `json:"applied_workflow"`

	// สถานะระดับงาน (ถ้าส่งมาเราจะคำนวณจาก steps ทับให้อยู่ดี)
	Status string `json:"status"`

	// step_name จะถูกคำนวณทับ ไม่ต้องพยายามตั้งค่าเอง
}

type PutTaskAppliedWorkflow struct {
	WorkFlowID   string `json:"workflow_id"`
	WorkFlowName string `json:"workflow_name"`
	Department   string `json:"department_id"`
	Description  string `json:"description"`
	// TotalHours   float64           `json:"total_hours"` // เซิร์ฟเวอร์จะคำนวณทับใหม่
	Steps []PutWorkflowStep `json:"steps"` // ต้องส่งทั้งชุด
	// Version      int               `json:"version"`
}

type PutWorkflowStep struct {
	StepID      string     `json:"step_id"` // ถ้าต้องการเก็บ id เดิม ส่งมาเลย
	StepName    string     `json:"step_name"`
	Description string     `json:"description"`
	Hours       float64    `json:"hours"`
	Order       int        `json:"order"`                  // จะถูก normalize และ reindex เป็น 1..N
	Status      string     `json:"status"`                 // todo|in_progress|skip|done
	StartedAt   *time.Time `json:"started_at,omitempty"`   // optional
	CompletedAt *time.Time `json:"completed_at,omitempty"` // optional
	Notes       string     `json:"notes"`
}
