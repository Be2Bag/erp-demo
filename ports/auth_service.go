package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type AuthService interface {
	Login(ctx context.Context, user dto.RequestLogin) (string, error)
	GetSessions(ctx context.Context, token string) (map[string]interface{}, error)
	ResetPassword(ctx context.Context, req dto.RequestResetPassword) error
	ConfirmResetPassword(ctx context.Context, req dto.RequestConfirmResetPassword) error
}
