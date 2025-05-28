package k0

import (
	"github.com/keylogme/keylogme-zero/storage"
)

type KeylogmeZeroConfig struct {
	Keylog  Config                `json:"keylog"`
	Storage storage.ConfigStorage `json:"storage"`
}

type Config struct {
	Devices        []DeviceInput        `json:"devices"`
	ShortcutGroups []ShortcutGroupInput `json:"shortcut_groups"`
	ShiftState     ShiftStateInput      `json:"shift_state"`
}
