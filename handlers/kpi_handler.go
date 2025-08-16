package handlers

import (
	"fmt"
	"strconv" // added
	"strings"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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

	kpi.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetKPITemplateList)
	kpi.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateKPITemplate)
	kpi.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetKPITemplateByID)
	kpi.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateKPITemplate)
	kpi.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteKPITemplate)

	kpi.Get("/evaluations", h.mdw.AuthCookieMiddleware(), h.GetKPIEvaluations)
	kpi.Post("/evaluations", h.mdw.AuthCookieMiddleware(), h.CreateKPIEvaluation)

	kpi.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetKPIStatistics)
}

// @Summary Get KPI Templates
// @Description Get a list of KPI Templates
// @Tags KPI
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search keyword"
// @Param department query string false "Department"
// @Param is_active query bool false "Is Active"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/kpi/list [get]
func (h *KPIHandler) GetKPITemplateList(c *fiber.Ctx) error {

	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid page parameter",
			MessageTH:  "พารามิเตอร์ page ไม่ถูกต้อง",
			Status:     "error",
		})
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid limit parameter",
			MessageTH:  "พารามิเตอร์ limit ไม่ถูกต้อง",
			Status:     "error",
		})
	}
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	search := strings.TrimSpace(c.Query("search", ""))
	dept := strings.TrimSpace(c.Query("department", ""))

	isActivePtr, parseErr := parseOptionalBool(c.Query("is_active", ""))
	if parseErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid is_active parameter",
			MessageTH:  "พารามิเตอร์ is_active ไม่ถูกต้อง",
			Status:     "error",
		})
	}

	q := dto.KPITemplateListQuery{
		Page:       page,
		Limit:      limit,
		Search:     search,
		Department: dept,
		IsActive:   isActivePtr,
	}

	list, errOnGetList := h.svc.ListKPITemplates(c.Context(), q)
	if errOnGetList != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve KPI Templates" + errOnGetList.Error(),
			MessageTH:  "ไม่สามารถดึงแม่แบบ KPI ได้",
			Status:     "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "KPI Templates retrieved successfully",
		MessageTH:  "ดึงแม่แบบ KPI สำเร็จ",
		Status:     "success",
		Data:       list,
	})
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
// @Router /v1/kpi/create [post]
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

// @Summary Get a KPI Template by ID
// @Description Get a KPI Template by ID
// @Tags KPI
// @Accept json
// @Produce json
// @Param id path string true "KPI Template ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/kpi/{id} [get]
func (h *KPIHandler) GetKPITemplateByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "id is required",
			MessageTH:  "ต้องระบุรหัส",
			Status:     "error",
		})
	}

	tpl, err := h.svc.GetKPITemplateByID(c.Context(), id)
	if err != nil {
		status := fiber.StatusInternalServerError
		msgEN := "Failed to get KPI Template"
		msgTH := "ไม่สามารถดึงแม่แบบ KPI ได้"
		if err == mongo.ErrNoDocuments {
			status = fiber.StatusNotFound
			msgEN = "KPI Template not found"
			msgTH = "ไม่พบแม่แบบ KPI"
		}
		return c.Status(status).JSON(dto.BaseResponse{
			StatusCode: status,
			MessageEN:  msgEN,
			MessageTH:  msgTH,
			Status:     "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "KPI Template retrieved successfully",
		MessageTH:  "ดึงแม่แบบ KPI สำเร็จ",
		Status:     "success",
		Data:       tpl,
	})
}

// @Summary Update a KPI Template
// @Description Update a KPI Template
// @Tags KPI
// @Accept json
// @Produce json
// @Param id path string true "KPI Template ID"
// @Param body body dto.KPITemplateDTO true "KPI Template Data"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/kpi/{id} [put]
func (h *KPIHandler) UpdateKPITemplate(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "id is required",
			MessageTH:  "ต้องระบุรหัส",
			Status:     "error",
		})
	}

	var body dto.KPITemplateDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
		})
	}

	tpl, err := h.svc.UpdateKPITemplate(c.Context(), id, body, claims)
	if err != nil {
		status := fiber.StatusInternalServerError
		msgEN := "Failed to update KPI Template: " + err.Error()
		msgTH := "ไม่สามารถอัปเดตแม่แบบ KPI ได้"

		switch {
		case err == mongo.ErrNoDocuments:
			status = fiber.StatusNotFound
			msgEN = "KPI Template not found"
			msgTH = "ไม่พบแม่แบบ KPI"
		case err.Error() == "items must not be empty":
			status = fiber.StatusBadRequest
			msgEN = "Failed to update KPI Template: items must not be empty"
			msgTH = "ไม่สามารถอัปเดตแม่แบบ KPI ได้: ต้องมีรายการอย่างน้อยหนึ่งรายการ"
		case err.Error() == "sum of weights must be 100":
			status = fiber.StatusBadRequest
			msgEN = "Failed to update KPI Template: sum of weights must be 100"
			msgTH = "ไม่สามารถอัปเดตแม่แบบ KPI ได้: ผลรวมน้ำหนักต้องเท่ากับ 100"
		case err.Error() == "template with the same name already exists in this department":
			status = fiber.StatusBadRequest
			msgEN = "Failed to update KPI Template: duplicate name in department"
			msgTH = "ไม่สามารถอัปเดตแม่แบบ KPI ได้: มีชื่อซ้ำในแผนก"
		}

		return c.Status(status).JSON(dto.BaseResponse{
			StatusCode: status,
			MessageEN:  msgEN,
			MessageTH:  msgTH,
			Status:     "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "KPI Template updated successfully",
		MessageTH:  "อัปเดตแม่แบบ KPI สำเร็จ",
		Status:     "success",
		Data:       tpl,
	})
}

// @Summary Delete a KPI Template
// @Description Delete a KPI Template
// @Tags KPI
// @Accept json
// @Produce json
// @Param id path string true "KPI Template ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/kpi/{id} [delete]
func (h *KPIHandler) DeleteKPITemplate(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "id is required",
			MessageTH:  "ต้องระบุรหัส",
			Status:     "error",
		})
	}

	if err := h.svc.DeleteKPITemplate(c.Context(), id); err != nil {
		status := fiber.StatusInternalServerError
		msgEN := "Failed to delete KPI Template"
		msgTH := "ไม่สามารถลบแม่แบบ KPI ได้"
		if err == mongo.ErrNoDocuments {
			status = fiber.StatusNotFound
			msgEN = "KPI Template not found"
			msgTH = "ไม่พบแม่แบบ KPI"
		}
		return c.Status(status).JSON(dto.BaseResponse{
			StatusCode: status,
			MessageEN:  msgEN,
			MessageTH:  msgTH,
			Status:     "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "KPI Template deleted successfully",
		MessageTH:  "ลบแม่แบบ KPI สำเร็จ",
		Status:     "success",
	})
}

// ตัวจัดการสำหรับการประเมิน KPI
func (h *KPIHandler) GetKPIEvaluations(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusNotImplemented,
		MessageEN:  "Not implemented",
		MessageTH:  "ยังไม่ถูกพัฒนา",
		Status:     "error",
	})
}

func (h *KPIHandler) CreateKPIEvaluation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusNotImplemented,
		MessageEN:  "Not implemented",
		MessageTH:  "ยังไม่ถูกพัฒนา",
		Status:     "error",
	})
}

// ตัวจัดการสำหรับสถิติ KPI
func (h *KPIHandler) GetKPIStatistics(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusNotImplemented,
		MessageEN:  "Not implemented",
		MessageTH:  "ยังไม่ถูกพัฒนา",
		Status:     "error",
	})
}

func parseOptionalBool(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
