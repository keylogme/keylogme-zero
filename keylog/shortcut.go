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

type shortcutsDetector struct {
	indexVal                int
	currPossibleShortcuts   []*Shortcut
	prevShortcutIDCompleted int64
	delayMS                 int64
	lastKeyTimestamp        time.Time
	Shortcuts               []Shortcut
}

func newShortcutsDetector(s []Shortcut) *shortcutsDetector {
	return &shortcutsDetector{
		Shortcuts:             s,
		indexVal:              0,
		currPossibleShortcuts: []*Shortcut{},
	}
}

func (sd *shortcutsDetector) Detect(kp string) int64 {
	if len(sd.currPossibleShortcuts) == 0 {
		for _, s := range sd.Shortcuts {
			if s.Values[0] == kp {
				sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, &s)
				sd.indexVal = 1
			}
		}
	} else {
		new_ps := []*Shortcut{}
		foundOnePossibleShortcutCompleted := int64(0)
		for _, ps := range sd.currPossibleShortcuts {
			if len((*ps).Values) <= sd.indexVal {
				continue
			}
			nextKeyShortcut := (*ps).Values[sd.indexVal]
			if nextKeyShortcut == kp {
				new_ps = append(new_ps, ps)
			}
			isLastKeyShortcut := len((*ps).Values) == sd.indexVal+1
			if nextKeyShortcut == kp && isLastKeyShortcut {
				foundOnePossibleShortcutCompleted = ps.Id
			}
		}
		if len(new_ps) == 1 && foundOnePossibleShortcutCompleted != 0 {
			// found only one possible shortcut
			sd.reset()
			return foundOnePossibleShortcutCompleted

		}
		if len(new_ps) == 0 {
			if sd.prevShortcutIDCompleted != 0 {
				prevShortcutId := sd.prevShortcutIDCompleted
				sd.reset()
				return prevShortcutId
			}
			sd.reset()
		} else {
			if foundOnePossibleShortcutCompleted != 0 {
				sd.prevShortcutIDCompleted = foundOnePossibleShortcutCompleted
			}
			sd.indexVal += 1
			sd.currPossibleShortcuts = new_ps
		}
	}
	return 0
}

func (sd *shortcutsDetector) reset() {
	sd.indexVal = 0
	sd.prevShortcutIDCompleted = 0
	sd.currPossibleShortcuts = []*Shortcut{}
}
