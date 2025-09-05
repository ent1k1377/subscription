package subscription

import (
	"github.com/ent1k1377/subscriptions/internal/domain"
	"github.com/ent1k1377/subscriptions/internal/transport/http/common"
)

// CreateSubscriptionRequest request структура для создания новой подписки
type CreateSubscriptionRequest struct {
	ServiceName string            `json:"service_name" example:"Netflix" binding:"required"`
	Price       int               `json:"price" example:"999" binding:"required,gte=0"`
	UserID      string            `json:"user_id" example:"f81d4fae-7dec-11d0-a765-00a0c91e6bf6" binding:"required,uuid"`
	StartDate   common.MonthYear  `json:"start_date" example:"01-2025" binding:"required"`
	EndDate     *common.MonthYear `json:"end_date,omitempty" example:"12-2025"`
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
