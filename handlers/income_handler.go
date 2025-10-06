package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type InComeHandler struct {
	svc ports.InComeService
	mdw *middleware.Middleware
}

func NewInComeHandler(s ports.InComeService, mdw *middleware.Middleware) *InComeHandler {
	return &InComeHandler{svc: s, mdw: mdw}
}

func (h *InComeHandler) InComeRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	inCome := versionOne.Group("in-come")

	inCome.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateInCome)
	inCome.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListInComes)
	inCome.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetInComeByID)
	inCome.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateInComeByID)
	inCome.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteInComeByID)

}

// @Summary Create In Come
// @Description Create a new income
// @Tags Income
// @Accept json
// @Produce json
// @Param request body dto.CreateIncomeDTO true "Create Income"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/in-come/create [post]
func (h *InComeHandler) CreateInCome(c *fiber.Ctx) error {

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

	var inCome dto.CreateIncomeDTO
	if err := c.BodyParser(&inCome); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateInCome(c.Context(), inCome, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create income" + err.Error(),
			MessageTH:  "สร้างรายได้ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Income created successfully",
		MessageTH:  "สร้างรายได้เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List In Comes
// @Description List all incomes with pagination and optional search
// @Tags Income
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10) maximum(100)
// @Param search query string false "Search term"
// @Param sortBy query string false "Field to sort by" default("created_at")
// @Param sortOrder query string false "Sort order" Enums(asc, desc) default("desc")
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/in-come/list [get]

func (h *InComeHandler) ListInComes(c *fiber.Ctx) error {

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

	var req dto.RequestListIncome
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

	list, err := h.svc.ListInComes(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list income" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการรายได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Success",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       list,
	})
}

// @Summary Get Income by ID
// @Description Get a single income by its ID
// @Tags Income
// @Accept json
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/in-come/{id} [get]

func (h *InComeHandler) GetInComeByID(c *fiber.Ctx) error {
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
	inComeID := c.Params("id")
	item, err := h.svc.GetIncomeByID(c.Context(), inComeID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get income",
			MessageTH:  "ไม่สามารถดึงรายได้",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Income not found",
			MessageTH:  "ไม่พบรายได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Success",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       item,
	})
}

// @Summary Update Income by ID
// @Description Update an existing income by its ID
// @Tags Income
// @Accept json
// @Produce json
// @Param id path string true "Income ID"
// @Param request body dto.UpdateIncomeDTO true "Update Income"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/in-come/{id} [put]
func (h *InComeHandler) UpdateInComeByID(c *fiber.Ctx) error {
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
	inComeID := c.Params("id")
	var body dto.UpdateIncomeDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdateInComeByID(c.Context(), inComeID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Income not found"
			MsgTH = "ไม่พบรายได้"
		}

		statusCode = fiber.StatusInternalServerError
		MsgEN = "Failed to update income" + errOnUpdate.Error()
		MsgTH = "อัปเดตไม่สำเร็จ"
	}

	return c.JSON(dto.BaseResponse{
		StatusCode: statusCode,
		MessageEN:  MsgEN,
		MessageTH:  MsgTH,
		Status:     "success",
		Data:       nil,
	})

}

// @Summary Delete Income by ID
// @Description Soft delete an income by its ID
// @Tags Income
// @Accept json
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/in-come/{id} [delete]
func (h *InComeHandler) DeleteInComeByID(c *fiber.Ctx) error {
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
	inComeID := c.Params("id")
	err = h.svc.DeleteInComeByInComeID(c.Context(), inComeID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete income",
			MessageTH:  "ลบไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Deleted",
		MessageTH:  "ลบแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
