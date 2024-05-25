package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/service"
	serviceMocks "github.com/kirillmc/auth/internal/service/mocks"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.UUID()
		role     = gofakeit.Number(0, 2)

		serviceErr = fmt.Errorf("error of service layer")

		req = &desc.CreateRequest{
			Username:        name,
			Email:           email,
			PasswordConfirm: password,
			Password:        password,
			Role:            desc.Role(role),
		}

		modelUser = &model.UserToCreate{
			Username:        name,
			Email:           email,
			ConfirmPassword: password,
			Password:        password,
			Role:            model.Role(role),
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{{
		name: "success create case",
		args: args{
			ctx: ctx,
			req: req,
		},
		want: res,
		err:  nil,
		userServiceMock: func(mc *minimock.Controller) service.UserService {
			mock := serviceMocks.NewUserServiceMock(mc)
			mock.CreateMock.Expect(ctx, modelUser).Return(id, nil)
			return mock
		},
	},
		{
			name: "service error create",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, modelUser).Return(0, serviceErr)
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

			newId, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newId)

		})
	}
}
