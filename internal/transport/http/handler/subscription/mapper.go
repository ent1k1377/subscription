package subscription

import (
	"github.com/google/uuid"
	"subscriptions/internal/domain"
	"time"
)

func ToCreateSubscriptionParams(request *CreateSubscriptionRequest) (*domain.CreateSubscriptionParams, error) {
	uuidParse, err := uuid.Parse(request.UserUUID)
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
