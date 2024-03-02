package user

import (
	"context"
	"github.com/kirillmc/auth/internal/converter"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	nUser, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	log.Printf("%v %v %v %v %v", nUser.Id, nUser.Name, nUser.Email, nUser.Role, nUser.CreatedAt)
	return converter.ToGetResponseFromService(nUser), nil
}
