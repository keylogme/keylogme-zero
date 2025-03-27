package keylog

import (
	"testing"
	"time"
)

func TestShiftStateDetectorMCU(t *testing.T) {
	config := getTestShiftStateConfig()
	ssd := NewShiftStateDetector(config)

	devId := "1"
	shiftKey := uint16(42)

	ev := getFakeEvent(devId, shiftKey, KeyPress)
	ssDect := ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	// press second key
	ev = getFakeEvent(devId, 2, KeyPress)
	ssDect = ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	ev = getFakeEvent(devId, shiftKey, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
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

	ev := getFakeEvent(devId, shiftKey, KeyPress)
	ssDect := ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}

	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	// press second key
	ev = getFakeEvent(devId, 2, KeyPress)
	ssDect = ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
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

	ev := getFakeEvent(devId, shiftKey, KeyPress)
	ssDect := ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}

	time.Sleep(config.ThresholdAuto.Duration / 2)
	// press second key
	ev = getFakeEvent(devId, 2, KeyPress)
	ssDect = ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// release second key
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	ev = getFakeEvent(devId, shiftKey, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
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
	ev := getFakeEvent(devId, shiftKey, KeyPress)
	ssd.handleKeyEvent(ev)
	if ssd.blockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	time.Sleep(config.ThresholdAuto.Duration / 2)
	ev = getFakeEvent(devId, 2, KeyPress)
	ssd.handleKeyEvent(ev)
	if ssd.blockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	// keyrelease
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssd.handleKeyEvent(ev)
	if !ssd.blockSaveKeylog() {
		t.Fatal("Block expected")
	}
	time.Sleep(config.ThresholdAuto.Duration + 1*time.Millisecond)
	ev = getFakeEvent(devId, shiftKey, KeyRelease)
	sd := ssd.handleKeyEvent(ev)
	if ssd.blockSaveKeylog() {
		t.Fatal("Block not expected")
	}
	if !sd.IsDetected() {
		t.Fatal("Detection expected")
	}
	if sd.Auto {
		t.Fatal("Auto not expected")
	}
	// press shift again
	ev = getFakeEvent(devId, shiftKey, KeyPress)
	ssd.handleKeyEvent(ev)
	if ssd.blockSaveKeylog() {
		t.Fatal("Block not expected")
	}
}
