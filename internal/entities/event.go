// Package entities contains objects for the application.
package entities

import (
	"fmt"
	"strings"
	"time"
)

// Event contains data for events.
type Event struct {
	EventID   int64      `json:"id"`
	EventType string     `json:"eventType"`
	UserID    int64      `json:"userID"`
	EventTime CustomTime `json:"eventTime"`
	Payload   string     `json:"payload"`
}

// CustomTime type for implementing the Marshaler and Unmarshaler interfaces.
type CustomTime struct {
	time.Time
}

const layout = "2006-01-02 15:04:05"

// UnmarshalJSON implements UnmarshalJSON method for unmarshalling custom time from JSON.
func (c *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		return
	}
	c.Time, err = time.Parse(layout, s)
	return
}

// MarshalJSON implements MarshalJSON method for marshalling custom time to JSON.
func (c CustomTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(layout))), nil
}
