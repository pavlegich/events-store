// Package entities contains objects for the application.
package entities

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Event contains data for events.
type Event struct {
	EventID   int64      `json:"eventID"`
	EventType string     `json:"eventType"`
	UserID    int64      `json:"userID"`
	EventTime CustomTime `json:"eventTime"`
	Payload   string     `json:"payload"`
}

// CustomTime type for implementing the Marshaler and Unmarshaler interfaces.
type CustomTime struct {
	time.Time
}

const Layout = "2006-01-02 15:04:05"

// UnmarshalJSON implements UnmarshalJSON method for unmarshalling custom time from JSON.
func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		return
	}
	t.Time, err = time.Parse(Layout, s)
	return
}

// MarshalJSON implements MarshalJSON method for marshalling custom time to JSON.
func (t CustomTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(Layout))), nil
}

// Scan implements Scan method for scanning the event time from the storage.
func (t *CustomTime) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	switch data := v.(type) {
	case time.Time:
		t.Time = data
		return nil
	case []byte:
		return json.Unmarshal(data, &t.Time)
	default:
		return fmt.Errorf("cannot scan type %t into Time", v)
	}
}
