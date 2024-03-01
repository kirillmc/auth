package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kirillmc/auth/internal/repository"
	"github.com/kirillmc/auth/internal/repository/user/converter"
	"github.com/kirillmc/auth/internal/repository/user/model"
	desc "github.com/kirillmc/auth/pkg/user_v1"
	"time"
)

// ТУТ ИМПЛЕМЕНТАЦИЯ МЕТОДОВ

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func (r repo) Create(ctx context.Context, req *desc.CreateRequest) (int64, error) {
	builder := sq.Insert(tableName).PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(req.Name, req.Email, genPassHash(req.Password), req.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r repo) Get(ctx context.Context, id int64) (*desc.GetResponse, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r repo) Update(ctx context.Context, req *desc.UpdateRequest) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(roleColumn, req.Role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": req.GetId()})
	if req.Name != nil {
		builder = builder.Set(nameColumn, req.Name.Value)
	}

	if req.Email != nil {
		builder = builder.Set(emailColumn, req.Email.Value)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r repo) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	builder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func genPassHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum(nil))
}
