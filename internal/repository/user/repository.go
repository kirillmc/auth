package user

import (
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
