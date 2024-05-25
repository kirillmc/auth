package config

import (
	"time"

	"github.com/joho/godotenv"
)

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type AccessTokenConfig interface {
	AccessTokenSecretKey() string
	AccessTokenExpiration() time.Duration
}

type RefreshTokenConfig interface {
	RefreshTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
}

type LoggerConfig interface {
	LogLevel() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
