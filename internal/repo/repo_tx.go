package repo

import (
	"context"
	"payment-sse/internal/env"
	verbErr "payment-sse/internal/error/verbose"
)

type transaction[T any] interface {
	Commit(c context.Context) verbErr.VerboseError
	Rollback(c context.Context) verbErr.VerboseError

	Tx() T
}

type TxRepo[T any] interface {
	Begin(c env.Env) (TxContext[T], verbErr.VerboseError)
}

type TxContext[T any] interface {
	env.Env
	transaction[T]
}
