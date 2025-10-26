package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReceivableHandler struct {
	svc ports.ReceivableService
	mdw *middleware.Middleware
}

func NewReceivableHandler(s ports.ReceivableService, mdw *middleware.Middleware) *ReceivableHandler {
	return &ReceivableHandler{svc: s, mdw: mdw}
}

func (h *ReceivableHandler) ReceivableRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	receivable := versionOne.Group("receivable")

	receivable.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateReceivable)
	receivable.Post("/record-receipt", h.mdw.AuthCookieMiddleware(), h.RecordReceipt)
	receivable.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListReceivables)
	receivable.Get("/summary", h.mdw.AuthCookieMiddleware(), h.SummaryReceivableByFilter)
	receivable.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetReceivableByID)
	receivable.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateReceivableByID)
	receivable.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteReceivableByID)

}

// @Summary      Create a new receivable
// @Description  Create a new receivable record
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param        receivable  body      dto.CreateReceivableDTO  true  "Receivable data"
// @Success      201  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/create [post]
func (h *ReceivableHandler) CreateReceivable(c *fiber.Ctx) error {

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

	var inCome dto.CreateReceivableDTO
	if err := c.BodyParser(&inCome); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateReceivable(c.Context(), inCome, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create receivable" + err.Error(),
			MessageTH:  "สร้างลูกหนี้ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Receivable created successfully",
		MessageTH:  "สร้างลูกหนี้เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary      List receivables
// @Description  Get a paginated list of receivables
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sortBy query string false "Sort by field" Enums(created_at, updated_at) default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc, desc) default(desc)
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/receivable/list [get]
func (h *ReceivableHandler) ListReceivables(c *fiber.Ctx) error {

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

	var req dto.RequestListReceivable
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

	list, err := h.svc.ListReceivables(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list receivables" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการลูกหนี้",
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

// @Summary      Get receivable by ID
// @Description  Get a receivable record by its ID
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receivable ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      404  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/{id} [get]
func (h *ReceivableHandler) GetReceivableByID(c *fiber.Ctx) error {
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
	receivableID := c.Params("id")
	item, err := h.svc.GetReceivableByID(c.Context(), receivableID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get receivable",
			MessageTH:  "ไม่สามารถดึงลูกหนี้",
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

// @Summary      Update receivable by ID
// @Description  Update a receivable record by its ID
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receivable ID"
// @Param        receivable  body      dto.UpdateReceivableDTO  true  "Receivable data to update"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      404  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/{id} [put]
func (h *ReceivableHandler) UpdateReceivableByID(c *fiber.Ctx) error {
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
	receivableID := c.Params("id")
	var body dto.UpdateReceivableDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdateReceivableByID(c.Context(), receivableID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Receivable not found"
			MsgTH = "ไม่พบลูกหนี้"
		}

		statusCode = fiber.StatusInternalServerError
		MsgEN = "Failed to update receivable" + errOnUpdate.Error()
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

// @Summary      Delete receivable by ID
// @Description  Delete a receivable record by its ID
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Receivable ID"
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/{id} [delete]
func (h *ReceivableHandler) DeleteReceivableByID(c *fiber.Ctx) error {
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
	receivableID := c.Params("id")
	err = h.svc.DeleteReceivableByID(c.Context(), receivableID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete receivable",
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

// @Summary      Summary receivables
// @Description  Get summary of receivables
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param bank_id query string false "Bank ID"
// @Param report query string true "Report type" Enums(day, month, all)
// @Success      200  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/summary [get]
func (h *ReceivableHandler) SummaryReceivableByFilter(c *fiber.Ctx) error {
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
	var req dto.RequestSummaryReceivable
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	summary, err := h.svc.SummaryReceivableByFilter(c.Context(), claims, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get receivable summary",
			MessageTH:  "ไม่สามารถดึงข้อมูลสรุปลูกหนี้",
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

// @Summary      Record receipt
// @Description  Record a receipt for a receivable
// @Tags         Receivables
// @Accept       json
// @Produce      json
// @Param        receipt  body      dto.RecordReceiptDTO  true  "Receipt data"
// @Success      201  {object}  dto.BaseResponse
// @Failure      400  {object}  dto.BaseResponse
// @Failure      401  {object}  dto.BaseResponse
// @Failure      500  {object}  dto.BaseResponse
// @Router       /v1/receivable/record-receipt [post]
func (h *ReceivableHandler) RecordReceipt(c *fiber.Ctx) error {
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
	var inCome dto.RecordReceiptDTO
	if err := c.BodyParser(&inCome); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	err = h.svc.RecordReceipt(c.Context(), inCome, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to record receipt" + err.Error(),
			MessageTH:  "บันทึกรายรับไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Receipt recorded successfully",
		MessageTH:  "บันทึกรายรับเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
