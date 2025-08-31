package subscription

import (
	"encoding/json"
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

func (my *MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(*my)
	formatted := t.Format("01-2006")
	return json.Marshal(formatted)
}

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type GetSubscriptionResponse struct {
	ID          string     `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}
