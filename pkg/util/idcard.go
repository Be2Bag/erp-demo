package util

import (
	"strings"
)

func MaskIDCard(id string) string {
	runes := []rune(id)
	l := len(runes)
	if l <= 6 {
		return id
	}
	maskCount := l - 6
	return strings.Repeat("X", maskCount) + string(runes[maskCount:])
}

func ValidateThaiID(id string) bool {
	if len(id) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		d := int(id[i] - '0')
		if d < 0 || d > 9 {
			return false
		}
		sum += d * (13 - i)
	}
	check := (11 - (sum % 11)) % 10
	last := id[12] - '0'
	return int(last) == check
}
