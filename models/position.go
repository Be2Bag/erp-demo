package models

import "time"

const CollectionPositions = "positions"

type Position struct {
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"` // วันที่สร้างข้อมูลนี้
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"` // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt    *time.Time `bson:"deleted_at" json:"deleted_at"` // วันที่ลบข้อมูล (soft delete)
	Note         *string    `bson:"note" json:"note"`
	PositionID   string     `bson:"position_id" json:"position_id"`     // รหัสตำแหน่งงาน (ไม่ซ้ำกัน)
	DepartmentID string     `bson:"department_id" json:"department_id"` // รหัสแผนก (FK ไปยัง Department)
	PositionName string     `bson:"position_name" json:"position_name"` // ชื่อตำแหน่งงาน
	Level        string     `bson:"level" json:"level"`                 // ระดับของตำแหน่งงาน
}
