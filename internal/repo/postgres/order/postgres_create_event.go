package pgOrd

import (
	"payment-sse/internal/repo"

	"github.com/jackc/pgx/v5"
)

func (r *postgresOrder) CreateEventOrder(
	c repo.TxContext[pgx.Tx],
	ce *repo.CreateEventOrder,
) repo.RepoError {
	tx := c.Tx()

	qry := `
	        insert into "event_order"
	            (id, order_id, status, created_at)
	        values
	            ($1, $2, $3, $4)
	       `

	_, err := tx.Exec(c, qry, ce.ID, ce.OrderID, ce.Status, ce.CreatedAt)
	if err != nil {
		return r.NewError(err, repo.Conflict)
	}

	return nil
}
