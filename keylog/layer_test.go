package keylog

import (
	"fmt"
	"testing"
	"time"
)

func getTestLayers() []Layer {
	return []Layer{
		{
			LayerId: 1,
			Codes:   []uint16{16, 17, 18}, // letters q, w, e
			ShiftStates: []ShortcutCodes{
				{
					Id:    "1_2",
					Codes: []uint16{42, 16}, // Q
					Type:  HoldShortcutType,
				},
			},
		},
		{
			LayerId: 2,
			Codes:   []uint16{2, 3, 4}, // numbers 1,2,3
			ShiftStates: []ShortcutCodes{
				{
					Id:    "4_5",
					Codes: []uint16{42, 2}, // !
					Type:  HoldShortcutType,
				},
			},
		},
	}
}

func TestChangeLayerSingleCodes(t *testing.T) {
	lsd := NewLayerDetector(getTestLayers(), 100*time.Millisecond)
	deviceId := "1"
	// first layer - press "q" and  "w"
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 3, KeyPress))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 3, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithShiftedCodesInMultipleLayers(t *testing.T) {
	lsd := NewLayerDetector(getTestLayers(), 100*time.Millisecond)
	deviceId := "1"
	// first layer - press "Q"
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	fmt.Println(lsd.GetCurrentLayerId())
	if lsd.GetCurrentLayerId() != 1 { // after first key , layer id is set
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change layer- use shifted code in second layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 2, KeyPress))
	if !ld.IsDetected() {
		t.Fatal("Detection expected ")
	}
	if lsd.GetCurrentLayerId() != 2 { //  layer id is changed
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 2, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
}
