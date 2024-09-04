package http

import (
	"fmt"
	"payment-sse/internal/config"
	"payment-sse/internal/protocol/http/dto/request"
)

type PaymentHelper struct {
	conf *config.Config
	t    *Testing
}

func newPaymentHelper(conf *config.Config, t *Testing) *PaymentHelper {
	return &PaymentHelper{conf, t}
}

func (s *PaymentHelper) SendPayment(pw *request.PaymentWebhook, status int) {
	s.t.Post(
		fmt.Sprintf("http://%s/webhooks/payments/orders", s.conf.Http.Addr()), pw, status)
}
