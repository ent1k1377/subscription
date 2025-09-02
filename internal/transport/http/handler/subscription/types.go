package subscription

import (
	"github.com/ent1k1377/subscriptions/internal/domain"
	"github.com/ent1k1377/subscriptions/internal/transport/http/common"
)

// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	// Название сервиса
	// Required true
	// Example: Yandex Plus
	ServiceName string            `json:"service_name"`
	Price       int               `json:"price"`
	UserID      string            `json:"user_id"`
	StartDate   common.MonthYear  `json:"start_date"`
	EndDate     *common.MonthYear `json:"end_date,omitempty"`
}

type GetSubscriptionResponse struct {
	ID          string            `json:"id"`
	ServiceName string            `json:"service_name"`
	Price       int               `json:"price"`
	UserID      string            `json:"user_id"`
	StartDate   common.MonthYear  `json:"start_date"`
	EndDate     *common.MonthYear `json:"end_date"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string            `json:"service_name"`
	Price       int               `json:"price"`
	EndDate     *common.MonthYear `json:"end_date,omitempty"`
}

type ListSubscriptionRequest struct {
	ServiceName *string `form:"service_name,omitempty"`
	UserID      *string `form:"user_id,omitempty"`
	Page        int     `form:"page"`
	Limit       int     `form:"limit"`
}

type ListSubscriptionResponse struct {
	Subscriptions []*domain.Subscription `json:"subscriptions"`
}

type TotalCostSubscriptionsRequest struct {
	ServiceName *string          `json:"service_name"`
	UserID      *string          `json:"user_id"`
	StartDate   common.MonthYear `json:"start_date"`
	EndDate     common.MonthYear `json:"end_date"`
}
