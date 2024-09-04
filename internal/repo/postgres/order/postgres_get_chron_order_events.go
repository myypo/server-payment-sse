package pgOrd

import (
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/env"
	"payment-sse/internal/repo"
	"payment-sse/internal/util"
	"time"

	"github.com/google/uuid"
)

// Could have intergated postgres replication instead, I think
func (r *postgresOrder) GetChronOrdersEvents(
	e env.Env,
	ordIds []uuid.UUID,
) (map[uuid.UUID]struct {
	UserID    uuid.UUID
	CreatedAt time.Time

	Events []domOrd.EventOrder
}, repo.RepoError,
) {
	qry := `
        select 
            o.id, o.user_id, o.created_at,
            array_agg((e.id, e.order_id, e.status, e.created_at) order by e.created_at) as orders
        from 
            "order" o
        left join
            "event_order" e 
        on
            o.id = e.order_id
        where
            o.id = any($1::uuid[])
        group by
            o.id
        `
	rows, err := r.db.Query(e, qry, ordIds)
	if err != nil {
		return nil, r.NewError(err)
	}
	defer rows.Close()

	orders := make(map[uuid.UUID]struct {
		UserID    uuid.UUID
		CreatedAt time.Time

		Events []domOrd.EventOrder
	})
	for rows.Next() {
		type rawEvent struct {
			ID      uuid.UUID
			OrderID uuid.UUID
			// HACK: despite implementing Scan OrderStatus can't be scanned here when used in a slice
			Status    string
			CreatedAt time.Time
		}
		type rawOrder struct {
			UserID    uuid.UUID
			CreatedAt time.Time

			Events []rawEvent
		}

		var raw rawOrder
		var id uuid.UUID
		err := rows.Scan(
			&id,
			&raw.UserID,
			&raw.CreatedAt,
			&raw.Events,
		)
		if err != nil {
			return nil, r.NewUnexpectedError(err)
		}
		e.LogDebug(raw)

		events, err := u.MapE(raw.Events, func(r rawEvent) (domOrd.EventOrder, error) {
			status, err := domOrd.OrderStatusFromString(r.Status)
			if err != nil {
				return domOrd.EventOrder{}, err
			}
			return domOrd.NewEventOrder(r.ID, r.OrderID, status, r.CreatedAt), nil
		})
		if err != nil {
			return nil, r.NewUnexpectedError(err)
		}

		orders[id] = struct {
			UserID    uuid.UUID
			CreatedAt time.Time

			Events []domOrd.EventOrder
		}{
			UserID:    raw.UserID,
			CreatedAt: raw.CreatedAt,

			Events: events,
		}
	}

	return orders, nil
}
