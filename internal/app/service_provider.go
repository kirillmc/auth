package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/closer"
	"github.com/kirillmc/auth/internal/config"
	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/repository"
	userRepo "github.com/kirillmc/auth/internal/repository/user"
	"github.com/kirillmc/auth/internal/service"
	userService "github.com/kirillmc/auth/internal/service/user"
	"log"
)

// содержит все зависимости, необходимые в рамках приложения
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool *pgxpool.Pool

	userRepository repository.UserRepository
	userService    service.UserService

	userImpl *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// если Get - в GO Get НЕ УКАЗЫВАЮТ: НЕ GetPGConfig, A PGConfig
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err) // делаем Log.Fatalf, чтобы не обрабатывать ошибку в другом месте
			// + инициализация происходит при старте приложения, поэтому если ошибка - можно и сервер уронить
			// можно кинуть panic()
		}
		s.pgConfig = pgConfig
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = grpcConfig
	}
	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN()) // получется каскадная инициализация
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = pool.Ping(ctx) // проверяем доступность, пингуя
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		// defer pgPool.Close - не отработает как в main, т.к. функция закончися почти сразу,
		// Поэтому:
		// Для закрытия реализуем closer - сущность, которая копит функции закрытия каких-либо ресурсов и вызывается в конце
		closer.Add(func() error { // т.к. метод Close у pool не возвращает ошибку - добавлена обертка,
			// возворащающая в качестве ошибки nil
			pool.Close()
			return nil
		})
		s.pgPool = pool
	}
	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.PgPool(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}
	return s.userImpl
}
