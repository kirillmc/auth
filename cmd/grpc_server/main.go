package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kirillmc/auth/internal/config"
	"github.com/kirillmc/auth/internal/config/env"
	desc "github.com/kirillmc/auth/pkg/user_v1"
)

const (
	tableName       = "users"
	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	p *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	buildInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(req.Name, req.Email, genPassHash(req.Password), req.Role).
		Suffix("RETURNING id")

	query, args, err := buildInsert.ToSql()
	if err != nil {
		//	log.Fatalf("failed to build query: %v", err)
		return nil, err
	}

	var userID int64
	err = s.p.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		//log.Fatalf("failed to insert note: %v", err)
		return nil, err
	}
	//pool.QueryRow // считать одну строку
	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	var id, role int64
	var name, email string
	var createdAt time.Time
	var updatedAt sql.NullTime
	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		//log.Fatalf("failed to build SELECT query: %v", err)
		return nil, err
	}
	err = s.p.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		//log.Fatalf("failed to SELECT user: %v", err)
		return nil, err
	}

	var upTime *timestamppb.Timestamp
	if updatedAt.Valid {
		upTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      desc.Role(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: upTime,
	}, nil
}
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: req.GetId()})
	if req.Role != nil {
		builderUpdate = builderUpdate.Set(roleColumn, req.Role.Value)
	}
	if req.Name != nil {
		builderUpdate = builderUpdate.Set(nameColumn, req.Name.Value)
	}
	if req.Email != nil {
		builderUpdate = builderUpdate.Set(emailColumn, req.Email.Value)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		//log.Fatalf("failed to build UPDATE query: %v", err)
		return nil, err
	}

	_, err = s.p.Exec(ctx, query, args...)
	if err != nil {
		//log.Fatalf("failed to update user")
		return nil, err
	}
	return nil, nil
}
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: req.GetId()})
	query, args, err := builderDelete.ToSql()
	if err != nil {
		//log.Fatalf("failed to build DELETE query: %v", err)
		return nil, err
	}

	_, err = s.p.Exec(ctx, query, args...)
	if err != nil {
		//log.Fatalf("failed to delete user with id %d", req.GetId())
		return nil, err
	}
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

	lis, err := net.Listen("tcp", grpcConfig.Address())
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
