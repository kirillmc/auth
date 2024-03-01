package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kirillmc/auth/internal/config"
	"github.com/kirillmc/auth/internal/config/env"
	"github.com/kirillmc/auth/internal/repository"
	"github.com/kirillmc/auth/internal/repository/user"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
	//p              *pgxpool.Pool //TODO: потом удалить
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("insered user with id: %d", id)
	//pool.QueryRow // считать одну строку
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	nUser, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	log.Printf("%v %v %v %v %v", nUser.Id, nUser.Name, nUser.Email, nUser.Role, nUser.CreatedAt)
	return nUser, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.userRepository.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("User %d updated", req.GetId())
	return nil, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userRepository.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("User %d was deleted", req.GetId())
	return nil, nil
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

	userRepo := user.NewRepository(pool)

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserV1Server(s, &server{userRepository: userRepo})
	log.Printf("server is listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//func genPassHash(pass string) string {
//	h := sha256.New()
//	h.Write([]byte(pass))
//	return fmt.Sprintf("%x", h.Sum(nil))
//}
