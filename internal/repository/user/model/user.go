package model

import (
	"database/sql"
	"time"

	"github.com/kirillmc/auth/internal/model"
)

type User struct {
	Id        int64        `db:"id"`
	Username  string       `db:"username"`
	Email     string       `ab:"email"`
	Role      model.Role   `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
