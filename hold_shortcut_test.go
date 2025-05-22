package keylog

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
)

func getFakeEvent(deviceId string, code uint16, keyevent keylogger.KeyEvent) keylogger.DeviceEvent {
	return keylogger.DeviceEvent{
		InputEvent: keylogger.InputEvent{
			Code:  code,
			Value: keyevent,
		},
		DeviceId: deviceId,
		ExecTime: time.Now(),
	}
}

func TestHoldShortcut_Detect(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{29, 47}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
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
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 51, keylogger.KeyRelease) // key not in shortcuts
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
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
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 46, keylogger.KeyRelease) // second key shortcut (id 1)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
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
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	// Now release and press again ctrl
	ev = getFakeEvent("1", 29, keylogger.KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectThreeKeys(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 56, 111}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 56, keylogger.KeyPress) // second key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 111, keylogger.KeyRelease) // last key press
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "1" {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_Aborted(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{29, 46}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 50, keylogger.KeyRelease) // rand key
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyRelease) // aborted (ctrl key released)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 46, keylogger.KeyRelease) // second key pressed but ctrl was released
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
	ds := newHoldShortcutDetector(sl, getAllHoldModifiers())

	ev := getFakeEvent("1", 29, keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 29, keylogger.KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	if len(ds.modPress) != 0 {
		t.Fatal("modPress should be empty")
	}
}
