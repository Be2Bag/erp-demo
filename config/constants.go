package config

// DefaultBankAccountIDs - รหัสบัญชีธนาคารที่ใช้เป็นค่าเริ่มต้นในระบบ
// ห้ามลบบัญชีเหล่านี้เพราะใช้สำหรับ auto-sync
var DefaultBankAccountIDs = struct {
	CompanyBank string // บัญชีบริษัท (ธนาคารกสิกรไทย)
	ShopBank    string // บัญชีร้านค้า (ธนาคารกรุงไทย)
}{
	CompanyBank: "307961ea-eb4f-4127-8e83-6eba0b8abbaf",
	ShopBank:    "d2791d28-4427-4857-9d38-1492110aaba3",
}

// DefaultTransactionCategoryIDs - รหัสหมวดหมู่รายการที่ใช้เป็นค่าเริ่มต้นในระบบ
// ห้ามลบหมวดหมู่เหล่านี้เพราะใช้สำหรับ auto-sync
var DefaultTransactionCategoryIDs = struct {
	CompanyExpense string // รายจ่ายบริษัท
	CompanyIncome  string // รายได้ใบงานบริษัท
}{
	CompanyExpense: "70e128e9-aef3-4699-83a1-7d34e1a1f342",
	CompanyIncome:  "ee1bbffd-aee7-4f1b-8c92-582d9449b0fd",
}

// IsProtectedBankAccount ตรวจสอบว่าเป็นบัญชีที่ห้ามลบหรือไม่
func IsProtectedBankAccount(bankID string) bool {
	return bankID == DefaultBankAccountIDs.CompanyBank ||
		bankID == DefaultBankAccountIDs.ShopBank
}

// IsProtectedTransactionCategory ตรวจสอบว่าเป็นหมวดหมู่ที่ห้ามลบหรือไม่
func IsProtectedTransactionCategory(categoryID string) bool {
	return categoryID == DefaultTransactionCategoryIDs.CompanyExpense ||
		categoryID == DefaultTransactionCategoryIDs.CompanyIncome
}

// PermissionConfig - กำหนด UserIDs ที่มีสิทธิ์พิเศษนอกเหนือจาก admin
var PermissionConfig = struct {
	VerifySignJob []string // UserIDs ที่สามารถ verify sign job ได้
}{
	VerifySignJob: []string{
		"ace8df6b-a785-491f-b33e-37085e7e8229",
	},
}

// HasPermission ตรวจสอบว่า userID มีสิทธิ์ในรายการหรือไม่
func HasPermission(allowedUserIDs []string, userID string) bool {
	for _, id := range allowedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}
