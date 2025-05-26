package shortcut

import (
	"fmt"
	"log"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/types"
)

type ShortcutsDetector struct {
	SeqDetector  SeqShortcutDetector
	HoldDetector HoldShortcutDetector
}

func MustGetNewShortcutsDetector(sgs []types.ShortcutGroupInput) *ShortcutsDetector {
	s, err := getShortcutsFromGroups(sgs)
	if err != nil {
		log.Fatalf("Error getting shortcuts from groups: %s", err.Error())
	}

	return &ShortcutsDetector{
		SeqDetector:  NewSeqShortcutDetector(s),
		HoldDetector: NewHoldShortcutDetector(s, getAllHoldModifiers()),
	}
}

// check duplicate ids
func getShortcutsFromGroups(s []types.ShortcutGroupInput) ([]types.ShortcutCodes, error) {
	scgIds := map[string]bool{}
	scIds := map[string]bool{}
	scs := []types.ShortcutCodes{}
	for _, sg := range s {
		if _, ok := scgIds[sg.Id]; ok {
			return []types.ShortcutCodes{}, fmt.Errorf("Repeated shortcut group id %s", sg.Id)
		}
		scgIds[sg.Id] = true
		for _, sc := range sg.Shortcuts {
			// check uniqueness shortcut ids
			if _, ok := scIds[sc.Id]; ok {
				return []types.ShortcutCodes{}, fmt.Errorf("Repeated shortcut id %s", sg.Id)
			}
			scIds[sc.Id] = true
			scs = append(scs, sc)
		}
	}
	return scs, nil
}

func (sd *ShortcutsDetector) HandleKeyEvent(ke keylogger.DeviceEvent) types.ShortcutDetected {
	sdect := sd.SeqDetector.handleKeyEvent(ke)
	if sdect.ShortcutId != "" {
		return sdect
	}
	return sd.HoldDetector.HandleKeyEvent(ke)
}
