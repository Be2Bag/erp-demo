package dto

import "time"

// ---------- Request DTO ----------
type CreatePositionDTO struct {
	DepartmentID string `json:"department_id"` // รหัสแผนก (FK ไปยัง Department)
	PositionName string `json:"position_name"` // ชื่อตำแหน่งงาน
	Level        string `json:"level"`         // ระดับของตำแหน่งงาน
}
type UpdatePositionDTO struct {
	DepartmentID string `json:"department_id"` // รหัสแผนก (FK ไปยัง Department)
	PositionName string `json:"position_name"` // ชื่อตำแหน่งงาน
	Level        string `json:"level"`         // ระดับของตำแหน่งงาน
	Note         string `json:"note"`
}

type RequestListPosition struct {
	Page       int    `query:"page"`          // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit      int    `query:"limit"`         // จำนวนรายการต่อหน้า
	Search     string `query:"search"`        // คำค้นหาสำหรับกรองข้อมูล
	Department string `query:"department_id"` // แผนก
	SortBy     string `query:"sort_by"`       // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder  string `query:"sort_order"`    // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

// ---------- Response DTO ----------
type PositionDTO struct {
	PositionID   string    `json:"position_id"`   // รหัสตำแหน่งงาน (ไม่ซ้ำกัน)
	ManagerName  string    `json:"manager_name"`  // ชื่อผู้จัดการ
	DepartmentID string    `json:"department_id"` // รหัสแผนก (FK ไปยัง Department)
	PositionName string    `json:"position_name"` // ชื่อตำแหน่งงาน
	Level        string    `json:"level"`         // ระดับของตำแหน่งงาน
	CreatedAt    time.Time `json:"created_at"`    // วันที่สร้าง
	UpdatedAt    time.Time `json:"updated_at"`    // วันที่อัปเดต
}
