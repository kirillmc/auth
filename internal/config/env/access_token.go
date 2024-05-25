package env

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/kirillmc/auth/internal/config"
)

const (
	accessTokenEnvName           = "ACCESS_TOKEN_SECRET_KEY"
	accessTokenExpirationEnvName = "ACCESS_TOKEN_EXPIRATION_IN_MINUTES"
)

type accessTokenConfig struct {
	secretKey      string
	timeExpiration time.Duration
}

func NewAccessTokenConfig() (config.AccessTokenConfig, error) {
	accessTokenSecretKey := os.Getenv(accessTokenEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("access token secret key not found")
	}

	accessTokenExpiration := os.Getenv(accessTokenExpirationEnvName)
	if len(accessTokenExpiration) == 0 {
		return nil, errors.New("access token expiration not found")
	}

	timeExpired, err := strconv.Atoi(accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &accessTokenConfig{
		secretKey:      accessTokenSecretKey,
		timeExpiration: time.Duration(timeExpired) * time.Minute,
	}, nil
}

func (cfg *accessTokenConfig) AccessTokenSecretKey() string {
	return cfg.secretKey
}

func (cfg *accessTokenConfig) AccessTokenExpiration() time.Duration {
	return cfg.timeExpiration
}
