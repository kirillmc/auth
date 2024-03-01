package repository

import (
	"context"
	desc "github.com/kirillmc/auth/pkg/user_v1"
)

// файл ТОЛЬКО для интерфейсов

type UserRepository interface {
	Create(ctx context.Context, req *desc.CreateRequest) (int64, error)
	Get(ctx context.Context, id int64) (*desc.GetResponse, error)
	Update(ctx context.Context, req *desc.UpdateRequest) error
	Delete(ctx context.Context, req *desc.DeleteRequest) error
}
