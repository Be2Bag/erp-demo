package handlers

import (
	"strings"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminHandler struct {
	svc ports.AdminService
	mdw *middleware.Middleware
}

func NewAdminHandler(s ports.AdminService, m *middleware.Middleware) *AdminHandler {
	return &AdminHandler{svc: s, mdw: m}
}

func (h *AdminHandler) AdminRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	admin := versionOne.Group("admin")

	admin.Put("/update-status-user", h.mdw.AuthCookieMiddleware(), h.UpdateStatusUser)
	admin.Put("/update-role-user", h.mdw.AuthCookieMiddleware(), h.UpdateRoleUser)
	admin.Put("/update-position-user", h.mdw.AuthCookieMiddleware(), h.UpdatePositionUser)

}

// @Summary Update User Status
// @Description สำหรับจัดการอนุมัติผู้ใช้ pending approved rejected
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.RequestUpdateUserStatus true "Request Update User Status"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/admin/update-status-user [put]
func (h *AdminHandler) UpdateStatusUser(c *fiber.Ctx) error {

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
	if claims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusForbidden,
			MessageEN:  "Forbidden",
			MessageTH:  "ห้ามเข้าถึง",
			Status:     "error",
			Data:       nil,
		})
	}

	var req dto.RequestUpdateUserStatus
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "ข้อมูลคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	errOnUpdate := h.svc.UpdateUserStatus(c.Context(), req)
	if errOnUpdate != nil {

		if errOnUpdate == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		if strings.Contains(errOnUpdate.Error(), "user status is not pending") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "User status is not pending",
				MessageTH:  "สถานะผู้ใช้ไม่ใช่ pending",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update user status",
			MessageTH:  "ไม่สามารถอัปเดตสถานะผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User status updated successfully",
		MessageTH:  "อัปเดตสถานะผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Update User Role
// @Description สำหรับจัดการอัปเดตบทบาทผู้ใช้
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.RequestUpdateUserRole true "Request Update User Role"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/admin/update-role-user [put]
func (h *AdminHandler) UpdateRoleUser(c *fiber.Ctx) error {

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
	if claims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusForbidden,
			MessageEN:  "Forbidden",
			MessageTH:  "ห้ามเข้าถึง",
			Status:     "error",
			Data:       nil,
		})
	}

	var req dto.RequestUpdateUserRole
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "ข้อมูลคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	errOnUpdate := h.svc.UpdateUserRole(c.Context(), req)
	if errOnUpdate != nil {

		if errOnUpdate == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update user role",
			MessageTH:  "ไม่สามารถอัปเดตบทบาทผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User role updated successfully",
		MessageTH:  "อัปเดตบทบาทผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Update User Position
// @Description สำหรับจัดการอัปเดตตำแหน่งผู้ใช้
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body dto.RequestUpdateUserPosition true "Request Update User Position"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/admin/update-position-user [put]
func (h *AdminHandler) UpdatePositionUser(c *fiber.Ctx) error {

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
	if claims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusForbidden,
			MessageEN:  "Forbidden",
			MessageTH:  "ห้ามเข้าถึง",
			Status:     "error",
			Data:       nil,
		})
	}

	var req dto.RequestUpdateUserPosition
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "ข้อมูลคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	errOnUpdate := h.svc.UpdateUserPosition(c.Context(), req)
	if errOnUpdate != nil {

		if errOnUpdate == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update user position",
			MessageTH:  "ไม่สามารถอัปเดตตำแหน่งผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User position updated successfully",
		MessageTH:  "อัปเดตตำแหน่งผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
