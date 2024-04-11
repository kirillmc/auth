package tests

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/repository"
	repositoryMocks "github.com/kirillmc/auth/internal/repository/mocks"
	userService "github.com/kirillmc/auth/internal/service/user"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UserToUpdate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.BeerName()
		email = gofakeit.Email()
		role  = gofakeit.Number(0, 2)

		repositoryErr = fmt.Errorf("error of repository layer")

		req = &model.UserToUpdate{
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
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success update case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "repository error update",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			service := userService.NewService(userRepositoryMock)

			err := service.Update(ctx, req)

			require.Equal(t, tt.err, err)
		})
	}
}
