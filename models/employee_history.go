package models

import "time"

const CollectionEmploymentHistories = "employment_histories"

type EmploymentHistory struct {
	FromDate       time.Time  `bson:"from_date" json:"from_date"` // วันที่เริ่มต้น
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at"`
	ToDate         *time.Time `bson:"to_date,omitempty" json:"to_date"`       // วันที่สิ้นสุด (nullable ถ้ายังทำอยู่)
	DeletedAt      *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"` // soft delete
	UserID         string     `bson:"user_id" json:"user_id"`                 // รหัสผู้ใช้ที่เกี่ยวข้อง
	PositionID     string     `bson:"position_id" json:"position_id"`         // ตำแหน่งในช่วงเวลานั้น
	DepartmentID   string     `bson:"department_id" json:"department_id"`     // แผนกในช่วงเวลานั้น
	EmploymentType string     `bson:"employment_type" json:"employment_type"` // ประเภทการจ้าง (เช่น full-time, intern)
	Note           string     `bson:"note,omitempty" json:"note,omitempty"`   // หมายเหตุ (ถ้ามี)
}
