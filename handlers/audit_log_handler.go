package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type auditLogHandler struct {
	svc ports.AuditLogService
	mdw *middleware.Middleware
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(s ports.AuditLogService, mdw *middleware.Middleware) *auditLogHandler {
	return &auditLogHandler{svc: s, mdw: mdw}
}

// AuditLogRoutes registers all audit log routes
func (h *auditLogHandler) AuditLogRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	auditLog := versionOne.Group("audit-log")

	auditLog.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListAuditLogs)
	auditLog.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetAuditLogByID)
}

// @Summary      List audit logs
// @Description  Get a paginated list of audit logs (Admin only)
// @Tags         AuditLog
// @Accept       json
// @Produce      json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search in path, email, full_name"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param user_id query string false "Filter by user ID"
// @Param action query string false "Filter by action" Enums(CREATE, READ, UPDATE, DELETE)
// @Param resource query string false "Filter by resource type"
// @Param method query string false "Filter by HTTP method" Enums(GET, POST, PUT, DELETE)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 403 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/audit-log/list [get]
func (h *auditLogHandler) ListAuditLogs(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
			Data:       nil,
		})
	}

	// Only admin can view all audit logs
	if claims.Role != "admin" && claims.Role != "superadmin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusForbidden,
			MessageEN:  "Access denied. Admin only.",
			MessageTH:  "ไม่มีสิทธิ์เข้าถึง เฉพาะผู้ดูแลระบบเท่านั้น",
			Status:     "error",
			Data:       nil,
		})
	}

	var req dto.RequestListAuditLog
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	result, err := h.svc.GetLogs(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get audit logs: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูล Audit Log",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Success",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       result,
	})
}

// @Summary      Get audit log by ID
// @Description  Get a single audit log record by its ID (Admin only)
// @Tags         AuditLog
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Audit Log ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      403  {object}  dto.BaseResponse
// @Failure      404  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/audit-log/{id} [get]
func (h *auditLogHandler) GetAuditLogByID(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
			Data:       nil,
		})
	}

	// Only admin can view audit logs
	if claims.Role != "admin" && claims.Role != "superadmin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusForbidden,
			MessageEN:  "Access denied. Admin only.",
			MessageTH:  "ไม่มีสิทธิ์เข้าถึง เฉพาะผู้ดูแลระบบเท่านั้น",
			Status:     "error",
			Data:       nil,
		})
	}

	logID := c.Params("id")
	log, err := h.svc.GetLogByID(c.Context(), logID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get audit log",
			MessageTH:  "ไม่สามารถดึงข้อมูล Audit Log",
			Status:     "error",
			Data:       nil,
		})
	}

	if log == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Audit log not found",
			MessageTH:  "ไม่พบ Audit Log",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Success",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       log,
	})
}
