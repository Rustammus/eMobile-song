package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Deps struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPool(ctx context.Context, d Deps) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
	maxAttempts := 5

	ctxPool, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pool, err = pgxpool.New(ctxPool, dsn)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for ; maxAttempts > 0; maxAttempts-- {
		if err = pool.Ping(ctxPing); err != nil {
			time.Sleep(time.Second)
			continue
		}
		return pool, nil
	}

	return nil, err
}
