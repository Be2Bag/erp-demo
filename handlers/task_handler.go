package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	svc ports.TaskService
	mdw *middleware.Middleware
}

func NewTaskHandler(s ports.TaskService, mdw *middleware.Middleware) *TaskHandler {
	return &TaskHandler{svc: s, mdw: mdw}
}

func (h *TaskHandler) TaskRoutes(router fiber.Router) {
	versionOne := router.Group("v1")
	tasks := versionOne.Group("tasks")

	tasks.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetTasks)
	tasks.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateTask)
	tasks.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetTaskByID)
	tasks.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateTask)
	tasks.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteTask)

	tasks.Put("/:id/workflow", h.mdw.AuthCookieMiddleware(), h.UpdateTaskWorkflow)

	tasks.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetTaskStatistics)
}

// ตัวจัดการงาน (Task Management Handlers)
func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับดึงข้อมูลงานทั้งหมด
	return nil
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
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

	var createTask dto.CreateTaskRequest
	if err := c.BodyParser(&createTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}

	err = h.svc.CreateTask(c.Context(), createTask, claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to create task" + err.Error(),
			MessageTH:  "สร้างงานไม่สำเร็จ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusCreated,
		MessageEN:  "Task created successfully",
		MessageTH:  "สร้างงานเรียบร้อยแล้ว",
		Status:     "success",
		Data:       nil,
	})
}

func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับดึงข้อมูลงานตามรหัส
	return nil
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับแก้ไขข้อมูลงาน
	return nil
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับลบงาน
	return nil
}

// ตัวจัดการเวิร์กโฟลว์ (Workflow Management Handler)
func (h *TaskHandler) UpdateTaskWorkflow(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับอัปเดตเวิร์กโฟลว์ของงาน
	return nil
}

// ตัวจัดการสถิติงาน (Task Statistics Handler)
func (h *TaskHandler) GetTaskStatistics(c *fiber.Ctx) error {
	// ฟังก์ชันสำหรับดึงข้อมูลสถิติงาน
	return nil
}
