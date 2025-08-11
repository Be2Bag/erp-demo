package dto

// JWTClaims represents the structure of JWT payload used in the system.
type JWTClaims struct {
	UserID       string `json:"UserID"`
	EmployeeCode string `json:"EmployeeCode"`
	Role         string `json:"Role"`
	TitleTH      string `json:"TitleTH"`
	FirstNameTH  string `json:"FirstNameTH"`
	LastNameTH   string `json:"LastNameTH"`
	Avatar       string `json:"Avatar"`
	Status       string `json:"Status"`
	Exp          int64  `json:"exp,omitempty"`
}
