package access

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/utils"
	"github.com/pkg/errors"
)

const (
	authorization = "authorization"
	authPrefix    = "Bearer "
)

func (s *serv) Check(ctx context.Context, endpointAddress string) error {
	accessConfig, err := env.NewAccessTokenConfig()
	if err != nil {
		return err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md[authorization]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(accessConfig.AccessTokenSecretKey()))
	if err != nil {
		return errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}
	return errors.New("access denied")
}
