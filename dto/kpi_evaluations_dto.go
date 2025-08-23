package dto

import "time"

type CreateKPIEvaluationRequest struct {
	JobID       string            `json:"job_id" binding:"required"`        // อ้างถึง SignJob
	TaskID      string            `json:"task_id,omitempty"`                // งานย่อย (ถ้ามี)
	KPIID       string            `json:"kpi_id" binding:"required"`        // อ้างถึง KPITemplate
	EvaluateeID string            `json:"evaluatee_id" binding:"required"`  // ผู้ถูกประเมิน
	Department  string            `json:"department_id" binding:"required"` // แผนก
	Scores      []KPIScoreRequest `json:"scores" binding:"required"`        // รายการคะแนนแต่ละ item
	Feedback    string            `json:"feedback,omitempty"`               // คอมเมนต์รวม (ถ้ามี)
}

type KPIScoreRequest struct {
	ItemID string `json:"item_id" binding:"required"` // อ้างถึง item ใน KPI template
	Score  int    `json:"score" binding:"required"`   // คะแนนที่ให้
	Notes  string `json:"notes,omitempty"`            // หมายเหตุเพิ่มเติม (ถ้ามี)
}

type RequestListKPIEvaluation struct {
	Page       int    `query:"page"`          // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`         // จำนวนรายการต่อหน้า
	Search     string `query:"search"`        // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department_id"` // แผนก
	SortBy     string `query:"sort_by"`       // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"`    // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

// <===================== Response ===============================>

type KPIEvaluationResponse struct {
	EvaluationID string             `json:"evaluation_id"`
	JobID        string             `json:"job_id"`
	TaskID       string             `json:"task_id,omitempty"`
	KPIID        string             `json:"kpi_id"`
	KPIName      string             `json:"kpi_name"`
	Version      int                `json:"version"`
	EvaluatorID  string             `json:"evaluator_id"`
	EvaluateeID  string             `json:"evaluatee_id"`
	Department   string             `json:"department_id"`
	Scores       []KPIScoreResponse `json:"scores"`
	TotalScore   int                `json:"total_score"`
	Feedback     string             `json:"feedback,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type KPIScoreResponse struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Weight   int    `json:"weight"`
	MaxScore int    `json:"max_score"`
	Score    int    `json:"score"`
	Notes    string `json:"notes,omitempty"`
}
