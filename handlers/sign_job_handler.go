package handlers

import (
	"errors"
	"strconv"

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
			MessageEN:  "Failed to create sign job",
			MessageTH:  "สร้างงานเซ็นชื่อไม่สำเร็จ",
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
// @Param size query int false "Page size (default 20)"
// @Param search query string false "Search text"
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

	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")

	if size > 100 {
		size = 100
	}
	if page < 1 {
		page = 1
	}

	list, err := h.svc.ListSignJobs(c.Context(), claims, page, size, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list sign jobs",
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
			MessageTH:  "ไม่พบนงาน",
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

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Sign job not found"
			MsgTH = "ไม่ใบพบงาน"
		}

		statusCode = fiber.StatusInternalServerError
		MsgEN = "Failed to update" + errOnUpdate.Error()
		MsgTH = "อัปเดตไม่สำเร็จ"
	}

	return c.JSON(dto.BaseResponse{
		StatusCode: statusCode,
		MessageEN:  MsgEN,
		MessageTH:  MsgTH,
		Status:     "success",
		Data:       nil,
	})

}

// @Summary Update Sign Job by ID
// @Description Update a sign job by its ID
// @Tags SignJob
// @Accept json
// @Produce json
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
