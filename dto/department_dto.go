package dto

import "time"

// ---------- Request DTO ----------
type CreateDepartmentDTO struct {
	DepartmentName string `json:"department_name"` // ชื่อแผนก
	ManagerID      string `json:"manager_id"`      // รหัสผู้จัดการแผนก (FK ไปยัง User)
}
type UpdateDepartmentDTO struct {
	DepartmentName string `json:"department_name"` // ชื่อแผนก
	ManagerID      string `json:"manager_id"`      // รหัสผู้จัดการแผนก (FK ไปยัง User)
}

type RequestListDepartment struct {
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
}

// ---------- Response DTO ----------
type DepartmentDTO struct {
	DepartmentID   string    `json:"department_id"`   // รหัสแผนก
	DepartmentName string    `json:"department_name"` // ชื่อแผนก
	ManagerID      string    `json:"manager_id"`      // รหัสผู้จัดการแผนก (FK ไปยัง User)
	ManagerName    string    `json:"manager_name"`    // ชื่อผู้จัดการแผนก
	CreatedAt      time.Time `json:"created_at"`      // วันที่สร้าง
	UpdatedAt      time.Time `json:"updated_at"`      // วันที่อัปเดต
}
