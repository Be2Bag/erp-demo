package handler

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type DropDownHandler struct {
	svc ports.DropDownService
}

func NewDropDownHandler(s ports.DropDownService) *DropDownHandler {
	return &DropDownHandler{svc: s}
}

func (h *DropDownHandler) DropDownRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	dropdown := versionOne.Group("dropdown")
	dropdown.Get("/position", h.GetPosition)
	dropdown.Get("/department", h.GetDepartment)

}

// @Summary Get all positions
// @Description ใช้สำหรับดึงข้อมูลตำแหน่งงานทั้งหมด
// @Tags dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetPositions}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/position [get]
func (h *DropDownHandler) GetPosition(c *fiber.Ctx) error {

	positions, errOnGetPositions := h.svc.GetPositions()
	if errOnGetPositions != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงตำแหน่งงานได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get positions successfully",
		MessageTH:  "ดึงตำแหน่งงานสำเร็จ",
		Status:     "success",
		Data:       positions,
	})
}

// @Summary Get all departments
// @Description ใช้สำหรับดึงข้อมูลแผนกทั้งหมด
// @Tags dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetDepartments}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/department [get]
func (h *DropDownHandler) GetDepartment(c *fiber.Ctx) error {
	departments, errOnGetDepartments := h.svc.GetDepartments()
	if errOnGetDepartments != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงแผนกได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get departments successfully",
		MessageTH:  "ดึงแผนกสำเร็จ",
		Status:     "success",
		Data:       departments,
	})
}
