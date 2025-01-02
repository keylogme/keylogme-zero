package keylog

import (
	"fmt"
	"log"
)

type ShortcutType string

const (
	SequentialShortcutType = "seq"
	HoldShortcutType       = "hold"
)

type ShortcutGroup struct {
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

func MustGetNewShortcutsDetector(sgs []ShortcutGroup) *shortcutsDetector {
	s, err := getShortcutsFromGroups(sgs)
	if err != nil {
		log.Fatalf("Error getting shortcuts from groups: %s", err.Error())
	}

	return &shortcutsDetector{
		SeqDetector:  newSeqShortcutDetector(s),
		HoldDetector: newHoldShortcutDetector(s),
	}
}

func getShortcutsFromGroups(s []ShortcutGroup) ([]ShortcutCodes, error) {
	var scIds map[string]bool
	var scs []ShortcutCodes
	for _, sg := range s {
		if _, ok := scIds[sg.Id]; ok {
			return []ShortcutCodes{}, fmt.Errorf("Repeated shortcut id %s", sg.Id)
		}
		scIds[sg.Id] = true
		scs = append(scs, sg.Shortcuts...)
	}
	return scs, nil
}

func (sd *shortcutsDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	sdect := sd.SeqDetector.handleKeyEvent(ke)
	if sdect.ShortcutId != "" {
		return sdect
	}
	return sd.HoldDetector.handleKeyEvent(ke)
}
