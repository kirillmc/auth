package access

import (
	"github.com/kirillmc/auth/internal/repository"
	def "github.com/kirillmc/auth/internal/service"
)

var _ def.AccessService = (*serv)(nil) //валидация имплементации интерфейса

type serv struct {
	accessRepository repository.AccessRepository
}

func NewService(accessRepository repository.AccessRepository) *serv {
	return &serv{accessRepository: accessRepository}
}
