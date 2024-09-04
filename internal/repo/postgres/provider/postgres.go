package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresProvider struct {
	pool *pgxpool.Pool
	conf Config
}

func NewPostgresProvider(
	conf Config,
) (*PostgresProvider, error) {
	pc, err := pgxpool.ParseConfig(conf.DBURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), pc)
	if err != nil {
		return nil, err
	}

	return &PostgresProvider{pool, conf}, nil
}

func (p *PostgresProvider) Conn() *pgxpool.Pool {
	return p.pool
}

func (p *PostgresProvider) MigrateUp() error {
	goose.SetDialect("postgres")
	return goose.RunContext(
		context.Background(),
		"up",
		stdlib.OpenDBFromPool(p.pool),
		p.conf.MigrationsPath,
	)
}
