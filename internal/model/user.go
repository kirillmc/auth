package model

import (
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
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
	Username *wrapperspb.StringValue
	Email    *wrapperspb.StringValue
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
