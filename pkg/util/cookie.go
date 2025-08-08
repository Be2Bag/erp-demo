package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(claims map[string]interface{}, secretKey string, expiration time.Duration) (string, error) {
	tokenClaims := jwt.MapClaims{}
	for k, v := range claims {
		tokenClaims[k] = v
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
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})
}
