package dto

import "time"

// ---------- Request DTO ----------.
type CreateBankAccountsDTO struct {
	BankName    string `json:"bank_name"`    // ชื่อธนาคาร
	AccountNo   string `json:"account_no"`   // เลขที่บัญชี
	AccountName string `json:"account_name"` // ชื่อบัญชี
}

type UpdateBankAccountsDTO = struct {
	BankName    string `json:"bank_name"`    // ชื่อธนาคาร
	AccountNo   string `json:"account_no"`   // เลขที่บัญชี
	AccountName string `json:"account_name"` // ชื่อบัญชี
	Note        string `json:"note"`         // หมายเหตุ
}

type RequestListBankAccounts struct {
	Search    string `query:"search"`     // คำค้นหาสำหรับกรองข้อมูล
	SortBy    string `query:"sort_by"`    // คอลัมน์ที่ต้องการเรียงลำดับ
	SortOrder string `query:"sort_order"` // ทิศทางการเรียงลำดับ (asc หรือ desc)
	Page      int    `query:"page"`       // หมายเลขหน้าที่ต้องการดึงข้อมูล
	Limit     int    `query:"limit"`      // จำนวนรายการต่อหน้า
}

// ---------- Response DTO ----------.
type BankAccountsDTO struct {
	CreatedAt   time.Time  `json:"created_at"`   // เวลาสร้าง
	UpdatedAt   time.Time  `json:"updated_at"`   // เวลาอัปเดตล่าสุด
	DeletedAt   *time.Time `json:"deleted_at"`   // เวลาเมื่อถูกลบ (soft delete)
	BankID      string     `json:"bank_id"`      // รหัสบัญชีธนาคาร
	BankName    string     `json:"bank_name"`    // ชื่อธนาคาร
	AccountNo   string     `json:"account_no"`   // เลขที่บัญชี
	AccountName string     `json:"account_name"` // ชื่อบัญชี
	CreatedBy   string     `json:"created_by"`   // ผู้สร้าง
}
