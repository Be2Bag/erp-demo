package services

import (
	"context"
	"fmt"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dropDownService struct {
	config       config.Config
	dropDownRepo ports.DropDownRepository
}

func NewDropDownService(cfg config.Config, dropDownRepo ports.DropDownRepository) ports.DropDownService {
	return &dropDownService{config: cfg, dropDownRepo: dropDownRepo}
}

func (s *dropDownService) GetPositions(ctx context.Context, departmentID string) ([]dto.ResponseGetPositions, error) {

	filter := bson.M{"deleted_at": nil, "department_id": departmentID}
	projection := bson.M{}

	positions, errOnGetPositions := s.dropDownRepo.GetPositions(ctx, filter, projection)
	if errOnGetPositions != nil {
		return nil, errOnGetPositions
	}

	if len(positions) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetPositions
	for _, position := range positions {
		response = append(response, dto.ResponseGetPositions{
			PositionID:   position.PositionID,
			PositionName: position.PositionName + " (" + position.Level + ")",
		})
	}

	return response, nil
}

func (s *dropDownService) GetDepartments(ctx context.Context) ([]dto.ResponseGetDepartments, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	departments, errOnGetDepartments := s.dropDownRepo.GetDepartments(ctx, filter, projection)
	if errOnGetDepartments != nil {
		return nil, errOnGetDepartments
	}

	if len(departments) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetDepartments
	for _, department := range departments {
		response = append(response, dto.ResponseGetDepartments{
			DepartmentID:   department.DepartmentID,
			DepartmentName: department.DepartmentName,
		})
	}

	return response, nil
}

func (s *dropDownService) GetProvinces(ctx context.Context) ([]dto.ResponseGetProvinces, error) {
	filter := bson.M{"deleted_at": nil}
	// if nameTH != "" {
	// 	filter["name_th"] = bson.M{"$regex": nameTH, "$options": "i"}
	// }
	projection := bson.M{}

	provinces, errOnGetProvinces := s.dropDownRepo.GetProvinces(ctx, filter, projection)
	if errOnGetProvinces != nil {
		return nil, errOnGetProvinces
	}

	if len(provinces) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetProvinces
	for _, province := range provinces {
		response = append(response, dto.ResponseGetProvinces{
			ProvinceID:   province.ID,
			ProvinceName: province.NameTH,
		})
	}

	return response, nil
}

func (s *dropDownService) GetDistricts(ctx context.Context, provinceID string) ([]dto.ResponseGetDistricts, error) {
	filter := bson.M{"deleted_at": nil, "province_id": provinceID}
	projection := bson.M{}

	districts, errOnGetDistricts := s.dropDownRepo.GetDistricts(ctx, filter, projection)
	if errOnGetDistricts != nil {
		return nil, errOnGetDistricts
	}

	if len(districts) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetDistricts
	for _, district := range districts {
		response = append(response, dto.ResponseGetDistricts{
			DistrictID:   district.ID,
			DistrictName: district.NameTH,
		})
	}

	return response, nil
}

func (s *dropDownService) GetSubDistricts(ctx context.Context, districtID string) ([]dto.ResponseGetSubDistricts, error) {
	filter := bson.M{"deleted_at": nil, "district_id": districtID}
	projection := bson.M{}

	subDistricts, errOnGetSubDistricts := s.dropDownRepo.GetSubDistricts(ctx, filter, projection)
	if errOnGetSubDistricts != nil {
		return nil, errOnGetSubDistricts
	}

	if len(subDistricts) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetSubDistricts
	for _, subDistrict := range subDistricts {
		response = append(response, dto.ResponseGetSubDistricts{
			SubDistrictID:   subDistrict.ID,
			SubDistrictName: subDistrict.NameTH,
			ZipCode:         subDistrict.ZipCode,
		})
	}

	return response, nil
}

func (s *dropDownService) GetSignTypes(ctx context.Context) ([]dto.ResponseGetSignTypes, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	signTypes, errOnGetSignTypes := s.dropDownRepo.GetSignTypes(ctx, filter, projection)
	if errOnGetSignTypes != nil {
		return nil, errOnGetSignTypes
	}

	if len(signTypes) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetSignTypes
	for _, signType := range signTypes {
		response = append(response, dto.ResponseGetSignTypes{
			TypeID: signType.TypeID,
			NameTH: signType.NameTH,
			NameEN: signType.NameEN,
		})
	}

	return response, nil
}

func (h *dropDownService) GetCustomerTypes(ctx context.Context) ([]dto.ResponseGetCustomerTypes, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	customerTypes, errOnGetCustomerTypes := h.dropDownRepo.GetCustomerTypes(ctx, filter, projection)
	if errOnGetCustomerTypes != nil {
		return nil, errOnGetCustomerTypes
	}

	if len(customerTypes) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetCustomerTypes
	for _, customerType := range customerTypes {
		response = append(response, dto.ResponseGetCustomerTypes{
			TypeID: customerType.TypeID,
			NameTH: customerType.NameTH,
			NameEN: customerType.NameEN,
		})
	}

	return response, nil
}

func (s *dropDownService) GetSignJobList(ctx context.Context, projectID string) ([]dto.ResponseGetSignList, error) {
	filter := bson.M{"project_id": projectID, "deleted_at": nil}
	projection := bson.M{}

	signJobs, errOnGetSignJobs := s.dropDownRepo.GetSignJobsList(ctx, filter, projection)
	if errOnGetSignJobs != nil {
		return nil, errOnGetSignJobs
	}

	if len(signJobs) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetSignList
	for _, signJob := range signJobs {
		response = append(response, dto.ResponseGetSignList{
			JobID:       signJob.JobID,
			ProjectName: signJob.ProjectName,
			JobName:     signJob.JobName,
			Content:     signJob.Content,
		})
	}

	return response, nil
}

func (s *dropDownService) GetProjectList(ctx context.Context) ([]dto.ResponseGetProjects, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	projects, errOnGetProjects := s.dropDownRepo.GetProjectsList(ctx, filter, projection)
	if errOnGetProjects != nil {
		return nil, errOnGetProjects
	}

	if len(projects) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetProjects
	for _, project := range projects {
		response = append(response, dto.ResponseGetProjects{
			ProjectID:   project.ProjectID,
			ProjectName: project.ProjectName,
		})
	}

	return response, nil
}

func (s *dropDownService) GetUserList(ctx context.Context) ([]dto.ResponseGetUsers, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{
		"user_id":       1,
		"title_th":      1,
		"first_name_th": 1,
		"last_name_th":  1,
		"full_name_th":  1,
	}

	users, errOnGetUsers := s.dropDownRepo.GetUsersList(ctx, filter, projection)
	if errOnGetUsers != nil {
		return nil, errOnGetUsers
	}

	if len(users) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	var response []dto.ResponseGetUsers
	for _, user := range users {
		response = append(response, dto.ResponseGetUsers{
			UserID:      user.UserID,
			TitleTH:     user.TitleTH,
			FirstNameTH: user.FirstNameTH,
			LastNameTH:  user.LastNameTH,
			FullNameTH:  fmt.Sprintf("%s %s %s", user.TitleTH, user.FirstNameTH, user.LastNameTH),
		})
	}

	return response, nil
}
