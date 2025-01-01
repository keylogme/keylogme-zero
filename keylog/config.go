package keylog

import "github.com/keylogme/zero-trust-logger/keylog/storage"

type Config struct {
	Devices   []DeviceInput   `json:"devices"`
	Shortcuts []ShortcutCodes `json:"shortcuts"`
}

type KeylogmeZeroConfig struct {
	Keylog  Config                `json:"keylog"`
	Storage storage.ConfigStorage `json:"storage"`
}
