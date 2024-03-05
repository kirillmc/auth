package app

import (
	"context"
	"google.golang.org/grpc"
)

// подвязываем инициализаторские штуки из service_provider к старту приложения

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context)
