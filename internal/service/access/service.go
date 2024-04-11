package access

import (
	"github.com/kirillmc/auth/internal/repository"
	def "github.com/kirillmc/auth/internal/service"
)

var _ def.AccessService = (*serv)(nil) //валидация имплементации интерфейса

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *serv {
	return &serv{userRepository: userRepository}
}
