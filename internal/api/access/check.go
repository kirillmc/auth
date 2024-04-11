package access

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	descAccess "github.com/kirillmc/auth/pkg/access_v1"
)

func (i *Implementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		log.Print("Error in api layer")
		return nil, err
	}

	return nil, nil
}
