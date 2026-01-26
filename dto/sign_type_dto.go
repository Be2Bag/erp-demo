package dto

import "time"

// ---------- Request DTO ----------
type CreateSignTypeDTO struct {
	NameTH string `json:"name_th" binding:"required"` // ชื่อภาษาไทย (บังคับ)
	NameEN string `json:"name_en" binding:"required"` // ชื่อภาษาอังกฤษ (บังคับ)
}

type UpdateSignTypeDTO struct {
	NameTH string `json:"name_th" binding:"required"` // ชื่อภาษาไทย (บังคับ)
	NameEN string `json:"name_en" binding:"required"` // ชื่อภาษาอังกฤษ (บังคับ)
}

type RequestListSignType struct {
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------
type SignTypeDTO struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	TypeID    string     `json:"type_id"`
	NameTH    string     `json:"name_th"`
	NameEN    string     `json:"name_en"`
	CreatedBy string     `json:"created_by"`
}
