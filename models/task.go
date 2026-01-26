package models

import (
	"time"
)

const CollectionTasks = "tasks" // ชื่อ collection ใน MongoDB สำหรับเก็บงาน (tasks)

type Tasks struct {
	StartDate time.Time `bson:"start_date" json:"start_date"` // วันที่เริ่มงาน
	EndDate   time.Time `bson:"end_date" json:"end_date"`     // วันที่สิ้นสุดงาน

	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`     // วันที่สร้าง
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`     // วันที่อัปเดตล่าสุด
	DeletedAt   *time.Time `bson:"deleted_at" json:"deleted_at"`     // วันที่ลบ (soft delete)
	TaskID      string     `bson:"task_id" json:"task_id"`           // รหัสงาน (UUID/unique)
	ProjectID   string     `bson:"project_id" json:"project_id"`     // รหัสโปรเจกต์
	ProjectName string     `bson:"project_name" json:"project_name"` // ชื่อโปรเจกต์
	JobID       string     `bson:"job_id" json:"job_id"`             // รหัสงาน
	JobName     string     `bson:"job_name" json:"job_name"`         // ชื่องาน
	Description string     `bson:"description" json:"description"`   // รายละเอียดงาน

	Department       string `bson:"department_id" json:"department_id"`         // แผนกที่เกี่ยวข้อง
	Assignee         string `bson:"assignee" json:"assignee"`                   // ผู้รับผิดชอบหลัก
	AssigneeName     string `bson:"assignee_name" json:"assignee_name"`         // ชื่อผู้รับผิดชอบหลัก
	AssigneeNickName string `bson:"assignee_nickname" json:"assignee_nickname"` // ชื่อเล่นผู้รับผิดชอบหลัก
	Importance       string `bson:"importance" json:"importance"`               // ความสำคัญ (low|medium|high)

	KPIID      string `bson:"kpi_id" json:"kpi_id"`           // รหัส KPI ที่เกี่ยวข้อง
	WorkFlowID string `bson:"workflow_id" json:"workflow_id"` // รหัส Workflow (อ้างอิง template/ค้นสถิติ)

	Status    string `bson:"status" json:"status"`         // สถานะปัจจุบันของงาน (todos|in_progress|skip|done)
	StepName  string `bson:"step_name" json:"step_name"`   // ชื่อขั้นตอนปัจจุบัน
	CreatedBy string `bson:"created_by" json:"created_by"` // ผู้สร้างงาน

	AppliedWorkflow TaskAppliedWorkflow `bson:"applied_workflow" json:"applied_workflow"` // Snapshot workflow ที่ใช้ในงานนี้

}

type TaskAppliedWorkflow struct {
	WorkFlowID   string             `bson:"workflow_id" json:"workflow_id"`     // รหัส Workflow (UUID)
	WorkFlowName string             `bson:"workflow_name" json:"workflow_name"` // ชื่อ Workflow
	Department   string             `bson:"department_id" json:"department_id"` // แผนกที่เกี่ยวข้อง
	Description  string             `bson:"description" json:"description"`     // รายละเอียดเพิ่มเติม
	Steps        []TaskWorkflowStep `bson:"steps" json:"steps"`                 // ลำดับขั้นตอนทั้งหมด
	TotalHours   float64            `bson:"total_hours" json:"total_hours"`     // ชั่วโมงรวม (แคชจากผลรวม step)
	Version      int                `bson:"version" json:"version"`             // เวอร์ชันของ template
}

type TaskWorkflowStep struct {
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`                         // วันที่สร้างขั้นตอนนี้
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`                         // วันที่อัปเดตขั้นตอนล่าสุด
	StartedAt   *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`     // เวลาที่เริ่ม (optional)
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"` // เวลาที่เสร็จ (optional)
	StepID      string     `bson:"step_id" json:"step_id"`                               // รหัส Step (UUID)
	StepName    string     `bson:"step_name" json:"step_name"`                           // ชื่อ Step
	Description string     `bson:"description,omitempty" json:"description,omitempty"`   // รายละเอียด (ไม่บังคับ)
	Status      string     `bson:"status" json:"status"`                                 // สถานะ (todo|in_progress|skip|done)
	Notes       string     `bson:"notes,omitempty" json:"notes,omitempty"`               // บันทึก/หมายเหตุ
	Hours       float64    `bson:"hours" json:"hours"`                                   // ชั่วโมงที่ใช้ (รองรับทศนิยม เช่น 0.5)
	Order       int        `bson:"order" json:"order"`                                   // ลำดับขั้นตอน (1..N)
}
