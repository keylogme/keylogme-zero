package keylog

import (
	"fmt"
	"testing"
)

func TestShortcutsDetector_Detect(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	scDetected := ds.Detect(1, "A") // rand key
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "B") // rand key
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "J") // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "S") // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "S") // second key makes it shortcut 2 not possible, then
	// shortcut 1 is expected
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	scDetected = ds.Detect(1, "A") // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_diffDevice_after_shortcut(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	scDetected := ds.Detect(1, "J") // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "S") // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(2, "S") // second key makes it shortcut 2 not possible, then
	// shortcut 1 is expected
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != 1 {
		t.Fatal("Device expected")
	}
	scDetected = ds.Detect(1, "A") // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_diffDevice_after_shortcut_2(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	scDetected := ds.Detect(1, "J") // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(1, "S") // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(2, "A") // shortcut id 2 is registered but this key is from other device
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != 1 {
		t.Fatal("Device expected")
	}
	scDetected = ds.Detect(1, "A") // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestShortcutsDetector_diffDevice_after_shortcut_3(t *testing.T) {
	sl := []Shortcut{
		// {ID: 1, Values: []string{"L_CTRL", "S"}, Type: SequentialShortcutType},
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	scDetected := ds.Detect(1, "J") // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(2, "J") // change of keyboard
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(2, "S")
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	fmt.Println(ds.currPossibleShortcuts)
	fmt.Printf("%#v\n", ds.prevShortcutDeviceDetected)
	scDetected = ds.Detect(2, "C") // third  key confirms detection shortcut
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != 2 {
		t.Fatal("Device expected")
	}
}

func TestShortcutsDetector_Detect_Not_Expected(t *testing.T) {
	sl := []Shortcut{
		{Id: 1, Values: []string{"J", "B", "H"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "S", "G"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	detected := ds.Detect(1, "J")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	// until here, 2 current possible shortcuts
	detected = ds.Detect(1, "B")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	// shortcut 2 should be dropped, only shortcut 1 is possible
	if len(ds.currPossibleShortcuts) != 1 {
		t.Fatal("Current possible shortcuts not expected")
	}
	// the third letter is G (shortcut 2) , but only shortcut 1 is possible
	detected = ds.Detect(1, "G")
	if detected.ShortcutId != 0 {
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
	detected := ds.Detect(1, "A")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "B")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "J")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "S")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "G")
	if detected.ShortcutId != 3 {
		t.Fatal("Detection expected")
	}
}

func TestShortcutsDetector_Detect_ReAttempt(t *testing.T) {
	sl := []Shortcut{
		{Id: 1, Values: []string{"J", "S"}, Type: SequentialShortcutType},
		{Id: 2, Values: []string{"J", "A"}, Type: SequentialShortcutType},
	}
	ds := newShortcutsDetector(sl)
	detected := ds.Detect(1, "J")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "B")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "J")
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect(1, "A")
	if detected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}
