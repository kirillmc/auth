package access

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	descAccess "github.com/kirillmc/auth/pkg/access_v1"
)

func (i *Implementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
