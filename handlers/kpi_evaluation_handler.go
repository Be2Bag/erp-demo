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
	kpiEvaluations.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)
	// kpiEvaluations.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluationByID)
	// kpiEvaluations.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateKPIEvaluation)
	// kpiEvaluations.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteKPIEvaluation)

	// kpiEvaluations.Get("/", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluations)
	// kpiEvaluations.Post("/", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)

	// kpiEvaluations.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetKPIStatistics)
}

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

func (h *KPIEvaluationHandler) CreateKPIEvaluation(c *fiber.Ctx) error {
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

	var createKPIEvaluation dto.CreateKPIEvaluationRequest
	if err := c.BodyParser(&createKPIEvaluation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateKPIEvaluation(c.Context(), createKPIEvaluation, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create KPI evaluation" + err.Error(),
			MessageTH:  "สร้างการประเมิน KPI ไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "KPI evaluation created successfully",
		MessageTH:  "สร้างการประเมิน KPI เรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}
