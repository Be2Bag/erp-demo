package models

import "time"

type KPITemplate struct {
	KPIID       string            `bson:"kpi_id" json:"kpi_id"`             // รหัส Template
	Name        string            `bson:"name" json:"name"`                 // ชื่อ Template
	Department  string            `bson:"department" json:"department"`     // แผนก
	Templates   []KPITemplateList `json:"templates"`                        // รายการ KPI
	TargetValue int               `bson:"target_value" json:"target_value"` // ค่าเป้าหมายรวม (100%)
	IsActive    bool              `bson:"is_active" json:"is_active"`
	CreatedBy   string            `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type KPITemplateList struct {
	KPIID       string     `bson:"kpi_id" json:"kpi_id"`
	Name        string     `bson:"name" json:"name"`
	Description string     `bson:"description" json:"description"`
	Category    string     `bson:"category" json:"category"`
	MaxScore    int        `bson:"max_score" json:"max_score"` // คะแนนเต็ม
	Value       int        `bson:"value" json:"value"`         // น้ำหนัก %
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// type KPIEvaluation struct {
// 	EvaluationID string     `bson:"evaluation_id" json:"evaluation_id"` // รหัสการประเมิน
// 	KPIID        string     `bson:"kpi_id" json:"kpi_id"`               // รหัส KPI ที่ประเมิน
// 	EmployeeID   string     `bson:"employee_id" json:"employee_id"`     // รหัสพนักงานที่ถูกประเมิน
// 	Period       string     `bson:"period" json:"period"`               // ช่วงเวลาการประเมิน
// 	Value        float64    `bson:"value" json:"value"`                 // ค่าที่วัดได้
// 	Score        float64    `bson:"score" json:"score"`                 // คะแนนที่คำนวณได้
// 	Comment      string     `bson:"comment" json:"comment"`             // ความคิดเห็นเพิ่มเติม
// 	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`       // วันที่สร้าง
// 	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`       // วันที่แก้ไขล่าสุด
// 	DeletedAt    *time.Time `bson:"deleted_at" json:"deleted_at"`       // วันที่ลบ (soft delete)
// }
