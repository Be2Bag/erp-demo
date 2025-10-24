package models

import "time"

const CollectionPaymentTransaction = "payment_transactions"

type PaymentTransaction struct {
	ID            string    `json:"id" bson:"_id,omitempty"`                                // รหัสเอกสาร
	PayableID     string    `json:"payable_id,omitempty" bson:"payable_id,omitempty"`       // ลิงก์ไปที่ Payable (ใช้เมื่อเป็นการจ่ายเงิน)
	ReceivableID  string    `json:"receivable_id,omitempty" bson:"receivable_id,omitempty"` // ลิงก์ไปที่ Receivable (ใช้เมื่อเป็นการรับเงิน)
	PaymentDate   time.Time `json:"payment_date" bson:"payment_date"`                       // วันที่ชำระเงิน
	AmountPaid    float64   `json:"amount_paid" bson:"amount_paid"`                         // จำนวนเงินที่ชำระในครั้งนี้
	PaymentMethod string    `json:"payment_method" bson:"payment_method"`                   // ช่องทางการชำระ: cash, bank_transfer, cheque, etc.
	Reference     string    `json:"reference" bson:"reference"`                             // เลขอ้างอิง เช่น เลขเช็ค / รหัสโอน
	Note          string    `json:"note,omitempty" bson:"note,omitempty"`                   // หมายเหตุเพิ่มเติม
	CreatedBy     string    `json:"created_by" bson:"created_by"`                           // ผู้บันทึกธุรกรรม
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`                           // วันที่บันทึก
}
