package shortcut

import (
	"testing"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/types"
)

func TestHoldShortcut_Detect(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: types.HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectMultiple_OnlyOne(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: types.HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 51, keylogger.KeyRelease) // key not in shortcuts
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

// test for copy/paste ;)
func TestHoldShortcut_DetectMultiple_Both(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{29, 47}, Type: types.HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 46, keylogger.KeyRelease) // second key shortcut (id 1)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

// test for paste paste paste
func TestHoldShortcut_DetectConsecutive(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{29, 47}, Type: types.HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	// Now release and press again ctrl
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyRelease)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectThreeKeys(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 56, 111}, Type: types.HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 56, keylogger.KeyPress) // second key shortcut (hold)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 111, keylogger.KeyRelease) // last key press
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_Aborted(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 50, keylogger.KeyRelease) // rand key
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyRelease) // aborted (ctrl key released)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent(
		"1",
		46,
		keylogger.KeyRelease,
	) // second key pressed but ctrl was released
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
}

func TestHoldShortcut_ModPressEmpty(t *testing.T) {
	sl := []types.ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: types.HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: types.HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := keylogger.GetFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = keylogger.GetFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = keylogger.GetFakeEvent("1", 29, keylogger.KeyRelease)
	scDetected = ds.HandleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	if len(ds.modPress) != 0 {
		t.Fatal("modPress should be empty")
	}
}
