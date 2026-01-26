package models

import "time"

const CollectionReceivable = "receivables" // เราเป็นลูกหนี้ (ต้องจ่ายเงิน)

type Receivable struct {
	IssueDate    time.Time  `json:"issue_date" bson:"issue_date"`                     // วันที่ออกใบแจ้งหนี้
	DueDate      time.Time  `json:"due_date" bson:"due_date"`                         // วันครบกำหนดชำระ
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`                     // วันที่สร้างข้อมูล
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at"`                     // วันที่แก้ไขล่าสุด
	DeletedAt    *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // วันที่ลบข้อมูล (ถ้ามี)
	IDReceivable string     `json:"id_receivable" bson:"id_receivable"`               // รหัสเอกสาร
	BankID       string     `json:"bank_id" bson:"bank_id"`                           // รหัสบัญชีธนาคารที่เกี่ยวข้อง
	Customer     string     `json:"customer" bson:"customer"`                         // ลูกค้า
	InvoiceNo    string     `json:"invoice_no" bson:"invoice_no"`                     // เลขที่ใบแจ้งหนี้
	Status       string     `json:"status" bson:"status"`                             // สถานะ: pending, paid, overdue, partial
	Phone        string     `json:"phone"`                                            // เบอร์โทรศัพท์ของพนักงาน
	Address      string     `json:"address"`                                          // ที่อยู่ของพนักงาน
	CreatedBy    string     `json:"created_by" bson:"created_by"`                     // ผู้สร้างข้อมูล
	Note         string     `json:"note" bson:"note"`                                 // หมายเหตุเพิ่มเติม
	JobID        string     `json:"job_id" bson:"job_id,omitempty"`                   // รหัสงาน
	Amount       float64    `json:"amount" bson:"amount"`                             // จำนวนเงินทั้งหมด
	Balance      float64    `json:"balance" bson:"balance"`                           // ยอดคงเหลือ
}
