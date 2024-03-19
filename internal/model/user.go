package model

import (
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kirillmc/auth/pkg/user_v1"
)

type User struct {
	Id        int64
	Name      string
	Email     string
	Role      user_v1.Role
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserToCreate struct {
	Name     string
	Email    string
	Password string
	Role     Role
}

type UserToUpdate struct {
	Id    int64
	Name  *wrapperspb.StringValue
	Email *wrapperspb.StringValue
	Role  Role
}

type Role int32

const (
	Role_UNKNOWN Role = 0
	Role_USER    Role = 1
	Role_ADMIN   Role = 2
)
