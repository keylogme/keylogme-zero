package keylog

import "github.com/keylogme/keylogme-zero/keylog/storage"

type Config struct {
	Devices        []DeviceInput   `json:"devices"`
	ShortcutGroups []ShortcutGroup `json:"shortcut_groups"`
}

type KeylogmeZeroConfigV1 struct {
	Keylog  Config                `json:"keylog"`
	Storage storage.ConfigStorage `json:"storage"`
}
