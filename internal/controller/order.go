package controller

import (
	"fmt"
	"payment-sse/internal/controller/service"
	dom "payment-sse/internal/domain"
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/env"
	repo "payment-sse/internal/repo"

	"github.com/google/uuid"
)

type OrderController[T any] struct {
	txRepo  repo.TxRepo[T]
	ordRepo repo.OrderRepo[T]

	streamEve service.EventStreamer
}

func NewOrderController[T any](
	txRepo repo.TxRepo[T],
	ordRepo repo.OrderRepo[T],
	streamEve service.EventStreamer,
) OrderController[T] {
	return OrderController[T]{txRepo, ordRepo, streamEve}
}

func (c *OrderController[T]) ListOrders(
	e env.Env,
	filt *domOrd.ListOrders,
) ([]domOrd.Order, dom.DomError) {
	orders, rerr := c.ordRepo.ListOrders(e, filt)
	if rerr != nil {
		e.LogUnexpectedError(rerr, "listing orders")
		return nil, dom.DomErrorFromVerbose(rerr, dom.Internal)
	}

	return orders, nil
}

func (c *OrderController[T]) PaymentWebhook(
	e env.Env,
	pw *domOrd.PaymentWebhook,
) dom.DomError {
	tx, verr := c.txRepo.Begin(e)
	if verr != nil {
		return dom.DomErrorFromVerbose(verr, dom.Internal)
	}

	ord, isNewOrd, rerr := c.ordRepo.EnsureOrderExists(
		tx,
		repo.EnsureOrderExistsFromPaymentHook(pw),
	)
	if rerr != nil {
		tx.Rollback(e)

		e.LogUnexpectedError(rerr, "ensuring order exists")
		return dom.DomErrorFromVerbose(rerr, dom.Internal)
	}
	if !isNewOrd && pw.Status == ord.Status {
		tx.Rollback(e)

		err := fmt.Errorf("the provided order event has already been processed")
		return dom.NewDomError(err, err, dom.Conflict)
	}

	if !domOrd.OrderStatusCompatible(ord.Status, pw.Status, ord.IsFinal) {
		tx.Rollback(e)

		err := fmt.Errorf(
			"the provided order status %s is incompatible with the current state %s",
			pw.Status,
			ord.Status,
		)
		return dom.NewDomError(err, err, dom.Gone)
	}

	if rerr := c.ordRepo.CreateEventOrder(tx, repo.CreateEventOrderFromPaymentHook(pw)); rerr != nil {
		tx.Rollback(e)

		if rerr.Type() == repo.Conflict {
			return dom.DomErrorFromVerbose(rerr, dom.Conflict)
		}

		e.LogUnexpectedError(rerr, "creating new event for order")
		return dom.DomErrorFromVerbose(rerr, dom.Internal)
	}

	tx.Commit(e)
	return nil
}

func (c *OrderController[T]) StreamEvents(
	e env.Env,
	ordId uuid.UUID,
) <-chan domOrd.PaymentEvent {
	return c.streamEve.NewStream(e, ordId)
}
