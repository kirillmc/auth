package access

import (
	"context"
	"log"
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
		log.Print("Error in service layer1")
		return err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("Error in service layer2")
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md[authorization]
	if !ok || len(authHeader) == 0 {
		log.Println("Error in service layer3")
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		log.Print("Error in service layer4")
		return errors.New("invalid authorization header format")
	}
	log.Printf("AuthHeader:%v\n", authHeader)
	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	log.Printf("Access token:%v\n", accessToken)

	claims, err := utils.VerifyToken(accessToken, []byte(accessConfig.AccessTokenSecretKey()))
	if err != nil {
		log.Print("Error in service layer5")
		log.Printf("Error:%v\n", err)
		return errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	log.Printf("MAP: %v,\n", accessibleMap)
	if err != nil {
		log.Print("Error in service layer6")
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
