package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PositionHandler struct {
	svc ports.PositionService
	mdw *middleware.Middleware
}

func NewPositionHandler(s ports.PositionService, mdw *middleware.Middleware) *PositionHandler {
	return &PositionHandler{svc: s, mdw: mdw}
}

func (h *PositionHandler) PositionRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	position := versionOne.Group("position")

	position.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetPositionList)
	position.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreatePosition)
	position.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetPositionByID)
	position.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdatePosition)
	position.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeletePosition)

}

// @Summary Create a new position
// @Description Create a new position
// @Tags Positions
// @Accept json
// @Produce json
// @Param position body dto.CreatePositionDTO true "Position data"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/position/create [post]
func (h *PositionHandler) CreatePosition(c *fiber.Ctx) error {
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

	var createPosition dto.CreatePositionDTO
	if err := c.BodyParser(&createPosition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreatePosition(c.Context(), createPosition, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create position" + err.Error(),
			MessageTH:  "สร้างตำแหน่งไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Position created successfully",
		MessageTH:  "สร้างตำแหน่งเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Update an existing position
// @Description Update an existing position
// @Tags Positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param position body dto.UpdatePositionDTO true "Position data"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/position/{id} [put]
func (h *PositionHandler) UpdatePosition(c *fiber.Ctx) error {
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

	var updatePosition dto.UpdatePositionDTO
	if err := c.BodyParser(&updatePosition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	positionID := c.Params("id")
	errOnuUpdate := h.svc.UpdatePosition(c.Context(), positionID, updatePosition, claims)
	if errOnuUpdate != nil {

		if errors.Is(errOnuUpdate, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "Position not found",
				MessageTH:  "ไม่พบตำแหน่ง",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update position" + errOnuUpdate.Error(),
			MessageTH:  "อัปเดตตำแหน่งไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Position updated successfully",
		MessageTH:  "อัปเดตตำแหน่งเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Delete a position by ID
// @Description Delete a position by ID
// @Tags Positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/position/{id} [delete]
func (h *PositionHandler) DeletePosition(c *fiber.Ctx) error {
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

	positionID := c.Params("id")
	err = h.svc.DeletePositionByID(c.Context(), positionID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete position" + err.Error(),
			MessageTH:  "ลบตำแหน่งไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Position deleted successfully",
		MessageTH:  "ลบตำแหน่งเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get a position by ID
// @Description Get a position by ID
// @Tags Positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 200 {object} dto.BaseResponse{data=dto.PositionDTO}
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/position/{id} [get]
func (h *PositionHandler) GetPositionByID(c *fiber.Ctx) error {
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

	positionID := c.Params("id")
	position, err := h.svc.GetPositionByID(c.Context(), positionID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get position" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลตำแหน่งได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if position == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Position not found",
			MessageTH:  "ไม่พบตำแหน่ง",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Position retrieved successfully",
		MessageTH:  "ดึงข้อมูลตำแหน่งเรียบร้อยแล้ว",
		Status:     "success",
		Data:       position,
	})
}

// @Summary Get a list of positions
// @Description Get a list of positions
// @Tags Positions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "ค้นหา position_name"
// @Param department_id query string false "Dropdown แผนก DPT001: แผนกออกแบบกราฟิก, DPT002: แผนกผลิต, DPT003: แผนกติดตั้ง, DPT004: แผนกบัญชี"
// @Param sort_by query string false "เรียงตาม created_at updated_at department_name position_name level"
// @Param sort_order query string false "เรียงลำดับ (asc เก่า→ใหม่ | desc ใหม่→เก่า (ค่าเริ่มต้น))"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/position/list [get]
func (h *PositionHandler) GetPositionList(c *fiber.Ctx) error {
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

	var req dto.RequestListPosition
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	positions, err := h.svc.GetPositionList(c.Context(), claims, req.Page, req.Limit, req.Search, req.Department, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve positions" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลตำแหน่งได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Positions retrieved successfully",
		MessageTH:  "ดึงข้อมูลตำแหน่งเรียบร้อยแล้ว",
		Status:     "success",
		Data:       positions,
	})
}
