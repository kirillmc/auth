package access

import (
	"github.com/kirillmc/auth/internal/service"
	desc "github.com/kirillmc/auth/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
