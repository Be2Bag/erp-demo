package handler

import (
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

	tasks.Get("/", h.mdw.AuthCookieMiddleware(), h.GetTasks)
	tasks.Post("/", h.mdw.AuthCookieMiddleware(), h.CreateTask)
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
	// ฟังก์ชันสำหรับสร้างงานใหม่
	return nil
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
