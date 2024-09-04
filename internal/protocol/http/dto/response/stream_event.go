package response

import (
	domOrd "payment-sse/internal/domain/order"
	"time"

	"github.com/google/uuid"
)

type StreamEvent struct {
	OrderID     uuid.UUID          `json:"order_id"`
	UserID      uuid.UUID          `json:"user_id"`
	OrderStatus domOrd.OrderStatus `json:"order_status"`
	IsFinal     bool               `json:"is_final"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func StreamEventFromDom(pay *domOrd.PaymentEvent, chinTimout time.Duration) *StreamEvent {
	return &StreamEvent{
		OrderID:     pay.OrderID,
		UserID:      pay.UserID,
		OrderStatus: pay.Event.Status,
		IsFinal:     pay.Event.Final(chinTimout),
		CreatedAt:   pay.CreatedAt,
		UpdatedAt:   pay.Event.CreatedAt,
	}
}
