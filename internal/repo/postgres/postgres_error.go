package pg

import (
	"errors"
	"payment-sse/internal/repo"

	"github.com/jackc/pgx/v5/pgconn"
)

func RepoErrorFromPostgres(err error, tableName string, expTypes ...repo.RepoErrorType,
) repo.RepoError {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return repo.NewUnexpectedRepoError(err, tableName)
	}

	typ := func() repo.RepoErrorType {
		switch pgErr.Code {
		case uniqueViolation:
			return repo.Conflict
		case doesNotExist:
			return repo.NotFound
		case constraintViolation:
			return repo.BadRequest
		}

		return repo.Internal
	}()

	return repo.NewRepoError(
		typ,
		err,
		tableName,
		pgErr.ConstraintName,
		expTypes...,
	)
}

func TypeFromError(err error) repo.RepoErrorType {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return repo.Internal
	}

	switch pgErr.Code {
	case uniqueViolation:
		return repo.Conflict
	case doesNotExist:
		return repo.NotFound
	case constraintViolation:
		return repo.BadRequest
	default:
		return repo.Internal
	}
}

func columnHumanReadable(pgErr pgconn.PgError) string {
	switch pgErr.ColumnName {
	case "event_id":
		return "event ID"
	}

	return pgErr.ColumnName
}

func constraintHumanReadable(pgErr pgconn.PgError) string {
	switch pgErr.ConstraintName {
	case "event_order_pkey":
		return "event ID has to be unique"
	case "order_pkey":
		return "order ID has to be unique"
	}

	return pgErr.ConstraintName
}

const (
	uniqueViolation     = "23505"
	doesNotExist        = "23503"
	constraintViolation = "23514"
)
