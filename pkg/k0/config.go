package k0

import (
	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/types"
)

type KeylogmeZeroConfig struct {
	Keylog  Config        `json:"keylog"`
	Storage ConfigStorage `json:"storage"`
}

type Config struct {
	Devices        []keylogger.DeviceInput    `json:"devices"`
	ShortcutGroups []types.ShortcutGroupInput `json:"shortcut_groups"`
	ShiftState     types.ShiftStateInput      `json:"shift_state"`
}
