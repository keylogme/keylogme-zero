package keylog

import (
	"testing"
	"time"
)

func TestShiftStateDetectorMCU(t *testing.T) {
	ssd := NewShiftStateDetector(getTestShiftStateConfig())

	devId := "1"
	shiftKey := uint16(42)

	ev := getFakeEvent(devId, shiftKey, KeyPress)
	ssDect := ssd.handleKeyEvent(ev)
	if ssDect.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
	if !ssDect.IsDetected() {
		t.Fatal("Detection expected")
	}
	if !ssDect.Auto {
		t.Fatal("Auto  expected")
	}
}

func TestShiftStateDetectorHuman(t *testing.T) {
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
	ev = getFakeEvent(devId, 2, KeyRelease)
	ssDect = ssd.handleKeyEvent(ev)
	if !ssDect.IsDetected() {
		t.Fatal("Detection expected")
	}
	if ssDect.Auto {
		t.Fatal("Not Auto expected")
	}
}
