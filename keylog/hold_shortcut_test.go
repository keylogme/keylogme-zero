package keylog

import "testing"

func getFakeEvent(deviceId string, code uint16, keyevent keyevent) DeviceEvent {
	return DeviceEvent{
		inputEvent{
			Type:  evKey,
			Code:  code,
			Value: int32(keyevent),
		},
		deviceId,
	}
}

func TestHoldShortcut_Detect(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 47}, Type: HoldShortcutType},
		{Id: 2, Codes: []uint16{29, 48}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectMultiple_OnlyOne(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 47}, Type: HoldShortcutType},
		{Id: 2, Codes: []uint16{29, 48}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 51, KeyRelease) // key not in shortcuts
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}

// test for copy/paste ;)
func TestHoldShortcut_DetectMultiple_Both(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 47}, Type: HoldShortcutType}, // copy
		{Id: 2, Codes: []uint16{29, 48}, Type: HoldShortcutType}, // paste
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, KeyRelease) // second key shortcut (id 1)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}

// test for paste paste paste
func TestHoldShortcut_DetectConsecutive(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 47}, Type: HoldShortcutType}, // copy
		{Id: 2, Codes: []uint16{29, 48}, Type: HoldShortcutType}, // paste
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 28, KeyRelease) // rand key
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
	ev = getFakeEvent("1", 48, KeyRelease) // second key shortcut (id 2)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 2 {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_DetectThreeKeys(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 56, 111}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 56, KeyPress) // second key shortcut (hold)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 111, KeyRelease) // last key press
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 1 {
		t.Fatal("Detection expected")
	}
}

func TestHoldShortcut_Aborted(t *testing.T) {
	sl := []ShortcutCodes{
		{Id: 1, Codes: []uint16{29, 47}, Type: HoldShortcutType},
	}
	ds := newHoldShortcutDetector(sl)

	ev := getFakeEvent("1", 29, KeyPress) // first key shortcut (hold)
	scDetected := ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 50, KeyRelease) // rand key
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 29, KeyRelease) // aborted (ctrl key released)
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent("1", 47, KeyRelease) // second key pressed but ctrl was released
	scDetected = ds.handleKeyEvent(ev)
	if scDetected.ShortcutId != 0 {
		t.Fatal("Detection not expected")
	}
}
