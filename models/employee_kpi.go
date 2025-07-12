package models

import "time"

// EmployeeKPI เก็บผลการประเมิน KPI ของพนักงานในแต่ละช่วงเวลา โดยเชื่อมโยงกับ KPIDefinition และ User ผ่าน FK
type EmployeeKPI struct {
	RecordID   string     `bson:"record_id" json:"record_id"`     // รหัสบันทึกของ KPI แต่ละรายการ
	KPIID      string     `bson:"kpi_id" json:"kpi_id"`           // รหัส KPI จาก KPIDefinition
	EmployeeID string     `bson:"employee_id" json:"employee_id"` // รหัสพนักงาน (FK ไปยัง User)
	Period     string     `bson:"period" json:"period"`           // ช่วงเวลาของการประเมิน (เช่น "2023-Q3")
	Value      float64    `bson:"value" json:"value"`             // ค่าที่วัดได้จริง
	Score      float64    `bson:"score" json:"score"`             // คะแนนที่คำนวณตามน้ำหนัก
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`   // วันที่สร้างบันทึก
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`   // วันที่แก้ไขบันทึกล่าสุด
	DeletedAt  *time.Time `bson:"deleted_at" json:"deleted_at"`   // วันที่ลบข้อมูล (soft delete)
}
