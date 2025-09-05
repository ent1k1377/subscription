package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ent1k1377/subscriptions/internal/database/postgres/repository"
	"github.com/ent1k1377/subscriptions/internal/domain"
	"github.com/ent1k1377/subscriptions/internal/transport/http/middleware"

	"github.com/google/uuid"
)

type Subscription struct {
	logger           *slog.Logger
	subscriptionRepo *repository.Subscription
}

func NewSubscription(baseLogger *slog.Logger, subscriptionRepo *repository.Subscription) *Subscription {
	logger := baseLogger.WithGroup("subscription service")

	return &Subscription{
		logger:           logger,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *Subscription) CreateSubscription(ctx context.Context, params *domain.CreateSubscriptionParams) error {
	logger := s.logger.With("request_id", ctx.Value(middleware.RequestIDKey).(string))

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

	logger.Info("Creating subscription",
		slog.Any("user_id", subscription.UserUUID),
		slog.String("service_name", subscription.ServiceName),
		slog.Int("price", subscription.Price),
	)

	err := s.subscriptionRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		logger.Error("Failed to create subscription",
			slog.String("error", err.Error()),
		)
		return err
	}

	logger.Info("Finish create subscription",
		slog.Any("user_id", subscription.UserUUID),
	)
	return nil
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
