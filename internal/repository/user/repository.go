package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillmc/auth/internal/client/db"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/repository"
	"github.com/kirillmc/auth/internal/repository/user/converter"
	modelRepo "github.com/kirillmc/auth/internal/repository/user/model"
	desc "github.com/kirillmc/auth/pkg/user_v1"
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
	db db.Client // Клиент вместо *pgx.Pool
}

func (r repo) Create(ctx context.Context, req *model.UserToCreate) (int64, error) {
	builder := sq.Insert(tableName).PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(req.Name, req.Email, genPassHash(req.Password), req.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	// Добавлено с db.Client
	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...) // Сканирует одну запись в user
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r repo) Update(ctx context.Context, req *model.UserToUpdate) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: req.Id})
	if req.Name != nil {
		builder = builder.Set(nameColumn, req.Name.Value)
	}

	if req.Email != nil {
		builder = builder.Set(emailColumn, req.Email.Value)
	}

	if req.Role != desc.Role_UNKNOWN {
		builder = builder.Set(roleColumn, req.Role)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": id})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func genPassHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum(nil))
}
