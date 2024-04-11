package tests

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/service"
	serviceMocks "github.com/kirillmc/auth/internal/service/mocks"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.BeerName()
		email = gofakeit.Email()
		role  = gofakeit.Number(0, 2)

		serviceErr = fmt.Errorf("error of service layer")

		req = &desc.UpdateRequest{
			Id: id,
			Name: &wrapperspb.StringValue{
				Value: name,
			},
			Email: &wrapperspb.StringValue{
				Value: email,
			},
			Role: desc.Role(role),
		}

		modelUserToUpdate = &model.UserToUpdate{
			Id: id,
			Username: &wrapperspb.StringValue{
				Value: name,
			},
			Email: &wrapperspb.StringValue{
				Value: email,
			},
			Role: model.Role(role),
		}
	)

	tests := []struct {
		name            string
		args            args
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success update case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelUserToUpdate).Return(nil)
				return mock
			},
		},
		{
			name: "service error update",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelUserToUpdate).Return(serviceErr)
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

			_, err := api.Update(ctx, req)

			require.Equal(t, tt.err, err)

		})

	}
}
