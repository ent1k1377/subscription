package common

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

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessfulResponse struct {
	Message string `json:"message"`
}
