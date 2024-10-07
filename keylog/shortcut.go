package keylog

import "time"

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

type shortcutDevice struct {
	Shortcut
	DeviceId int64
}

type ShortcutDetected struct {
	ShortcutId int64
	DeviceId   int64
}

type shortcutsDetector struct {
	indexVal                   int
	currPossibleShortcuts      []shortcutDevice
	prevShortcutDeviceDetected ShortcutDetected
	delayMS                    int64
	lastKeyTimestamp           time.Time
	Shortcuts                  []Shortcut
}

func newShortcutsDetector(s []Shortcut) *shortcutsDetector {
	return &shortcutsDetector{
		Shortcuts:             s,
		indexVal:              0,
		currPossibleShortcuts: []shortcutDevice{},
	}
}

func (sd *shortcutsDetector) handleFirstKey(deviceId int64, kp string) {
	for _, s := range sd.Shortcuts {
		if s.Values[0] == kp {
			scd := shortcutDevice{Shortcut: s, DeviceId: deviceId}
			sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, scd)
			sd.indexVal = 1
		}
	}
}

func (sd *shortcutsDetector) handleChangeOfDevice(deviceId int64, kp string) ShortcutDetected {
	if len(sd.currPossibleShortcuts) > 0 && sd.currPossibleShortcuts[0].DeviceId != deviceId {
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

func (sd *shortcutsDetector) checkPossibleShortcuts(
	deviceId int64,
	kp string,
) ([]shortcutDevice, ShortcutDetected) {
	new_ps := []shortcutDevice{}
	foundOnePossibleShortcutCompleted := new(ShortcutDetected)
	for _, ps := range sd.currPossibleShortcuts {
		if len(ps.Values) <= sd.indexVal {
			continue
		}
		nextKeyShortcut := ps.Values[sd.indexVal]
		if nextKeyShortcut == kp && ps.DeviceId == deviceId {
			// if nextKeyShortcut == kp {
			new_ps = append(new_ps, ps)
		}
		isLastKeyShortcut := len((ps).Values) == sd.indexVal+1
		if nextKeyShortcut == kp && isLastKeyShortcut && ps.DeviceId == deviceId {
			// if nextKeyShortcut == kp && isLastKeyShortcut {
			foundOnePossibleShortcutCompleted.DeviceId = ps.DeviceId
			foundOnePossibleShortcutCompleted.ShortcutId = ps.Id
		}
	}
	return new_ps, *foundOnePossibleShortcutCompleted
}

func (sd *shortcutsDetector) Detect(deviceId int64, kp string) ShortcutDetected {
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

func (sd *shortcutsDetector) reset() {
	sd.indexVal = 0
	sd.prevShortcutDeviceDetected = *new(ShortcutDetected)
	sd.currPossibleShortcuts = []shortcutDevice{}
}
