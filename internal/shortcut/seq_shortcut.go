package shortcut

import (
	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/types"
)

type shortcutDevice struct {
	types.ShortcutCodes
	DeviceId string
}

type SeqShortcutDetector struct {
	shortcuts                  []types.ShortcutCodes
	indexVal                   int
	currPossibleShortcuts      []shortcutDevice
	prevShortcutDeviceDetected types.ShortcutDetected
	// delayMS                    int64
	// lastKeyTimestamp           time.Time
}

func NewSeqShortcutDetector(shortcuts []types.ShortcutCodes) SeqShortcutDetector {
	shortsSeq := []types.ShortcutCodes{}
	for _, s := range shortcuts {
		if s.Type == types.SequentialShortcutType && len(s.Codes) > 1 {
			shortsSeq = append(shortsSeq, s)
		}
	}
	return SeqShortcutDetector{
		shortcuts:             shortsSeq,
		indexVal:              0,
		currPossibleShortcuts: []shortcutDevice{},
	}
}

func (sd *SeqShortcutDetector) handleKeyEvent(ke keylogger.DeviceEvent) types.ShortcutDetected {
	if ke.KeyRelease() {
		return sd.Detect(ke.DeviceId, ke.Code)
	}
	return *new(types.ShortcutDetected)
}

func (sd *SeqShortcutDetector) Detect(deviceId string, kp uint16) types.ShortcutDetected {
	if sdet := sd.handleChangeOfDevice(deviceId, kp); sdet.ShortcutId != "" {
		return sdet
	}
	if len(sd.currPossibleShortcuts) == 0 {
		sd.handleFirstKey(deviceId, kp)
	} else {
		newPossibleShortcuts, shortcutCompleted := sd.checkPossibleShortcuts(deviceId, kp)
		if len(newPossibleShortcuts) == 1 && shortcutCompleted.ShortcutId != "" {
			// found only one possible shortcut
			sd.reset()
			return shortcutCompleted

		}
		if len(newPossibleShortcuts) == 0 {
			if sd.prevShortcutDeviceDetected.ShortcutId != "" {
				output := sd.prevShortcutDeviceDetected
				sd.reset()
				return output
			}
			sd.reset()
		} else {
			if shortcutCompleted.ShortcutId != "" {
				sd.prevShortcutDeviceDetected = shortcutCompleted
			}
			sd.indexVal += 1
			sd.currPossibleShortcuts = newPossibleShortcuts
		}
	}
	return *new(types.ShortcutDetected)
}

func (sd *SeqShortcutDetector) handleFirstKey(deviceId string, kp uint16) {
	for _, s := range sd.shortcuts {
		if s.Codes[0] == kp {
			scd := shortcutDevice{ShortcutCodes: s, DeviceId: deviceId}
			sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, scd)
			sd.indexVal = 1
		}
	}
}

func (sd *SeqShortcutDetector) handleChangeOfDevice(
	deviceId string,
	kp uint16,
) types.ShortcutDetected {
	if len(sd.currPossibleShortcuts) > 0 &&
		sd.currPossibleShortcuts[0].DeviceId != deviceId {
		if sd.prevShortcutDeviceDetected.ShortcutId != "" {
			output := sd.prevShortcutDeviceDetected
			sd.reset()
			sd.handleFirstKey(deviceId, kp)
			return output
		}
		sd.reset()
	}
	return *new(types.ShortcutDetected)
}

func (sd *SeqShortcutDetector) checkPossibleShortcuts(
	deviceId string,
	kp uint16,
) ([]shortcutDevice, types.ShortcutDetected) {
	new_ps := []shortcutDevice{}
	foundOnePossibleShortcutCompleted := new(types.ShortcutDetected)
	for _, ps := range sd.currPossibleShortcuts {
		if len(ps.Codes) <= sd.indexVal {
			continue
		}
		nextKeyShortcut := ps.Codes[sd.indexVal]
		if nextKeyShortcut == kp && ps.DeviceId == deviceId {
			// if nextKeyShortcut == kp {
			new_ps = append(new_ps, ps)
		}
		isLastKeyShortcut := len((ps).Codes) == sd.indexVal+1
		if nextKeyShortcut == kp && isLastKeyShortcut && ps.DeviceId == deviceId {
			// if nextKeyShortcut == kp && isLastKeyShortcut {
			foundOnePossibleShortcutCompleted.DeviceId = ps.DeviceId
			foundOnePossibleShortcutCompleted.ShortcutId = ps.Id
		}
	}
	return new_ps, *foundOnePossibleShortcutCompleted
}

func (sd *SeqShortcutDetector) reset() {
	sd.indexVal = 0
	sd.prevShortcutDeviceDetected = *new(types.ShortcutDetected)
	sd.currPossibleShortcuts = []shortcutDevice{}
}
