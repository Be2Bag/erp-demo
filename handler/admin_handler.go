package handler

import (
	"strings"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminHandler struct {
	svc ports.AdminService
}

func NewAdminHandler(s ports.AdminService) *AdminHandler {
	return &AdminHandler{svc: s}
}

func (h *AdminHandler) AdminRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	admin := versionOne.Group("admin")

	admin.Put("/update-status-user", h.UpdateStatusUser)

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

	err := h.svc.UpdateUserStatus(c.Context(), req)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		if strings.Contains(err.Error(), "user status is not pending") {
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
