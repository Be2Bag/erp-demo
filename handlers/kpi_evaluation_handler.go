package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type KPIEvaluationHandler struct {
	svc ports.KPIEvaluationService
	mdw *middleware.Middleware
}

func NewKPIEvaluationHandler(s ports.KPIEvaluationService, mdw *middleware.Middleware) *KPIEvaluationHandler {
	return &KPIEvaluationHandler{svc: s, mdw: mdw}
}

func (h *KPIEvaluationHandler) KPIEvaluationRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	kpiEvaluations := versionOne.Group("kpi-evaluations")

	kpiEvaluations.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluationList)
	// kpiEvaluations.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)
	// kpiEvaluations.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluationByID)
	kpiEvaluations.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateKPIEvaluation)
	// kpiEvaluations.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteKPIEvaluation)

	// kpiEvaluations.Get("/", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluations)
	// kpiEvaluations.Post("/", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)

	// kpiEvaluations.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetKPIStatistics)
}

// @Summary List KPI Evaluations
// @Description Get a list of KPI evaluations
// @Tags KPI Evaluations
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search term"
// @Param department_id query string false "Department ID"
// @Param sort_by query string false "Sort by"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {object} dto.BaseResponse{data=[]dto.KPIEvaluationResponse}
// @Router /v1/kpi-evaluations/list [get]
func (h *KPIEvaluationHandler) GetKPIEvaluationList(c *fiber.Ctx) error {

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

	var req dto.RequestListKPIEvaluation
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

	tasks, errOnGetTasks := h.svc.ListKPIEvaluation(c.Context(), claims, req.Page, req.Limit, req.Search, req.Department, req.SortBy, req.SortOrder)
	if errOnGetTasks != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  errOnGetTasks.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "OK",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       tasks,
	})
}

// @Summary Update KPI Evaluation
// @Description Update an existing KPI evaluation
// @Tags KPI Evaluations
// @Accept json
// @Produce json
// @Param id path string true "KPI Evaluation ID"
// @Param request body dto.UpdateKPIEvaluationRequest true "Update KPI Evaluation Request"
// @Success 200 {object} dto.BaseResponse
// @Router /v1/kpi-evaluations/{id} [put]
func (h *KPIEvaluationHandler) UpdateKPIEvaluation(c *fiber.Ctx) error {
	var req dto.UpdateKPIEvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	evaluationID := c.Params("id")
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

	if err := h.svc.UpdateKPIEvaluation(c.Context(), evaluationID, req, claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถอัปเดตการประเมิน KPI ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "OK",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
