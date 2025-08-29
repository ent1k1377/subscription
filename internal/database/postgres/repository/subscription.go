package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"subscriptions/internal/domain"
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
