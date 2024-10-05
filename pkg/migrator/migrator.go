package migrator

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Deps struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Source   string
}

type Migrator struct {
	migrator *migrate.Migrate
}

func (m Migrator) Up() error {
	return m.migrator.Up()
}

func (m Migrator) Close() (error, error) {
	return m.migrator.Close()
}

func NewMigrator(d Deps) (*Migrator, error) {
	url := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		d.Username, d.Password, d.Host, d.Port, d.Database)

	m, err := migrate.New(d.Source, url)
	if err != nil {
		return nil, err
	}

	return &Migrator{m}, nil
}
