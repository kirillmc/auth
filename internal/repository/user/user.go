package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/auth/internal/repository/user/converter"
	modelRepo "github.com/kirillmc/auth/internal/repository/user/model"
	"github.com/kirillmc/platform_common/pkg/db"
)

// ТУТ ИМПЛЕМЕНТАЦИЯ МЕТОДОВ

func (r *repo) Create(ctx context.Context, req *model.UserToCreate) (int64, error) {
	hashPass, err := genPassHash(req.Password)
	if err != nil {
		return 0, err
	}

	builder := sq.Insert(tableName).PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(req.Username, req.Email, hashPass, req.Role).
		Suffix(returnId)

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

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
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

func (r *repo) Update(ctx context.Context, req *model.UserToUpdate) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: req.Id})

	if req.Username != nil {
		builder = builder.Set(nameColumn, req.Username.Value)
	}

	if req.Email != nil {
		builder = builder.Set(emailColumn, req.Email.Value)
	}

	if req.Role != model.RoleUnknown {
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

func (r *repo) Delete(ctx context.Context, id int64) error {
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

func genPassHash1(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))

	return fmt.Sprintf("%x", h.Sum(nil))
}
func genPassHash(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}
