package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	svc ports.ProjectService
	mdw *middleware.Middleware
}

func NewProjectHandler(s ports.ProjectService, mdw *middleware.Middleware) *ProjectHandler {
	return &ProjectHandler{svc: s, mdw: mdw}
}

func (h *ProjectHandler) ProjectRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	project := versionOne.Group("project")

	project.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetProjectList)
	project.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateProject)
	project.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetProjectByID)
	project.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateProject)
	project.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteProject)

}

// @Summary Create a new project
// @Description Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body dto.CreateProjectDTO true "Project data"
// @Success 201 {object} dto.BaseResponse{data=dto.ProjectDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/project/create [post]
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
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

	var createProject dto.CreateProjectDTO
	if err := c.BodyParser(&createProject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateProject(c.Context(), createProject, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create project" + err.Error(),
			MessageTH:  "สร้างโปรเจกต์ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Project created successfully",
		MessageTH:  "สร้างโปรเจกต์เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get a list of projects
// @Description Get a list of projects
// @Tags Projects
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "ค้นหา project_name "
// @Param sort_by query string false "เรียงตาม created_at updated_at project_name"
// @Param sort_order query string false "เรียงลำดับ (asc เก่า→ใหม่ | desc ใหม่→เก่า (ค่าเริ่มต้น))"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/project/list [get]
func (h *ProjectHandler) GetProjectList(c *fiber.Ctx) error {
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

	var req dto.RequestListProject
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

	projects, err := h.svc.ListProject(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve projects" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลโปรเจกต์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Projects retrieved successfully",
		MessageTH:  "ดึงข้อมูลโปรเจกต์เรียบร้อยแล้ว",
		Status:     "success",
		Data:       projects,
	})
}

// @Summary Get a project by ID
// @Description Get a project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} dto.BaseResponse{data=dto.ProjectDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/project/{id} [get]
func (h *ProjectHandler) GetProjectByID(c *fiber.Ctx) error {
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

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid project ID",
			MessageTH:  "รหัสโปรเจกต์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	project, err := h.svc.GetProjectByID(c.Context(), projectID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve project" + err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลโปรเจกต์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Project retrieved successfully",
		MessageTH:  "ดึงข้อมูลโปรเจกต์เรียบร้อยแล้ว",
		Status:     "success",
		Data:       project,
	})
}

// @Summary Update a project by ID
// @Description Update a project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param body body dto.UpdateProjectDTO true "Project data"
// @Success 200 {object} dto.BaseResponse{data=dto.ProjectDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/project/{id} [put]
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
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

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid project ID",
			MessageTH:  "รหัสโปรเจกต์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	var updateDTO dto.UpdateProjectDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.UpdateProjectByID(c.Context(), projectID, updateDTO, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update project" + err.Error(),
			MessageTH:  "ไม่สามารถอัปเดตโปรเจกต์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Project updated successfully",
		MessageTH:  "อัปเดตโปรเจกต์เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Delete a project by ID
// @Description Delete a project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
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

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid project ID",
			MessageTH:  "รหัสโปรเจกต์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if err := h.svc.DeleteProjectByID(c.Context(), projectID, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete project" + err.Error(),
			MessageTH:  "ไม่สามารถลบโปรเจกต์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Project deleted successfully",
		MessageTH:  "ลบโปรเจกต์เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
