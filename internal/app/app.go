package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/kirillmc/auth/internal/config"
	localInterceptor "github.com/kirillmc/auth/internal/interceptor"
	"github.com/kirillmc/auth/internal/logger"
	"github.com/kirillmc/auth/internal/metric"
	"github.com/kirillmc/auth/internal/tracing"
	descAccess "github.com/kirillmc/auth/pkg/access_v1"
	descAuth "github.com/kirillmc/auth/pkg/auth_v1"
	descUser "github.com/kirillmc/auth/pkg/user_v1"
	_ "github.com/kirillmc/auth/statik"
	"github.com/kirillmc/platform_common/pkg/closer"
	"github.com/kirillmc/platform_common/pkg/interceptor"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
)

const (
	SERVICE_PEM = "tls/service.pem"
	SERVICE_KEY = "tls/service.key"
	serviceName = "auth-server-service"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// подвязываем инициализаторские штуки из service_provider к старту приложения

type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll() // не блокирующий вызов
		closer.Wait()     // блокирующий
		// делает graceful shutdown - но не сработает при log.Fatal или os.Exit
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal("failed to run gRPC server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal("failed to run HTTP server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			logger.Fatal("failed to run Prometheus server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatal("failed to run Swagger server", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initPrometheusServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.InitLoggerWithLogLevel(a.serviceProvider.LoggerConfig().LogLevel())

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile(SERVICE_PEM, SERVICE_KEY)
	if err != nil {
		logger.Fatal("failed to load TLS keys", zap.Error(err))
	}

	tracing.Init(serviceName)

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				localInterceptor.LogInterceptor,
				localInterceptor.MetricsInterceptor,
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				interceptor.ValidateInerceptor,
			),
		),
	)

	reflection.Register(a.grpcServer) // рефлексия вкл для постмана

	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImplementation(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImplementation(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImplementation(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initPrometheusServer(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		logger.Fatal("failed to init metrics", zap.Any("error", err))
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    a.serviceProvider.PrometheusConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("gRPC server is running on: ", zap.String("address", a.serviceProvider.GRPCConfig().Address()))

	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info("HTTP server is running on: ", zap.String("address", a.serviceProvider.HTTPConfig().Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Info("Prometheus server is running on: ", zap.String("address", a.serviceProvider.PrometheusConfig().Address()))

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Info("Swagger server is running on: ", zap.String("address", a.serviceProvider.SwaggerConfig().Address()))

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Serving swagger", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Info("Opening swagger at", zap.String("path", path))

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Write swagger file", zap.String("path", path))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file", zap.String("path", path))
	}
}
