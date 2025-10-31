package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseHandler struct {
	svc ports.ExpenseService
	mdw *middleware.Middleware
}

func NewExpenseHandler(s ports.ExpenseService, mdw *middleware.Middleware) *ExpenseHandler {
	return &ExpenseHandler{svc: s, mdw: mdw}
}

func (h *ExpenseHandler) ExpenseRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	expense := versionOne.Group("expense")

	expense.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateExpense)
	expense.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListExpenses)
	expense.Get("/summary", h.mdw.AuthCookieMiddleware(), h.SummaryExpenseByFilter)
	expense.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetExpenseByID)
	expense.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateExpenseByID)
	expense.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteExpenseByID)

}

// @Summary Create Expense
// @Description Create a new expense
// @Tags Expense
// @Accept json
// @Produce json
// @Param request body dto.CreateExpenseDTO true "Create Expense"
// @Success 201 {object} dto.BaseResponse{data=dto.ExpenseDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/create [post]
func (h *ExpenseHandler) CreateExpense(c *fiber.Ctx) error {

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

	var expense dto.CreateExpenseDTO
	if err := c.BodyParser(&expense); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateExpense(c.Context(), expense, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create expense" + err.Error(),
			MessageTH:  "สร้างรายจ่ายไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Expense created successfully",
		MessageTH:  "สร้างรายจ่ายเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List Expenses
// @Description List all expenses with pagination and optional search
// @Tags Expense
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10) maximum(100)
// @Param search query string false "Search term"
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order" Enums(asc, desc) default("desc")
// @Param transaction_category_id query string false "Transaction Category ID"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param bank_id query string false "Bank ID"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/list [get]
func (h *ExpenseHandler) ListExpenses(c *fiber.Ctx) error {

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

	var req dto.RequestListExpense
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

	list, err := h.svc.ListExpenses(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.TransactionCategoryID, req.StartDate, req.EndDate, req.BankID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list expenses" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการรายจ่าย",
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

// @Summary Get Expense by ID
// @Description Get a single expense by its ID
// @Tags Expense
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} dto.BaseResponse{data=dto.ExpenseDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/{id} [get]
func (h *ExpenseHandler) GetExpenseByID(c *fiber.Ctx) error {
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
	expenseID := c.Params("id")
	item, err := h.svc.GetExpenseByID(c.Context(), expenseID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get expense",
			MessageTH:  "ไม่สามารถดึงรายจ่าย",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Expense not found",
			MessageTH:  "ไม่พบรายจ่าย",
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

// @Summary Update Expense by ID
// @Description Update an existing expense by its ID
// @Tags Expense
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Param request body dto.UpdateExpenseDTO true "Update Expense"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/{id} [put]
func (h *ExpenseHandler) UpdateExpenseByID(c *fiber.Ctx) error {
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
	expenseID := c.Params("id")
	var body dto.UpdateExpenseDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdateExpenseByID(c.Context(), expenseID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Expense not found"
			MsgTH = "ไม่พบรายจ่าย"
		}

		statusCode = fiber.StatusInternalServerError
		MsgEN = "Failed to update expense" + errOnUpdate.Error()
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

// @Summary Delete Expense by ID
// @Description Soft delete an expense by its ID
// @Tags Expense
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/{id} [delete]
func (h *ExpenseHandler) DeleteExpenseByID(c *fiber.Ctx) error {
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
	expenseID := c.Params("id")
	err = h.svc.DeleteExpenseByID(c.Context(), expenseID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete expense",
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

// @Summary Expense Summary by Filter
// @Description Get summary of expense for today, this month, and all time
// @Tags Expense
// @Accept json
// @Produce json
// @Param bank_id query string false "Bank ID to filter expenses"
// @Success 200 {object} dto.BaseResponse{data=dto.ExpenseSummaryDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/expense/summary [get]
func (h *ExpenseHandler) SummaryExpenseByFilter(c *fiber.Ctx) error {
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
	var req dto.RequestExpenseSummary
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	summary, err := h.svc.SummaryExpenseByFilter(c.Context(), claims, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get expense summary",
			MessageTH:  "ไม่สามารถดึงข้อมูลสรุปรายจ่าย",
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
