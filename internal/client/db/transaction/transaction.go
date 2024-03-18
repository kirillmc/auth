package transaction

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/kirillmc/auth/internal/client/db"
	"github.com/kirillmc/auth/internal/client/db/pg"
)

// TRANSACTION MANAGER

type manager struct {
	db db.Transactor
}

// NewTransactionManager - функция, которая создает новый менеджер транзакций, который
// удовлетворяет интерфейсу db.TxManager
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

// transaction - основная функция, которая выполняет указанный пользователем обработчик в транзакции
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// Если это вложенная транзакция, пропускаем инициализацию новой транзакции и выполняем обработчик
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx) // Проверка от вложенных транзакций
	if ok {
		return fn(ctx)
	}

	// Стартуем новую транзакцию.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {

	}

	// Кладем транзакцию в контекст.
	ctx = pg.MakeContextTx(ctx, tx)
}

func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevl: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
