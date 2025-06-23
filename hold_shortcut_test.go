package k0

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/types"
)

var (
	ALL_CODES          = types.GetAllCodes()
	CTRL_CODES         = types.GetCtrlCodes()
	ALT_CODES          = types.GetAltCodes()
	SHIFT_CODES        = types.GetShiftCodes()
	ALL_MODIFIER_CODES = types.GetAllModifierCodes()
)

func getFakeEvent(deviceId string, code uint16, keyevent keylogger.KeyEvent) DeviceEvent {
	return DeviceEvent{
		InputEvent: keylogger.InputEvent{
			Time: time.Now(),
			Code: code,
			Type: keyevent,
		},
		DeviceId: deviceId,
	}
}

func TestHoldShortcut_Detect(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{CTRL_CODES[0], 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
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
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{CTRL_CODES[0], 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
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
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{CTRL_CODES[0], 47}, Type: HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
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
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType}, // copy
		{Id: "2", Codes: []uint16{CTRL_CODES[0], 47}, Type: HoldShortcutType}, // paste
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", 28, keylogger.KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
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
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
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
		{Id: "1", Codes: []uint16{CTRL_CODES[0], ALT_CODES[0], 111}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", ALT_CODES[0], keylogger.KeyPress) // second key shortcut (hold)
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
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 50, keylogger.KeyRelease) // rand key
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyRelease) // aborted (ctrl key released)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent(
		"1",
		46,
		keylogger.KeyRelease,
	) // second key pressed but ctrl was released
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
}

func TestHoldShortcut_ModPressEmpty(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: "1", Codes: []uint16{CTRL_CODES[0], 46}, Type: HoldShortcutType},
		{Id: "2", Codes: []uint16{CTRL_CODES[0], 47}, Type: HoldShortcutType},
	}
	ds := NewHoldShortcutDetector(sl, ALL_MODIFIER_CODES)

	ev := getFakeEvent("1", CTRL_CODES[0], keylogger.KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, keylogger.KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "2" {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", CTRL_CODES[0], keylogger.KeyRelease)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != "" {
		t.Fatal("Detection not expected")
	}
	if len(ds.modPress) != 0 {
		t.Fatal("modPress should be empty")
	}
}
