package dto

// Request

type RequestLogin struct {
	Email    string `json:"email"`    // อีเมลของผู้ใช้
	Password string `json:"password"` // รหัสผ่าน
}

type JWTClaims struct {
	UserID       string `json:"user_id"`       // User ID
	EmployeeCode string `json:"employee_code"` // Employee code
	Role         string `json:"role"`          // User role
	TitleTH      string `json:"title_th"`      // Title in Thai
	FirstNameTH  string `json:"first_name_th"` // First name in Thai
	LastNameTH   string `json:"last_name_th"`  // Last name in Thai
	Avatar       string `json:"avatar"`        // User avatar URL
	Status       string `json:"status"`        // User status (e.g., active, inactive)
}

type RequestResetPassword struct {
	Email       string `json:"email"`        // อีเมลของผู้ใช้
	RedirectURL string `json:"redirect_url"` // URL ที่จะเปลี่ยนเส้นทางหลังจากรีเซ็ตรหัสผ่าน
}

type RequestConfirmResetPassword struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
