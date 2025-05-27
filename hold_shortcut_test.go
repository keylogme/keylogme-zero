package k0

import (
	"testing"
)

func TestHoldShortcut_Detect(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectMultiple_OnlyOne(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 51, KeyRelease) // key not in shortcuts
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

// test for copy/paste ;)
func TestHoldShortcut_DetectMultiple_Both(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 46, KeyRelease) // second key shortcut (id 1)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

// test for paste paste paste
func TestHoldShortcut_DetectConsecutive(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	// Now release and press again ctrl
	ev = GetFakeEvent("1", 29, KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectThreeKeys(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 56, 111}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 56, KeyPress) // second key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 111, KeyRelease) // last key press
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_Aborted(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 50, KeyRelease) // rand key
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 29, KeyRelease) // aborted (ctrl key released)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent(
		"1",
		46,
		KeyRelease,
	) // second key pressed but ctrl was released
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
}

func TestHoldShortcut_ModPressEmpty(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := GetFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = GetFakeEvent("1", 47, KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = GetFakeEvent("1", 29, KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	if len(ds.modPress) != 0 {
		t.Fatal("modPress should be empty")
	}
}
