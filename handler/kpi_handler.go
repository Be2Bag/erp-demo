package handler

import (
	"fmt"
	"strings"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type KPIHandler struct {
	svc ports.KPIService
	mdw *middleware.Middleware
}

func NewKPIHandler(s ports.KPIService, mdw *middleware.Middleware) *KPIHandler {
	return &KPIHandler{svc: s, mdw: mdw}
}

func (h *KPIHandler) KPIRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	kpi := versionOne.Group("kpi")

	kpi.Get("/templates", h.mdw.AuthCookieMiddleware(), h.GetKPITemplates)
	kpi.Post("/templates", h.mdw.AuthCookieMiddleware(), h.CreateKPITemplate)
	kpi.Get("/templates/:id", h.mdw.AuthCookieMiddleware(), h.GetKPITemplateByID)
	kpi.Put("/templates/:id", h.mdw.AuthCookieMiddleware(), h.UpdateKPITemplate)
	kpi.Delete("/templates/:id", h.mdw.AuthCookieMiddleware(), h.DeleteKPITemplate)

	kpi.Get("/evaluations", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluations)
	kpi.Post("/evaluations", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)

	kpi.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetKPIStatistics)
}

// ตัวจัดการสำหรับการจัดการแม่แบบ KPI
func (h *KPIHandler) GetKPITemplates(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับดึงข้อมูลแม่แบบ KPI
	return nil
}

// @Summary Create a new KPI Template
// @Description Create a new KPI Template
// @Tags KPI
// @Accept json
// @Produce json
// @Param template body dto.KPITemplateDTO true "KPI Template"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/kpi/templates [post]
func (h *KPIHandler) CreateKPITemplate(c *fiber.Ctx) error {

	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}

	var template dto.KPITemplateDTO
	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	if err := h.svc.CreateKPITemplate(c.Context(), template, claims); err != nil {
		var (
			statusCode int
			messageEN  string
			messageTH  string
		)

		switch {
		case err.Error() == "items must not be empty":
			statusCode = fiber.StatusBadRequest
			messageEN = "Failed to create KPI Template: items must not be empty"
			messageTH = "ไม่สามารถสร้างแม่แบบ KPI ได้: ต้องมีรายการอย่างน้อยหนึ่งรายการ"
		case err.Error() == "sum of weights must be 100, got "+fmt.Sprint(template.TotalWeight):
			statusCode = fiber.StatusBadRequest
			messageEN = "Failed to create KPI Template: sum of weights must be 100, got " + fmt.Sprint(template.TotalWeight)
			messageTH = "ไม่สามารถสร้างแม่แบบ KPI ได้: ผลรวมของน้ำหนักต้องเท่ากับ 100"
		case err.Error() == "template with the same name already exists in this department":
			statusCode = fiber.StatusBadRequest
			messageEN = "Failed to create KPI Template: template with the same name already exists in this department"
			messageTH = "ไม่สามารถสร้างแม่แบบ KPI ได้: มีแม่แบบชื่อเดียวกันในแผนกนี้แล้ว"
		default:
			if strings.Contains(err.Error(), "items[") && (strings.Contains(err.Error(), "max_score") || strings.Contains(err.Error(), "weight")) {
				statusCode = fiber.StatusBadRequest
				messageEN = "Failed to create KPI Template: " + err.Error()
				messageTH = "ไม่สามารถสร้างแม่แบบ KPI ได้: ข้อมูลรายการไม่ถูกต้อง"
			} else {
				statusCode = fiber.StatusInternalServerError
				messageEN = "Failed to create KPI Template: " + err.Error()
				messageTH = "ไม่สามารถสร้างแม่แบบ KPI ได้"
			}
		}

		return c.Status(statusCode).JSON(dto.BaseResponse{
			StatusCode: statusCode,
			MessageEN:  messageEN,
			MessageTH:  messageTH,
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "KPI Template created successfully",
		MessageTH:  "สร้างแม่แบบ KPI สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

func (h *KPIHandler) GetKPITemplateByID(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับดึงข้อมูลแม่แบบ KPI ตามรหัส
	return nil
}

func (h *KPIHandler) UpdateKPITemplate(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับอัปเดตแม่แบบ KPI
	return nil
}

func (h *KPIHandler) DeleteKPITemplate(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับลบแม่แบบ KPI
	return nil
}

// ตัวจัดการสำหรับการประเมิน KPI
func (h *KPIHandler) GetKPIEvaluations(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับดึงข้อมูลการประเมิน KPI
	return nil
}

func (h *KPIHandler) CreateKPIEvaluation(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับสร้างการประเมิน KPI ใหม่
	return nil
}

// ตัวจัดการสำหรับสถิติ KPI
func (h *KPIHandler) GetKPIStatistics(c *fiber.Ctx) error {
	// การนำไปใช้สำหรับดึงข้อมูลสถิติ KPI
	return nil
}
