package handler

import (
	"payment-sse/internal/controller"
	dom "payment-sse/internal/domain"
	"payment-sse/internal/env"
	"payment-sse/internal/protocol/http/dto/response"
	httpErr "payment-sse/internal/protocol/http/error"
	"payment-sse/internal/protocol/http/middleware"
	"payment-sse/internal/protocol/http/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler[T any] struct {
	domConf dom.Config

	ordCont controller.OrderController[T]
	log     *zap.Logger

	tsServ *service.TranslationService
}

func NewHandler[T any](
	domConf dom.Config,

	ordCont controller.OrderController[T],
	log *zap.Logger,
) *Handler[T] {
	tsServ := service.NewTranslationService()
	return &Handler[T]{
		domConf,
		ordCont,
		log,
		&tsServ,
	}
}

func (h *Handler[T]) Route(e *gin.Engine) {
	{
		orders := e.Group("/orders")
		orders.GET("", h.listOrders)
		orders.GET("/:id/events", middleware.NewStreamHeadersMiddleware(), h.streamOrderEvents)
	}

	{
		wh := e.Group("/webhooks")
		{
			pay := wh.Group("/payments")
			pay.POST("/orders", h.paymentWebhook)
		}
	}
}

func (h *Handler[T]) envFromGin(g *gin.Context) env.Env {
	return env.NewEnv(g, h.log, &h.domConf)
}

func succSend(g *gin.Context, statusCode int, data any) {
	g.JSON(statusCode, response.NewSuccessResponse(data))
}

func errSend(g *gin.Context, err *httpErr.HttpError) {
	g.JSON(err.Code, response.NewErrorResponse(err))
}

func (h *Handler[T]) bindUri(
	g *gin.Context,
	uri any,
) *httpErr.HttpError {
	if err := g.ShouldBindUri(uri); err != nil {
		return h.tsServ.TranslateEN(err)
	}
	return nil
}

func (h *Handler[T]) bindQuery(
	g *gin.Context,
	qry any,
) *httpErr.HttpError {
	if err := g.ShouldBindQuery(qry); err != nil {
		return h.tsServ.TranslateEN(err)
	}
	return nil
}

func (h *Handler[T]) bindJSON(
	g *gin.Context,
	req any,
) *httpErr.HttpError {
	if err := g.ShouldBindJSON(req); err != nil {
		return h.tsServ.TranslateEN(err)
	}
	return nil
}
