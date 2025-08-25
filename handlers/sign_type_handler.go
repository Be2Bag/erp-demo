package handlers

import (
	"errors"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignTypeHandler struct {
	svc ports.SignTypeService
	mdw *middleware.Middleware
}

func NewSignTypeHandler(s ports.SignTypeService, mdw *middleware.Middleware) *SignTypeHandler {
	return &SignTypeHandler{svc: s, mdw: mdw}
}

func (h *SignTypeHandler) SignTypeRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	signType := versionOne.Group("sign-type")

	signType.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateSignType)
	signType.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListSignTypes)
	signType.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetSignTypeByID)
	signType.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateSignTypeByID)
	signType.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteSignTypeByID)

}

// @Summary Create Sign Type
// @Description Create a new sign type
// @Tags SignType
// @Accept json
// @Produce json
// @Param request body dto.CreateSignTypeDTO true "Create Sign Type"
// @Success 201 {object} dto.BaseResponse{data=dto.CreateSignTypeDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-type/create [post]
func (h *SignTypeHandler) CreateSignType(c *fiber.Ctx) error {

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

	var signType dto.CreateSignTypeDTO
	if err := c.BodyParser(&signType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateSignType(c.Context(), signType, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create sign type" + err.Error(),
			MessageTH:  "สร้างประเภทงานไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Sign type created successfully",
		MessageTH:  "สร้างประเภทงานเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List Sign Types
// @Description Get a list of sign types
// @Tags SignType
// @Accept json
// @Produce json
// @Param request query dto.RequestListSignType true "List Sign Types"
// @Success 200 {object} dto.BaseResponse{data=[]dto.SignTypeDTO}
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-type/list [get]
func (h *SignTypeHandler) ListSignTypes(c *fiber.Ctx) error {

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

	var req dto.RequestListSignType
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

	list, err := h.svc.ListSignTypes(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list sign types" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการประเภทงานได้",
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

// @Summary Get Sign Type by ID
// @Description Get a sign type by its ID
// @Tags SignType
// @Accept json
// @Produce json
// @Param id path string true "Sign Type ID"
// @Success 200 {object} dto.BaseResponse{data=dto.SignTypeDTO}
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-type/{id} [get]
func (h *SignTypeHandler) GetSignTypeByID(c *fiber.Ctx) error {
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
	signTypeID := c.Params("id")
	item, err := h.svc.GetSignTypeByTypeID(c.Context(), signTypeID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get sign type",
			MessageTH:  "ไม่สามารถดึงประเภทงานได้",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Sign type not found",
			MessageTH:  "ไม่พบบประเภทงาน",
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

// @Summary Update Sign Type by ID
// @Description Update a sign type by its ID
// @Tags SignType
// @Accept json
// @Produce json
// @Param id path string true "Sign Type ID"
// @Param request body dto.UpdateSignTypeDTO true "Update Sign Type"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-type/{id} [put]
func (h *SignTypeHandler) UpdateSignTypeByID(c *fiber.Ctx) error {
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
	signTypeID := c.Params("id")
	var body dto.UpdateSignTypeDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid payload",
			MessageTH:  "ข้อมูลไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	errOnUpdate := h.svc.UpdateSignTypeByTypeID(c.Context(), signTypeID, body, claims)
	statusCode := fiber.StatusOK
	MsgEN := "Updated"
	MsgTH := "อัปเดตแล้ว"

	if errOnUpdate != nil {

		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			statusCode = fiber.StatusNotFound
			MsgEN = "Sign type not found"
			MsgTH = "ไม่พบประเภทงาน"
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

// @Summary Delete Sign Type by ID
// @Description Delete a sign type by its ID
// @Tags SignType
// @Accept json
// @Produce json
// @Param id path string true "Sign Type ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/sign-type/{id} [delete]
func (h *SignTypeHandler) DeleteSignTypeByID(c *fiber.Ctx) error {
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
	signTypeID := c.Params("id")
	err = h.svc.DeleteSignTypeByTypeID(c.Context(), signTypeID, claims)
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
