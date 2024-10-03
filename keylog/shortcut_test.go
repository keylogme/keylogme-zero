package keylog

import (
	"testing"
)

func TestShortcutsDetector_Detect(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	scId := ds.Detect("A") // rand key
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
	scId = ds.Detect("B") // rand key
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
	scId = ds.Detect("J") // first key shortcut
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
	scId = ds.Detect("S") // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
	scId = ds.Detect("S") // second key makes it shortcut 2 not possible, then
	// shortcut 1 is expected
	if scId != 1 {
		t.Fatal("Detection expected")
	}
	scId = ds.Detect("A") // third  key but no detection
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_Detect_Not_Expected(t *testing.T) {
	sl := []Shortcut{
		{Id: 1, Values: []string{"J", "B", "H"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "G"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	detected := ds.Detect("J")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	// until here, 2 current possible shortcuts
	detected = ds.Detect("B")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	// shortcut 2 should be dropped, only shortcut 1 is possible
	if len(ds.currPossibleShortcuts) != 1 {
		t.Fatal("Current possible shortcuts not expected")
	}
	// the third letter is G (shortcut 2) , but only shortcut 1 is possible
	detected = ds.Detect("G")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_Detect_Multiple_Possible(t *testing.T) {
	sl := []Shortcut{
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "A"}, Type: SequentialShortcutType},
		{Id: 3, Values: []string{"J", "S", "G"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	detected := ds.Detect("A")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("B")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("J")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("S")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("G")
	if detected != 3 {
		t.Fatal("Detection expected")
	}
}

func TestShortcutsDetector_Detect_ReAttempt(t *testing.T) {
	sl := []Shortcut{
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	detected := ds.Detect("J")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("B")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("J")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("A")
	if detected != 2 {
		t.Fatal("Detection expected")
	}
}
