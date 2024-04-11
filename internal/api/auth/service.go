package auth

import (
	"github.com/kirillmc/auth/internal/service"
	desc "github.com/kirillmc/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
