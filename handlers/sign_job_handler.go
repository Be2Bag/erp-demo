package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignJobHandler struct {
	svc ports.SignJobService
	mdw *middleware.Middleware
}

func NewSignJobHandler(s ports.SignJobService, mdw *middleware.Middleware) *SignJobHandler {
	return &SignJobHandler{svc: s, mdw: mdw}
}

func (h *SignJobHandler) SignJobRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	signJob := versionOne.Group("sign-job")

	signJob.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateSignJob)
	signJob.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListSignJobs)
	signJob.Put("/verify/:id", h.mdw.AuthCookieMiddleware(), h.VerifySignJob)
	signJob.Put("/confirm/:id", h.mdw.AuthCookieMiddleware(), h.ConfirmSignJob)
	signJob.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetSignJobByID)
	signJob.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateSignJobByID)
	signJob.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteSignJobByID)

}

// @Summary Create Sign Job
// @Description Create a new sign job
// @Tags SignJob
// @Accept json
// @Produce json
// @Param request body dto.CreateSignJobDTO true "Create Sign Job"
// @Success 201 {object} dto.BaseResponse{data=dto.SignJobDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/create [post]
func (h *SignJobHandler) CreateSignJob(c *fiber.Ctx) error {

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

	var signJob dto.CreateSignJobDTO
	if err := c.BodyParser(&signJob); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateSignJob(c.Context(), signJob, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create sign job: " + err.Error(),
			MessageTH:  "สร้างงานไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Sign job created successfully",
		MessageTH:  "สร้างงานเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List Sign Jobs
// @Description List sign jobs (search & pagination)
// @Tags SignJob
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 10)"
// @Param search query string false "ค้นหาด้วย ชื่อโปรเจกต์, ชื่องาน,ชื่อบริษัท,ชื่อผู้ติดต่อ "
// @Param status query string false "สถานะงาน in_progress, done"
// @Param sort_by query string false "เรียงตาม created_at updated_at due_date job_name project_name company_name status price_thb quantity"
// @Param sort_order query string false "เรียงลำดับ (asc เก่า→ใหม่ | desc ใหม่→เก่า (ค่าเริ่มต้น))"
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/list [get]
func (h *SignJobHandler) ListSignJobs(c *fiber.Ctx) error {

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

	var req dto.RequestListSignJobs
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

	list, err := h.svc.ListSignJobs(c.Context(), claims, req.Page, req.Limit, req.Search, req.Status, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list sign jobs: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการงานได้",
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

// @Summary Get Sign Job by ID
// @Description Get a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
// @Param id path string true "Sign Job ID"
// @Success 200 {object} dto.BaseResponse{data=dto.SignJobDTO}
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/{id} [get]
func (h *SignJobHandler) GetSignJobByID(c *fiber.Ctx) error {
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
	jobID := c.Params("id")
	item, err := h.svc.GetSignJobByJobID(c.Context(), jobID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get sign job",
			MessageTH:  "ไม่สามารถดึงงานได้",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Sign job not found",
			MessageTH:  "ไม่พบงาน",
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

// @Summary Update Sign Job by ID
// @Description Update a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
// @Param id path string true "Sign Job ID"
// @Param body body dto.UpdateSignJobDTO true "Update Sign Job"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/{id} [put]
func (h *SignJobHandler) UpdateSignJobByID(c *fiber.Ctx) error {
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
	jobID := c.Params("id")
	var body dto.UpdateSignJobDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdateSignJobByJobID(c.Context(), jobID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"
	status := "success"

	if errOnUpdate != nil {
		status = "error"
		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Sign job not found"
			MsgTH = "ไม่พบใบงาน"
		} else {
			statusCode = fiber.StatusInternalServerError
			MsgEN = "Failed to update: " + errOnUpdate.Error()
			MsgTH = "อัปเดตไม่สำเร็จ"
		}
	}

	return c.Status(statusCode).JSON(dto.BaseResponse{
		StatusCode: statusCode,
		MessageEN:  MsgEN,
		MessageTH:  MsgTH,
		Status:     status,
		Data:       nil,
	})

}

// @Summary Update Sign Job by ID
// @Description Update a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
// @Param id path string true "Sign Job ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/{id} [delete]
func (h *SignJobHandler) DeleteSignJobByID(c *fiber.Ctx) error {
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
	jobID := c.Params("id")
	err = h.svc.DeleteSignJobByJobID(c.Context(), jobID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete",
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

// @Summary Verify Sign Job
// @Description Verify a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
// @Param id path string true "Sign Job ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/verify/{id} [put]
func (h *SignJobHandler) VerifySignJob(c *fiber.Ctx) error {
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
	jobID := c.Params("id")
	err = h.svc.VerifySignJob(c.Context(), jobID, claims)
	if err != nil {

		statusCode := fiber.StatusInternalServerError
		MsgEN := "Failed to verify"
		MsgTH := "ส่งงานไม่สำเร็จ"

		if errors.Is(err, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Sign job not found"
			MsgTH = "ไม่พบใบงาน"
		} else if err.Error() == "can not verify" {
			statusCode = fiber.StatusBadRequest
			MsgEN = err.Error()
			MsgTH = "ไม่สามารถยืนยันงานได้ เนื่องจากมีงานที่กำลังดำเนินการอยู่ในระบบ"
		} else if err.Error() == "no tasks found for this job" {
			statusCode = fiber.StatusBadRequest
			MsgEN = err.Error()
			MsgTH = "ไม่พบงานใดๆ ที่จัดการงานของใบงานนี้"
		} else if err.Error() == "only admin can verify" {
			statusCode = fiber.StatusForbidden
			MsgEN = err.Error()
			MsgTH = "มีเพียงแอดมินเท่านั้นที่สามารถยืนยันได้"
		}

		return c.Status(statusCode).JSON(dto.BaseResponse{
			StatusCode: statusCode,
			MessageEN:  MsgEN,
			MessageTH:  MsgTH,
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Verified",
		MessageTH:  "ตรวจสอบแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Confirm Sign Job
// @Description Confirm a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
// @Param id path string true "Sign Job ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-job/confirm/{id} [put]
func (h *SignJobHandler) ConfirmSignJob(c *fiber.Ctx) error {
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
	jobID := c.Params("id")
	err = h.svc.ConfirmSignJob(c.Context(), jobID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to confirm",
			MessageTH:  "ยืนยันงานไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Confirmed",
		MessageTH:  "ยืนยันแล้ว สามารถเริ่มงานได้",
		Status:     "success",
		Data:       nil,
	})
}
