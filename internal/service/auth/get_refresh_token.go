package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/utils"
)

func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	refreshConfig, err := env.NewRefreshTokenConfig()
	if err != nil {
		return "", err
	}

	claims, err := utils.VerifyToken(oldRefreshToken, []byte(refreshConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(
		model.UserForToken{
			Username: claims.Username,
			//TODO:  Нужно ли лезть в базу или роль можно также всзять из calims?
			Role: claims.Role,
		},
		[]byte(refreshConfig.RefreshTokenSecretKey()),
		refreshConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
