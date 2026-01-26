package dto

import "time"

// ---------- Request DTO ----------
type CreateProjectDTO struct {
	ProjectName string `json:"project_name"`
}

type UpdateProjectDTO = struct {
	ProjectName string `json:"project_name"`
	Note        string `json:"note"`
}

type RequestListProject struct {
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------
type ProjectDTO struct {
	CreatedAt   time.Time  `json:"created_at"` // เวลาสร้าง
	UpdatedAt   time.Time  `json:"updated_at"` // เวลาอัปเดตล่าสุด
	DeletedAt   *time.Time `json:"deleted_at"` // เวลาเมื่อถูกลบ (soft delete)
	ProjectID   string     `json:"project_id"`
	ProjectName string     `json:"project_name"`
	CreatedBy   string     `json:"created_by"` // ผู้สร้าง
}
