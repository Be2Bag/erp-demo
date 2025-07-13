package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authService struct {
	config   config.Config
	authRepo ports.AuthRepository
	userRepo ports.UserRepository
}

func NewAuthService(cfg config.Config, authRepo ports.AuthRepository, userRepo ports.UserRepository) ports.AuthService {
	return &authService{config: cfg, authRepo: authRepo, userRepo: userRepo}
}

func (s *authService) Login(ctx context.Context, user dto.RequestLogin) (string, error) {

	filter := bson.M{
		"email":      user.Email,
		"deleted_at": nil,
	}

	projection := bson.M{
		"user_id":       1,
		"username":      1,
		"email":         1,
		"password":      1,
		"role":          1,
		"status":        1,
		"employee_code": 1,
		"title_th":      1,
		"first_name_th": 1,
		"last_name_th":  1,
		"avatar":        1,
	}

	userData, errOnGetUserData := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetUserData != nil {
		return "", errOnGetUserData
	}
	if userData == nil {
		return "", mongo.ErrNoDocuments
	}

	hashedPassword := util.HashPassword(user.Password, s.config.Hash.Salt)

	if userData[0].Password != hashedPassword {
		return "", fmt.Errorf("invalid password")
	}

	claims := map[string]interface{}{
		"UserID":       userData[0].UserID,
		"Username":     userData[0].Username,
		"EmployeeCode": userData[0].EmployeeCode,
		"Role":         userData[0].Role,
		"TitleTH":      userData[0].TitleTH,
		"FirstNameTH":  userData[0].FirstNameTH,
		"LastNameTH":   userData[0].LastNameTH,
		"Avatar":       userData[0].Avatar,
		"Status":       userData[0].Status,
	}

	token, err := util.GenerateJWTToken(claims, s.config.JWT.SecretKey, 50000*time.Second) // 5 minutes expiration
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *authService) GetSessions(ctx context.Context, token string) (map[string]interface{}, error) {

	claims, errOnVerifyJWTToken := util.VerifyJWTToken(token, s.config.JWT.SecretKey)
	if errOnVerifyJWTToken != nil {
		return nil, errOnVerifyJWTToken
	}

	return claims, nil
}
