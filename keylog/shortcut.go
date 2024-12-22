package keylog

type ShortcutType string

const (
	SequentialShortcutType = "seq"
	TogetherShortcutType   = "tog"
)

type Shortcut struct {
	Id     int64        `json:"id"`
	Values []string     `json:"values"`
	Type   ShortcutType `json:"type"`
}

type ShortcutCodes struct {
	Id    int64        `json:"id"`
	Codes []uint16     `json:"codes"`
	Type  ShortcutType `json:"type"`
}

type shortcutDevice struct {
	Shortcut
	DeviceId string
}

type ShortcutDetected struct {
	ShortcutId int64
	DeviceId   string
}

type shortcutsDetector struct {
	SeqDetector  seqShortcutDetector
	HoldDetector holdShortcutDetector
}

func NewShortcutsDetector(s []Shortcut) *shortcutsDetector {
	return &shortcutsDetector{
		SeqDetector:  newSeqShortcutDetector(s),
		HoldDetector: newHoldShortcutDetector([]ShortcutCodes{}),
	}
}

func (sd *shortcutsDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	sdect := sd.SeqDetector.handleKeyEvent(ke)
	if sdect.ShortcutId != 0 {
		return sdect
	}
	return sd.HoldDetector.handleKeyEvent(ke)
}
