package service

import (
	"github.com/google/uuid"
	"subscriptions/internal/database/postgres/repository"
	"subscriptions/internal/domain"
	"time"
)

type Subscription struct {
	subscriptionRepo *repository.Subscription
}

func NewSubscription(subscriptionRepo *repository.Subscription) *Subscription {
	return &Subscription{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *Subscription) CreateSubscription(params *domain.CreateSubscriptionParams) error {
	subscription := &domain.Subscription{
		UUID:        uuid.New(),
		ServiceName: params.ServiceName,
		Price:       params.Price,
		UserUUID:    params.UserUUID,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	return s.subscriptionRepo.CreateSubscription(subscription)
}
