package dto

import "time"

// ---------- Request DTO ----------

type CreatePayableDTO struct {
	BankAccount BankAccountPayable `json:"bank_account" bson:"bank_account"` // ข้อมูลบัญชีธนาคารที่เกี่ยวข้อง
	Supplier    string             `json:"supplier"`                         // ชื่อผู้จำหน่าย
	BankID      string             `json:"bank_id"`                          // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	PurchaseNo  string             `json:"purchase_no"`                      // เลขที่ใบสั่งซื้อ
	InvoiceNo   string             `json:"invoice_no"`                       // เลขที่ใบแจ้งหนี้
	IssueDate   string             `json:"issue_date"`                       // วันที่ออกใบแจ้งหนี้
	DueDate     string             `json:"due_date"`                         // วันที่ครบกำหนดชำระ
	PaymentRef  string             `json:"payment_ref"`                      // เลขที่อ้างอิงการชำระเงิน
	Phone       string             `json:"phone" bson:"phone"`               // เบอร์โทรศัพท์ผู้ขาย / เจ้าหนี้
	Address     string             `json:"address" bson:"address"`           // ที่อยู่ผู้ขาย / เจ้าหนี้
	Note        string             `json:"note,omitempty"`                   // หมายเหตุ
	Items       []ReceiptItemDTO   `json:"items,omitempty"`                  // รายการสินค้า/บริการ
	Amount      float64            `json:"amount"`                           // จำนวนเงิน
	Balance     float64            `json:"balance"`                          // ยอดคงเหลือ
}

type UpdatePayableDTO struct {
	BankAccount BankAccountPayable `json:"bank_account" bson:"bank_account"` // ข้อมูลบัญชีธนาคารที่เกี่ยวข้อง
	Supplier    string             `json:"supplier,omitempty"`               // ชื่อผู้จำหน่าย
	BankID      string             `json:"bank_id,omitempty"`                // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	PurchaseNo  string             `json:"purchase_no,omitempty"`            // เลขที่ใบสั่งซื้อ
	InvoiceNo   string             `json:"invoice_no,omitempty"`             // เลขที่ใบแจ้งหนี้
	IssueDate   string             `json:"issue_date,omitempty"`             // วันที่ออกใบแจ้งหนี้
	DueDate     string             `json:"due_date,omitempty"`               // วันที่ครบกำหนดชำระ
	Status      string             `json:"status,omitempty"`                 // สถานะ
	PaymentRef  string             `json:"payment_ref,omitempty"`            // เลขที่อ้างอิงการชำระเงิน
	Phone       string             `json:"phone" bson:"phone"`               // เบอร์โทรศัพท์ผู้ขาย / เจ้าหนี้
	Address     string             `json:"address" bson:"address"`           // ที่อยู่ผู้ขาย / เจ้าหนี้
	Note        string             `json:"note,omitempty"`                   // หมายเหตุ
	Items       []ReceiptItemDTO   `json:"items,omitempty"`                  // รายการสินค้า/บริการ
	Amount      float64            `json:"amount,omitempty"`                 // จำนวนเงิน
	Balance     float64            `json:"balance,omitempty"`                // ยอดคงเหลือ
}

type RequestListPayable struct {
	Search    string `query:"search"`     // คำค้นหา
	SortBy    string `query:"sort_by"`    // เรียงตาม
	SortOrder string `query:"sort_order"` // ลำดับการเรียง
	Status    string `query:"status"`     // สถานะ
	Supplier  string `query:"supplier"`   // ผู้จำหน่าย
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	BankID    string `query:"bank_id"` // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	Page      int    `query:"page"`    // หน้า
	Limit     int    `query:"limit"`   // จำนวนต่อหน้า
}

type RequestSummaryPayable struct {
	BankID    string `query:"bank_id"` // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// ---------- Response DTO ----------

type PayableDTO struct {
	IssueDate    time.Time               `json:"issue_date"`                       // วันที่ออกใบแจ้งหนี้
	DueDate      time.Time               `json:"due_date"`                         // วันที่ครบกำหนดชำระ
	BankAccount  BankAccountPayable      `json:"bank_account" bson:"bank_account"` // ข้อมูลบัญชีธนาคารที่เกี่ยวข้อง
	IDPayable    string                  `json:"id_payable"`                       // รหัสเจ้าหนี้
	BankID       string                  `json:"bank_id"`                          // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	BankName     string                  `json:"bank_name"`                        // ชื่อบัญชีธนาคาร
	Supplier     string                  `json:"supplier"`                         // ชื่อผู้จำหน่าย
	PurchaseNo   string                  `json:"purchase_no"`                      // เลขที่ใบสั่งซื้อ
	InvoiceNo    string                  `json:"invoice_no"`                       // เลขที่ใบแจ้งหนี้
	Status       string                  `json:"status"`                           // สถานะ
	PaymentRef   string                  `json:"payment_ref"`                      // เลขที่อ้างอิงการชำระเงิน
	Phone        string                  `json:"phone" bson:"phone"`               // เบอร์โทรศัพท์ผู้ขาย / เจ้าหนี้
	Address      string                  `json:"address" bson:"address"`           // ที่อยู่ผู้ขาย / เจ้าหนี้
	Note         string                  `json:"note"`                             // หมายเหตุ
	Items        []ReceiptItemDTO        `json:"items,omitempty"`                  // รายการสินค้า/บริการ
	Transactions []PaymentTransactionDTO `json:"transactions"`                     // รายการชำระเงิน
	Amount       float64                 `json:"amount"`                           // จำนวนเงิน
	Balance      float64                 `json:"balance"`                          // ยอดคงเหลือ
}

type PaymentTransactionDTO struct {
	PaymentDate     time.Time `json:"payment_date"`     // วันที่รับ/จ่ายเงิน
	CreatedAt       time.Time `json:"created_at"`       // วันที่บันทึก
	UpdatedAt       time.Time `json:"updated_at"`       // วันที่แก้ไขล่าสุด
	IDTransaction   string    `json:"id_transaction"`   // รหัสเอกสาร (เช่น PAY-2024-001-001)
	BankID          string    `json:"bank_id"`          // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	RefInvoiceNo    string    `json:"ref_invoice_no"`   // อ้างอิงใบแจ้งหนี้ (InvoiceNo)
	TransactionType string    `json:"transaction_type"` // ประเภท: receivable (ลูกหนี้) หรือ payable (เจ้าหนี้)
	PaymentMethod   string    `json:"payment_method"`   // วิธีการชำระเงิน (cash, transfer, cheque)
	PaymentRef      string    `json:"payment_ref"`      // หมายเลขอ้างอิง (เลขสลิป, เช็ค, ใบเสร็จ)
	Note            string    `json:"note"`             // หมายเหตุเพิ่มเติม
	CreatedBy       string    `json:"created_by"`       // ผู้บันทึก
	Amount          float64   `json:"amount"`           // จำนวนเงินที่จ่าย/รับในครั้งนี้
}

type PayableSummaryDTO struct {
	TotalAmount  float64 `json:"total_amount"`  // ยอดรวมทั้งหมด
	TotalDue     float64 `json:"total_due"`     // ยอดคงค้าง
	OverdueCount int     `json:"overdue_count"` // จำนวนรายการเกินกำหนด
}

type BankAccountPayable struct {
	BankName    string `bson:"bank_name" json:"bank_name"`
	AccountNo   string `bson:"account_no" json:"account_no"`
	AccountName string `bson:"account_name" json:"account_name"`
}
