package pgOrd

import (
	"fmt"
	dom "payment-sse/internal/domain"
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/env"
	"payment-sse/internal/repo"
	"payment-sse/internal/util"
	"strings"
)

func (r *postgresOrder) ListOrders(
	e env.Env,
	filt *domOrd.ListOrders,
) ([]domOrd.Order, repo.RepoError) {
	base := `
        select
            o.id, o.user_id, e.status, is_order_final(e.status, e.created_at, $1) as is_final, o.created_at, e.created_at
        from
            "order" o
        left join (
            select distinct on (order_id)
                e.order_id,
                e.status,
                e.created_at
            from event_order e
            order by e.order_id, e.created_at desc
        ) e on o.id = e.order_id
	    `
	args := []any{e.DomConf().PaymentConfirmationIn}
	where := make([]string, 0)

	if statuses, ok := filt.StatusesOrFinal.Left(); ok {
		where = append(where, fmt.Sprintf("status = any ($%v)", len(args)+1))
		args = append(
			args,
			u.Map(statuses, func(s domOrd.OrderStatus) string { return s.String() }),
		)
	}
	if isFinal, ok := filt.StatusesOrFinal.Right(); ok {
		if isFinal {
			where = append(
				where,
				fmt.Sprintf("is_order_final(e.status, e.created_at, $1)"),
			)
		} else {
			where = append(
				where,
				fmt.Sprintf("not is_order_final(e.status, e.created_at, $1)"),
			)
		}
	}
	if userId, ok := filt.UserID.Some(); ok {
		where = append(where, fmt.Sprintf("o.user_id = $%v", len(args)+1))
		args = append(args, userId)
	}

	if len(where) > 0 {
		base = base + " where " + strings.Join(where, " and ")
	}

	base = base + " group by o.id, e.status, e.created_at"

	base = base + fmt.Sprintf(" order by $%v %v", len(args)+1, filt.SortOrder.String())
	byTime := func() string {
		if filt.SortByTime == dom.CreatedAt {
			return fmt.Sprintf("o.%s", filt.SortByTime.String())
		}
		if filt.SortByTime == dom.UpdatedAt {
			return fmt.Sprintf("e.%s", filt.SortByTime.String())
		}
		return fmt.Sprintf("o.%s", filt.SortByTime.String())
	}()
	args = append(args, byTime)

	base = base + fmt.Sprintf(" limit $%v", len(args)+1)
	args = append(args, filt.Limit)
	base = base + fmt.Sprintf(" offset $%v", len(args)+1)
	args = append(args, filt.Offset)

	e.LogDebug(base)
	rows, err := r.db.Query(e, base, args...)
	if err != nil {
		return nil, r.NewError(err)
	}
	defer rows.Close()

	orders := make([]domOrd.Order, 0)
	for rows.Next() {
		var ord domOrd.Order
		err := rows.Scan(
			&ord.ID,
			&ord.UserID,
			&ord.Status,
			&ord.IsFinal,
			&ord.CreatedAt,
			&ord.UpdatedAt,
		)
		if err != nil {
			return nil, r.NewUnexpectedError(err)
		}

		orders = append(orders, ord)
	}

	return orders, nil
}
