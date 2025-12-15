package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type DepartmentHandler struct {
	svc ports.DepartmentService
	mdw *middleware.Middleware
}

func NewDepartmentHandler(s ports.DepartmentService, mdw *middleware.Middleware) *DepartmentHandler {
	return &DepartmentHandler{svc: s, mdw: mdw}
}

func (h *DepartmentHandler) DepartmentRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	department := versionOne.Group("department")

	department.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetDepartmentList)
	department.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateDepartment)
	department.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetDepartmentByID)
	department.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateDepartment)
	department.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteDepartment)

}

// @Summary Create a new department
// @Description Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Param department body dto.CreateDepartmentDTO true "Department data"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/department/create [post]
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
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

	var createDepartment dto.CreateDepartmentDTO
	if err := c.BodyParser(&createDepartment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateDepartment(c.Context(), createDepartment, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create department: " + err.Error(),
			MessageTH:  "สร้างแผนกไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Department created successfully",
		MessageTH:  "สร้างแผนกเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Update an existing department
// @Description Update an existing department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param department body dto.UpdateDepartmentDTO true "Department data"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/department/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
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

	var updateDepartment dto.UpdateDepartmentDTO
	if err := c.BodyParser(&updateDepartment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	departmentID := c.Params("id")
	err = h.svc.UpdateDepartment(c.Context(), departmentID, updateDepartment, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update department: " + err.Error(),
			MessageTH:  "อัปเดตแผนกไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Department updated successfully",
		MessageTH:  "อัปเดตแผนกเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Delete a department
// @Description Delete a department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/department/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
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

	departmentID := c.Params("id")
	err = h.svc.DeleteDepartmentByID(c.Context(), departmentID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete department: " + err.Error(),
			MessageTH:  "ลบแผนกไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Department deleted successfully",
		MessageTH:  "ลบแผนกเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get a department by ID
// @Description Get a department by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse{data=dto.DepartmentDTO}
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/department/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *fiber.Ctx) error {
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

	departmentID := c.Params("id")
	department, err := h.svc.GetDepartmentByID(c.Context(), departmentID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get department: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลแผนกได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if department == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Department not found",
			MessageTH:  "ไม่พบแผนก",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Department retrieved successfully",
		MessageTH:  "ดึงข้อมูลแผนกเรียบร้อยแล้ว",
		Status:     "success",
		Data:       department,
	})
}

// @Summary Get a list of departments
// @Description Get a list of departments
// @Tags Departments
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "ค้นหา department_name"
// @Param sort_by query string false "เรียงตาม created_at updated_at department_name"
// @Param sort_order query string false "เรียงลำดับ (asc เก่า→ใหม่ | desc ใหม่→เก่า (ค่าเริ่มต้น))"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/department/list [get]
func (h *DepartmentHandler) GetDepartmentList(c *fiber.Ctx) error {
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

	var req dto.RequestListDepartment
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

	departments, err := h.svc.GetDepartmentList(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve departments: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลแผนกได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Projects retrieved successfully",
		MessageTH:  "ดึงข้อมูลแผนกเรียบร้อยแล้ว",
		Status:     "success",
		Data:       departments,
	})
}
