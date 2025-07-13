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

func (s *dropDownService) GetPositions() ([]dto.ResponseGetPositions, error) {

	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	positions, errOnGetPositions := s.dropDownRepo.GetPositions(context.Background(), filter, projection)
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

func (s *dropDownService) GetDepartments() ([]dto.ResponseGetDepartments, error) {
	filter := bson.M{"deleted_at": nil}
	projection := bson.M{}

	departments, errOnGetDepartments := s.dropDownRepo.GetDepartments(context.Background(), filter, projection)
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
