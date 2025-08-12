package models

import "time"

const CollectionDepartments = "departments"

type Department struct {
	DepartmentID   string     `bson:"department_id" json:"department_id"`     // รหัสแผนก (ไม่ซ้ำกัน)
	DepartmentName string     `bson:"department_name" json:"department_name"` // ชื่อแผนก
	ManagerID      string     `bson:"manager_id" json:"manager_id"`           // รหัสผู้จัดการแผนก (FK ไปยัง User)
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`           // วันที่สร้างข้อมูลนี้
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at"`           // วันที่แก้ไขข้อมูลล่าสุด
	DeletedAt      *time.Time `bson:"deleted_at" json:"deleted_at"`           // วันที่ลบข้อมูล (soft delete)
}
