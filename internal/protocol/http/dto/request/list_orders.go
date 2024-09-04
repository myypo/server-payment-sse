package request

import (
	dom "payment-sse/internal/domain"
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/util"

	"github.com/google/uuid"
)

type ListOrders struct {
	Status  CommaSepArray      `form:"status"   binding:"omitempty"`
	IsFinal u.Maybe[bool]      `form:"is_final" binding:"omitempty"`
	UserID  u.Maybe[uuid.UUID] `form:"user_id"  binding:"omitempty,uuid"`

	Limit     u.Maybe[uint]   `form:"limit"      binding:"omitempty"`
	Offset    u.Maybe[uint]   `form:"offset"     binding:"omitempty"`
	SortBy    u.Maybe[string] `form:"sort_by"    binding:"omitempty"`
	SortOrder u.Maybe[string] `form:"sort_order" binding:"omitempty"`
}

func ListOrdersFromRequest(req *ListOrders) (*domOrd.ListOrders, dom.DomError) {
	return domOrd.NewListOrders(
		req.Status.Values(),
		req.IsFinal,
		req.UserID,

		req.Limit,
		req.Offset,
		req.SortBy,
		req.SortOrder,
	)
}
