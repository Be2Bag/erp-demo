package dto

import "time"

type CreateKPIEvaluationRequest struct {
	JobID       string            `json:"job_id" binding:"required"`
	TaskID      string            `json:"task_id,omitempty"`
	KPIID       string            `json:"kpi_id" binding:"required"`
	Version     int               `json:"version" binding:"required"`
	EvaluatorID string            `json:"evaluator_id" binding:"required"`
	EvaluateeID string            `json:"evaluatee_id" binding:"required"`
	Department  string            `json:"department_id" binding:"required"`
	Scores      []KPIScoreRequest `json:"scores" binding:"required"`
	Feedback    string            `json:"feedback,omitempty"`
}

type KPIScoreRequest struct {
	ItemID string `json:"item_id" binding:"required"`
	Score  int    `json:"score" binding:"required"`
	Notes  string `json:"notes,omitempty"`
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
