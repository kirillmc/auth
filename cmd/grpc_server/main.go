package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

const (
	grpcPort = 50051

	// Поработать с конфигом и изменить
	dbDSN = "host=localhost port=50321 dbname=users user=users-user password=users-password sslmode=disable"
)

type server struct {
	desc.UnimplementedUserV1Server
	p *pgxpool.Pool
}

// Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
// Name            string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
//
//	Email           string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
//	Password        string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
//	PasswordConfirm string `protobuf:"bytes,4,opt,name=password_confirm,json=passwordConfirm,proto3" json:"password_confirm,omitempty"`
//	Role            Role
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	buildInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(req.Name, req.Email, genPassHash(req.Password), req.Role).
		Suffix("RETURNING id")

	query, args, err := buildInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	err = s.p.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}
	//pool.QueryRow // считать одну строку
	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      "AAA",
		Email:     "AAA",
		Role:      desc.Role_USER,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("id of updated user: %d", req.GetId())
	return nil, nil
}
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("id of deleted user: %d", req.GetId())
	return nil, nil
}

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserV1Server(s, &server{p: pool})
	log.Printf("server is listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func genPassHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum(nil))
}
