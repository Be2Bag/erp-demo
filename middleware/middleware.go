package middleware

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	JWT config.JWTConfig
}

func NewMiddleware(jwtCfg config.JWTConfig) *Middleware {
	return &Middleware{JWT: jwtCfg}
}

// Replace previous standalone AuthCookieMiddleware with method version.
func (m *Middleware) AuthCookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("auth_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "No auth token cookie found",
				MessageTH:  "ไม่พบคุกกี้ auth token",
				Status:     "error",
				Data:       nil,
			})
		}
		if m.JWT.SecretKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "JWT secret not configured",
				MessageTH:  "ไม่ได้ตั้งค่า JWT secret",
				Status:     "error",
				Data:       nil,
			})
		}

		claims, err := util.VerifyAndParseJWTClaims(token, m.JWT.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Invalid or expired token",
				MessageTH:  "โทเคนไม่ถูกต้องหรือหมดอายุ",
				Status:     "error",
				Data:       nil,
			})
		}

		c.Locals("auth_claims", claims)
		return c.Next()
	}
}

func GetClaims(c *fiber.Ctx) (*dto.JWTClaims, error) {
	claims, ok := c.Locals("auth_claims").(*dto.JWTClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}

func (m *Middleware) TimeoutMiddleware(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return c.Status(fiber.StatusRequestTimeout).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusRequestTimeout,
				MessageEN:  "Request timed out",
				MessageTH:  "คำขอใช้เวลานานเกินไป",
				Status:     "error",
				Data:       nil,
			})
		}
	}
}
