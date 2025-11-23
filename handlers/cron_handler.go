package handlers

import (
	"github.com/Be2Bag/erp-demo/cron"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/gofiber/fiber/v2"
)

type CronHandler struct {
	statusChecker *cron.StatusChecker
	middleware    *middleware.Middleware
}

func NewCronHandler(statusChecker *cron.StatusChecker, middleware *middleware.Middleware) *CronHandler {
	return &CronHandler{
		statusChecker: statusChecker,
		middleware:    middleware,
	}
}

// RunStatusCheck
// @Summary รัน cronjob ตรวจสอบสถานะ Payable/Receivable ทันที
// @Description รันการตรวจสอบและอัปเดตสถานะของ Payable และ Receivable ทันที (ไม่ต้องรอถึงเวลา 00:00 น.)
// @Tags Cron
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "สำเร็จ"
// @Failure 500 {object} map[string]interface{} "เกิดข้อผิดพลาด"
// @Router /cron/status-check [post]
func (h *CronHandler) RunStatusCheck(c *fiber.Ctx) error {
	if err := h.statusChecker.RunNow(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "รัน cronjob ไม่สำเร็จ",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "รัน cronjob สำเร็จ - ตรวจสอบ log เพื่อดูรายละเอียด",
	})
}

// CronRoutes กำหนด routes สำหรับ Cron
func (h *CronHandler) CronRoutes(r fiber.Router) {
	cronGroup := r.Group("/cron")

	// ต้อง authenticate ก่อนเรียกใช้ (ใช้เฉพาะ Admin)
	cronGroup.Post("/status-check", h.middleware.AuthCookieMiddleware(), h.RunStatusCheck)
}
