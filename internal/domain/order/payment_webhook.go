package domOrd

import (
	dom "payment-sse/internal/domain"
	"time"

	"github.com/google/uuid"
)

type PaymentWebhook struct {
	OrderID uuid.UUID
	EventID uuid.UUID
	UserID  uuid.UUID
	Status  OrderStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPaymentWebhook(
	ordId, eventId, userId uuid.UUID,
	rawStatus string,

	createdAt, updatedAt time.Time,
) (*PaymentWebhook, dom.DomError) {
	status, err := OrderStatusFromString(rawStatus)
	if err != nil {
		return nil, dom.NewBadRequest(err)
	}

	return &PaymentWebhook{
		OrderID: ordId,
		EventID: eventId,
		UserID:  userId,
		Status:  status,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
