package dto

import "time"

// ---------- Request DTO ----------

type CreatePayableDTO struct {
	Supplier   string  `json:"supplier"`              // ชื่อผู้จำหน่าย
	BankID     string  `json:"bank_id"`               // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	PurchaseNo string  `json:"purchase_no"`           // เลขที่ใบสั่งซื้อ
	InvoiceNo  string  `json:"invoice_no"`            // เลขที่ใบแจ้งหนี้
	IssueDate  string  `json:"issue_date"`            // วันที่ออกใบแจ้งหนี้
	DueDate    string  `json:"due_date"`              // วันที่ครบกำหนดชำระ
	Amount     float64 `json:"amount"`                // จำนวนเงิน
	Balance    float64 `json:"balance"`               // ยอดคงเหลือ
	PaymentRef string  `json:"payment_ref,omitempty"` // เลขที่อ้างอิงการชำระเงิน
	Note       string  `json:"note,omitempty"`        // หมายเหตุ
}

type UpdatePayableDTO struct {
	Supplier   string  `json:"supplier,omitempty"`    // ชื่อผู้จำหน่าย
	BankID     string  `json:"bank_id,omitempty"`     // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	PurchaseNo string  `json:"purchase_no,omitempty"` // เลขที่ใบสั่งซื้อ
	InvoiceNo  string  `json:"invoice_no,omitempty"`  // เลขที่ใบแจ้งหนี้
	IssueDate  string  `json:"issue_date,omitempty"`  // วันที่ออกใบแจ้งหนี้
	DueDate    string  `json:"due_date,omitempty"`    // วันที่ครบกำหนดชำระ
	Amount     float64 `json:"amount,omitempty"`      // จำนวนเงิน
	Balance    float64 `json:"balance,omitempty"`     // ยอดคงเหลือ
	Status     string  `json:"status,omitempty"`      // สถานะ
	PaymentRef string  `json:"payment_ref,omitempty"` // เลขที่อ้างอิงการชำระเงิน
	Note       string  `json:"note,omitempty"`        // หมายเหตุ
}

type RequestListPayable struct {
	Page      int    `query:"page"`       // หน้า
	Limit     int    `query:"limit"`      // จำนวนต่อหน้า
	Search    string `query:"search"`     // คำค้นหา
	SortBy    string `query:"sort_by"`    // เรียงตาม
	SortOrder string `query:"sort_order"` // ลำดับการเรียง
	Status    string `query:"status"`     // สถานะ
	Supplier  string `query:"supplier"`   // ผู้จำหน่าย
}

type RequestSummaryPayable struct {
	BankID string `query:"bank_id"` // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	Report string `query:"report"`  // รายงานประเภท day | month | all
}

// ---------- Response DTO ----------

type PayableDTO struct {
	IDPayable  string    `json:"id_payable"`  // รหัสเจ้าหนี้
	BankID     string    `json:"bank_id"`     // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	Supplier   string    `json:"supplier"`    // ชื่อผู้จำหน่าย
	PurchaseNo string    `json:"purchase_no"` // เลขที่ใบสั่งซื้อ
	InvoiceNo  string    `json:"invoice_no"`  // เลขที่ใบแจ้งหนี้
	IssueDate  time.Time `json:"issue_date"`  // วันที่ออกใบแจ้งหนี้
	DueDate    time.Time `json:"due_date"`    // วันที่ครบกำหนดชำระ
	Amount     float64   `json:"amount"`      // จำนวนเงิน
	Balance    float64   `json:"balance"`     // ยอดคงเหลือ
}

type PayableSummaryDTO struct {
	TotalAmount  float64 `json:"total_amount"`  // ยอดรวมทั้งหมด
	TotalDue     float64 `json:"total_due"`     // ยอดคงค้าง
	OverdueCount int     `json:"overdue_count"` // จำนวนรายการเกินกำหนด
}
