package auth

import (
	"context"

	"github.com/kirillmc/auth/internal/converter"
	descAuth "github.com/kirillmc/auth/pkg/auth_v1"
)

func (i *Implementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, converter.ToUserToLoginFromDescAuth(req))
	if err != nil {
		return nil, err
	}

	return &descAuth.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
