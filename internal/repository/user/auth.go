package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/platform_common/pkg/db"
)

// ТУТ ИМПЛЕМЕНТАЦИЯ МЕТОДОВ

func (r *repo) GetRole(ctx context.Context, userName string) (model.Role, error) {
	builder := sq.Select(roleColumn).PlaceholderFormat(sq.Dollar).From(tableName).Where(sq.Eq{nameColumn: userName}).Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.RoleUnknown, err
	}

	q := db.Query{
		Name:     "user_repository.GetRole",
		QueryRaw: query,
	}

	var role model.Role

	//err=r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role)
	err = r.db.DB().ScanOneContext(ctx, &role, q, args...)
	if err != nil {
		return model.RoleUnknown, err
	}

	return role, nil
}

func (r *repo) GetHashPass(ctx context.Context, userName string) (string, error) {
	builder := sq.Select(passwordColumn).PlaceholderFormat(sq.Dollar).From(tableName).Where(sq.Eq{nameColumn: userName}).Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "user_repository.GetHashPass",
		QueryRaw: query,
	}

	var hashPass string

	//err=r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role)
	err = r.db.DB().ScanOneContext(ctx, &hashPass, q, args...)
	if err != nil {
		return "", err
	}

	return hashPass, nil
}
