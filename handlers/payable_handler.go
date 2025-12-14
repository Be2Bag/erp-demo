package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PayableHandler struct {
	svc ports.PayableService
	mdw *middleware.Middleware
}

func NewPayableHandler(s ports.PayableService, mdw *middleware.Middleware) *PayableHandler {
	return &PayableHandler{svc: s, mdw: mdw}
}

func (h *PayableHandler) PayableRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	payable := versionOne.Group("payable")

	payable.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreatePayable)
	payable.Post("/record-payment", h.mdw.AuthCookieMiddleware(), h.RecordPayment)
	payable.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListPayables)
	payable.Get("/summary", h.mdw.AuthCookieMiddleware(), h.SummaryPayableByFilter)
	payable.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetPayableByID)
	payable.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdatePayableByID)
	payable.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeletePayableByID)

}

// @Summary Create a new payable
// @Description Create a new payable record
// @Tags Payables
// @Accept json
// @Produce json
// @Param payable body dto.CreatePayableDTO true "Payable to create"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/create [post]
func (h *PayableHandler) CreatePayable(c *fiber.Ctx) error {

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

	var payable dto.CreatePayableDTO
	if err := c.BodyParser(&payable); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreatePayable(c.Context(), payable, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create payable" + err.Error(),
			MessageTH:  "สร้างเจ้าหนี้ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Payable created successfully",
		MessageTH:  "สร้างเจ้าหนี้เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List payables with pagination and filtering
// @Description Retrieve a paginated list of payables with optional search and filtering
// @Tags Payables
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field" Enums(created_at, updated_at) default(created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param status query string false "Filter by status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param bank_id query string false "Bank ID"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/list [get]
func (h *PayableHandler) ListPayables(c *fiber.Ctx) error {

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

	var req dto.RequestListPayable
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

	list, err := h.svc.ListPayables(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.Status, req.StartDate, req.EndDate, req.BankID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list payables" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการเจ้าหนี้",
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

// @Summary Get payable by ID
// @Description Retrieve a payable record by its ID
// @Tags Payables
// @Accept json
// @Produce json
// @Param id path string true "Payable ID"
// @Success 200 {object} dto.BaseResponse{data=dto.PayableDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/{id} [get]
func (h *PayableHandler) GetPayableByID(c *fiber.Ctx) error {
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
	payableID := c.Params("id")
	item, err := h.svc.GetPayableByID(c.Context(), payableID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get payable",
			MessageTH:  "ไม่สามารถดึงเจ้าหนี้",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Payable not found",
			MessageTH:  "ไม่พบเจ้าหนี้",
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

// @Summary Update payable by ID
// @Description Update a payable record by its ID
// @Tags Payables
// @Accept json
// @Produce json
// @Param id path string true "Payable ID"
// @Param payable body dto.UpdatePayableDTO true "Payable data to update"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/{id} [put]
func (h *PayableHandler) UpdatePayableByID(c *fiber.Ctx) error {
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
	payableID := c.Params("id")
	var body dto.UpdatePayableDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdatePayableByID(c.Context(), payableID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Payable not found"
			MsgTH = "ไม่พบเจ้าหนี้"
		}

		statusCode = fiber.StatusInternalServerError
		MsgEN = "Failed to update payable" + errOnUpdate.Error()
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

// @Summary Delete payable by ID
// @Description Soft delete a payable record by its ID
// @Tags Payables
// @Accept json
// @Produce json
// @Param id path string true "Payable ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/{id} [delete]
func (h *PayableHandler) DeletePayableByID(c *fiber.Ctx) error {
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
	payableID := c.Params("id")
	err = h.svc.DeletePayableByID(c.Context(), payableID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete payable",
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

// @Summary Get payable summary by filter
// @Description Retrieve a summary of payables based on filters
// @Tags Payables
// @Accept json
// @Produce json
// @Param bank_id query string false "Bank ID"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param report query string true "Report type" Enums(day, month, all)
// @Success 200 {object} dto.BaseResponse{data=dto.PayableSummaryDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/summary [get]
func (h *PayableHandler) SummaryPayableByFilter(c *fiber.Ctx) error {
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
	var req dto.RequestSummaryPayable
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ค้นหาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	summary, err := h.svc.SummaryPayableByFilter(c.Context(), claims, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get payable summary",
			MessageTH:  "ไม่สามารถดึงข้อมูลสรุปเจ้าหนี้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Success",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       summary,
	})
}

// @Summary Record a payment for a payable
// @Description Record a payment transaction and update the payable balance/status accordingly
// @Tags Payables
// @Accept json
// @Produce json
// @Param payment body dto.RecordPaymentDTO true "Payment details"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/payable/record-payment [post]
func (h *PayableHandler) RecordPayment(c *fiber.Ctx) error {
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
	var req dto.RecordPaymentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	err = h.svc.RecordPayment(c.Context(), req, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to record payment" + err.Error(),
			MessageTH:  "บันทึกการชำระเงินไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Payment recorded successfully",
		MessageTH:  "บันทึกการชำระเงินสำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
