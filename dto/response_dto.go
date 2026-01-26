package dto

type BaseResponse struct {
	Data       any    `json:"data"`
	MessageTH  string `json:"message_th" example:"สำเร็จ"`
	MessageEN  string `json:"message_en" example:"Success"`
	Status     string `json:"status" example:"success"`
	StatusCode int    `json:"status_code" example:"200"`
}
