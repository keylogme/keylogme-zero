package types

type ShortcutType string

const (
	SequentialShortcutType = "seq"
	HoldShortcutType       = "hold"
)

type ShortcutGroupInput struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Shortcuts []ShortcutCodes `json:"shortcuts"`
}

type ShortcutCodes struct {
	Id    string       `json:"id"`
	Name  string       `json:"name"`
	Codes []uint16     `json:"codes"`
	Type  ShortcutType `json:"type"`
}

type ShortcutDetected struct {
	ShortcutId string
	DeviceId   string
}

func (sd ShortcutDetected) IsDetected() bool {
	return sd.ShortcutId != ""
}
