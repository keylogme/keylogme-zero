package keylog

import (
	"slices"
)

type holdShortcutDetector struct {
	shortcuts []ShortcutCodes
	modifiers []uint16
	modPress  []uint16
}

func getCtrlKeys() []uint16 {
	return []uint16{29, 97}
}

func getShiftKeys() []uint16 {
	return []uint16{42, 54}
}

func getAltKeys() []uint16 {
	return []uint16{56, 100}
}

func getAllHoldModifiers() []uint16 {
	return slices.Concat(getCtrlKeys(), getShiftKeys(), getAltKeys())
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

func (hd *holdShortcutDetector) isHolded() bool {
	return len(hd.modPress) > 0
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
	if ke.Type == evKey && ke.KeyRelease() && hd.isHolded() {
		return hd.detect(ke.DeviceId, ke.Code)
	}
	if ke.Type == evKey && ke.KeyPress() &&
		slices.Contains(hd.modifiers, ke.Code) &&
		!slices.Contains(hd.modPress, ke.Code) {
		hd.modPress = append(hd.modPress, ke.Code)
	}
	return ShortcutDetected{}
}

func (hd *holdShortcutDetector) detect(deviceId string, code uint16) ShortcutDetected {
	// cleanup old modifiers
	hd.modPress = slices.DeleteFunc(hd.modPress, func(v uint16) bool {
		return v == code
	})
	if len(hd.modPress) == 0 {
		return ShortcutDetected{}
	}
	tempCodes := slices.Clone(hd.modPress)
	tempCodes = append(tempCodes, code)
	slices.Sort(tempCodes)

	// slog.Info(fmt.Sprintf("detect %s %d \n", deviceId, code))
	for _, s := range hd.shortcuts {
		// slog.Info(fmt.Sprintf("shortcut %v  vs  %v\n", s, tempCodes))
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
