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

type UpdateSubscriptionParams struct {
	ServiceName string
	Price       int
	EndDate     *time.Time
}

type ListSubscriptionParams struct {
	ServiceName string
	UserID      uuid.UUID
	Page        int
	Limit       int
}
