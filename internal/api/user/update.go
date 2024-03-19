package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kirillmc/auth/internal/converter"
	desc "github.com/kirillmc/auth/pkg/user_v1"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUserModelUpdateFromDesc(req))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
