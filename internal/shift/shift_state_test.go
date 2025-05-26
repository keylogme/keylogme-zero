package shift

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
)

func TestShiftStateDetectorMCU(t *testing.T) {
	config := getTestShiftStateConfig()
	ssd := NewShiftStateDetector(config)

	devId := "1"
	shiftKey := uint16(42)

	ev := keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyPress)
	ssDect := ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	// press second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyPress)
	ssDect = ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyRelease)
	ssDect = ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	ev = keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyRelease)
	ssDect = ssd.HandleKeyEvent(ev)
	if !ssDect.IsDetected() {
		t.Fatal("Detection expected")
	}
	if !ssDect.Auto {
		t.Fatal("Auto  expected")
	}
}

func TestShiftStateDetectorHuman_PressSlow(t *testing.T) {
	config := getTestShiftStateConfig()
	ssd := NewShiftStateDetector(config)

	devId := "1"
	shiftKey := uint16(42)

	ev := keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyPress)
	ssDect := ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}

	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	// press second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyPress)
	ssDect = ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyRelease)
	ssDect = ssd.HandleKeyEvent(ev)
	if !ssDect.IsDetected() {
		t.Fatal("Detection expected")
	}

	if ssDect.Auto {
		t.Fatal("Not Auto expected")
	}
}

func TestShiftStateDetectorHuman_ReleaseSlow(t *testing.T) {
	config := getTestShiftStateConfig()
	ssd := NewShiftStateDetector(config)

	devId := "1"
	shiftKey := uint16(42)

	ev := keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyPress)
	ssDect := ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}

	time.Sleep(config.ThresholdAuto.Duration / 2)
	// press second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyPress)
	ssDect = ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyRelease)
	ssDect = ssd.HandleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	ev = keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyRelease)
	ssDect = ssd.HandleKeyEvent(ev)
	if !ssDect.IsDetected() {
		t.Fatal("Detection expected")
	}

	if ssDect.Auto {
		t.Fatal("Not Auto expected")
	}
}

func TestShiftStateBlockKeylog(t *testing.T) {
	config := getTestShiftStateConfig()
	ssd := NewShiftStateDetector(config)

	devId := "1"
	shiftKey := uint16(42)

	// keypress
	ev := keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyPress)
	ssd.HandleKeyEvent(ev)
	if ssd.BlockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyPress)
	ssd.HandleKeyEvent(ev)
	if ssd.BlockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	// keyrelease
	ev = keylogger.GetFakeEvent(devId, 2, keylogger.KeyRelease)
	ssd.HandleKeyEvent(ev)
	if !ssd.BlockSaveKeylog() {
		t.Fatal("Block expected")
	}
	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	ev = keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyRelease)
	sd := ssd.HandleKeyEvent(ev)
	if ssd.BlockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	if !sd.IsDetected() {
		t.Fatal("Detection expected")
	}
	if sd.Auto {
		t.Fatal("Auto not expected")
	}
	// press shift again
	ev = keylogger.GetFakeEvent(devId, shiftKey, keylogger.KeyPress)
	ssd.HandleKeyEvent(ev)
	if ssd.BlockSaveKeylog() {
		t.Fatal("Block not expected")
	}
}
