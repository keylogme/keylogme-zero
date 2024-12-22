package keylog

type Config struct {
	Devices   []DeviceInput `json:"devices"`
	Shortcuts []Shortcut    `json:"shortcuts"`
}
