package internal

type ShortcutType string

const (
	SequentialShortcutType = "seq"
	TogetherShortcutType   = "tog"
)

type Shortcut struct {
	ID     int64        `json:"id"`
	Values []string     `json:"values"`
	Type   ShortcutType `json:"type"`
}

type ShortcutsDetector struct {
	indexVal              int
	currPossibleShortcuts []*Shortcut
	Shortcuts             []Shortcut
}

func NewShortcutsDetector(s []Shortcut) *ShortcutsDetector {
	return &ShortcutsDetector{
		Shortcuts:             s,
		indexVal:              0,
		currPossibleShortcuts: []*Shortcut{},
	}
}

func (sd *ShortcutsDetector) Detect(kp string) int64 {
	if len(sd.currPossibleShortcuts) == 0 {
		for _, s := range sd.Shortcuts {
			if s.Values[0] == kp {
				sd.currPossibleShortcuts = append(sd.currPossibleShortcuts, &s)
				sd.indexVal = 1
			}
		}
	} else {
		new_ps := []*Shortcut{}
		for _, ps := range sd.currPossibleShortcuts {
			if (*ps).Values[sd.indexVal] == kp {
				new_ps = append(new_ps, ps)
			}
			if (*ps).Values[sd.indexVal] == kp && len((*ps).Values) == sd.indexVal+1 {
				sd.reset()
				return ps.ID
			}
		}
		if len(new_ps) == 0 {
			sd.reset()
		} else {
			sd.indexVal += 1
		}
	}
	return 0
}

func (sd *ShortcutsDetector) reset() {
	sd.indexVal = 0
	sd.currPossibleShortcuts = []*Shortcut{}
}
