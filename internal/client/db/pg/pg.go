package pg

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kirillmc/auth/internal/client/db"
)

type pg struct {
	dbc *pgxpool.Pool
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	//TODO implement me
	panic("implement me")
}

func (p *pg) Ping(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (p pg) Close() {
	//TODO implement me
	panic("implement me")
}

func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {

}
