package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	svc    ports.UserService
	upload ports.UpLoadService
	mdw    *middleware.Middleware // added
}

func NewUserHandler(s ports.UserService, upload ports.UpLoadService, mdw *middleware.Middleware) *UserHandler {
	return &UserHandler{svc: s, upload: upload, mdw: mdw}
}

func (h *UserHandler) UserRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	user := versionOne.Group("user")

	user.Post("/create", h.CreateUser)
	user.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetAllUser)
	user.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetUserByID)
	user.Put("/documents", h.mdw.AuthCookieMiddleware(), h.UpdateDocuments)
	user.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateUserByID)
	user.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteUserByID)

}

// @Summary Create a new user
// @Description ใช้สำหรับสร้างผู้ใช้ใหม่ โดยจะไม่สามารถสร้างผู้ใช้ที่มีบัตรประชาชนซ้ำได้
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "อีเมล"
// @Param password formData string true "รหัสผ่าน"
// @Param title_th formData string false "คำนำหน้าชื่อ (TH)"
// @Param title_en formData string false "คำนำหน้าชื่อ (EN)"
// @Param first_name_th formData string true "ชื่อ (TH)"
// @Param last_name_th formData string true "นามสกุล (TH)"
// @Param first_name_en formData string false "ชื่อ (EN)"
// @Param last_name_en formData string false "นามสกุล (EN)"
// @Param id_card formData string true "เลขบัตรประชาชน (13 หลัก)"
// @Param avatar formData file false "ไฟล์รูปโปรไฟล์"
// @Param phone formData string false "เบอร์โทร"
// @Param employee_code formData string false "รหัสพนักงาน"
// @Param gender formData string false "เพศ"
// @Param birth_date formData string false "วันเกิด (YYYY-MM-DD)"
// @Param hire_date formData string false "วันที่เริ่มงาน (YYYY-MM-DD)"
// @Param position_id formData string false "ตำแหน่งงาน"
// @Param department_id formData string false "แผนก"
// @Param employment_type formData string false "ประเภทพนักงาน"
// @Param address_line1 formData string false "ที่อยู่บรรทัดที่ 1"
// @Param address_line2 formData string false "ที่อยู่บรรทัดที่ 2"
// @Param subdistrict formData string false "ตำบล/แขวง"
// @Param district formData string false "อำเภอ/เขต"
// @Param province formData string false "จังหวัด"
// @Param postal_code formData string false "รหัสไปรษณีย์"
// @Param country formData string false "ประเทศ"
// @Param bank_name formData string false "ชื่อธนาคาร"
// @Param account_no formData string false "เลขบัญชี"
// @Param account_name formData string false "ชื่อบัญชี"
// @Success 201 {object} dto.BaseSuccess201ResponseSwagger
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/create [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user dto.RequestCreateUser
	if _, err := c.MultipartForm(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid form-data payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.TitleTH = c.FormValue("title_th")
	user.TitleEN = c.FormValue("title_en")
	user.FirstNameTH = c.FormValue("first_name_th")
	user.LastNameTH = c.FormValue("last_name_th")
	user.FirstNameEN = c.FormValue("first_name_en")
	user.LastNameEN = c.FormValue("last_name_en")
	user.IDCard = c.FormValue("id_card")

	checkIDCard := util.ValidateThaiID(user.IDCard)
	if !checkIDCard {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid ID card format",
			MessageTH:  "รูปแบบบัตรประชาชนไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		// ไม่มีไฟล์แนบมา → ข้ามอัปโหลด
		if errors.Is(err, fasthttp.ErrMissingFile) {
			user.Avatar = ""
		} else {
			// error อื่นๆ ที่ควรแจ้งกลับ
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "Failed to parse uploaded file: " + err.Error(),
				MessageTH:  "ไม่สามารถอัปโหลดไฟล์ได้ กรุณาอัปโหลดไฟล์รูปภาพ",
				Status:     "error",
				Data:       nil,
			})
		}
	} else {
		// มีการแนบไฟล์มา → ดำเนินการอัปโหลด
		if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "Failed to create temporary directory: " + err.Error(),
				MessageTH:  "ไม่สามารถสร้างโฟลเดอร์ชั่วคราวได้",
				Status:     "error",
				Data:       nil,
			})
		}

		tempFilePath := fmt.Sprintf("./tmp/%s", fileHeader.Filename)
		if err := c.SaveFile(fileHeader, tempFilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "Failed to save uploaded file: " + err.Error(),
				MessageTH:  "ไม่สามารถบันทึกไฟล์ที่อัปโหลดได้",
				Status:     "error",
				Data:       nil,
			})
		}
		defer os.Remove(tempFilePath)

		ext := filepath.Ext(fileHeader.Filename)
		uuid := uuid.New().String()
		newName := fmt.Sprintf("%s/%s%s", "avatars", uuid, ext)

		errOnUpload := h.upload.UploadFileCloudflare(c.Context(), tempFilePath, newName)
		if errOnUpload != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "Failed to upload file to storage: " + errOnUpload.Error(),
				MessageTH:  "ไม่สามารถอัปโหลดไฟล์ไปยังที่เก็บข้อมูลได้",
				Status:     "error",
				Data:       nil,
			})
		}

		url, errOnGetURL := h.upload.GetFileURLCloudflare(c.Context(), dto.RequestGetFile{
			Folder: "avatars",
			File:   uuid + ext,
		})
		if errOnGetURL != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "Failed to get file URL: " + errOnGetURL.Error(),
				MessageTH:  "ไม่สามารถดึง URL ของไฟล์ได้",
				Status:     "error",
				Data:       nil,
			})
		}

		user.Avatar = url
	}
	user.Phone = c.FormValue("phone")
	user.EmployeeCode = c.FormValue("employee_code")
	user.Gender = c.FormValue("gender")

	if birthDate := c.FormValue("birth_date"); birthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", birthDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "Invalid birth_date format",
				MessageTH:  "รูปแบบวันเกิดไม่ถูกต้อง",
				Status:     "error",
				Data:       nil,
			})
		}
		user.BirthDate = parsedDate
	}
	if hireDate := c.FormValue("hire_date"); hireDate != "" {
		parsedDate, err := time.Parse("2006-01-02", hireDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "Invalid hire_date format",
				MessageTH:  "รูปแบบวันที่เริ่มงานไม่ถูกต้อง",
				Status:     "error",
				Data:       nil,
			})
		}
		user.HireDate = parsedDate
	}
	user.PositionID = c.FormValue("position_id")
	user.DepartmentID = c.FormValue("department_id")
	user.EmploymentType = c.FormValue("employment_type")

	user.Address = dto.Address{
		AddressLine1: c.FormValue("address_line1"),
		AddressLine2: c.FormValue("address_line2"),
		Subdistrict:  c.FormValue("subdistrict"),
		District:     c.FormValue("district"),
		Province:     c.FormValue("province"),
		PostalCode:   c.FormValue("postal_code"),
		Country:      c.FormValue("country"),
	}

	user.BankInfo = dto.BankInfo{
		BankName:    c.FormValue("bank_name"),
		AccountNo:   c.FormValue("account_no"),
		AccountName: c.FormValue("account_name"),
	}

	errOnCreateUser := h.svc.Create(context.Background(), user)
	if errOnCreateUser != nil {

		if strings.Contains(errOnCreateUser.Error(), "user with ID card") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusBadRequest,
				MessageEN:  "User with this ID card already exists: " + errOnCreateUser.Error(),
				MessageTH:  "มีผู้ใช้ที่มีบัตรประชาชนนี้อยู่แล้ว",
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
// @Param search query string false "Search first_name_th last_name_th first_name_en last_name_en"
// @Param status query string false "Filter by user status (e.g., pending, approved, rejected)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc|desc)"
// @Success 200 {object} dto.BaseSuccessPaginationResponseSwagger
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/list [get]
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
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload: " + err.Error(),
			MessageTH:  "ข้อมูลที่ส่งไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	updatedUser, err := h.svc.UpdateUserByID(context.Background(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update user: " + err.Error(),
			MessageTH:  "ไม่สามารถอัปเดตผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
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
				MessageEN:  "User not found" + err.Error(),
				MessageTH:  "ไม่พบผู้ใช้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to delete user" + err.Error(),
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

// @Summary Update user documents
// @Description ใช้สำหรับอัปเดตเอกสารของผู้ใช้
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "User ID"
// @Param type formData string true "Document type (avatars = รูปโปรไฟล์ , idcards = หลักฐานสำเนาบัตรประชาชน, graduation = หลักฐานการจบการศึกษา, transcript = หลักฐานการศึกษา, resume = หลักฐานการสมัครงาน, health = หลักฐานการตรวจสุขภาพ, military = หลักฐานการผ่านการเกณฑ์ทหาร, criminal = หลักฐานการตรวจประวัติอาชญากรรม, other = โฟลเดอร์อัปโหลดทั่วไป)"
// @Param file formData file true "Document file"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseError400ResponseSwagger
// @Failure 500 {object} dto.BaseError500ResponseSwagger
// @Router /v1/user/documents [put]
func (h *UserHandler) UpdateDocuments(c *fiber.Ctx) error {

	var req dto.RequestUpdateDocuments

	req.UserID = c.FormValue("user_id")
	req.Type = c.FormValue("type")
	allowedTypes := map[string]bool{
		"avatars":    true,
		"idcards":    true,
		"graduation": true,
		"transcript": true,
		"resume":     true,
		"health":     true,
		"military":   true,
		"criminal":   true,
		"other":      true,
	}
	if !allowedTypes[req.Type] {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid document type",
			MessageTH:  "ประเภทเอกสารไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Failed to parse uploaded file" + err.Error(),
			MessageTH:  "ไม่สามารถแยกไฟล์ที่อัปโหลดได้",
			Status:     "error",
			Data:       nil,
		})
	}

	// Ensure the ./tmp/ directory exists
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create temporary directory: " + err.Error(),
			MessageTH:  "ไม่สามารถสร้างโฟลเดอร์ชั่วคราวได้",
			Status:     "error",
			Data:       nil,
		})
	}

	tempFilePath := fmt.Sprintf("./tmp/%s", fileHeader.Filename)
	if err := c.SaveFile(fileHeader, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to save uploaded file" + err.Error(),
			MessageTH:  "ไม่สามารถบันทึกไฟล์ที่อัปโหลดได้",
			Status:     "error",
			Data:       nil,
		})
	}
	defer os.Remove(tempFilePath) // Clean up the temporary file

	ext := filepath.Ext(fileHeader.Filename) // .jpg .png .pdf
	uuid := uuid.New().String()
	newName := fmt.Sprintf("%s/%s%s", req.Type, uuid, ext)

	errOnUpload := h.upload.UploadFileCloudflare(c.Context(), tempFilePath, newName)
	if errOnUpload != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to upload file to storage: " + errOnUpload.Error(),
			MessageTH:  "ไม่สามารถอัปโหลดไฟล์ไปยังที่เก็บข้อมูลได้",
			Status:     "error",
			Data:       nil,
		})
	}

	url, errOnGetURL := h.upload.GetFileURLCloudflare(c.Context(), dto.RequestGetFile{
		Folder: req.Type,
		File:   uuid + ext,
	})
	if errOnGetURL != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get file URL: " + errOnGetURL.Error(),
			MessageTH:  "ไม่สามารถดึง URL ของไฟล์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	req.Name = uuid + ext
	req.FileURL = url

	_, errOnUpdate := h.svc.UpdateDocuments(context.Background(), req)
	if errOnUpdate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update documents" + errOnUpdate.Error(),
			MessageTH:  "ไม่สามารถอัปเดตเอกสารได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Documents updated successfully",
		MessageTH:  "อัปเดตเอกสารสำเร็จ",
		Status:     "success",
		Data:       url,
	})
}
