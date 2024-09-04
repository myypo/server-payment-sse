package http

import (
	"payment-sse/internal/controller"
	dom "payment-sse/internal/domain"
	"payment-sse/internal/protocol"
	"payment-sse/internal/protocol/http/handler"
	"payment-sse/internal/protocol/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Http[T any] struct {
	conf Config

	handler *handler.Handler[T]
	log     *zap.Logger
}

func NewHttp[T any](
	conf Config,
	domConf dom.Config,
	log *zap.Logger,

	ordCont controller.OrderController[T],
) (protocol.Protocol, error) {
	handler := handler.NewHandler(domConf, ordCont, log)

	return &Http[T]{
		conf,

		handler,
		log,
	}, nil
}

func (p *Http[T]) Listen() error {
	e := gin.Default()

	e.Use(middleware.NewLoggerMiddleware(p.log))

	p.handler.Route(e)

	return e.Run(p.conf.Addr())
}
