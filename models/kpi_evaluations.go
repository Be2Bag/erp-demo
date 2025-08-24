package models

import "time"

const CollectionKPIEvaluations = "kpi_evaluations"

type KPIEvaluation struct {
	EvaluationID string     `bson:"evaluation_id" json:"evaluation_id"`         // UUID
	ProjectID    string     `bson:"project_id" json:"project_id"`               // อ้างถึง Project
	JobID        string     `bson:"job_id" json:"job_id"`                       // อ้างถึง SignJob
	TaskID       string     `bson:"task_id,omitempty" json:"task_id,omitempty"` // ถ้ามีงานย่อย
	KPIID        string     `bson:"kpi_id" json:"kpi_id"`                       // อ้างถึง KPITemplate
	Version      int        `bson:"version" json:"version"`                     // ใช้ version ของ KPI template ตอนนั้น
	EvaluatorID  string     `bson:"evaluator_id" json:"evaluator_id"`           // ใครประเมิน
	EvaluateeID  string     `bson:"evaluatee_id" json:"evaluatee_id"`           // ใครถูกประเมิน (เช่น assignee)
	Department   string     `bson:"department_id" json:"department_id"`         // แผนก
	Scores       []KPIScore `bson:"scores" json:"scores"`                       // รายการคะแนนแต่ละ item
	TotalScore   float32    `bson:"total_score" json:"total_score"`             // รวมคะแนน
	Feedback     string     `bson:"feedback" json:"feedback"`                   // คอมเมนต์รวม
	IsEvaluated  bool       `bson:"is_evaluated" json:"is_evaluated"`           // ประเมินแล้วหรือยัง
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" bson:"deleted_at"` // วันที่ลบ (ถ้ามี)
}

type KPIScore struct {
	ItemID   string `bson:"item_id" json:"item_id"` // อ้างถึง item ใน KPI template
	Name     string `bson:"name" json:"name"`       // สำเนาชื่อ item
	Category string `bson:"category" json:"category"`
	Weight   int    `bson:"weight" json:"weight"` // weight ตอนนั้น
	MaxScore int    `bson:"max_score" json:"max_score"`
	Score    int    `bson:"score" json:"score"` // คะแนนที่ให้จริง
	Notes    string `bson:"notes,omitempty" json:"notes"`
}
