package dto

type BaseSuccess201ResponseSwagger struct {
	Data       any    `json:"data"`
	MessageTH  string `json:"message_th" example:"สร้างผู้ใช้สำเร็จ"`
	MessageEN  string `json:"message_en" example:"User created successfully"`
	Status     string `json:"status" example:"success"`
	StatusCode int    `json:"status_code" example:"201"`
}

type BaseError400ResponseSwagger struct {
	Data       any    `json:"data"`
	MessageTH  string `json:"message_th" example:"ข้อมูลไม่ถูกต้อง กรุณาตรวจสอบ"`
	MessageEN  string `json:"message_en" example:"Invalid request"`
	Status     string `json:"status" example:"error"`
	StatusCode int    `json:"status_code" example:"400"`
}

type BaseError401ResponseSwagger struct {
	Data       any    `json:"data"`
	MessageTH  string `json:"message_th" example:"การเข้าถึงถูกปฏิเสธ"`
	MessageEN  string `json:"message_en" example:"Unauthorized"`
	Status     string `json:"status" example:"error"`
	StatusCode int    `json:"status_code" example:"401"`
}

type BaseError500ResponseSwagger struct {
	Data       any    `json:"data"`
	MessageTH  string `json:"message_th" example:"เกิดข้อผิดพลาดในระบบ"`
	MessageEN  string `json:"message_en" example:"Internal server error"`
	Status     string `json:"status" example:"error"`
	StatusCode int    `json:"status_code" example:"500"`
}

type BaseSuccessPaginationResponseSwagger struct {
	MessageTH  string     `json:"message_th" example:"สำเร็จ"`
	MessageEN  string     `json:"message_en" example:"Success"`
	Status     string     `json:"status" example:"success"`
	Data       Pagination `json:"data"`
	StatusCode int        `json:"status_code" example:"200"`
}
