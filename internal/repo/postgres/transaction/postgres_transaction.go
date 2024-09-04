package pgTx

import (
	"context"
	"payment-sse/internal/env"
	verbErr "payment-sse/internal/error/verbose"
	"payment-sse/internal/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgTx struct {
	env.Env
	tx pgx.Tx
}

func (t *pgTx) Commit(c context.Context) verbErr.VerboseError {
	if err := t.tx.Commit(c); err != nil {
		return verbErr.DefaultVerboseError(err)
	}

	return nil
}

func (t *pgTx) Rollback(c context.Context) verbErr.VerboseError {
	if err := t.tx.Rollback(c); err != nil {
		return verbErr.DefaultVerboseError(err)
	}

	return nil
}

func (t *pgTx) Tx() pgx.Tx {
	return t.tx
}

func newPgTx(e env.Env, tx pgx.Tx) repo.TxContext[pgx.Tx] {
	return &pgTx{e, tx}
}

type pgTxProvider struct {
	pool *pgxpool.Pool
}

func (p *pgTxProvider) Begin(e env.Env) (repo.TxContext[pgx.Tx], verbErr.VerboseError) {
	tx, err := p.pool.BeginTx(e, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, verbErr.DefaultVerboseError(err)
	}

	return newPgTx(e, tx), nil
}

func NewPGTransaction(pool *pgxpool.Pool) repo.TxRepo[pgx.Tx] {
	return &pgTxProvider{pool}
}
