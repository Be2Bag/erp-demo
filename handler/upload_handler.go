package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpLoadHandler struct {
	svc ports.UpLoadService
}

func NewUpLoadHandler(s ports.UpLoadService) *UpLoadHandler {
	return &UpLoadHandler{svc: s}
}

func (h *UpLoadHandler) UpLoadRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	upload := versionOne.Group("upload")
	upload.Post("/file", h.Upload)
	upload.Get("/list/:key", h.GetListFile)
	upload.Get("/file", h.GetFile)
}

func (h *UpLoadHandler) Upload(c *fiber.Ctx) error {

	// ดึงพารามิเตอร์โฟลเดอร์จาก query เช่น ?folder=avatars
	folder := c.Query("folder")
	if folder == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Folder parameter is required",
			MessageTH:  "ต้องระบุพารามิเตอร์ folder",
			Status:     "error",
			Data:       nil,
		})
	}
	// ตรวจสอบว่า folder อยู่ในรายการที่อนุญาต
	allowedFolders := map[string]bool{
		"avatars":    true, //รูปโปรไฟล์
		"idcards":    true, //หลักฐานสำเนาบัตรประชาชน
		"graduation": true, //หลักฐานการจบการศึกษา
		"transcript": true, //หลักฐานการศึกษา
		"resume":     true, //หลักฐานการสมัครงาน
		"health":     true, //หลักฐานการตรวจสุขภาพ
		"military":   true, //หลักฐานการผ่านการเกณฑ์ทหาร
		"criminal":   true, //หลักฐานการตรวจประวัติอาชญากรรม
		"other":      true, //โฟลเดอร์อัปโหลดทั่วไป
	}
	if !allowedFolders[folder] {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid folder name",
			MessageTH:  "ชื่อโฟลเดอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	// Parse the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Failed to parse uploaded file",
			MessageTH:  "ไม่สามารถแยกไฟล์ที่อัปโหลดได้",
			Status:     "error",
			Data:       nil,
		})
	}

	// Save the file temporarily
	tempFilePath := fmt.Sprintf("./temp/%s", fileHeader.Filename)
	if err := c.SaveFile(fileHeader, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to save uploaded file",
			MessageTH:  "ไม่สามารถบันทึกไฟล์ที่อัปโหลดได้",
			Status:     "error",
			Data:       nil,
		})
	}
	defer os.Remove(tempFilePath) // Clean up the temporary file

	ext := filepath.Ext(fileHeader.Filename) // .jpg .png .pdf
	uuid := uuid.New().String()              // Generate a unique name for the file
	newName := fmt.Sprintf("%s/%s%s", folder, uuid, ext)

	// Upload the file to storage
	err = h.svc.UploadFile(c.Context(), tempFilePath, newName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to upload file to storage: " + err.Error(),
			MessageTH:  "ไม่สามารถอัปโหลดไฟล์ไปยังที่เก็บข้อมูลได้",
			Status:     "error",
			Data:       nil,
		})
	}

	url, errOnGetURL := h.svc.GetFileURL(c.Context(), dto.RequestGetFile{
		Folder: folder,
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

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "File uploaded successfully",
		MessageTH:  "ไฟล์อัปโหลดสำเร็จ",
		Status:     "success",
		Data:       url,
	})
}

func (h *UpLoadHandler) GetListFile(c *fiber.Ctx) error {

	key := c.Params("key")

	list, err := h.svc.ListFiles(c.Context(), key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to list files: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงรายการไฟล์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "File list retrieved successfully",
		MessageTH:  "ดึงรายการไฟล์สำเร็จ",
		Status:     "success",
		Data:       list,
	})
}

func (h *UpLoadHandler) GetFile(c *fiber.Ctx) error {

	var req dto.RequestGetFile
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid query parameters",
			MessageTH:  "พารามิเตอร์ไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	url, err := h.svc.GetFileURL(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to get file: " + err.Error(),
			MessageTH:  "ไม่สามารถดึงไฟล์ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "File list retrieved successfully",
		MessageTH:  "ดึงรายการไฟล์สำเร็จ",
		Status:     "success",
		Data:       url,
	})
}
