package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

func (s *Subscription) UpdateSubscription(uuid uuid.UUID, params *domain.UpdateSubscriptionParams) error {
	ctx := context.Background()
	query := `UPDATE subscriptions SET service_name=$1, price=$2, end_date=$3 WHERE id = $4`
	_, err := s.pool.Exec(ctx, query, params.ServiceName, params.Price, params.EndDate, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subscription) DeleteSubscription(uuid uuid.UUID) error {
	ctx := context.Background()
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subscription) ListSubscriptions(params *domain.ListSubscriptionParams) ([]*domain.Subscription, error) {
	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args := s.buildListSubscriptionsQuery(params)
	fmt.Println(query)
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]*domain.Subscription, 0)
	for rows.Next() {
		var subscription domain.Subscription
		_ = rows.Scan(
			&subscription.UUID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserUUID,
			&subscription.StartDate,
			&subscription.EndDate,
		)
		subscriptions = append(subscriptions, &subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *Subscription) buildListSubscriptionsQuery(params *domain.ListSubscriptionParams) (string, []any) {
	var conditions []string
	var args []any
	pos := 1

	if params.ServiceName != nil && *params.ServiceName != "" {
		conditions = append(conditions, "service_name ILIKE $"+strconv.Itoa(pos))
		args = append(args, *params.ServiceName)
		pos++
	}

	if params.UserID != nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(pos))
		args = append(args, *params.UserID)
		pos++
	}

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " LIMIT $" + strconv.Itoa(pos) + " OFFSET $" + strconv.Itoa(pos+1)
	args = append(args, params.Page, params.Limit)

	return query, args
}

func (s *Subscription) TotalCostSubscriptions(params *domain.TotalCostSubscriptionsParams) (int, error) {
	ctx := context.Background()
	var totalCost int
	fmt.Println(params)
	query, args := s.buildTotalCostSubscriptionsQuery(params)
	fmt.Println(query)
	err := s.pool.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, err
	}

	return totalCost, nil
}

func (s *Subscription) buildTotalCostSubscriptionsQuery(params *domain.TotalCostSubscriptionsParams) (string, []any) {
	var conditions []string
	var args []any
	pos := 1

	if params.ServiceName != nil && *params.ServiceName != "" {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(pos))
		args = append(args, *params.ServiceName)
		pos++
	}

	if params.UserID != nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(pos))
		args = append(args, *params.UserID)
		pos++
	}

	conditions = append(conditions, "start_date >= $"+strconv.Itoa(pos))
	args = append(args, params.StartDate)
	pos++

	conditions = append(conditions, "(end_date IS NULL OR end_date <= $"+strconv.Itoa(pos)+")")
	args = append(args, params.EndDate)
	pos++

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE `
	if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}

	return query, args
}
