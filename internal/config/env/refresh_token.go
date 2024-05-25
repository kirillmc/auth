package env

import (
	"os"
	"strconv"
	"time"

	"github.com/kirillmc/auth/internal/config"
	"github.com/pkg/errors"
)

const (
	refreshTokenEnvName           = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION_IN_MINUTES"
)

type refreshTokenConfig struct {
	secretKey      string
	timeExpiration time.Duration
}

func NewRefreshTokenConfig() (config.RefreshTokenConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refresh token secret key not found")
	}

	refreshTokenExpiration := os.Getenv(refreshTokenExpirationEnvName)
	if len(refreshTokenExpiration) == 0 {
		return nil, errors.New("refresh token expiration not found")
	}

	timeExpired, err := strconv.Atoi(refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &refreshTokenConfig{
		secretKey:      refreshTokenSecretKey,
		timeExpiration: time.Duration(timeExpired) * time.Minute,
	}, nil
}

func (cfg *refreshTokenConfig) RefreshTokenSecretKey() string {
	return cfg.secretKey
}

func (cfg *refreshTokenConfig) RefreshTokenExpiration() time.Duration {
	return cfg.timeExpiration
}
