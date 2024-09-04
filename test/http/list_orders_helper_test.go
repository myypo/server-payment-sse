package http

import (
	"fmt"
	"net/http"
	"payment-sse/internal/config"
	domOrd "payment-sse/internal/domain/order"
)

type ListOrdersHelper struct {
	conf *config.Config
	t    *Testing
}

func newListOrdersHelper(conf *config.Config, t *Testing) *ListOrdersHelper {
	return &ListOrdersHelper{conf, t}
}

func (s *ListOrdersHelper) ListOrders(qry string) []domOrd.Order {
	var orders []domOrd.Order
	s.t.Get(
		fmt.Sprintf("http://%s/orders?%s", s.conf.Http.Addr(), qry),
		&orders,
		http.StatusOK,
	)
	return orders
}
