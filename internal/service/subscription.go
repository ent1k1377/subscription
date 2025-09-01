package service

import (
	"subscriptions/internal/database/postgres/repository"
	"subscriptions/internal/domain"
	"time"

	"github.com/google/uuid"
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

func (s *Subscription) GetSubscription(uuid uuid.UUID) (*domain.Subscription, error) {
	return s.subscriptionRepo.GetSubscription(uuid)
}

func (s *Subscription) UpdateSubscription(uuid uuid.UUID, params *domain.UpdateSubscriptionParams) error {
	return s.subscriptionRepo.UpdateSubscription(uuid, params)
}

func (s *Subscription) DeleteSubscription(uuid uuid.UUID) error {
	return s.subscriptionRepo.DeleteSubscription(uuid)
}

func (s *Subscription) ListSubscriptions(params *domain.ListSubscriptionParams) ([]*domain.Subscription, error) {
	return s.subscriptionRepo.ListSubscriptions(params)
}

func (s *Subscription) TotalCostSubscriptions(params *domain.TotalCostSubscriptionsParams) (int, error) {
	return s.subscriptionRepo.TotalCostSubscriptions(params)
}
