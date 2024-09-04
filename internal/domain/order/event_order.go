package domOrd

import (
	"time"

	"github.com/google/uuid"
)

type EventOrder struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Status  OrderStatus

	CreatedAt time.Time
}

func NewEventOrder(eveId, ordId uuid.UUID, status OrderStatus, createdAt time.Time) EventOrder {
	return EventOrder{ID: eveId, OrderID: ordId, Status: status, CreatedAt: createdAt}
}

func (e *EventOrder) Final(chinTimeout time.Duration) bool {
	switch e.Status {
	case ChangedMyMind, Failed, GiveMyMoneyBack:
		return true
	case Chinazes:
		if time.Since(e.CreatedAt) >= chinTimeout {
			return true
		}
		return false
	default:
		return false
	}
}

type PaymentEvent struct {
	OrderID   uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time

	Event EventOrder
}

func NewPaymentEvent(
	ordId, userId uuid.UUID,
	createdAt time.Time,
	event EventOrder,
) PaymentEvent {
	return PaymentEvent{OrderID: ordId, UserID: userId, CreatedAt: createdAt, Event: event}
}
