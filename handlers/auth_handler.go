package handlers

import (
	"time"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	svc ports.AuthService
}

func NewAuthHandler(s ports.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

func (h *AuthHandler) AuthRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	auth := versionOne.Group("auth")

	auth.Post("/login", h.Login)
	auth.Get("/sessions", h.GetSessions)
	auth.Post("/logout", h.Logout)
	auth.Post("/reset", h.ResetPassword)
	auth.Post("/confirm-reset", h.ConfirmResetPassword)

}

// @Summary User login
// @Description ใช้สำหรับเข้าสู่ระบบผู้ใช้
// @Tags auth
// @Accept json
// @Produce json
// @Param auth body dto.RequestLogin true "User login payload"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var auth dto.RequestLogin
	if err := c.BodyParser(&auth); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	token, errOnGetToken := h.svc.Login(c.Context(), auth)
	if errOnGetToken != nil {

		if errOnGetToken == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Invalid email or password",
				MessageTH:  "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
				Status:     "error",
				Data:       nil,
			})
		}

		if errOnGetToken.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Invalid email or password",
				MessageTH:  "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Login failed: " + errOnGetToken.Error(),
			MessageTH:  "การเข้าสู่ระบบล้มเหลว: " + errOnGetToken.Error(),
			Status:     "error",
			Data:       nil,
		})

	}

	util.SetSessionCookie(c, "auth_token", token, 50000*time.Second) // 5 minutes expiration

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Login successful",
		MessageTH:  "เข้าสู่ระบบสำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Get user sessions
// @Description ใช้สำหรับดึงข้อมูลคุกกี้ auth token ของผู้ใช้
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseError401ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/auth/sessions [get]
func (h *AuthHandler) GetSessions(c *fiber.Ctx) error {
	cookie := c.Cookies("auth_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "No auth token cookie found",
			MessageTH:  "ไม่พบคุกกี้ auth token",
			Status:     "error",
			Data:       nil,
		})
	}

	claims, err := h.svc.GetSessions(c.Context(), cookie)
	if err != nil {

		if err.Error() == "token expired" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Token expired",
				MessageTH:  "คุกกี้หมดอายุ",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Invalid token",
			MessageTH:  "คุกกี้ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Sessions retrieved successfully",
		MessageTH:  "ดึงข้อมูลเซสชันสำเร็จ",
		Status:     "success",
		Data:       claims,
	})
}

// @Summary User logout
// @Description ใช้สำหรับออกจากระบบผู้ใช้
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {

	cookie := c.Cookies("auth_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "No auth token cookie found",
			MessageTH:  "ไม่พบคุกกี้ auth token",
			Status:     "error",
			Data:       nil,
		})
	}

	util.DeleteCookie(c, "auth_token")

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Logout successful",
		MessageTH:  "ออกจากระบบสำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Reset user password
// @Description ใช้สำหรับรีเซ็ตรหัสผ่านของผู้ใช้
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body dto.RequestResetPassword true "Reset password payload"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/auth/reset [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var reset dto.RequestResetPassword
	if err := c.BodyParser(&reset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err := h.svc.ResetPassword(c.Context(), reset)
	if err != nil {

		if err.Error() == mongo.ErrNoDocuments.Error() {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "User not found",
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to reset password: " + err.Error(),
			MessageTH:  "ไม่สามารถรีเซ็ตรหัสผ่านได้: " + err.Error(),
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Password reset successfully",
		MessageTH:  "รีเซ็ตรหัสผ่านสำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Confirm reset password
// @Description ใช้สำหรับยืนยันการรีเซ็ตรหัสผ่านของผู้ใช้ token จะหมดอายุภายใน 15 นาที
// @Tags auth
// @Accept json
// @Produce json
// @Param confirm body dto.RequestConfirmResetPassword true "Confirm reset password payload"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 401 {object} dto.BaseError401ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/auth/confirm-reset [post]
func (h *AuthHandler) ConfirmResetPassword(c *fiber.Ctx) error {
	var req dto.RequestConfirmResetPassword
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err := h.svc.ConfirmResetPassword(c.Context(), req)
	if err != nil {
		if err.Error() == "token expired" || err.Error() == "invalid token" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Invalid or expired token",
				MessageTH:  "โทเค็นไม่ถูกต้องหรือหมดอายุ",
				Status:     "error",
				Data:       nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to confirm reset password: " + err.Error(),
			MessageTH:  "ยืนยันการรีเซ็ตรหัสผ่านล้มเหลว: " + err.Error(),
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Password has been reset successfully",
		MessageTH:  "รีเซ็ตรหัสผ่านสำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
