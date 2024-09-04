package repo

import (
	"payment-sse/internal/domain/order"
	"payment-sse/internal/env"
	"time"

	"github.com/google/uuid"
)

type OrderRepo[T any] interface {
	// Inserts a new order, or, if order with such ID already exists returns the old order,
	// bool signals if an insert happened
	EnsureOrderExists(c TxContext[T], co *EnsureOrderExists) (*domOrd.Order, bool, RepoError)
	ListOrders(e env.Env, filt *domOrd.ListOrders) ([]domOrd.Order, RepoError)

	CreateEventOrder(c TxContext[T], ce *CreateEventOrder) RepoError
	GetChronOrdersEvents(
		e env.Env,
		ordIds []uuid.UUID,
	) (map[uuid.UUID]struct {
		UserID    uuid.UUID
		CreatedAt time.Time

		Events []domOrd.EventOrder
	}, RepoError)
}

type EnsureOrderExists struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Status domOrd.OrderStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateEventOrder struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Status  domOrd.OrderStatus

	CreatedAt time.Time
}

func EnsureOrderExistsFromPaymentHook(
	pw *domOrd.PaymentWebhook,
) *EnsureOrderExists {
	return &EnsureOrderExists{
		ID:     pw.OrderID,
		UserID: pw.UserID,
		Status: pw.Status,

		CreatedAt: pw.CreatedAt,
		UpdatedAt: pw.UpdatedAt,
	}
}

func CreateEventOrderFromPaymentHook(
	pw *domOrd.PaymentWebhook,
) *CreateEventOrder {
	return &CreateEventOrder{
		ID:      pw.EventID,
		OrderID: pw.OrderID,
		Status:  pw.Status,

		// Count order update as the time of the event creation
		CreatedAt: pw.UpdatedAt,
	}
}
