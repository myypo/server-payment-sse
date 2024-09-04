package domOrd

import (
	"encoding/json"
	"fmt"
)

type OrderStatus int

const (
	// Missable events
	CoolOrderCreated OrderStatus = iota
	SbuVerificationPending
	ConfirmedByMayor
	// Potentially final
	ChangedMyMind
	Failed
	Chinazes
	GiveMyMoneyBack
)

func OrderStatusFromString(s string) (OrderStatus, error) {
	switch s {
	case "cool_order_created":
		return CoolOrderCreated, nil
	case "sbu_verification_pending":
		return SbuVerificationPending, nil
	case "confirmed_by_mayor":
		return ConfirmedByMayor, nil
	case "changed_my_mind":
		return ChangedMyMind, nil
	case "failed":
		return Failed, nil
	case "chinazes":
		return Chinazes, nil
	case "give_my_money_back":
		return GiveMyMoneyBack, nil
	default:
		return -1, fmt.Errorf("unknown order status provided: %s", s)
	}
}

func (s OrderStatus) String() string {
	return [...]string{
		"cool_order_created",
		"sbu_verification_pending",
		"confirmed_by_mayor",
		"changed_my_mind",
		"failed",
		"chinazes",
		"give_my_money_back",
	}[s]
}

func OrderStatusCompatible(old OrderStatus, new OrderStatus, isFinal bool) bool {
	// Means we are creating a new order
	if new == old {
		return true
	}
	// Can receive non-final events in any order, including the missed ones
	if new <= ConfirmedByMayor {
		return true
	}
	if new == Chinazes && old == GiveMyMoneyBack {
		return true
	}

	if !isFinal {
		if new == Chinazes {
			return true
		}

		// Can change mind on when it is neither final or chinazes
		if new == ChangedMyMind && old != Chinazes {
			return true
		}
		// Can get money back on non-final chinazes
		if new == GiveMyMoneyBack {
			return true
		}
		// Don't allow failed if the most recent event was chinazes
		if new == Failed && old != Chinazes {
			return true
		}
	}

	return false
}

func (os *OrderStatus) Scan(value any) error {
	if value == nil {
		return fmt.Errorf("null value")
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type")
	}

	orderStatus, err := OrderStatusFromString(str)
	if err != nil {
		return err
	}

	*os = orderStatus
	return nil
}

func (os *OrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var err error
	*os, err = OrderStatusFromString(s)
	if err != nil {
		return err
	}

	return nil
}

func (m OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
