package keylog

import (
	"fmt"
	"testing"
)

func TestSeqShortcut_Detect(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 31, 30}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	scDetected := ds.Detect("1", 30) // rand key
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 48) // rand key
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 36) // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 31) // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 31) // second key makes it shortcut 2 not possible, then
	// shortcut 1 is expected
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	scDetected = ds.Detect("1", 30) // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestSeqShortcut_diffDevice_after_shortcut(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 31, 30}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	scDetected := ds.Detect("1", 36) // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 31) // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("2", 31) // second key makes it shortcut 2 not possible, then
	// shortcut 1 is expected
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != "1" {
		t.Fatal("Device expected")
	}
	scDetected = ds.Detect("1", 30) // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestSeqShortcut_diffDevice_after_shortcut_2(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 31, 30}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	scDetected := ds.Detect("1", 36) // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("1", 31) // second key shortcut, shortcut 1 not detected yet
	// because shortcut 2 is still possible
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect(
		"2",
		30,
	) // shortcut id 2 is registered but this key is from other device
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != "1" {
		t.Fatal("Device expected")
	}
	scDetected = ds.Detect("1", 30) // third  key but no detection
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestSeqShortcut_diffDevice_after_shortcut_3(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 31, 30}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	scDetected := ds.Detect("1", 36) // first key shortcut
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("2", 36) // change of keyboard
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	scDetected = ds.Detect("2", 31)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	fmt.Println(ds.currPossibleShortcuts)
	fmt.Printf("%#v\n", ds.prevShortcutDeviceDetected)
	scDetected = ds.Detect("2", 46) // third  key confirms detection shortcut
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	if scDetected.DeviceId != "2" {
		t.Fatal("Device expected")
	}
}

func TestSeqShortcut_Detect_Not_Expected(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 48, 35}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 31, 34}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	detected := ds.Detect("1", 36)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	// until here, 2 current possible shortcuts
	detected = ds.Detect("1", 48)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	// shortcut 2 should be dropped, only shortcut 1 is possible
	if len(ds.currPossibleShortcuts) != 1 {
		t.Fatal("Current possible shortcuts not expected")
	}
	// the third letter is G (shortcut 2) , but only shortcut 1 is possible
	detected = ds.Detect("1", 34)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}

func TestSeqShortcut_Detect_Multiple_Possible(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 30}, Type: SequentialShortcutType},
		{Id: 3, Codes: []uint16{36, 31, 34}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	detected := ds.Detect("1", 30)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 48)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 36)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 31)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 34)
	if detected.ShortcutId != 3 {
		t.Fatal("Detection expected")
	}
}

func TestSeqShortcut_Detect_ReAttempt(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{36, 31}, Type: SequentialShortcutType},
		{Id: 2, Codes: []uint16{36, 30}, Type: SequentialShortcutType},
	}
	ds := newSeqShortcutDetector(sl)
	detected := ds.Detect("1", 36)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 48)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 36)
	if detected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	detected = ds.Detect("1", 30)
	if detected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}
