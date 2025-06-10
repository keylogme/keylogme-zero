package k0

import (
	"fmt"
	"log"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
)

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

type shortcutsDetector struct {
	seqDetector  seqShortcutDetector
	holdDetector holdShortcutDetector
}

func MustGetNewShortcutsDetector(sgs []ShortcutGroupInput) *shortcutsDetector {
	s, err := getShortcutsFromGroups(sgs)
	if err != nil {
		log.Fatalf("Error getting shortcuts from groups: %s", err.Error())
	}

	return &shortcutsDetector{
		seqDetector:  NewSeqShortcutDetector(s),
		holdDetector: NewHoldShortcutDetector(s, keylogger.GetAllModifierCodes()),
	}
}

// check duplicate ids
func getShortcutsFromGroups(s []ShortcutGroupInput) ([]ShortcutCodes, error) {
	scgIds := map[string]bool{}
	scIds := map[string]bool{}
	scs := []ShortcutCodes{}
	for _, sg := range s {
		if _, ok := scgIds[sg.Id]; ok {
			return []ShortcutCodes{}, fmt.Errorf("repeated shortcut group id %s", sg.Id)
		}
		scgIds[sg.Id] = true
		for _, sc := range sg.Shortcuts {
			// check uniqueness shortcut ids
			if _, ok := scIds[sc.Id]; ok {
				return []ShortcutCodes{}, fmt.Errorf("repeated shortcut id %s", sg.Id)
			}
			scIds[sc.Id] = true
			scs = append(scs, sc)
		}
	}
	return scs, nil
}

func (sd *shortcutsDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	sdect := sd.seqDetector.handleKeyEvent(ke)
	if sdect.ShortcutId != "" {
		return sdect
	}
	return sd.holdDetector.handleKeyEvent(ke)
}
