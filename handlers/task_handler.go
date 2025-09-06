package handlers

import (
	"errors"
	"strings"

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
	// versionTwo := router.Group("v2")

	tasks := versionOne.Group("tasks")
	// tasksV2 := versionTwo.Group("tasks")

	tasks.Get("/list", h.mdw.AuthCookieMiddleware(), h.GetListTasks)
	tasks.Post("/create", h.mdw.AuthCookieMiddleware(), h.CreateTask)
	tasks.Get("/:id", h.mdw.AuthCookieMiddleware(), h.GetTaskByID)
	// tasks.Put("/:id", h.mdw.AuthCookieMiddleware(), h.UpdateTask)
	tasks.Put("/:id", h.mdw.AuthCookieMiddleware(), h.PutTaskV2)
	tasks.Delete("/:id", h.mdw.AuthCookieMiddleware(), h.DeleteTask)
	tasks.Put("/:task_id/steps/:step_id", h.mdw.AuthCookieMiddleware(), h.UpdateStepStatusNote)

	// tasksV2.Put("/:id", h.mdw.AuthCookieMiddleware(), h.PutTaskV2)

}

// @Summary Get list of tasks
// @Description Get list of tasks with pagination and filtering
// @Tags Tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 10)"
// @Param search query string false "ค้นหาด้วย project_name หรือ job_name"
// @Param department_id query string false "Dropdown แผนก DPT001: แผนกออกแบบกราฟิก, DPT002: แผนกผลิต, DPT003: แผนกติดตั้ง, DPT004: แผนกบัญชี"
// @Param status query string false "สถานะ (todo|in_progress|done) (ค่าเริ่มต้น: todo, in_progress)"
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

	var req dto.RequestListTask
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

	tasks, errOnGetTasks := h.svc.GetListTasks(c.Context(), claims, req.Page, req.Limit, req.Search, req.Department, req.SortBy, req.SortOrder, req.Status)
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

// // @Summary Update a task
// // @Description Update a task
// // @Tags Tasks
// // @Accept json
// // @Produce json
// // @Param id path string true "Task ID"
// // @Param request body dto.UpdateTaskRequest true "Update Task Request"
// // @Success 200 {object} dto.BaseResponse
// // @Failure 400 {object} dto.BaseResponse
// // @Failure 404 {object} dto.BaseResponse
// // @Failure 500 {object} dto.BaseResponse
// // @Router /v1/tasks/{id} [put]
// func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	var req dto.UpdateTaskRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
// 			StatusCode: fiber.StatusBadRequest,
// 			MessageEN:  "Invalid request body",
// 			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
// 			Status:     "error",
// 			Data:       nil,
// 		})
// 	}
// 	claims, err := middleware.GetClaims(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
// 			StatusCode: fiber.StatusUnauthorized,
// 			MessageEN:  "Unauthorized",
// 			MessageTH:  "ไม่ได้รับอนุญาต",
// 			Status:     "error",
// 		})
// 	}
// 	errOnUpdate := h.svc.UpdateTask(c.Context(), id, req, claims.UserID)
// 	if errOnUpdate != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
// 			StatusCode: fiber.StatusBadRequest,
// 			MessageEN:  errOnUpdate.Error(),
// 			MessageTH:  "ไม่สามารถอัปเดต Task ได้",
// 			Status:     "error",
// 			Data:       nil,
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
// 		StatusCode: fiber.StatusOK,
// 		MessageEN:  "Task updated successfully",
// 		MessageTH:  "อัปเดต Task สำเร็จ",
// 		Status:     "success",
// 		Data:       nil,
// 	})
// }

// @Summary Delete a task
// @Description Delete a task by ID
// @Tags Tasks
// @Accept json
// @Param id path string true "Task ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
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
	id := c.Params("id")
	err = h.svc.DeleteTask(c.Context(), id, claims)
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

// @Summary Update step status
// @Description Update step status by task ID and step ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID"
// @Param step_id path string true "Step ID"
// @Param request body dto.UpdateStepStatusNoteRequest true "Update Step Status Request"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/{task_id}/steps/{step_id} [put]
// handlers/task.go
func (h *TaskHandler) UpdateStepStatusNote(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusUnauthorized,
			MessageEN:  "Unauthorized",
			MessageTH:  "ไม่ได้รับอนุญาต",
			Status:     "error",
		})
	}

	taskID := c.Params("task_id")
	stepID := c.Params("step_id")

	var req dto.UpdateStepStatusNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request payload",
			MessageTH:  "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			Status:     "error",
		})
	}
	// ป้องกันส่งค่าว่างทั้งคู่
	if req.Status == nil && (req.Notes == nil || strings.TrimSpace(*req.Notes) == "") {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Nothing to update",
			MessageTH:  "ไม่มีข้อมูลสำหรับอัปเดต",
			Status:     "error",
		})
	}

	errOnUpdate := h.svc.UpdateStepStatus(c.Context(), taskID, stepID, req, claims)
	if errOnUpdate != nil {
		if errors.Is(errOnUpdate, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "Task or step not found",
				MessageTH:  "ไม่พบนงานหรือขั้นตอน",
				Status:     "error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusInternalServerError,
			MessageEN:  "Failed to update step: " + errOnUpdate.Error(),
			MessageTH:  "อัปเดตขั้นตอนไม่สำเร็จ",
			Status:     "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Step updated",
		MessageTH:  "อัปเดตสำเร็จ",
		Status:     "success",
		Data:       nil, // คืนข้อมูลที่อัปเดตแล้ว (ดู struct ด้านล่าง)
	})
}

// @Summary Update task
// @Description Update task by ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body dto.UpdateTaskPutRequest true "Update Task Request (PUT)"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.BaseResponse
// @Failure 404 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/tasks/{id} [put]
func (h *TaskHandler) PutTaskV2(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateTaskPutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  "Invalid request body",
			MessageTH:  "รูปแบบคำขอไม่ถูกต้อง",
			Status:     "error",
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

	if err := h.svc.ReplaceTask(c.Context(), id, req, claims.UserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusBadRequest,
			MessageEN:  err.Error(),
			MessageTH:  "ไม่สามารถอัปเดต Task (PUT) ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Task replaced successfully",
		MessageTH:  "แทนที่ Task สำเร็จ",
		Status:     "success",
		Data:       nil,
	})
}
