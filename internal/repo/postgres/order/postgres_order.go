package pgOrd

import (
	"payment-sse/internal/repo"
	pg "payment-sse/internal/repo/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresOrder struct {
	db *pgxpool.Pool
}

const orderTableName = "order"

func NewOrderPostgres(db *pgxpool.Pool) repo.OrderRepo[pgx.Tx] {
	return &postgresOrder{db}
}

func (r *postgresOrder) NewUnexpectedError(techErr error) repo.RepoError {
	return repo.NewUnexpectedRepoError(techErr, orderTableName)
}

func (r *postgresOrder) NewError(techErr error, expTypes ...repo.RepoErrorType) repo.RepoError {
	return pg.RepoErrorFromPostgres(techErr, orderTableName, expTypes...)
}
