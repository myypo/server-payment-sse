package domOrd

import (
	"fmt"
	"payment-sse/internal/domain"
	"payment-sse/internal/util"

	"github.com/google/uuid"
)

type ListOrders struct {
	StatusesOrFinal *u.Either[[]OrderStatus, bool]
	UserID          u.Maybe[uuid.UUID]

	dom.List
}

func NewListOrders(
	mbStatuses []string,
	mbIsFinal u.Maybe[bool],
	mbUserId u.Maybe[uuid.UUID],
	mbLimit u.Maybe[uint],
	mbOffset u.Maybe[uint],
	mbSortBy u.Maybe[string],
	mbSortOrder u.Maybe[string],
) (*ListOrders, dom.DomError) {
	statusOrFinal, err := func() (*u.Either[[]OrderStatus, bool], dom.DomError) {
		if len(mbStatuses) != 0 && mbIsFinal.IsSome() {
			err := fmt.Errorf("can not filter orders by both `status` and `is final`")
			return nil, dom.NewDomError(err, err, dom.BadRequest)
		}

		if len(mbStatuses) > 0 {
			mbStatuses, err := u.MapE(
				mbStatuses,
				func(str string) (OrderStatus, error) { return OrderStatusFromString(str) },
			)
			if err != nil {
				return nil, dom.NewDomError(err, err, dom.BadRequest)
			}

			return u.EitherFromLeft[[]OrderStatus, bool](mbStatuses), nil
		}

		if mbIsFinal, ok := mbIsFinal.Some(); ok {
			return u.EitherFromRight[[]OrderStatus](mbIsFinal), nil
		}

		err := fmt.Errorf("either order `status` filters or `is final` filter has to be provided")
		return nil, dom.NewDomError(err, err, dom.BadRequest)
	}()
	if err != nil {
		return nil, err
	}

	sortByTime, err := func() (dom.SortByTime, dom.DomError) {
		strSortBy, ok := mbSortBy.Some()
		if !ok {
			return dom.CreatedAt, nil
		}

		sortByTime, err := dom.SortByTimeFromString(strSortBy)
		if err != nil {
			return -1, dom.NewDomError(err, err, dom.BadRequest)
		}
		return sortByTime, nil
	}()
	if err != nil {
		return nil, err
	}

	sortOrder, err := func() (dom.SortOrder, dom.DomError) {
		strSortOrder, ok := mbSortOrder.Some()
		if !ok {
			return dom.Descending, nil
		}

		sortOrder, err := dom.SortOrderFromString(strSortOrder)
		if err != nil {
			return -1, dom.NewDomError(err, err, dom.BadRequest)
		}
		return sortOrder, nil
	}()
	if err != nil {
		return nil, err
	}

	return &ListOrders{
		StatusesOrFinal: statusOrFinal,
		UserID:          mbUserId,

		List: dom.List{
			Limit:      u.OR(mbLimit, 10),
			Offset:     u.OR(mbOffset, 0),
			SortByTime: sortByTime,
			SortOrder:  sortOrder,
		},
	}, nil
}
