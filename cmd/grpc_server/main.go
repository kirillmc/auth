package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	userAPI "github.com/kirillmc/auth/internal/api/user"
	"github.com/kirillmc/auth/internal/config"
	"github.com/kirillmc/auth/internal/config/env"
	userRepo "github.com/kirillmc/auth/internal/repository/user"
	userService "github.com/kirillmc/auth/internal/service/user"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	//Считываем environment variables
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failded to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := userRepo.NewRepository(pool)
	userService := userService.NewService(userRepo)

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserV1Server(s, userAPI.NewImplementation(userService))
	log.Printf("server is listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
