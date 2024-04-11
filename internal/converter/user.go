package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kirillmc/auth/internal/model"
	descAuth "github.com/kirillmc/auth/pkg/auth_v1"
	desc "github.com/kirillmc/auth/pkg/user_v1"
)

func ToGetResponseFromService(user *model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserModelCreateFromDesc(user *desc.CreateRequest) *model.UserToCreate {
	return &model.UserToCreate{
		Username:        user.Username,
		Email:           user.Email,
		Role:            model.Role(user.Role),
		Password:        user.Password,
		ConfirmPassword: user.PasswordConfirm,
	}
}

func ToUserModelUpdateFromDesc(user *desc.UpdateRequest) *model.UserToUpdate {
	return &model.UserToUpdate{
		Id:       user.Id,
		Username: user.Name,
		Email:    user.Email,
		Role:     model.Role(user.Role),
	}
}

func ToUserToLoginFromDescAuth(user *descAuth.LoginRequest) *model.UserToLogin {
	return &model.UserToLogin{
		Username: user.Username,
		Password: user.Password,
	}
}
