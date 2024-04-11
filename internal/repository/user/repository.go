package user

import (
	"github.com/kirillmc/auth/internal/repository"
	"github.com/kirillmc/platform_common/pkg/db"
)

// ТУТ ИМПЛЕМЕНТАЦИЯ МЕТОДОВ

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "username"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"

	returnId = "RETURNING id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
