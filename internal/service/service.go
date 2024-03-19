package service

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, req *model.UserToCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req *model.UserToUpdate) error
}
