package crud

import (
	"context"
	"eMobile/internal/config"
	"eMobile/pkg/client/postgres"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Pagination struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// GetPool
// Get connection pool wrapped in Client interface
func GetPool(conf *config.Config, logger logging.Logger) Client {
	pool, err := postgres.NewPool(context.TODO(), postgres.Deps{
		Username: conf.Storage.Username,
		Password: conf.Storage.Password,
		Host:     conf.Storage.Host,
		Port:     conf.Storage.Port,
		Database: conf.Storage.Database,
	})
	if err != nil {
		logger.Fatalf("Can't crate connection Pool. Abort start app. \n Error: %s", err.Error())
	} else {
		logger.Infof("Connected to database: %s", conf.Storage.Database)
	}
	return pool
}
