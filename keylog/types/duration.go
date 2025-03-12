package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration wraps time.Duration to allow custom unmarshaling
type Duration struct {
	time.Duration
}

// UnmarshalJSON implements custom unmarshaling for Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		// Assume the number is in milliseconds
		d.Duration = time.Duration(value) * time.Millisecond
	case string:
		// Parse using time.ParseDuration
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		d.Duration = duration
	default:
		return fmt.Errorf("invalid duration format")
	}
	return nil
}
