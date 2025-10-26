package dto // แพ็กเกจ dto

type RecordReceiptDTO struct { // โครงสร้าง DTO สำหรับบันทึกรับชำระ
	ReceivableID  string  `json:"receivable_id"`  // รหัสรายการลูกหนี้ที่รับชำระ
	BankID        string  `json:"bank_id"`        // รหัสบัญชีธนาคารที่รับชำระ
	Amount        float64 `json:"amount"`         // จำนวนเงินที่รับชำระ
	PaymentDate   string  `json:"payment_date"`   // วันที่รับชำระ (รูปแบบ YYYY-MM-DD)
	PaymentMethod string  `json:"payment_method"` // วิธีชำระเงิน เช่น เงินสด โอน บัตร
	PaymentRef    string  `json:"payment_ref"`    // เลขอ้างอิงการชำระ เช่น สลิป/เลขที่ธุรกรรม
	Note          string  `json:"note"`           // หมายเหตุเพิ่มเติม
} // จบโครงสร้าง
