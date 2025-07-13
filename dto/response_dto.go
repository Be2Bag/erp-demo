package dto

type BaseResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	MessageTH  string `json:"message_th" example:"สำเร็จ"`
	MessageEN  string `json:"message_en" example:"Success"`
	Status     string `json:"status" example:"success"`
	Data       any    `json:"data"`
}
