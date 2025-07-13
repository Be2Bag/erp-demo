package ports

import (
	"context"

	"github.com/Be2Bag/erp-demo/dto"
)

type AuthService interface {
	Login(ctx context.Context, user dto.RequestLogin) (string, error)
	GetSessions(ctx context.Context, token string) (map[string]interface{}, error)
}
