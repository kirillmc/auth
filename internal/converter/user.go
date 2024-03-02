package converter

import (
	"github.com/kirillmc/auth/internal/model"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateRequestFromService(user *model.UserToCreate) *desc.CreateRequest {
	return &desc.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func ToGetResponseFromService(user *model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserModelCreateFromDesc(user *desc.CreateRequest) *model.UserToCreate {
	return &model.UserToCreate{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}
}

func ToUserModelUpdateFromDesc(user *desc.UpdateRequest) *model.UserToUpdate {
	return &model.UserToUpdate{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
