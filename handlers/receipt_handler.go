package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type receiptHandler struct {
	svc ports.ReceiptService
	mdw *middleware.Middleware
}

func NewReceiptHandler(s ports.ReceiptService, mdw *middleware.Middleware) *receiptHandler {
	return &receiptHandler{svc: s, mdw: mdw}
}

func (h *receiptHandler) ReceiptRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	receipt := versionOne.Group("receipt")

	receipt.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateReceipt)
	receipt.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListReceipts)
	receipt.Get("/summary", h.mdw.AuthCookieMiddleware(), h.SummaryReceiptByFilter)
	receipt.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetReceiptByID)
	receipt.Post("/:id/confirm", h.mdw.AuthCookieMiddleware(), h.ConfirmReceiptByID)
	receipt.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteReceiptByID)

}

// @Summary      Create a new receipt
// @Description  Create a new receipt record
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param        receipt  body      dto.CreateReceiptDTO  true  "Receipt data"
// @Success      201  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receipt/create [post]
func (h *receiptHandler) CreateReceipt(c *fiber.Ctx) error {

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

	var receipt dto.CreateReceiptDTO
	if err := c.BodyParser(&receipt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateReceipt(c.Context(), receipt, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create receipt" + err.Error(),
			MessageTH:  "สร้างใบเสร็จไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Receipt created successfully",
		MessageTH:  "สร้างใบเสร็จเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary      List receipts
// @Description  Get a paginated list of receipts
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field" Enums(created_at, updated_at) default(created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param status query string false "Filter by status" Enums(paid, pending)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param bill_type query string false "Filter by bill type" Enums(quotation, delivery_note, receipt)
// @Param type_receipt query string false "Filter by type receipt" Enums(company, shop)
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/receipt/list [get]
func (h *receiptHandler) ListReceipts(c *fiber.Ctx) error {

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

	var req dto.RequestListReceipt
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

	list, err := h.svc.ListReceipts(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.Status, req.StartDate, req.EndDate, req.BillType, req.TypeReceipt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list receipts" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการใบเสร็จ",
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

// @Summary      Get receipts by ID
// @Description  Get a receipt record by its ID
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receipt ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      404  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receipt/{id} [get]
func (h *receiptHandler) GetReceiptByID(c *fiber.Ctx) error {
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
	receiptID := c.Params("id")
	item, err := h.svc.GetReceiptByID(c.Context(), receiptID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get receipt",
			MessageTH:  "ไม่สามารถดึงใบเสร็จ",
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

// @Summary      Delete receipt by ID
// @Description  Delete a receipt record by its ID
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receipt ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receipt/{id} [delete]
func (h *receiptHandler) DeleteReceiptByID(c *fiber.Ctx) error {
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
	receiptID := c.Params("id")
	err = h.svc.DeleteReceiptByID(c.Context(), receiptID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete receipt",
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

// @Summary      Summary receipts
// @Description  Get summary of receipts
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param report query string true "Report type" Enums(day, month, all) default(all)
// @Param type_receipt query string false "Type receipt" Enums(company, shop)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receipt/summary [get]
func (h *receiptHandler) SummaryReceiptByFilter(c *fiber.Ctx) error {
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
	var req dto.RequestSummaryReceipt
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	summary, err := h.svc.SummaryReceiptByFilter(c.Context(), claims, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get receipt summary",
			MessageTH:  "ไม่สามารถดึงข้อมูลสรุปใบเสร็จ",
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

// @Summary      Confirm receipt by ID
// @Description  Confirm a receipt record by its ID
// @Tags         Receipts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receipt ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receipt/{id}/confirm [post]
func (h *receiptHandler) ConfirmReceiptByID(c *fiber.Ctx) error {
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
	receiptID := c.Params("id")
	err = h.svc.ConfirmReceiptByID(c.Context(), receiptID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to confirm receipt",
			MessageTH:  "ยืนยันไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Receipt confirmed successfully",
		MessageTH:  "ยืนยันใบเสร็จเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
