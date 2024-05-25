package auth

import (
	"context"

	descAuth "github.com/kirillmc/auth/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	return &descAuth.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
