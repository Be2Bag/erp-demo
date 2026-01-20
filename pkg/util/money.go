package util

import "math"

// Round2 ปัดเศษจำนวนเงินเป็น 2 ตำแหน่งทศนิยม (สตางค์)
// ใช้เพื่อป้องกัน floating-point precision error ในการคำนวณเงิน
func Round2(val float64) float64 {
	return math.Round(val*100) / 100
}

// IsZeroBalance ตรวจสอบว่ายอดเงินถือว่าเป็น 0 หรือไม่
// (น้อยกว่า 0.01 หรือ 1 สตางค์ ถือว่าเป็น 0)
func IsZeroBalance(val float64) bool {
	return val < 0.01
}

// IsPositiveAmount ตรวจสอบว่ามียอดเงินที่มากกว่าศูนย์หรือไม่
// (ต้องมากกว่าหรือเท่ากับ 0.01 หรือ 1 สตางค์)
func IsPositiveAmount(val float64) bool {
	return val >= 0.01
}
