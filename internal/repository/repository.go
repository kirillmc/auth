package repository

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
)

// файл ТОЛЬКО для интерфейсов

type UserRepository interface {
	Create(ctx context.Context, req *model.UserToCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, req *model.UserToUpdate) error
	Delete(ctx context.Context, id int64) error

	GetRole(ctx context.Context, userName string) (model.Role, error)
	GetHashPass(ctx context.Context, userName string) (string, error)

	GetAccessibleRoles(ctx context.Context) (map[string]model.Role, error)
}
