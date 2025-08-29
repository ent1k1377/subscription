package subscription

import (
	"strings"
	"time"
)

type MonthYear time.Time

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	
	*my = MonthYear(t)
	
	return nil
}

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserUUID    string     `json:"user_uuid"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}
