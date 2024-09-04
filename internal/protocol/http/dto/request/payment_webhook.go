package request

import (
	dom "payment-sse/internal/domain"
	domOrd "payment-sse/internal/domain/order"
	"time"

	"github.com/google/uuid"
)

type PaymentWebhook struct {
	OrderID   uuid.UUID `json:"order_id"   binding:"required,uuid"`
	EventID   uuid.UUID `json:"event_id"   binding:"required,uuid"`
	UserID    uuid.UUID `json:"user_id"    binding:"required,uuid"`
	Status    string    `json:"status"     binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
}

func PaymentWebhookFromRequest(req *PaymentWebhook) (*domOrd.PaymentWebhook, dom.DomError) {
	return domOrd.NewPaymentWebhook(
		req.OrderID,
		req.EventID,
		req.UserID,
		req.Status,

		req.CreatedAt,
		req.UpdatedAt,
	)
}
