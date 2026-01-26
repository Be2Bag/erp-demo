package models

import "time"

const CollectionUserTaskStats = "user_task_stats"

type UserTaskStats struct {
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
	KPI       UserTaskKPI `json:"kpi,omitempty" bson:"kpi,omitempty"` // KPI (ถ้ามี)

	DeletedAt    *time.Time     `json:"deleted_at" bson:"deleted_at"`
	UserID       string         `json:"user_id" bson:"user_id"`                       // อ้างอิงผู้ใช้
	DepartmentID string         `json:"department_id,omitempty" bson:"department_id"` // ช่วย filter รายแผนก (เลือกใส่)
	Totals       UserTaskTotals `json:"totals" bson:"totals"`                         // ตัวเลขรวม
}

// จำนวนงานรวมตามสถานะ
type UserTaskTotals struct {
	Assigned   int `json:"assigned" bson:"assigned"`         // งานทั้งหมดที่เป็น assignee อยู่
	Open       int `json:"open" bson:"open"`                 // todo + in_progress
	InProgress int `json:"in_progress" bson:"in_progress"`   // งานสถานะ in_progress
	Completed  int `json:"completed" bson:"completed"`       // งานสถานะ done
	Skipped    int `json:"skipped,omitempty" bson:"skipped"` // ถ้าต้องการนับงานที่ถูกข้าม (option)
}

// KPI ของผู้ใช้ (เก็บเป็นตัวเลข 0..1 แล้วไปฟอร์แมตเป็น % ที่ชั้นตอบกลับ)
type UserTaskKPI struct {
	Score            *float64   `json:"score,omitempty" bson:"score,omitempty"`                           // เช่น 1.0 = 100%
	LastCalculatedAt *time.Time `json:"last_calculated_at,omitempty" bson:"last_calculated_at,omitempty"` // เวลาอัปเดต KPI ล่าสุด
}

/*
แนะนำดัชนีใน MongoDB:

db.user_task_stats.createIndex({ user_id: 1, period: 1 }, { unique: true });
*/
