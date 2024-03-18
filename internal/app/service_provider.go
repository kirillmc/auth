package app

import (
	"context"
	"log"

	"github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/client/db"
	"github.com/kirillmc/auth/internal/client/db/pg"
	"github.com/kirillmc/auth/internal/client/db/transaction"
	"github.com/kirillmc/auth/internal/closer"
	"github.com/kirillmc/auth/internal/config"
	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/repository"
	userRepo "github.com/kirillmc/auth/internal/repository/user"
	"github.com/kirillmc/auth/internal/service"
	userService "github.com/kirillmc/auth/internal/service/user"
)

// содержит все зависимости, необходимые в рамках приложения
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("faailed to create db client: %v", err)
		}
		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) PgPool(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN()) // получется каскадная инициализация
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = cl.DB().Ping(ctx) // проверяем доступность, пингуя
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}

		// defer pgPool.Close - не отработает как в main, т.к. функция закончися почти сразу,
		// Поэтому:
		// Для закрытия реализуем closer - сущность, которая копит функции закрытия каких-либо ресурсов и вызывается в конце
		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.PgPool(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}
	return s.userImpl
}
