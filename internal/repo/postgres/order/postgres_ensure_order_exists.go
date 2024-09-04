package pgOrd

import (
	"payment-sse/internal/domain/order"
	"payment-sse/internal/repo"

	"github.com/jackc/pgx/v5"
)

func (r *postgresOrder) EnsureOrderExists(
	c repo.TxContext[pgx.Tx],
	co *repo.EnsureOrderExists,
) (*domOrd.Order, bool, repo.RepoError) {
	tx := c.Tx()

	qry := `
        with ins_order as (
            insert into "order"
                (id, user_id, created_at)
            values
                ($1, $2, $3)
            on conflict (id)
                do nothing
            returning
                id, user_id, $4::order_status as status, created_at, $5::timestamptz as updated_at, is_order_final($4, $5, $6) as is_final, true as ins_happened
        )
        select 
            o.id, o.user_id, o.status, o.created_at, o.updated_at, o.is_final, o.ins_happened
        from 
            "ins_order" o
        left join (
            select distinct on (order_id)
                e.order_id,
                e.status,
                e.created_at as updated_at
            from event_order e
            order by e.order_id, e.created_at desc
        ) e on o.id = e.order_id
        union select
            o.id, o.user_id, e.status, o.created_at, e.updated_at, is_order_final(e.status, e.updated_at, $6) as is_final, false as ins_happened
        from 
            "order" o
        left join (
            select distinct on (order_id)
                e.order_id,
                e.status,
                e.created_at as updated_at
            from event_order e
            order by e.order_id, e.created_at desc
        ) e on o.id = e.order_id
        where
            o.id = $1
	       `
	var insHappened bool
	var ord domOrd.Order
	err := tx.QueryRow(c, qry, co.ID, co.UserID, co.CreatedAt, co.Status, co.UpdatedAt, c.DomConf().PaymentConfirmationIn).
		Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.CreatedAt, &ord.UpdatedAt, &ord.IsFinal, &insHappened)
	if err != nil {
		return nil, false, r.NewUnexpectedError(err)
	}

	return &ord, insHappened, nil
}
