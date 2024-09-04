package handler

import (
	"io"
	"payment-sse/internal/protocol/http/dto/request"
	"payment-sse/internal/protocol/http/dto/response"
	httpErr "payment-sse/internal/protocol/http/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler[T]) streamOrderEvents(g *gin.Context) {
	var uri request.SelectorUUID
	if err := h.bindUri(g, &uri); err != nil {
		errSend(g, err)
		return
	}

	ordId, err := uuid.Parse(uri.ID)
	if err != nil {
		errSend(g, httpErr.NewBadRequest(err))
		return
	}
	e := h.envFromGin(g)
	rx := h.ordCont.StreamEvents(e, ordId)

	g.Stream(func(w io.Writer) bool {
		if pay, ok := <-rx; ok {
			g.SSEvent(
				"event_order",
				response.StreamEventFromDom(&pay, e.DomConf().PaymentConfirmationIn),
			)
			return true
		}
		h.envFromGin(g).LogDebug("closed stream!!!")
		return false
	})
}
