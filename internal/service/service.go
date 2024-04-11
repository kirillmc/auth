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

type AccessService interface {
	Check(ctx context.Context, endpointAddress string) error
}

type AuthService interface {
	Login(ctx context.Context, req *model.UserToLogin) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}
