package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	UUID        uuid.UUID
	ServiceName string
	Price       int
	UserUUID    uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateSubscriptionParams struct {
	ServiceName string
	Price       int
	UserUUID    uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}
