package auth

import (
	"context"
	"errors"

	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/utils"
	"github.com/kirillmc/platform_common/pkg/verify_password"
)

func (s *serv) Login(ctx context.Context, req *model.UserToLogin) (string, error) {
	// Сверяем хэши пароля
	hashPass, err := s.userRepository.GetHashPass(ctx, req.Username)
	if err != nil {
		return "", err
	}

	role, err := s.userRepository.GetRole(ctx, req.Username)
	if err != nil {
		return "", err
	}

	if !verify_password.VerifyPassword(hashPass, req.Password) {
		return "", errors.New("Wrong password!!!")
	}

	refreshConfig, err := env.NewRefreshTokenConfig()
	if err != nil {
		return "", err
	}

	refreshToken, err := utils.GenerateToken(
		model.UserForToken{
			Username: req.Username,
			Role:     role,
		},
		[]byte(refreshConfig.RefreshTokenSecretKey()),
		refreshConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
