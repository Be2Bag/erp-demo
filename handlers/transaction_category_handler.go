package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type TransactionCategoryHandler struct {
	svc ports.TransactionCategoryService
	mdw *middleware.Middleware
}

func NewTransactionCategoryHandler(s ports.TransactionCategoryService, mdw *middleware.Middleware) *TransactionCategoryHandler {
	return &TransactionCategoryHandler{svc: s, mdw: mdw}
}

func (h *TransactionCategoryHandler) TransactionCategoryRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	transactionCategory := versionOne.Group("transaction-category")

	transactionCategory.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetTransactionCategoryList)
	transactionCategory.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateTransactionCategory)
	transactionCategory.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetTransactionCategoryByID)
	transactionCategory.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateTransactionCategory)
	transactionCategory.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteTransactionCategory)

}

// @Summary Create Transaction Category
// @Description สร้างหมวดหมู่รายการ (Transaction Category) ใหม่
// @Tags Transaction Category
// @Accept json
// @Produce json
// @Param body body dto.CreateTransactionCategoryDTO true "Transaction Category data"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/transaction-category/create [post]
func (h *TransactionCategoryHandler) CreateTransactionCategory(c *fiber.Ctx) error {
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

	var createCategory dto.CreateTransactionCategoryDTO
	if err := c.BodyParser(&createCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateTransactionCategory(c.Context(), createCategory, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create category" + err.Error(),
			MessageTH:  "สร้างหมวดหมู่ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Category created successfully",
		MessageTH:  "สร้างหมวดหมู่เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get Transaction Category List
// @Description ดึงรายการหมวดหมู่รายการ (Transaction Category) ทั้งหมด
// @Tags Transaction Category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field" default("created_at")
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param type query string false "Type filter Enums(income, expense) default(income)"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/transaction-category/list [get]
func (h *TransactionCategoryHandler) GetTransactionCategoryList(c *fiber.Ctx) error {
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

	var req dto.RequestListTransactionCategory
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

	categories, err := h.svc.ListTransactionCategory(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.Type)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve categories" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Categories retrieved successfully",
		MessageTH:  "ดึงข้อมูลหมวดหมู่เรียบร้อยแล้ว",
		Status:     "success",
		Data:       categories,
	})
}

// @Summary Get Transaction Category by ID
// @Description ดึงข้อมูลหมวดหมู่รายการ (Transaction Category) โดยใช้ ID
// @Tags Transaction Category
// @Accept json
// @Produce json
// @Param id path string true "Transaction Category ID"
// @Success 200 {object} dto.BaseResponse{data=dto.TransactionCategoryDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/transaction-category/{id} [get]
func (h *TransactionCategoryHandler) GetTransactionCategoryByID(c *fiber.Ctx) error {

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

	transactionCategoryByID := c.Params("id")
	if transactionCategoryByID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid category ID",
			MessageTH:  "รหัสหมวดหมู่ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	category, err := h.svc.GetTransactionCategoryByID(c.Context(), transactionCategoryByID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve category" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if category == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Category not found",
			MessageTH:  "ไม่พบหมวดหมู่",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Category retrieved successfully",
		MessageTH:  "ดึงข้อมูลหมวดหมู่เรียบร้อยแล้ว",
		Status:     "success",
		Data:       category,
	})
}

// @Summary Update Transaction Category
// @Description อัปเดตข้อมูลหมวดหมู่รายการ (Transaction Category) โดยใช้ ID
// @Tags Transaction Category
// @Accept json
// @Produce json
// @Param id path string true "Transaction Category ID"
// @Param body body dto.UpdateTransactionCategoryDTO true "Updated Transaction Category data"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/transaction-category/{id} [put]
func (h *TransactionCategoryHandler) UpdateTransactionCategory(c *fiber.Ctx) error {
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

	transactionCategoryByID := c.Params("id")
	if transactionCategoryByID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid category ID",
			MessageTH:  "รหัสหมวดหมู่ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	var updateDTO dto.UpdateTransactionCategoryDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.UpdateTransactionCategoryByID(c.Context(), transactionCategoryByID, updateDTO, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update category" + err.Error(),
			MessageTH:  "ไม่สามารถอัปเดตหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Category updated successfully",
		MessageTH:  "อัปเดตหมวดหมู่เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Delete Transaction Category
// @Description ลบหมวดหมู่รายการ (Transaction Category) โดยใช้ ID
// @Tags Transaction Category
// @Accept json
// @Produce json
// @Param id path string true "Transaction Category ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/transaction-category/{id} [delete]
func (h *TransactionCategoryHandler) DeleteTransactionCategory(c *fiber.Ctx) error {
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

	transactionCategoryID := c.Params("id")
	if transactionCategoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid project ID",
			MessageTH:  "รหัสโปรเจกต์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.DeleteTransactionCategoryByID(c.Context(), transactionCategoryID, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete category" + err.Error(),
			MessageTH:  "ไม่สามารถลบหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Category deleted successfully",
		MessageTH:  "ลบหมวดหมู่เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
