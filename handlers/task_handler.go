package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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

	tasks.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetListTasks)
	tasks.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateTask)
	tasks.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetTaskByID)
	tasks.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateTask)
	tasks.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteTask)

	tasks.Put("/:id/workflow", h.mdw.AuthCookieMiddleware(), h.UpdateTaskWorkflow)

	tasks.Get("/stats", h.mdw.AuthCookieMiddleware(), h.GetTaskStatistics)
}

// @Summary Get list of tasks
// @Description Get list of tasks with pagination and filtering
// @Tags Tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 10)"
// @Param search query string false "ค้นหาด้วย project_name หรือ job_name"
// @Param department query string false "Dropdown แผนก DPT001: แผนกออกแบบกราฟิก, DPT002: แผนกผลิต, DPT003: แผนกติดตั้ง, DPT004: แผนกบัญชี"
// @Param sort_by query string false "เรียงตาม created_at updated_at project_name job_name  (ค่าเริ่มต้น: created_at)"
// @Param sort_order query string false "เรียงลำดับ (asc เก่า→ใหม่ | desc ใหม่→เก่า (ค่าเริ่มต้น))"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/list [get]
func (h *TaskHandler) GetListTasks(c *fiber.Ctx) error {

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

	var req dto.RequestListWorkflow
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

	tasks, errOnGetTasks := h.svc.GetListTasks(c.Context(), claims, req.Page, req.Limit, req.Search, req.Department, req.SortBy, req.SortOrder)
	if errOnGetTasks != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  errOnGetTasks.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "OK",
		MessageTH:  "สำเร็จ",
		Status:     "success",
		Data:       tasks,
	})
}

// @Summary Create a new task
// @Description Create a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body dto.CreateTaskRequest true "Create Task Request"
// @Success 201 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 401 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/create [post]
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

// @Summary Get a task by ID
// @Description Get a task by ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "id is required",
			MessageTH:  "ต้องระบุรหัส",
			Status:     "error",
		})
	}

	task, err := h.svc.GetTaskByID(c.Context(), id)
	if err != nil {
		status := fiber.StatusInternalServerError
		msgEN := "Failed to get Task"
		msgTH := "ไม่สามารถดึงข้อมูลงานได้"
		if err == mongo.ErrNoDocuments {
			status = fiber.StatusNotFound
			msgEN = "Task not found"
			msgTH = "ไม่พบข้อมูลงาน"
		}
		return c.Status(status).JSON(dto.BaseResponse{
			StatusCode: status,
			MessageEN:  msgEN,
			MessageTH:  msgTH,
			Status:     "error",
		})
	}
	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "Task not found",
			MessageTH:  "ไม่พบข้อมูลงาน",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Task retrieved successfully",
		MessageTH:  "ดึงข้อมูลงานสำเร็จ",
		Status:     "success",
		Data:       task,
	})
}

// @Summary Update a task
// @Description Update a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body dto.UpdateTaskRequest true "Update Task Request"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
			Data:       nil,
		})
	}
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}
	errOnUpdate := h.svc.UpdateTask(c.Context(), id, req, claims.UserID)
	if errOnUpdate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  errOnUpdate.Error(),
			MessageTH:  "ไม่สามารถอัปเดต Task ได้",
			Status:     "error",
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Task updated successfully",
		MessageTH:  "อัปเดต Task สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
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
