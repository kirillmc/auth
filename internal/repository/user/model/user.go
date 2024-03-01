package model

import (
	"database/sql"
	"github.com/kirillmc/auth/pkg/user_v1"
	"time"
)

//	Id        int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
//	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
//	Email     string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
//	Role      Role                   `protobuf:"varint,4,opt,name=role,proto3,enum=user_v1.Role" json:"role,omitempty"`
//	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
//	UpdatedAt *timestamppb.Timestamp

type User struct {
	Id        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `ab:"email"`
	Role      user_v1.Role `db:"role"` //TODO: Спросить: нужно ли enum делать локально в этом файле/пакете или такого импорта достаточно?
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

//type UserToCreate struct {
//	Name     string       `db:"name"`
//	Email    string       `db:"email"`
//	Password string       `db:"password"`
//	Role     user_v1.Role `db:"role"`
//}
