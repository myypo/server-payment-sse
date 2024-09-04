package env

import (
	"context"
	"fmt"

	dom "payment-sse/internal/domain"
	verbErr "payment-sse/internal/error/verbose"

	"go.uber.org/zap"
)

type Env interface {
	context.Context
	LogUnexpectedError(err verbErr.VerboseError, while string)
	LogDebug(vs ...any)
	DomConf() *dom.Config
}

type env struct {
	context.Context
	*zap.Logger
	*dom.Config
}

func NewEnv(c context.Context, log *zap.Logger, conf *dom.Config) Env {
	return &env{
		Context: c,
		Logger:  log,
		Config:  conf,
	}
}

func (e *env) LogUnexpectedError(err verbErr.VerboseError, when string) {
	e.Error(fmt.Sprintf("Unexpected error occurred when %s", when), zap.Error(err.Verbose()))
}

func (e *env) LogDebug(vs ...any) {
	e.Debug(fmt.Sprintf("Debug log %s", vs))
}

func (e *env) DomConf() *dom.Config {
	return e.Config
}
