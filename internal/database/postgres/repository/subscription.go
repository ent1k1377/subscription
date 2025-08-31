package repository

import (
	"context"
	"subscriptions/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscription struct {
	pool *pgxpool.Pool
}

func NewSubscription(pool *pgxpool.Pool) *Subscription {
	return &Subscription{pool: pool}
}

func (s *Subscription) CreateSubscription(subscription *domain.Subscription) error {
	ctx := context.Background()
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.pool.Exec(ctx, query,
		subscription.UUID.String(),
		subscription.ServiceName,
		subscription.Price,
		subscription.UserUUID.String(),
		subscription.StartDate,
		subscription.EndDate,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Subscription) GetSubscription(uuid uuid.UUID) (*domain.Subscription, error) {
	ctx := context.Background()
	var subscription domain.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	err := s.pool.QueryRow(ctx, query, uuid).Scan(
		&subscription.UUID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserUUID,
		&subscription.StartDate,
		&subscription.EndDate,
	)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}
