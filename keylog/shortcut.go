package keylog

type ShortcutType string

const (
	SequentialShortcutType = "seq"
	HoldShortcutType       = "hold"
)

type ShortcutCodes struct {
	Id    string       `json:"id"`
	Name  string       `json:"name"`
	Codes []uint16     `json:"codes"`
	Type  ShortcutType `json:"type"`
}

type shortcutDevice struct {
	ShortcutCodes
	DeviceId string
}

type ShortcutDetected struct {
	ShortcutId string
	DeviceId   string
}

type shortcutsDetector struct {
	SeqDetector  seqShortcutDetector
	HoldDetector holdShortcutDetector
}

func NewShortcutsDetector(s []ShortcutCodes) *shortcutsDetector {
	return &shortcutsDetector{
		SeqDetector:  newSeqShortcutDetector(s),
		HoldDetector: newHoldShortcutDetector(s),
	}
}

func (sd *shortcutsDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	sdect := sd.SeqDetector.handleKeyEvent(ke)
	if sdect.ShortcutId != "" {
		return sdect
	}
	return sd.HoldDetector.handleKeyEvent(ke)
}
