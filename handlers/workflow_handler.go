package handlers

import (
	"strconv"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type WorkFlowHandler struct {
	svc ports.WorkFlowService
	mdw *middleware.Middleware
}

func NewWorkFlowHandler(s ports.WorkFlowService, mdw *middleware.Middleware) *WorkFlowHandler {
	return &WorkFlowHandler{svc: s, mdw: mdw}
}

func (h *WorkFlowHandler) WorkFlowRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	workFlow := versionOne.Group("workflow")
	workFlow.Get("/list", h.ListWorkflows)
	workFlow.Post("/create", h.mdw.AuthCookieMiddleware(), h.createWorkflow)
	workFlow.Get("/:id", h.GetWorkflowByID)
	workFlow.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateWorkflow)
	workFlow.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteWorkflow)
}

// @Summary Create a new workflow
// @Description Create a new workflow
// @Tags Workflow
// @Accept json
// @Produce json
// @Param request body dto.CreateWorkflowTemplateDTO true "Create Workflow Template"
// @Success 201 {object} dto.BaseResponse{data=dto.WorkflowTemplateDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Router /v1/workflow/create [post]
func (h *WorkFlowHandler) createWorkflow(c *fiber.Ctx) error {
	var req dto.CreateWorkflowTemplateDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}

	out, err := h.svc.CreateWorkflowTemplate(c.Context(), req, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถสร้าง Workflow ได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Workflow created successfully",
		MessageTH:  "สร้าง Workflow สำเร็จ",
		Status:     "success",
		Data:       out,
	})
}

// @Summary Get a workflow by ID
// @Description Get a workflow by ID
// @Tags Workflow
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Success 200 {object} dto.BaseResponse{data=dto.WorkflowTemplateDTO}
// @Failure 404 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Router /v1/workflow/{id} [get]
func (h *WorkFlowHandler) GetWorkflowByID(c *fiber.Ctx) error {
	id := c.Params("id")
	out, err := h.svc.GetWorkflowTemplateByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Workflow not found",
			MessageTH:  "ไม่พบ Workflow",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "OK",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       out,
	})
}

// @Summary List all workflows
// @Description List all workflows
// @Tags Workflow
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param department query string false "Department"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} dto.BaseResponse{data=[]dto.WorkflowTemplateDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Router /v1/workflow/list [get]
func (h *WorkFlowHandler) ListWorkflows(c *fiber.Ctx) error {
	search := c.Query("search")
	dept := c.Query("department")
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	sort := c.Query("sort", "updated_at:desc")

	items, total, err := h.svc.ListWorkflowTemplates(c.Context(), search, dept, page, limit, sort)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลได้",
			Status:     "error",
			Data:       nil,
		})
	}

	var totalPages int
	if limit > 0 {
		totalPages = int((total + limit - 1) / limit)
	} else {
		totalPages = 0
	}

	data := fiber.Map{
		"items": items,
		"pagination": dto.Pagination{
			Page:       int(page),
			Size:       int(limit),
			TotalCount: int(total),
			TotalPages: totalPages,
		},
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "OK",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       data,
	})
}

// @Summary Update a workflow
// @Description Update a workflow
// @Tags Workflow
// @Accept json
// @Produce json
// @Param id path string true "Workflow ID"
// @Param request body dto.UpdateWorkflowTemplateDTO true "Update Workflow Template"
// @Success 200 {object} dto.BaseResponse{data=dto.WorkflowTemplateDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Router /v1/workflow/{id} [put]
func (h *WorkFlowHandler) UpdateWorkflow(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateWorkflowTemplateDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}
	out, err := h.svc.UpdateWorkflowTemplate(c.Context(), id, req, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถอัปเดต Workflow ได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Workflow updated successfully",
		MessageTH:  "อัปเดต Workflow สำเร็จ",
		Status:     "success",
		Data:       out,
	})
}

// @Summary Delete a workflow
// @Description Delete a workflow
// @Tags Workflow
// @Param id path string true "Workflow ID"
// @Success 204 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Router /v1/workflow/{id} [delete]
func (h *WorkFlowHandler) DeleteWorkflow(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteWorkflowTemplate(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถลบ Workflow ได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Workflow deleted successfully",
		MessageTH:  "ลบ Workflow สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
