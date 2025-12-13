package handlers

import (
	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type BankAccountsHandler struct {
	svc ports.BankAccountsService
	mdw *middleware.Middleware
}

func NewBankAccountsHandler(s ports.BankAccountsService, mdw *middleware.Middleware) *BankAccountsHandler {
	return &BankAccountsHandler{svc: s, mdw: mdw}
}

func (h *BankAccountsHandler) BankAccountsRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	bankAccounts := versionOne.Group("bank-accounts")

	bankAccounts.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateBankAccount)
	bankAccounts.Get("/list", h.mdw.AuthCookieMiddleware(), h.ListBankAccounts)
	bankAccounts.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetBankAccountByID)
	bankAccounts.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateBankAccountByID)
	bankAccounts.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteBankAccountByID)

}

// @Summary Create Bank Account
// @Description Create a new bank account
// @Tags BankAccounts
// @Accept json
// @Produce json
// @Param bankAccount body dto.CreateBankAccountsDTO true "Bank Account Payload"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/bank-accounts/create [post]
func (h *BankAccountsHandler) CreateBankAccount(c *fiber.Ctx) error {

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

	var bankAccount dto.CreateBankAccountsDTO
	if err := c.BodyParser(&bankAccount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateBankAccount(c.Context(), bankAccount, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create bank account" + err.Error(),
			MessageTH:  "สร้างบัญชีธนาคารไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Bank account created successfully",
		MessageTH:  "สร้างบัญชีธนาคารเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Update Bank Account by ID
// @Description Update a bank account's details by its ID
// @Tags BankAccounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Param bankAccount body dto.UpdateBankAccountsDTO true "Bank Account Update Payload"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/bank-accounts/{id} [put]
func (h *BankAccountsHandler) UpdateBankAccountByID(c *fiber.Ctx) error {
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
	bankID := c.Params("id")
	if bankID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Bank ID is required",
			MessageTH:  "ต้องระบุรหัสบัญชีธนาคาร",
			Status:     "error",
			Data:       nil,
		})
	}
	var updateDTO dto.UpdateBankAccountsDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	err = h.svc.UpdateBankAccountID(c.Context(), bankID, updateDTO, claims)
	if err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "Bank account not found",
				MessageTH:  "ไม่พบบัญชีธนาคาร",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update bank account",
			MessageTH:  "อัปเดตบัญชีธนาคารไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Bank account updated successfully",
		MessageTH:  "อัปเดตบัญชีธนาคารเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary Delete Bank Account by ID
// @Description Soft delete a bank account by its ID
// @Tags BankAccounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/bank-accounts/{id} [delete]
func (h *BankAccountsHandler) DeleteBankAccountByID(c *fiber.Ctx) error {
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
	bankID := c.Params("id")
	if bankID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Bank ID is required",
			MessageTH:  "ต้องระบุรหัสบัญชีธนาคาร",
			Status:     "error",
			Data:       nil,
		})
	}

	if config.IsProtectedBankAccount(bankID) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "This bank account cannot be deleted",
			MessageTH:  "ห้ามลบบัญชีนี้",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.DeleteBankAccountByID(c.Context(), bankID, claims)
	if err != nil {
		if err == fiber.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "Bank account not found",
				MessageTH:  "ไม่พบบัญชีธนาคาร",
				Status:     "error",
				Data:       nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete bank account",
			MessageTH:  "ลบบัญชีธนาคารไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Bank account deleted successfully",
		MessageTH:  "ลบบัญชีธนาคารเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

// @Summary List Bank Accounts
// @Description Retrieve a paginated list of bank accounts with optional search and sorting
// @Tags BankAccounts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search term to filter bank accounts"
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc or desc)" default("desc")
// @Success 200 {object} dto.BaseResponse{data=dto.Pagination}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/bank-accounts/list [get]
func (h *BankAccountsHandler) ListBankAccounts(c *fiber.Ctx) error {
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

	var req dto.RequestListBankAccounts
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

	list, err := h.svc.ListBankAccounts(c.Context(), claims, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list bank accounts" + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการบัญชีธนาคารได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "List bank accounts successfully",
		MessageTH:  "ดึงรายการบัญชีธนาคารสำเร็จ",
		Status:     "success",
		Data:       list,
	})
}

// @Summary Get Bank Account by ID
// @Description Retrieve a bank account by its ID
// @Tags BankAccounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 200 {object} dto.BaseResponse{data=dto.BankAccountsDTO}
// @Failure 400 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/bank-accounts/{id} [get]
func (h *BankAccountsHandler) GetBankAccountByID(c *fiber.Ctx) error {
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
	bankID := c.Params("id")
	item, err := h.svc.GetListBankAccountByBankID(c.Context(), bankID, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get bank account",
			MessageTH:  "ไม่สามารถดึงบัญชีธนาคารได้",
			Status:     "error",
			Data:       nil,
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Bank account not found",
			MessageTH:  "ไม่พบบัญชีธนาคาร",
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
