package user

import (
	"context"
	"log"

	"github.com/kirillmc/auth/internal/converter"
	desc "github.com/kirillmc/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserModelCreateFromDesc(req))
	if err != nil {
		return nil, err
	}
	log.Printf("inserted user with id: %d", id)
	//pool.QueryRow // считать одну строку
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
