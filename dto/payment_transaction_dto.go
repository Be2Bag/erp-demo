package dto

type RecordPaymentDTO struct { // โครงสร้างข้อมูลสำหรับบันทึกการชำระเงิน.
	PayableID     string  `json:"payable_id"`     // ไอดีของเจ้าหนี้/รายการที่ต้องชำระ
	BankID        string  `json:"bank_id"`        // รหัสบัญชีธนาคารที่รับชำระ
	PaymentDate   string  `json:"payment_date"`   // วันที่ชำระเงิน (เช่น 2025-10-26)
	PaymentMethod string  `json:"payment_method"` // วิธีการชำระเงิน (เช่น โอนเงิน/เงินสด/เช็ค)
	PaymentRef    string  `json:"payment_ref"`    // เลขอ้างอิงการชำระเงิน/หลักฐาน
	Note          string  `json:"note"`           // หมายเหตุเพิ่มเติม
	Amount        float64 `json:"amount"`         // จำนวนเงินที่ชำระ
} // ปิดโครงสร้างข้อมูล
