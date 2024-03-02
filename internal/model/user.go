package model

import (
	"database/sql"
	"github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
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
	Role     user_v1.Role
}

type UserToUpdate struct {
	Id    int64
	Name  *wrapperspb.StringValue
	Email *wrapperspb.StringValue
	Role  user_v1.Role
}
