package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	svc ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) UserRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	user := versionOne.Group("user")

	user.Post("/", h.CreateUser)
	user.Get("/", h.GetAllUser)
	user.Get("/:id", h.GetUserByID)
	user.Put("/:id", h.UpdateUserByID)
	user.Delete("/:id", h.DeleteUserByID)

}

// @Summary Create a new user
// @Description ใช้สำหรับสร้างผู้ใช้ใหม่ โดยจะไม่สามารถสร้างผู้ใช้ที่มีบัตรประชาชนซ้ำได้
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.RequestCreateUser true "User create payload"
// @Success 201 {object} dto.BaseSuccess201ResponseSwagger
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user dto.RequestCreateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	err := h.svc.Create(context.Background(), user)
	if err != nil {

		if strings.Contains(err.Error(), "user with ID card") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "User with this ID card already exists",
				MessageTH:  "มีผู้ใช้ที่มีบัตรประชาชนนี้อยู่แล้ว",
				Status:     "error",
				Data:       nil,
			})
		}

		if strings.Contains(err.Error(), "invalid ID card format") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "Invalid ID card format",
				MessageTH:  "รูปแบบบัตรประชาชนไม่ถูกต้อง",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถสร้างผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "User created successfully",
		MessageTH:  "สร้างผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get all users
// @Description ใช้สำหรับดึงรายการผู้ใช้งานแบบแบ่งหน้า
// @Tags user
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc|desc)"
// @Success 200 {object} dto.BaseSuccessPaginationResponseSwagger
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user [get]
func (h *UserHandler) GetAllUser(c *fiber.Ctx) error {

	var req dto.RequestGetUserAll
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	users, err := h.svc.GetAll(context.Background(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve users",
			MessageTH:  "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Users retrieved successfully",
		MessageTH:  "ดึงข้อมูลผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       users,
	})
}

// @Summary Get user by ID
// @Description ใช้สำหรับดึงข้อมูลผู้ใช้ตาม ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "ID is required",
			MessageTH:  "ต้องระบุ ID",
			Status:     "error",
			Data:       nil,
		})
	}

	user, err := h.svc.GetByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to retrieve user",
			MessageTH:  "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "User not found",
			MessageTH:  "ไม่พบผู้ใช้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User retrieved successfully",
		MessageTH:  "ดึงข้อมูลผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       user,
	})
}

// @Summary Update user by ID
// @Description ใช้สำหรับอัปเดตข้อมูลผู้ใช้ตาม ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.RequestUpdateUser true "User update payload"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/{id} [put]
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "ID is required",
			MessageTH:  "ต้องระบุ ID",
			Status:     "error",
			Data:       nil,
		})
	}

	var req dto.RequestUpdateUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	updatedUser, err := h.svc.UpdateUserByID(context.Background(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User updated successfully",
		MessageTH:  "อัปเดตผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       updatedUser,
	})
}

// @Summary Delete user by ID
// @Description ใช้สำหรับลบผู้ใช้ตาม ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/{id} [delete]
func (h *UserHandler) DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "ID is required",
			MessageTH:  "ต้องระบุ ID",
			Status:     "error",
			Data:       nil,
		})
	}

	err := h.svc.DeleteUserByID(context.Background(), id)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete user",
			MessageTH:  "ไม่สามารถลบผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "User deleted successfully",
		MessageTH:  "ลบผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
