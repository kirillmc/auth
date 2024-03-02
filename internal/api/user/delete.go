package user

import (
	"context"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	log.Printf("User %d was deleted", req.GetId())
	return nil, nil
}
