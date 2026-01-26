package services

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
	authRepo ports.AuthRepository
	userRepo ports.UserRepository
	config   config.Config
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
	// ตรวจสอบความยาว slice ป้องกัน panic จากการเข้าถึง index 0 เมื่อไม่มีข้อมูล
	if len(userData) == 0 {
		return "", mongo.ErrNoDocuments
	}

	hashedPassword := util.HashPassword(user.Password, s.config.Hash.Salt)

	if userData[0].Password != hashedPassword {
		return "", fmt.Errorf("invalid password")
	}

	claims := dto.JWTClaims{
		UserID:       userData[0].UserID,
		Email:        userData[0].Email,
		EmployeeCode: userData[0].EmployeeCode,
		Role:         userData[0].Role,
		TitleTH:      userData[0].TitleTH,
		FirstNameTH:  userData[0].FirstNameTH,
		LastNameTH:   userData[0].LastNameTH,
		Avatar:       userData[0].Avatar,
		Status:       userData[0].Status,
	}

	token, err := util.GenerateJWTToken(claims, s.config.JWT.SecretKey, 30*24*time.Hour) // 30 days expiration
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

func (s *authService) ResetPassword(ctx context.Context, req dto.RequestResetPassword) error {

	filter := bson.M{"email": req.Email, "deleted_at": nil}
	projection := bson.M{}
	userData, errOnGetUserData := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetUserData != nil {
		return errOnGetUserData
	}
	if len(userData) == 0 {
		return mongo.ErrNoDocuments
	}

	claims := dto.JWTClaims{
		UserID:       userData[0].UserID,
		Email:        userData[0].Email,
		EmployeeCode: userData[0].EmployeeCode,
		Role:         userData[0].Role,
		TitleTH:      userData[0].TitleTH,
		FirstNameTH:  userData[0].FirstNameTH,
		LastNameTH:   userData[0].LastNameTH,
		Avatar:       userData[0].Avatar,
		Status:       userData[0].Status,
	}

	token, err := util.GenerateJWTToken(claims, s.config.JWT.SecretKey, 15*time.Minute) // 15 minutes expiration
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", req.RedirectURL, token)

	emailCfg := util.EmailConfig{
		Host:     s.config.Email.Host,
		Port:     s.config.Email.Port,
		Username: s.config.Email.Username,
		Password: s.config.Email.Password,
		From:     s.config.Email.From,
	}

	if err := util.SendResetPasswordEmail(emailCfg, req.Email, resetLink); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

func (s *authService) ConfirmResetPassword(ctx context.Context, req dto.RequestConfirmResetPassword) error {

	claims, err := util.VerifyJWTToken(req.Token, s.config.JWT.SecretKey)
	if err != nil {
		return err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return fmt.Errorf("token payload missing email")
	}

	hashed := util.HashPassword(req.NewPassword, s.config.Hash.Salt)

	filter := bson.M{"email": email, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"password": hashed, "updated_at": time.Now()}}
	_, errOnUpdate := s.userRepo.UpdateUserByFilter(ctx, filter, update)
	if errOnUpdate != nil {
		if errOnUpdate == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to update password: %w", errOnUpdate)
	}

	return nil
}
