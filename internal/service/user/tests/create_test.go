package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/repository"
	repositoryMocks "github.com/kirillmc/auth/internal/repository/mocks"
	userService "github.com/kirillmc/auth/internal/service/user"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UserToCreate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.UUID()
		role     = gofakeit.Number(0, 2)

		repositoryErr = fmt.Errorf("error of repository layer")

		req = &model.UserToCreate{
			Username:        name,
			Email:           email,
			Password:        password,
			ConfirmPassword: password,
			Role:            model.Role(role),
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{{
		name: "success create case",
		args: args{
			ctx: ctx,
			req: req,
		},
		want: id,
		err:  nil,
		userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
			mock := repositoryMocks.NewUserRepositoryMock(mc)
			mock.CreateMock.Expect(ctx, req).Return(id, nil)
			return mock
		},
	},
		{
			name: "repository error create",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repositoryErr)
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

			newId, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newId)

		})
	}
}
