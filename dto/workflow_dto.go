package dto

import "time"

// ---------- Request DTO ----------
type CreateWorkflowTemplateDTO struct {
	Name        string                  `json:"name"`        // ชื่อ Template
	Department  string                  `json:"department"`  // แผนก (dropdown)
	Description string                  `json:"description"` // คำอธิบายการใช้งาน
	Steps       []CreateWorkflowStepDTO `json:"steps"`       // ขั้นตอนการทำงาน
}

type CreateWorkflowStepDTO struct {
	Name        string  `json:"name"`                  // ชื่อขั้นตอน
	Description string  `json:"description,omitempty"` // คำอธิบาย (ไม่บังคับ)
	Hours       float64 `json:"hours"`                 // ชั่วโมง (รองรับทศนิยม)
	Order       int     `json:"order"`                 // ลำดับ (1,2,3,...)
}

// ---------- Response DTO ----------
type WorkflowTemplateDTO struct {
	TemplateID  string            `json:"template_id"`
	Name        string            `json:"name"`
	Department  string            `json:"department"`
	Description string            `json:"description"`
	TotalHours  float64           `json:"total_hours"` // ผลรวมชั่วโมงจากทุก step
	Steps       []WorkflowStepDTO `json:"steps"`
	IsActive    bool              `json:"is_active"`
	Version     int               `json:"version"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type WorkflowStepDTO struct {
	StepID      string    `json:"step_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Hours       float64   `json:"hours"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
