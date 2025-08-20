package services

import (
	"context"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type adminService struct {
	config    config.Config
	adminRepo ports.AdminRepository
	authRepo  ports.AuthRepository
	userRepo  ports.UserRepository
}

func NewAdminService(cfg config.Config, adminRepo ports.AdminRepository, authRepo ports.AuthRepository, userRepo ports.UserRepository) ports.AdminService {
	return &adminService{config: cfg, adminRepo: adminRepo, authRepo: authRepo, userRepo: userRepo}
}

func (s *adminService) UpdateUserStatus(ctx context.Context, req dto.RequestUpdateUserStatus) error {

	filter := bson.M{"user_id": req.UserID, "deleted_at": nil}
	projection := bson.M{}

	users, errOnGetUser := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetUser != nil {
		return errOnGetUser
	}

	if len(users) == 0 {
		return mongo.ErrNoDocuments
	}

	// if users[0].Status != "pending" {
	// 	return fmt.Errorf("user status is not pending, current status: %s", users[0].Status)
	// }

	update := bson.M{"$set": bson.M{"status": req.Status, "updated_at": time.Now()}}

	if req.Status == "deleted" {
		update = bson.M{"$set": bson.M{"status": req.Status, "deleted_at": time.Now()}}
	}

	_, errOnUpdateStatus := s.userRepo.UpdateUserByFilter(ctx, filter, update)
	if errOnUpdateStatus != nil {
		return errOnUpdateStatus
	}

	return nil
}

func (s *adminService) UpdateUserRole(ctx context.Context, req dto.RequestUpdateUserRole) error {

	filter := bson.M{"user_id": req.UserID, "deleted_at": nil}
	projection := bson.M{}

	users, errOnGetUser := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetUser != nil {
		return errOnGetUser
	}

	if len(users) == 0 {
		return mongo.ErrNoDocuments
	}

	update := bson.M{"$set": bson.M{"role": req.Role, "note": req.Note, "updated_at": time.Now()}}

	_, errOnUpdateRole := s.userRepo.UpdateUserByFilter(ctx, filter, update)
	if errOnUpdateRole != nil {
		return errOnUpdateRole
	}

	return nil
}

func (s *adminService) UpdateUserPosition(ctx context.Context, req dto.RequestUpdateUserPosition) error {

	filter := bson.M{"user_id": req.UserID, "deleted_at": nil}
	projection := bson.M{}

	users, errOnGetUser := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetUser != nil {
		return errOnGetUser
	}

	if len(users) == 0 {
		return mongo.ErrNoDocuments
	}

	update := bson.M{"$set": bson.M{"department_id": req.DepartmentID, "position_id": req.PositionID, "note": req.Note, "updated_at": time.Now()}}

	_, errOnUpdatePosition := s.userRepo.UpdateUserByFilter(ctx, filter, update)
	if errOnUpdatePosition != nil {
		return errOnUpdatePosition
	}

	return nil
}
