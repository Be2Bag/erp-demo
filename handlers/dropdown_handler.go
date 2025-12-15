package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/middleware"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type DropDownHandler struct {
	svc ports.DropDownService
	mdw *middleware.Middleware
}

func NewDropDownHandler(s ports.DropDownService, m *middleware.Middleware) *DropDownHandler {
	return &DropDownHandler{svc: s, mdw: m}
}

func (h *DropDownHandler) DropDownRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	dropdown := versionOne.Group("dropdown")

	dropdown.Get("/department", h.GetDepartment)
	dropdown.Get("/project", h.mdw.AuthCookieMiddleware(), h.GetProject)
	dropdown.Get("/user", h.mdw.AuthCookieMiddleware(), h.GetUserAll)
	dropdown.Get("/user/:id", h.mdw.AuthCookieMiddleware(), h.GetUser)
	dropdown.Get("/province", h.GetProvince)
	dropdown.Get("/sign-type", h.GetSignType)
	dropdown.Get("/sign-job-list/:id", h.mdw.AuthCookieMiddleware(), h.GetSignJobList)
	dropdown.Get("/position/:id", h.GetPosition)
	dropdown.Get("/district/:id", h.GetDistrict)
	dropdown.Get("/subdistrict/:id", h.GetSubDistrict)
	dropdown.Get("/customer-type", h.GetCustomerTypes)
	dropdown.Get("/kpi/:id", h.mdw.AuthCookieMiddleware(), h.GetKPI)
	dropdown.Get("/workflow/:id", h.mdw.AuthCookieMiddleware(), h.GetWorkflow)
	dropdown.Get("/category", h.mdw.AuthCookieMiddleware(), h.GetCategory)
	dropdown.Get("/transaction-category/:types", h.mdw.AuthCookieMiddleware(), h.GetTransactionCategory)
	dropdown.Get("/bank-accounts", h.mdw.AuthCookieMiddleware(), h.GetBankAccountsCategory)
}

// @Summary Get all positions
// @Description ใช้สำหรับดึงข้อมูลตำแหน่งงานทั้งหมด
// @Tags Dropdown
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetPositions}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/position/{id} [get]
func (h *DropDownHandler) GetPosition(c *fiber.Ctx) error {
	department_ID := c.Params("id")

	positions, errOnGetPositions := h.svc.GetPositions(c.Context(), department_ID)
	if errOnGetPositions != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงตำแหน่งงานได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(positions) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No positions found",
			MessageTH:  "ไม่พบข้อมูลตำแหน่งงาน",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get positions successfully",
		MessageTH:  "ดึงตำแหน่งงานสำเร็จ",
		Status:     "success",
		Data:       positions,
	})
}

// @Summary Get all departments
// @Description ใช้สำหรับดึงข้อมูลแผนกทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetDepartments}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/department [get]
func (h *DropDownHandler) GetDepartment(c *fiber.Ctx) error {
	departments, errOnGetDepartments := h.svc.GetDepartments(c.Context())
	if errOnGetDepartments != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงแผนกได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(departments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No departments found",
			MessageTH:  "ไม่พบข้อมูลแผนก",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get departments successfully",
		MessageTH:  "ดึงแผนกสำเร็จ",
		Status:     "success",
		Data:       departments,
	})
}

// @Summary Get all provinces
// @Description ใช้สำหรับดึงข้อมูลจังหวัดทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetProvinces}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/province [get]
func (h *DropDownHandler) GetProvince(c *fiber.Ctx) error {
	provinces, errOnGetProvinces := h.svc.GetProvinces(c.Context())
	if errOnGetProvinces != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงจังหวัดได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(provinces) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No provinces found",
			MessageTH:  "ไม่พบข้อมูลจังหวัด",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get provinces successfully",
		MessageTH:  "ดึงจังหวัดสำเร็จ",
		Status:     "success",
		Data:       provinces,
	})
}

// @Summary Get all districts by province ID
// @Description ใช้สำหรับดึงข้อมูลอำเภอทั้งหมดตามรหัสจังหวัด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "Province ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetDistricts}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/district/{id} [get]
func (h *DropDownHandler) GetDistrict(c *fiber.Ctx) error {
	provinceID := c.Params("id")
	districts, errOnGetDistricts := h.svc.GetDistricts(c.Context(), provinceID)
	if errOnGetDistricts != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงอำเภอได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(districts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No districts found",
			MessageTH:  "ไม่พบข้อมูลอำเภอ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get districts successfully",
		MessageTH:  "ดึงอำเภอสำเร็จ",
		Status:     "success",
		Data:       districts,
	})
}

// @Summary Get all sub-districts by district ID
// @Description ใช้สำหรับดึงข้อมูลตำบลทั้งหมดตามรหัสอำเภอ
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "District ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetSubDistricts}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/subdistrict/{id} [get]
func (h *DropDownHandler) GetSubDistrict(c *fiber.Ctx) error {
	districtID := c.Params("id")
	subDistricts, errOnGetSubDistricts := h.svc.GetSubDistricts(c.Context(), districtID)
	if errOnGetSubDistricts != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงตำบลได้",
			Status:     "error",
			Data:       nil,
		})
	}
	if len(subDistricts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No sub-districts found",
			MessageTH:  "ไม่พบข้อมูลตำบล",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get sub-districts successfully",
		MessageTH:  "ดึงตำบลสำเร็จ",
		Status:     "success",
		Data:       subDistricts,
	})
}

// @Summary Get all sign types
// @Description ใช้สำหรับดึงข้อมูลประเภทงานทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetSignTypes}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/sign-type [get]
func (h *DropDownHandler) GetSignType(c *fiber.Ctx) error {
	signTypes, errOnGetSignTypes := h.svc.GetSignTypes(c.Context())
	if errOnGetSignTypes != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงประเภทงานได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(signTypes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No sign types found",
			MessageTH:  "ไม่พบข้อมูลประเภทงาน",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get sign types successfully",
		MessageTH:  "ดึงประเภทงานสำเร็จ",
		Status:     "success",
		Data:       signTypes,
	})
}

// @Summary Get all customer types
// @Description ใช้สำหรับดึงข้อมูลประเภทลูกค้าทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetCustomerTypes}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/customer-type [get]
func (h *DropDownHandler) GetCustomerTypes(c *fiber.Ctx) error {
	customerTypes, errOnGetCustomerTypes := h.svc.GetCustomerTypes(c.Context())
	if errOnGetCustomerTypes != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงประเภทลูกค้าได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(customerTypes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No customer types found",
			MessageTH:  "ไม่พบข้อมูลประเภทลูกค้า",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get customer types successfully",
		MessageTH:  "ดึงประเภทลูกค้าสำเร็จ",
		Status:     "success",
		Data:       customerTypes,
	})
}

// @Summary Get all sign jobs
// @Description ใช้สำหรับดึงข้อมูลใบงานทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetSignList}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/sign-job-list/{id} [get]
func (h *DropDownHandler) GetSignJobList(c *fiber.Ctx) error {

	projectID := c.Params("id")
	signJobs, errOnGetSignJobs := h.svc.GetSignJobList(c.Context(), projectID)
	if errOnGetSignJobs != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงใบงานได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(signJobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No sign jobs found",
			MessageTH:  "ไม่พบข้อมูลใบงาน",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get sign jobs successfully",
		MessageTH:  "ดึงใบงานสำเร็จ",
		Status:     "success",
		Data:       signJobs,
	})
}

// @Summary Get all projects
// @Description ใช้สำหรับดึงข้อมูลโครงการทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetProjects}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/project [get]
func (h *DropDownHandler) GetProject(c *fiber.Ctx) error {
	projects, errOnGetProjects := h.svc.GetProjectList(c.Context())
	if errOnGetProjects != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงข้อมูลโครงการได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No projects found",
			MessageTH:  "ไม่พบข้อมูลโครงการ",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get projects successfully",
		MessageTH:  "ดึงข้อมูลโครงการสำเร็จ",
		Status:     "success",
		Data:       projects,
	})
}

// @Summary Get all users
// @Description ใช้สำหรับดึงข้อมูลผู้ใช้ทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetUsers}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/user/{id} [get]
func (h *DropDownHandler) GetUser(c *fiber.Ctx) error {

	departmentID := c.Params("id")
	users, errOnGetUsers := h.svc.GetUserList(c.Context(), departmentID)
	if errOnGetUsers != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No users found",
			MessageTH:  "ไม่พบข้อมูลผู้ใช้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get users successfully",
		MessageTH:  "ดึงข้อมูลผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       users,
	})
}

// @Summary Get all KPIs
// @Description ใช้สำหรับดึงข้อมูล KPI ทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.KPITemplateDTO}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/kpi/{id} [get]
func (h *DropDownHandler) GetKPI(c *fiber.Ctx) error {

	departmentID := c.Params("id")

	kpis, errOnGetKPI := h.svc.GetKPI(c.Context(), departmentID)
	if errOnGetKPI != nil {

		if errOnGetKPI == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "No KPIs found",
				MessageTH:  "ไม่สามารถดึงข้อมูล KPI ได้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  errOnGetKPI.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูล KPI ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(kpis) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No KPIs found",
			MessageTH:  "ไม่พบข้อมูล KPI",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get KPI successfully",
		MessageTH:  "ดึงข้อมูล KPI สำเร็จ",
		Status:     "success",
		Data:       kpis,
	})
}

// @Summary Get all workflows
// @Description ใช้สำหรับดึงข้อมูล Workflow ทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetWorkflows}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/workflow/{id} [get]
func (h *DropDownHandler) GetWorkflow(c *fiber.Ctx) error {

	departmentID := c.Params("id")

	workflows, errOnGetWorkflows := h.svc.GetWorkflows(c.Context(), departmentID)
	if errOnGetWorkflows != nil {

		if errOnGetWorkflows == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "No workflows found",
				MessageTH:  "ไม่สามารถดึงข้อมูล Workflow ได้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงข้อมูล Workflow ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(workflows) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No workflows found",
			MessageTH:  "ไม่พบข้อมูล Workflow",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get workflows successfully",
		MessageTH:  "ดึงข้อมูล Workflow สำเร็จ",
		Status:     "success",
		Data:       workflows,
	})
}

// @Summary Get all categories
// @Description ใช้สำหรับดึงข้อมูลหมวดหมู่ทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.CategoryDTO}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/category [get]
func (h *DropDownHandler) GetCategory(c *fiber.Ctx) error {

	categorys, errOnGetCategorys := h.svc.GetCategorys(c.Context())
	if errOnGetCategorys != nil {

		if errOnGetCategorys == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "No categories found",
				MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  errOnGetCategorys.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(categorys) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No categories found",
			MessageTH:  "ไม่พบข้อมูลหมวดหมู่",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get categories successfully",
		MessageTH:  "ดึงข้อมูลหมวดหมู่สำเร็จ",
		Status:     "success",
		Data:       categorys,
	})
}

// @Summary Get all users
// @Description ใช้สำหรับดึงข้อมูลผู้ใช้ทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetUsers}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/user [get]
func (h *DropDownHandler) GetUserAll(c *fiber.Ctx) error {

	users, errOnGetUsers := h.svc.GetUserListAll(c.Context())
	if errOnGetUsers != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  fiber.ErrBadGateway.Message,
			MessageTH:  "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No users found",
			MessageTH:  "ไม่พบข้อมูลผู้ใช้",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get users successfully",
		MessageTH:  "ดึงข้อมูลผู้ใช้สำเร็จ",
		Status:     "success",
		Data:       users,
	})
}

// @Summary Get all transaction categories
// @Description ใช้สำหรับดึงข้อมูลหมวดหมู่รายการเคลื่อนไหวทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Param types path string true "Types"
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetTransactionCategorys}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/transaction-category/{types} [get]
func (h *DropDownHandler) GetTransactionCategory(c *fiber.Ctx) error {

	types := c.Params("types")

	transactionCategorys, errOnGetTransactionCategorys := h.svc.GetTransactionCategory(c.Context(), types)
	if errOnGetTransactionCategorys != nil {

		if errOnGetTransactionCategorys == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "No transaction categories found",
				MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  errOnGetTransactionCategorys.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลหมวดหมู่ได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(transactionCategorys) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No transaction categories found",
			MessageTH:  "ไม่พบข้อมูลหมวดหมู่",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get transaction categories successfully",
		MessageTH:  "ดึงข้อมูลหมวดหมู่รายการสำเร็จ",
		Status:     "success",
		Data:       transactionCategorys,
	})
}

// @Summary Get all bank accounts
// @Description ใช้สำหรับดึงข้อมูลบัญชีธนาคารทั้งหมด
// @Tags Dropdown
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]dto.ResponseGetBankAccounts}
// @Failure 502 {object} dto.BaseResponse
// @Failure 500 {object} dto.BaseResponse
// @Router /v1/dropdown/bank-accounts [get]
func (h *DropDownHandler) GetBankAccountsCategory(c *fiber.Ctx) error {

	bankAccountsCategorys, errOnGetBankAccountsCategorys := h.svc.GetBankAccountsList(c.Context())
	if errOnGetBankAccountsCategorys != nil {

		if errOnGetBankAccountsCategorys == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusNotFound,
				MessageEN:  "No bank accounts found",
				MessageTH:  "ไม่สามารถดึงข้อมูลบัญชีธนาคารได้",
				Status:     "error",
				Data:       nil,
			})
		}

		return c.Status(fiber.ErrBadGateway.Code).JSON(dto.BaseResponse{
			StatusCode: fiber.ErrBadGateway.Code,
			MessageEN:  errOnGetBankAccountsCategorys.Error(),
			MessageTH:  "ไม่สามารถดึงข้อมูลบัญชีธนาคารได้",
			Status:     "error",
			Data:       nil,
		})
	}

	if len(bankAccountsCategorys) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.BaseResponse{
			StatusCode: fiber.StatusNotFound,
			MessageEN:  "No bank accounts categories found",
			MessageTH:  "ไม่พบข้อมูลหมวดหมู่บัญชีธนาคาร",
			Status:     "error",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get bank accounts categories successfully",
		MessageTH:  "ดึงข้อมูลหมวดหมู่บัญชีธนาคารสำเร็จ",
		Status:     "success",
		Data:       bankAccountsCategorys,
	})
}
