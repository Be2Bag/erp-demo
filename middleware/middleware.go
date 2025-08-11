package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func ExecWithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := ExecWithTimeout(c.UserContext(), timeout, func(ctx context.Context) error {
			c.SetUserContext(ctx)
			return c.Next()
		})

		if err != nil {
			// timeout
			if errors.Is(err, context.DeadlineExceeded) {
				return c.Status(fiber.StatusGatewayTimeout).JSON(dto.BaseResponse{
					StatusCode: fiber.StatusGatewayTimeout,
					MessageEN:  "Request timed out",
					MessageTH:  "หมดเวลาการร้องขอ",
					Status:     "error",
					Data:       nil,
				})
			}
			// ถ้าเป็น error ของ Fiber (เช่น 404)
			if e, ok := err.(*fiber.Error); ok {
				// ปรับ MessageTH ตามความเหมาะสม
				return c.Status(e.Code).JSON(dto.BaseResponse{
					StatusCode: e.Code,
					MessageEN:  e.Message,
					MessageTH:  "ไม่พบข้อมูล",
					Status:     "error",
					Data:       nil,
				})
			}
			// error อื่นๆ
			return err
		}

		return nil
	}
}

// Add struct to hold dependencies and a constructor.
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

		claims, err := util.VerifyJWTToken(token, m.JWT.SecretKey)
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
