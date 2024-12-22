package keylog

type Config struct {
	Devices   []DeviceInput   `json:"devices"`
	Shortcuts []ShortcutCodes `json:"shortcuts"`
}
