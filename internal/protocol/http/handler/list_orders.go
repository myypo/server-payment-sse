package handler

import (
	"net/http"
	"payment-sse/internal/protocol/http/dto/request"
	httpErr "payment-sse/internal/protocol/http/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler[T]) listOrders(g *gin.Context) {
	var qry request.ListOrders
	if err := h.bindQuery(g, &qry); err != nil {
		errSend(g, err)
		return
	}

	listOrd, err := request.ListOrdersFromRequest(&qry)
	if err != nil {
		errSend(g, httpErr.NewBadRequest(err))
		return

	}

	resp, err := h.ordCont.ListOrders(h.envFromGin(g), listOrd)
	if err != nil {
		errSend(g, httpErr.HttpErrorFromDom(err))
		return
	}

	succSend(g, http.StatusOK, resp)
}
