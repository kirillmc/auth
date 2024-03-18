package user

import (
	"github.com/kirillmc/auth/internal/client/db"
	"github.com/kirillmc/auth/internal/repository"
	def "github.com/kirillmc/auth/internal/service"
)

var _ def.UserService = (*serv)(nil) //валидация имплементации интерфейса

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) *serv {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
