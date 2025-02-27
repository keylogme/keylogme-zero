package keylog

import "github.com/keylogme/keylogme-zero/v1/keylog/storage"

type Config struct {
	Devices        []DeviceInput   `json:"devices"`
	ShortcutGroups []ShortcutGroup `json:"shortcut_groups"`
}

type KeylogmeZeroConfig struct {
	Keylog  Config                `json:"keylog"`
	Storage storage.ConfigStorage `json:"storage"`
}
