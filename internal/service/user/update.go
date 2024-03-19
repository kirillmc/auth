package user

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, req *model.UserToUpdate) error {
	err := s.userRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
