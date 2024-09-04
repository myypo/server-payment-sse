package handler

import (
	"net/http"
	"payment-sse/internal/protocol/http/dto/request"
	httpErr "payment-sse/internal/protocol/http/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler[T]) paymentWebhook(g *gin.Context) {
	var body request.PaymentWebhook
	if err := h.bindJSON(g, &body); err != nil {
		errSend(g, err)
		return
	}

	payHook, err := request.PaymentWebhookFromRequest(&body)
	if err != nil {
		errSend(g, httpErr.NewBadRequest(err))
		return
	}

	err = h.ordCont.PaymentWebhook(h.envFromGin(g), payHook)
	if err != nil {
		errSend(g, httpErr.HttpErrorFromDom(err))
		return
	}

	succSend(g, http.StatusOK, nil)
}
