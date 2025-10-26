package models

import "time"

const CollectionPayable = "payables" // เราเป็นเจ้าหนี้ (รอเก็บเงิน)

type Payable struct {
	IDPayable  string     `json:"id_payable" bson:"id_payable"`                     // รหัสเอกสาร
	BankID     string     `json:"bank_id" bson:"bank_id"`                           // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	Supplier   string     `json:"supplier" bson:"supplier"`                         // ชื่อผู้ขาย / เจ้าหนี้
	PurchaseNo string     `json:"purchase_no" bson:"purchase_no"`                   // เลขที่ใบสั่งซื้อ / ใบรับของ
	InvoiceNo  string     `json:"invoice_no" bson:"invoice_no"`                     // เลขที่ใบแจ้งหนี้จากผู้ขาย
	IssueDate  time.Time  `json:"issue_date" bson:"issue_date"`                     // วันที่ออกใบแจ้งหนี้
	DueDate    time.Time  `json:"due_date" bson:"due_date"`                         // วันครบกำหนดชำระ
	Amount     float64    `json:"amount" bson:"amount"`                             // จำนวนเงินทั้งหมดในใบแจ้งหนี้
	Balance    float64    `json:"balance" bson:"balance"`                           // ยอดคงเหลือที่ยังไม่ได้ชำระ
	Status     string     `json:"status" bson:"status"`                             // สถานะ: pending, paid, overdue, partial
	PaymentRef string     `json:"payment_ref" bson:"payment_ref"`                   // เลขอ้างอิงการชำระเงิน (เช่น เลขที่เช็ค, โอนเงิน, เอกสารภายใน)
	Note       string     `json:"note" bson:"note"`                                 // หมายเหตุเพิ่มเติม
	CreatedBy  string     `json:"created_by" bson:"created_by"`                     // ผู้สร้างข้อมูล
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`                     // วันที่สร้างข้อมูล
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`                     // วันที่แก้ไขล่าสุด
	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // วันที่ลบข้อมูล (ถ้ามี)
}

// pending → ยังไม่จ่ายเลย
// partial → จ่ายแล้วบางส่วน
// paid → จ่ายครบแล้ว
// overdue → เลยกำหนดจ่ายแล้วยังไม่จ่ายครบ
