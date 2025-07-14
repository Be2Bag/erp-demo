package dto

type RequestUpdateUserStatus struct {
	UserID string `json:"user_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	Status string `json:"status" validate:"required,oneof=pending approved rejected cancelled" example:"approved"`
}
