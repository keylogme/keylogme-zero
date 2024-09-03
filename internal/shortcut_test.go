package internal

import (
	"testing"
)

func TestShortcutsDetector_Detect(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{ID: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{ID: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := NewShortcutsDetector(sl)
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
	scId = ds.Detect("S") // second key shortcut => ID=1
	if scId != 1 {
		t.Fatal("Detection expected")
	}
	scId = ds.Detect("S") // second key but no detection of ID=1
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
	scId = ds.Detect("A") // third  key but no detection of ID=2
	if scId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_Detect_Not_Expected(t *testing.T) {
	sl := []Shortcut{
		{ID: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
	}
	ds := NewShortcutsDetector(sl)
	detected := ds.Detect("J")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("B")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("S")
	if detected != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_Detect_Multiple_Possible(t *testing.T) {
	sl := []Shortcut{
		{ID: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{ID: 2, Values: []string{"J", "A"}, Type: SequentialShortcutType},
	}
	ds := NewShortcutsDetector(sl)
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
	detected = ds.Detect("A")
	if detected != 2 {
		t.Fatal("Detection expected")
	}
}

func TestShortcutsDetector_Detect_ReAttempt(t *testing.T) {
	sl := []Shortcut{
		{ID: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{ID: 2, Values: []string{"J", "A"}, Type: SequentialShortcutType},
	}
	ds := NewShortcutsDetector(sl)
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
