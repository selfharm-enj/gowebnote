// Адаптер для работы через pgxpool.Pool.
package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Создание нового pgpool.Pool. Panic если произойдет ошибка.
func MustNewPool(ctx context.Context, c *Config) *pgxpool.Pool {
	const op = "pkg.adapters.postgres.pool.New"
	pool, err := pgxpool.New(ctx, c.getConnString())
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		panic(err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		panic(err)
	}
	slog.Debug("db connected", "op", op)
	return pool
}
