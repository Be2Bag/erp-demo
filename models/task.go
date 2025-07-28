package models

import "time"

type Task struct {
	TaskID         string         `bson:"task_id" json:"task_id"`                 // รหัสงาน (ไม่ซ้ำกัน)
	ProjectName    string         `bson:"project_name" json:"project_name"`       // ชื่อโปรเจกต์
	Title          string         `bson:"title" json:"title"`                     // ชื่องาน
	Description    string         `bson:"description" json:"description"`         // รายละเอียดงาน
	Department     string         `bson:"department" json:"department"`           // แผนกที่เกี่ยวข้อง
	AssigneeID     string         `bson:"assignee_id" json:"assignee_id"`         // รหัสผู้รับผิดชอบงาน (FK ไปยัง User)
	Priority       string         `bson:"priority" json:"priority"`               // ลำดับความสำคัญ (เช่น high, medium, low)
	Status         string         `bson:"status" json:"status"`                   // สถานะงาน (เช่น pending, in_progress, completed)
	StartDate      time.Time      `bson:"start_date" json:"start_date"`           // วันที่เริ่มงาน
	EndDate        time.Time      `bson:"end_date" json:"end_date"`               // วันที่สิ้นสุดงาน
	Progress       int            `bson:"progress" json:"progress"`               // ความคืบหน้า (%)
	KPIScore       float64        `bson:"kpi_score" json:"kpi_score"`             // คะแนน KPI
	EstimatedHours int            `bson:"estimated_hours" json:"estimated_hours"` // ชั่วโมงที่คาดว่าจะใช้
	ActualHours    int            `bson:"actual_hours" json:"actual_hours"`       // ชั่วโมงที่ใช้จริง
	Budget         float64        `bson:"budget" json:"budget"`                   // งบประมาณที่ตั้งไว้
	ActualCost     float64        `bson:"actual_cost" json:"actual_cost"`         // ค่าใช้จ่ายจริง
	WorkflowSteps  []WorkflowStep `bson:"workflow_steps" json:"workflow_steps"`   // ขั้นตอนการทำงาน
	CreatedBy      string         `bson:"created_by" json:"created_by"`           // รหัสผู้สร้างงาน
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`           // วันที่สร้างงาน
	UpdatedAt      time.Time      `bson:"updated_at" json:"updated_at"`           // วันที่แก้ไขงานล่าสุด
	DeletedAt      *time.Time     `bson:"deleted_at" json:"deleted_at"`           // วันที่ลบงาน (soft delete)
}

type WorkflowStep struct {
	StepID      string    `bson:"step_id" json:"step_id"`           // รหัสขั้นตอน
	Title       string    `bson:"title" json:"title"`               // ชื่อขั้นตอน
	Description string    `bson:"description" json:"description"`   // รายละเอียดขั้นตอน
	Order       int       `bson:"order" json:"order"`               // ลำดับขั้นตอน
	Status      string    `bson:"status" json:"status"`             // สถานะขั้นตอน (pending, completed)
	CompletedAt time.Time `bson:"completed_at" json:"completed_at"` // วันที่เสร็จสิ้นขั้นตอน
}
