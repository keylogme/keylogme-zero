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

func (sd *shortcutsDetector) Detect(deviceId int64, kp string) ShortcutDetected {
	if len(sd.currPossibleShortcuts) == 0 {
		for _, s := range sd.Shortcuts {
			if s.Values[0] == kp {
				scd := shortcutDevice{Shortcut: s, DeviceId: deviceId}
				sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, scd)
				sd.indexVal = 1
			}
		}
	} else {
		new_ps := []shortcutDevice{}
		foundOnePossibleShortcutCompleted := new(ShortcutDetected)
		for _, ps := range sd.currPossibleShortcuts {
			if len(ps.Values) <= sd.indexVal {
				continue
			}
			nextKeyShortcut := ps.Values[sd.indexVal]
			if nextKeyShortcut == kp && ps.DeviceId == deviceId {
				new_ps = append(new_ps, ps)
			}
			isLastKeyShortcut := len((ps).Values) == sd.indexVal+1
			if nextKeyShortcut == kp && isLastKeyShortcut && ps.DeviceId == deviceId {
				foundOnePossibleShortcutCompleted.DeviceId = ps.DeviceId
				foundOnePossibleShortcutCompleted.ShortcutId = ps.Id
			}
		}
		if len(new_ps) == 1 && foundOnePossibleShortcutCompleted.ShortcutId != 0 {
			// found only one possible shortcut
			sd.reset()
			return *foundOnePossibleShortcutCompleted

		}
		if len(new_ps) == 0 {
			if sd.prevShortcutDeviceDetected.ShortcutId != 0 {
				output := sd.prevShortcutDeviceDetected
				sd.reset()
				return output
			}
			sd.reset()
		} else {
			if foundOnePossibleShortcutCompleted.ShortcutId != 0 {
				sd.prevShortcutDeviceDetected = *foundOnePossibleShortcutCompleted
			}
			sd.indexVal += 1
			sd.currPossibleShortcuts = new_ps
		}
	}
	return *new(ShortcutDetected)
}

func (sd *shortcutsDetector) reset() {
	sd.indexVal = 0
	sd.prevShortcutDeviceDetected = *new(ShortcutDetected)
	sd.currPossibleShortcuts = []shortcutDevice{}
}
