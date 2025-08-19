package handlers

import (
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
)

type DropDownHandler struct {
	svc ports.DropDownService
}

func NewDropDownHandler(s ports.DropDownService) *DropDownHandler {
	return &DropDownHandler{svc: s}
}

func (h *DropDownHandler) DropDownRoutes(router fiber.Router) {

	versionOne := router.Group("v1")
	dropdown := versionOne.Group("dropdown")

	dropdown.Get("/department", h.GetDepartment)
	dropdown.Get("/project", h.GetProject)
	dropdown.Get("/sign-job-list/:id", h.GetSignJobList)
	dropdown.Get("/province", h.GetProvince)
	dropdown.Get("/sign-type", h.GetSignType)
	dropdown.Get("/position/:id", h.GetPosition)
	dropdown.Get("/district/:id", h.GetDistrict)
	dropdown.Get("/subdistrict/:id", h.GetSubDistrict)
	dropdown.Get("/customer-type", h.GetCustomerTypes)

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

	return c.Status(fiber.StatusOK).JSON(dto.BaseResponse{
		StatusCode: fiber.StatusOK,
		MessageEN:  "Get projects successfully",
		MessageTH:  "ดึงข้อมูลโครงการสำเร็จ",
		Status:     "success",
		Data:       projects,
	})
}
