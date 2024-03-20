package user

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, req *model.UserToCreate) (int64, error) {
	id, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	return id, nil
}
