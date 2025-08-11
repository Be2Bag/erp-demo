package dto

// Request

type RequestLogin struct {
	Email    string `json:"email"`    // อีเมลของผู้ใช้
	Password string `json:"password"` // รหัสผ่าน
}

type RequestResetPassword struct {
	Email       string `json:"email" example:"example@mail.com"`                              // อีเมลของผู้ใช้
	RedirectURL string `json:"redirect_url" example:"https://erp-demo-frontend.onrender.com"` // URL ที่จะเปลี่ยนเส้นทางหลังจากรีเซ็ตรหัสผ่าน
}

type RequestConfirmResetPassword struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
