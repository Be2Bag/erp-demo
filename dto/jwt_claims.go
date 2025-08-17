package dto

// JWTClaims represents the structure of JWT payload used in the system.
type JWTClaims struct {
	UserID       string `json:"user_id"`
	EmployeeCode string `json:"employee_code"`
	Role         string `json:"role"`
	TitleTH      string `json:"title_th"`
	FirstNameTH  string `json:"first_name_th"`
	LastNameTH   string `json:"last_name_th"`
	Avatar       string `json:"avatar"`
	Status       string `json:"status"`
	Exp          int64  `json:"exp,omitempty"`
}
