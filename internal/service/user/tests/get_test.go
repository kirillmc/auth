package tests

import (
	"context"
	"database/sql"
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

func TestGet(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		id  int64
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

		repositoryErr = fmt.Errorf("error of repository layer")

		res = &model.User{
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
	)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success get case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "repository error get",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repositoryErr)
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

			resUser, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resUser)
		})
	}
}
