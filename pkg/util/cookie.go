package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(claims dto.JWTClaims, secretKey string, expiration time.Duration) (string, error) {
	tokenClaims := jwt.MapClaims{
		"user_id":       claims.UserID,
		"email":         claims.Email,
		"employee_code": claims.EmployeeCode,
		"role":          claims.Role,
		"title_th":      claims.TitleTH,
		"first_name_th": claims.FirstNameTH,
		"last_name_th":  claims.LastNameTH,
		"avatar":        claims.Avatar,
		"status":        claims.Status,
	}
	if expiration > 0 {
		tokenClaims["exp"] = time.Now().Add(expiration).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(secretKey))
}

func VerifyJWTToken(tokenStr, secretKey string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func SetSessionCookie(c *fiber.Ctx, name, value string, duration time.Duration) {
	domain := ""
	secure := false

	// ถ้าไม่ใช่ localhost ให้ตั้ง domain และ secure
	if !strings.Contains(c.Hostname(), "localhost") {
		domain = ".rkp-media.com"
		secure = true
	}

	cookie := &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "None",
		Domain:   domain,
		Path:     "/",
	}
	c.Cookie(cookie)
}

func DeleteCookie(ctx *fiber.Ctx, name string) {
	domain := ""
	secure := false

	// ถ้าไม่ใช่ localhost ให้ตั้ง domain และ secure
	if !strings.Contains(ctx.Hostname(), "localhost") {
		domain = ".rkp-media.com"
		secure = true
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "None",
		Domain:   domain,
	})
}

func VerifyAndParseJWTClaims(tokenStr, secretKey string) (*dto.JWTClaims, error) {
	raw, err := VerifyJWTToken(tokenStr, secretKey)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var claims dto.JWTClaims
	if err := json.Unmarshal(b, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}
