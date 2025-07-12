package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/Be2Bag/erp-demo/dto"
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
