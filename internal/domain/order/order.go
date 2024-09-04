package domOrd

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Status  OrderStatus
	IsFinal bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
