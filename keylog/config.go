package keylog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/keylogme/keylogme-zero/v1/keylog/storage"
)

type KeylogmeZeroConfig struct {
	Keylog  Config                `json:"keylog"`
	Storage storage.ConfigStorage `json:"storage"`
}

type Config struct {
	Devices        []DeviceInput   `json:"devices"`
	ShortcutGroups []ShortcutGroup `json:"shortcut_groups"`
	ShiftState     ShiftState      `json:"shift_state"`
	Layers         []Layer         `json:"layers"`
}

type ShiftState struct {
	ThresholdAuto Duration `json:"threshold_auto"`
}

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
