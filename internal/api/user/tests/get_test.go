package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/service"
	serviceMocks "github.com/kirillmc/auth/internal/service/mocks"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = gofakeit.Number(0, 2)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("error of service layer")

		req = &desc.GetRequest{
			Id: id,
		}

		serviceRes = &model.User{
			Id:        id,
			Username:  name,
			Email:     email,
			Role:      model.Role(role),
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}

		res = &desc.GetResponse{
			Id:        id,
			Username:  name,
			Email:     email,
			Role:      desc.Role(role),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success get case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(serviceRes, nil)
				return mock
			},
		},
		{
			name: "service error get",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			resU, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resU)
		})
	}
}
