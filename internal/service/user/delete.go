package user

import (
	"context"

	"github.com/kirillmc/auth/internal/logger"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger.Info("user was deleted successfully")
	return nil
}
