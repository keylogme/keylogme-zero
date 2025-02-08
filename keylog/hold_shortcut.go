package keylog

import (
	"slices"
)

type holdShortcutDetector struct {
	shortcuts []ShortcutCodes
	modifiers []uint16
	modPress  []uint16
}

func getAllHoldModifiers() []uint16 {
	return []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt
}

// detects hold shortcuts like Ctrl+C, Ctrl+Alt+Del, Shift+C
func newHoldShortcutDetector(shortcuts []ShortcutCodes, modifiers []uint16) holdShortcutDetector {
	hsd := holdShortcutDetector{
		shortcuts: []ShortcutCodes{},
		modifiers: modifiers,
		modPress:  []uint16{},
	}
	hsd.setShortcuts(shortcuts)
	return hsd
}

func (hd *holdShortcutDetector) setShortcuts(shortcuts []ShortcutCodes) {
	newS := []ShortcutCodes{}
	for _, s := range shortcuts {
		if s.Type == HoldShortcutType && len(s.Codes) > 1 {
			// sort codes (Important for detect function)
			slices.Sort(s.Codes)
			newS = append(newS, s)
		}
	}
	hd.shortcuts = newS
}

func (hd *holdShortcutDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	if ke.Type == evKey && ke.KeyPress() && slices.Contains(hd.modifiers, ke.Code) {
		hd.modPress = append(hd.modPress, ke.Code)
	}
	if ke.Type == evKey && ke.KeyRelease() && len(hd.modPress) > 0 {
		return hd.detect(ke.DeviceId, ke.Code)
	}
	return ShortcutDetected{}
}

func (hd *holdShortcutDetector) detect(deviceId string, code uint16) ShortcutDetected {
	// cleanup old modifiers
	hd.modPress = slices.DeleteFunc(hd.modPress, func(v uint16) bool {
		if v == code {
			return true
		}
		return false
	})
	tempCodes := slices.Clone(hd.modPress)
	tempCodes = append(tempCodes, code)
	slices.Sort(tempCodes)

	// slog.Info(fmt.Sprintf("detect %s %d \n", deviceId, code))
	for _, s := range hd.shortcuts {
		isEqual := slices.Equal(tempCodes, s.Codes)
		if isEqual {
			return ShortcutDetected{
				ShortcutId: s.Id,
				DeviceId:   deviceId,
			}
		}
	}
	return ShortcutDetected{}
}

// func (hd *holdShortcutDetector) reset() {
// 	hd.modPress = []uint16{}
// }
