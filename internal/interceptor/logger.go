package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kirillmc/auth/internal/logger"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("request", req))
		return res, err
	}

	logger.Info("request success", zap.String("method", info.FullMethod), zap.Any("request", res), zap.Any("response", res), zap.Any("duration", time.Since(now)))

	return res, err
}
