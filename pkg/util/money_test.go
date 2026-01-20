package util

import (
	"testing"
)

func TestRound2(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		// กรณีปกติ
		{"simple round", 100.5, 100.5},
		{"no change needed", 100.00, 100.00},

		// กรณี floating-point precision error
		{"0.1 + 0.2 = 0.3", 0.1 + 0.2, 0.30},
		{"10.0 - 9.9 = 0.1", 10.0 - 9.9, 0.10},
		{"100.10 - 100.00 = 0.10", 100.10 - 100.00, 0.10},

		// กรณี VAT 7%
		{"VAT: 1000 * 0.07 = 70", 1000 * 0.07, 70.00},
		{"VAT: 999.99 * 0.07", 999.99 * 0.07, 70.00},

		// กรณีปัดเศษ
		{"round up", 100.125, 100.13},
		{"round down", 100.124, 100.12},
		{"round half", 100.115, 100.12},

		// กรณีค่าเล็กมาก
		{"very small", 0.001, 0.00},
		{"one satang", 0.01, 0.01},
		{"near zero", 0.009, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round2(tt.input)
			if result != tt.expected {
				t.Errorf("Round2(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsZeroBalance(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"zero", 0, true},
		{"very small positive", 0.001, true},
		{"just under threshold", 0.009, true},
		{"at threshold", 0.01, false},
		{"above threshold", 0.02, false},
		{"negative", -0.01, true},
		{"floating point near zero", 10.0 - 9.9 - 0.1, true}, // ~1e-16
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsZeroBalance(tt.input)
			if result != tt.expected {
				t.Errorf("IsZeroBalance(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPositiveAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"zero", 0, false},
		{"very small positive", 0.001, false},
		{"just under threshold", 0.009, false},
		{"at threshold (1 satang)", 0.01, true},
		{"above threshold", 0.02, true},
		{"100 baht", 100.00, true},
		{"negative", -0.01, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPositiveAmount(tt.input)
			if result != tt.expected {
				t.Errorf("IsPositiveAmount(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ทดสอบ scenario จริง: การคำนวณยอดชำระ
func TestPaymentCalculation(t *testing.T) {
	// สถานการณ์: ยอด 1000 บาท จ่าย 999.90 บาท
	amount := 1000.00
	paid := 999.90

	balance := Round2(amount - paid)
	if balance != 0.10 {
		t.Errorf("Balance should be 0.10, got %v", balance)
	}
	if !IsPositiveAmount(balance) {
		t.Errorf("Balance 0.10 should be positive amount")
	}

	// จ่ายอีก 0.10 บาท
	balance = Round2(balance - 0.10)
	if balance != 0.00 {
		t.Errorf("Final balance should be 0.00, got %v", balance)
	}
	if !IsZeroBalance(balance) {
		t.Errorf("Balance 0.00 should be zero balance")
	}
}

// ทดสอบ scenario: floating-point precision error ที่เกิดขึ้นจริง
func TestFloatingPointScenario(t *testing.T) {
	// สถานการณ์ที่เคยเกิดบัญชี: 0.1 + 0.2 ไม่เท่ากับ 0.3
	result := 0.1 + 0.2
	if result == 0.3 {
		t.Log("Direct comparison works (unexpected in float64)")
	} else {
		t.Logf("Floating point issue: 0.1 + 0.2 = %.17f", result)
	}

	// แต่หลังจากปัดเศษจะถูกต้อง
	rounded := Round2(result)
	if rounded != 0.30 {
		t.Errorf("After Round2, should be 0.30, got %v", rounded)
	}
}
