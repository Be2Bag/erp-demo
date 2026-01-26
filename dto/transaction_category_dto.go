package dto

import "time"

// ---------- Request DTO ----------

// สำหรับสร้าง TransactionCategory ใหม่
type CreateTransactionCategoryDTO struct {
	Type                      string `json:"type" validate:"required"`                         // ประเภทหมวดหมู่ (เช่น รายรับ, รายจ่าย)
	TransactionCategoryNameTH string `json:"transaction_category_name_th" validate:"required"` // ชื่อหมวดหมู่ภาษาไทย
	Description               string `json:"description" validate:"required"`                  // รายละเอียด
}

// สำหรับอัปเดต Category
type UpdateTransactionCategoryDTO struct {
	Type                      *string `json:"type,omitempty"`                         // ประเภทหมวดหมู่
	TransactionCategoryNameTH *string `json:"transaction_category_name_th,omitempty"` // ชื่อหมวดหมู่ภาษาไทย
	Description               *string `json:"description,omitempty"`                  // รายละเอียด
	Note                      *string `json:"note,omitempty"`                         // หมายเหตุเพิ่มเติม
}

// สำหรับ query list category
type RequestListTransactionCategory struct {
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Type      string `query:"type"`       // ประเภทหมวดหมู่ (เช่น รายรับ, รายจ่าย)
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------

type TransactionCategoryDTO struct {
	CreatedAt                 time.Time  `json:"created_at"`                   // วันเวลาที่สร้าง
	UpdatedAt                 time.Time  `json:"updated_at"`                   // วันเวลาที่แก้ไขล่าสุด
	DeletedAt                 *time.Time `json:"deleted_at"`                   // วันเวลาที่ลบ (ถ้ามี)
	Note                      *string    `json:"note,omitempty"`               // หมายเหตุเพิ่มเติม
	TransactionCategoryID     string     `json:"transaction_category_id"`      // UUID
	Type                      string     `json:"type"`                         // ประเภทหมวดหมู่
	TransactionCategoryNameTH string     `json:"transaction_category_name_th"` // ชื่อหมวดหมู่ภาษาไทย
	Description               string     `json:"description"`                  // รายละเอียด
	CreatedBy                 string     `json:"created_by"`                   // ผู้สร้าง
}
