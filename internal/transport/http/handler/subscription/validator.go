package subscription

import (
	"fmt"

	"github.com/google/uuid"
)

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e *ValidationError) Add(field, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string]string)
	}
	e.Errors[field] = message
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

func ValidateCreateSubscriptionRequest(req CreateSubscriptionRequest) error {
	var vErr ValidationError

	if req.ServiceName == "" {
		vErr.Add("service_name", "required")
	}
	if req.Price < 0 {
		vErr.Add("price", "must be greater than or equal to zero")
	}

	_, err := uuid.Parse(req.UserID)
	if err != nil {
		vErr.Add("user_id", "must be a valid UUID")
	}

	if vErr.HasErrors() {
		return err
	}

	return nil
}
