package subscription

import (
	"subscriptions/internal/domain"
	"time"

	"github.com/google/uuid"
)

func ToCreateSubscriptionParams(request *CreateSubscriptionRequest) (*domain.CreateSubscriptionParams, error) {
	uuidParse, err := uuid.Parse(request.UserID)
	if err != nil {
		return nil, err
	}

	startDate := time.Time(request.StartDate)
	var endDate *time.Time
	if request.EndDate != nil {
		endTime := time.Time(*request.EndDate)
		endDate = &endTime
	}

	return &domain.CreateSubscriptionParams{
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserUUID:    uuidParse,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func ToGetSubscriptionResponse(subscription *domain.Subscription) *GetSubscriptionResponse {
	startDate := MonthYear(subscription.StartDate)
	var endDate *MonthYear
	if subscription.EndDate != nil {
		endTime := MonthYear(*subscription.EndDate)
		endDate = &endTime
	}

	return &GetSubscriptionResponse{
		ID:          subscription.UUID.String(),
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserUUID.String(),
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func ToUpdateSubscriptionParams(request *UpdateSubscriptionRequest) *domain.UpdateSubscriptionParams {
	var endDate *time.Time
	if request.EndDate != nil {
		endTime := time.Time(*request.EndDate)
		endDate = &endTime
	}

	return &domain.UpdateSubscriptionParams{
		ServiceName: request.ServiceName,
		Price:       request.Price,
		EndDate:     endDate,
	}
}

func ToListSubscriptionParams(request *ListSubscriptionRequest) (*domain.ListSubscriptionParams, error) {
	var userID *uuid.UUID
	if request.UserID != nil {
		id, err := uuid.Parse(*request.UserID)
		if err != nil {
			return nil, err
		}
		userID = &id
	}

	return &domain.ListSubscriptionParams{
		ServiceName: request.ServiceName,
		UserID:      userID,
		Page:        request.Page,
		Limit:       request.Limit,
	}, nil
}

func ToListSubscriptionResponse(subscriptions []*domain.Subscription) *ListSubscriptionResponse {
	return &ListSubscriptionResponse{
		Subscriptions: subscriptions,
	}
}
