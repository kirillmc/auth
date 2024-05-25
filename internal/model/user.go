package model

import (
	"database/sql"
	"time"

	"github.com/kirillmc/platform_common/pkg/nillable"
)

type User struct {
	Id        int64
	Username  string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserToCreate struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	Role            Role
}

type UserToUpdate struct {
	Id       int64
	Username nillable.NilString
	Email    nillable.NilString
	Role     Role
}

type UserToLogin struct {
	Username string
	Password string
}

type UserForToken struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type Role int32

const (
	RoleUnknown Role = 0
	RoleUser    Role = 1
	RoleAdmin   Role = 2
)
