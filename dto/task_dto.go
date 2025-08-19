package dto

// ===== Task Create Request =====
type CreateTaskRequest struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	JobName     string `json:"job_name"`
	Description string `json:"description,omitempty"`
	Department  string `json:"department"`
	Assignee    string `json:"assignee"`
	Importance  string `json:"importance"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	KPIID       string `json:"kpi_id,omitempty"`
	WorkflowID  string `json:"workflow_id"`

	ExtraSteps []ExtraStepRequest `json:"extra_steps,omitempty"`
}

// ===== Extra step ที่ FE จะส่งเพิ่ม =====
type ExtraStepRequest struct {
	StepName    string         `json:"step_name"`
	Description string         `json:"description,omitempty"`
	Hours       float64        `json:"hours"`
	Insert      InsertPosition `json:"insert"`
}

type InsertPosition struct {
	Position     string  `json:"position"` // start|end|before|after
	AnchorStepID *string `json:"anchor_step_id,omitempty"`
}
