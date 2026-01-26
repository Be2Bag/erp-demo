package dto

import "time"

// ---------- Request DTO ----------
type CreateWorkflowTemplateDTO struct {
	WorkFlowName string                  `json:"workflow_name"` // ชื่อ Template
	Department   string                  `json:"department_id"` // แผนก (dropdown)
	Description  string                  `json:"description"`   // คำอธิบายการใช้งาน
	Steps        []CreateWorkflowStepDTO `json:"steps"`         // ขั้นตอนการทำงาน
}

type CreateWorkflowStepDTO struct {
	StepName    string  `json:"step_name"`             // ชื่อ Template
	Description string  `json:"description,omitempty"` // คำอธิบาย (ไม่บังคับ)
	Hours       float64 `json:"hours"`                 // ชั่วโมง (รองรับทศนิยม)
	Order       int     `json:"order"`                 // ลำดับ (1,2,3,...)
}

// Partial update payload (use pointer fields)
type UpdateWorkflowTemplateDTO struct {
	Steps        *[]CreateWorkflowStepDTO `json:"steps,omitempty"`
	WorkFlowName string                   `json:"workflow_name"` // ชื่อ Template
	Department   string                   `json:"department_id,omitempty"`
	Description  string                   `json:"description,omitempty"`
}

type RequestListWorkflow struct {
	Search     string `query:"search"`        // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department_id"` // แผนก
	SortBy     string `query:"sort_by"`       // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"`    // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page       int    `query:"page"`          // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`         // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------
type WorkflowTemplateDTO struct {
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	WorkFlowID   string            `json:"workflow_id"`
	WorkFlowName string            `json:"workflow_name"` // ชื่อ Template
	Department   string            `json:"department_id"`
	Description  string            `json:"description"`
	CreatedBy    string            `json:"created_by"`
	Steps        []WorkflowStepDTO `json:"steps"`
	TotalHours   float64           `json:"total_hours"` // ผลรวมชั่วโมงจากทุก step
	Version      int               `json:"version"`
	IsActive     bool              `json:"is_active"`
}

type WorkflowStepDTO struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StepID      string    `json:"step_id"`
	StepName    string    `json:"step_name"`
	Description string    `json:"description,omitempty"`
	Hours       float64   `json:"hours"`
	Order       int       `json:"order"`
}
