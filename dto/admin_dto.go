package dto

type RequestUpdateUserStatus struct {
	UserID string `json:"user_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	Status string `json:"status" validate:"required,oneof=pending approved rejected" example:"approved"`
}

type RequestUpdateUserRole struct {
	UserID string `json:"user_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	Role   string `json:"role" validate:"required,oneof=admin user" example:"admin"`
	Note   string `json:"note"  example:"ทดสอบปรับบทบาท"`
}

type RequestUpdateUserPosition struct {
	UserID       string `json:"user_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	PositionID   string `json:"position_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	DepartmentID string `json:"department_id" validate:"required" example:"50f7a957-8c2c-4a76-88ed-7c247471f28f"`
	Note         string `json:"note"  example:"ทดสอบปรับตำแหน่ง"`
}
