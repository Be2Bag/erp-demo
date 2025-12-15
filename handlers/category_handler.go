package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	svc ports.CategoryService
	mdw *middleware.Middleware
}

func NewCategoryHandler(s ports.CategoryService, mdw *middleware.Middleware) *CategoryHandler {
	return &CategoryHandler{svc: s, mdw: mdw}
}

func (h *CategoryHandler) CategoryRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	category := versionOne.Group("category")

	category.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetCategoryList)
	category.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateCategory)
	category.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetCategoryByID)
	category.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateCategory)
	category.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteCategory)

}

// @Summary Create a new category
// @Description Create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param body body dto.CreateCategoryDTO true "Category data"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/category/create [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
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

	var createCategory dto.CreateCategoryDTO
	if err := c.BodyParser(&createCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateCategory(c.Context(), createCategory, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create category: " + err.Error(),
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

// @Summary Get a list of categories
// @Description Get a list of categories
// @Tags Category
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/category/list [get]
func (h *CategoryHandler) GetCategoryList(c *fiber.Ctx) error {
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

	var req dto.RequestListCategory
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

	categories, err := h.svc.ListCategory(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve categories: " + err.Error(),
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

// @Summary Get a category by ID
// @Description Get a category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/category/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {

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

	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid category ID",
			MessageTH:  "รหัสหมวดหมู่ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	category, err := h.svc.GetCategoryByID(c.Context(), categoryID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve category: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
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

// @Summary Update a category
// @Description Update a category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param body body dto.UpdateCategoryDTO true "Category data"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/category/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
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

	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid category ID",
			MessageTH:  "รหัสหมวดหมู่ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	var updateDTO dto.UpdateCategoryDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.UpdateCategoryByID(c.Context(), categoryID, updateDTO, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update category: " + err.Error(),
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

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/category/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
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

	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid project ID",
			MessageTH:  "รหัสโปรเจกต์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.DeleteCategoryByID(c.Context(), categoryID, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete category: " + err.Error(),
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
