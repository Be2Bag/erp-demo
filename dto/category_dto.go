package dto

import "time"

// ---------- Request DTO ----------

// สำหรับสร้าง Category ใหม่
type CreateCategoryDTO struct {
	DepartmentID   string `json:"department_id" validate:"required"`    // รหัสแผนก (FK ไปยัง Department)
	CategoryNameTH string `json:"category_name_th" validate:"required"` // ชื่อหมวดหมู่ภาษาไทย
	CategoryNameEN string `json:"category_name_en" validate:"required"` // ชื่อหมวดหมู่ภาษาอังกฤษ
	Description    string `json:"description" validate:"required"`      // รายละเอียด
}

// สำหรับอัปเดต Category
type UpdateCategoryDTO struct {
	DepartmentID   *string `json:"department_id,omitempty"`    // รหัสแผนก
	CategoryNameTH *string `json:"category_name_th,omitempty"` // ชื่อหมวดหมู่ภาษาไทย
	CategoryNameEN *string `json:"category_name_en,omitempty"` // ชื่อหมวดหมู่ภาษาอังกฤษ
	Description    *string `json:"description,omitempty"`      // รายละเอียด
	Note           *string `json:"note,omitempty"`             // หมายเหตุเพิ่มเติม
}

// สำหรับ query list category
type RequestListCategory struct {
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------

type CategoryDTO struct {
	CreatedAt      time.Time  `json:"created_at"`       // วันเวลาที่สร้าง
	UpdatedAt      time.Time  `json:"updated_at"`       // วันเวลาที่แก้ไขล่าสุด
	DeletedAt      *time.Time `json:"deleted_at"`       // วันเวลาที่ลบ (ถ้ามี)
	Note           *string    `json:"note,omitempty"`   // หมายเหตุเพิ่มเติม
	CategoryID     string     `json:"category_id"`      // UUID
	DepartmentID   string     `json:"department_id"`    // รหัสแผนก
	CategoryNameTH string     `json:"category_name_th"` // ชื่อหมวดหมู่ภาษาไทย
	CategoryNameEN string     `json:"category_name_en"` // ชื่อหมวดหมู่ภาษาอังกฤษ
	Description    string     `json:"description"`      // รายละเอียด
	CreatedBy      string     `json:"created_by"`       // ผู้สร้าง (user_id หรือ username)
}
