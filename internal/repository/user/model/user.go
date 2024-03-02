package model

import (
	"database/sql"
	"github.com/kirillmc/auth/pkg/user_v1"
	"time"
)

type User struct {
	Id        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `ab:"email"`
	Role      user_v1.Role `db:"role"` //TODO: Спросить: нужно ли enum делать локально в этом файле/пакете или такого импорта достаточно?
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
