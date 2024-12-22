package keylog

import (
	"time"
)

type seqShortcutDetector struct {
	shortcuts                  []ShortcutCodes
	indexVal                   int
	currPossibleShortcuts      []shortcutDevice
	prevShortcutDeviceDetected ShortcutDetected
	delayMS                    int64
	lastKeyTimestamp           time.Time
}

func newSeqShortcutDetector(shortcuts []ShortcutCodes) seqShortcutDetector {
	shortsSeq := []ShortcutCodes{}
	for _, s := range shortcuts {
		if s.Type == SequentialShortcutType {
			shortsSeq = append(shortsSeq, s)
		}
	}
	return seqShortcutDetector{
		shortcuts:             shortsSeq,
		indexVal:              0,
		currPossibleShortcuts: []shortcutDevice{},
	}
}

func (sd *seqShortcutDetector) handleKeyEvent(ke DeviceEvent) ShortcutDetected {
	if ke.Type == evKey && ke.KeyRelease() {
		return sd.Detect(ke.DeviceId, ke.Code)
	}
	return *new(ShortcutDetected)
}

func (sd *seqShortcutDetector) Detect(deviceId string, kp uint16) ShortcutDetected {
	if sdet := sd.handleChangeOfDevice(deviceId, kp); sdet.ShortcutId != 0 {
		return sdet
	}
	if len(sd.currPossibleShortcuts) == 0 {
		sd.handleFirstKey(deviceId, kp)
	} else {
		newPossibleShortcuts, shortcutCompleted := sd.checkPossibleShortcuts(deviceId, kp)
		if len(newPossibleShortcuts) == 1 && shortcutCompleted.ShortcutId != 0 {
			// found only one possible shortcut
			sd.reset()
			return shortcutCompleted

		}
		if len(newPossibleShortcuts) == 0 {
			if sd.prevShortcutDeviceDetected.ShortcutId != 0 {
				output := sd.prevShortcutDeviceDetected
				sd.reset()
				return output
			}
			sd.reset()
		} else {
			if shortcutCompleted.ShortcutId != 0 {
				sd.prevShortcutDeviceDetected = shortcutCompleted
			}
			sd.indexVal += 1
			sd.currPossibleShortcuts = newPossibleShortcuts
		}
	}
	return *new(ShortcutDetected)
}

func (sd *seqShortcutDetector) handleFirstKey(deviceId string, kp uint16) {
	for _, s := range sd.shortcuts {
		if s.Codes[0] == kp {
			scd := shortcutDevice{ShortcutCodes: s, DeviceId: deviceId}
			sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, scd)
			sd.indexVal = 1
		}
	}
}

func (sd *seqShortcutDetector) handleChangeOfDevice(deviceId string, kp uint16) ShortcutDetected {
	if len(sd.currPossibleShortcuts) > 0 &&
		sd.currPossibleShortcuts[0].DeviceId != deviceId {
		if sd.prevShortcutDeviceDetected.ShortcutId != 0 {
			output := sd.prevShortcutDeviceDetected
			sd.reset()
			sd.handleFirstKey(deviceId, kp)
			return output
		}
		sd.reset()
	}
	return *new(ShortcutDetected)
}

func (sd *seqShortcutDetector) checkPossibleShortcuts(
	deviceId string,
	kp uint16,
) ([]shortcutDevice, ShortcutDetected) {
	new_ps := []shortcutDevice{}
	foundOnePossibleShortcutCompleted := new(ShortcutDetected)
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

func (sd *seqShortcutDetector) reset() {
	sd.indexVal = 0
	sd.prevShortcutDeviceDetected = *new(ShortcutDetected)
	sd.currPossibleShortcuts = []shortcutDevice{}
}
