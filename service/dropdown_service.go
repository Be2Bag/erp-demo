package service

import (
	"context"

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

func (s *dropDownService) GetPositions(ctx context.Context) ([]dto.ResponseGetPositions, error) {

	filter := bson.M{"deleted_at": nil}
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
