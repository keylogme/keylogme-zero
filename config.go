package keylog

import (
	"github.com/keylogme/keylogme-zero/storage"
	"github.com/keylogme/keylogme-zero/types"
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
	ThresholdAuto types.Duration `json:"threshold_auto"`
}
