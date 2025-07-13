package dto

type BaseSuccess201ResponseSwagger struct {
	StatusCode int    `json:"status_code" example:"201"`
	MessageTH  string `json:"message_th" example:"สร้างผู้ใช้สำเร็จ"`
	MessageEN  string `json:"message_en" example:"User created successfully"`
	Status     string `json:"status" example:"success"`
	Data       any    `json:"data"`
}

type BaseError400ResponseSwagger struct {
	StatusCode int    `json:"status_code" example:"400"`
	MessageTH  string `json:"message_th" example:"ข้อมูลไม่ถูกต้อง กรุณาตรวจสอบ"`
	MessageEN  string `json:"message_en" example:"Invalid request"`
	Status     string `json:"status" example:"error"`
	Data       any    `json:"data"`
}

type BaseError401ResponseSwagger struct {
	StatusCode int    `json:"status_code" example:"401"`
	MessageTH  string `json:"message_th" example:"การเข้าถึงถูกปฏิเสธ"`
	MessageEN  string `json:"message_en" example:"Unauthorized"`
	Status     string `json:"status" example:"error"`
	Data       any    `json:"data"`
}

type BaseError500ResponseSwagger struct {
	StatusCode int    `json:"status_code" example:"500"`
	MessageTH  string `json:"message_th" example:"เกิดข้อผิดพลาดในระบบ"`
	MessageEN  string `json:"message_en" example:"Internal server error"`
	Status     string `json:"status" example:"error"`
	Data       any    `json:"data"`
}

type BaseSuccessPaginationResponseSwagger struct {
	StatusCode int        `json:"status_code" example:"200"`
	MessageTH  string     `json:"message_th" example:"สำเร็จ"`
	MessageEN  string     `json:"message_en" example:"Success"`
	Status     string     `json:"status" example:"success"`
	Data       Pagination `json:"data"`
}
