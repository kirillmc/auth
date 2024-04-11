package user

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
	"github.com/pkg/errors"
)

func (s *serv) Create(ctx context.Context, req *model.UserToCreate) (int64, error) {
	if req.ConfirmPassword != req.Password {
		return 0, errors.New("Passwords are not equal")
	}

	id, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	return id, nil
}
