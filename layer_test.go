package keylog

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

func getTestLayers() []Layer {
	return []Layer{
		{
			LayerId: 1,
			Codes:   []uint16{16, 17, 18}, // letters q, w, e
			ShiftedCodes: ShiftedCodes{
				ShiftCode: 42,           // shift key
				Codes:     []uint16{16}, // Q
			},
		},
		{
			LayerId: 2,
			Codes:   []uint16{2, 3, 4}, // numbers 1,2,3
			ShiftedCodes: ShiftedCodes{
				ShiftCode: 42,          // shift key
				Codes:     []uint16{2}, // !
			},
		},
	}
}

func getTestLayersCodesEmpty() []Layer {
	return []Layer{
		{
			LayerId: 1,
			Codes:   []uint16{16, 17, 18}, // letters q, w, e
			ShiftedCodes: ShiftedCodes{
				ShiftCode: 42,           // shift key
				Codes:     []uint16{16}, // Q
			},
		},
		{
			LayerId: 2,
			Codes:   []uint16{2},
			ShiftedCodes: ShiftedCodes{
				ShiftCode: 42,          // shift key
				Codes:     []uint16{2}, // !
			},
		},
	}
}

func getTestShiftStateConfig() ShiftState {
	return ShiftState{
		ThresholdAuto: types.Duration{Duration: 100 * time.Millisecond},
	}
}

func TestChangeLayerSingleCodes(t *testing.T) {
	lsd := NewLayerDetector(getTestLayers(), getTestShiftStateConfig())
	deviceId := "1"
	// first layer - press "q" and  "w"
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
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
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 3, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 2 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithShiftedCodesInMultipleLayers(t *testing.T) {
	lsd := NewLayerDetector(getTestLayers(), getTestShiftStateConfig())
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
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 { // after Shift+Q press , layer id = 1 is set
		t.Fatal("Layer id incorrect")
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
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 2, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected ")
	}
	if lsd.GetCurrentLayerId() != 2 { //  layer id is changed
		t.Fatal("Layer id incorrect")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithShiftedCodesInMultipleLayers_CodesEmpty(t *testing.T) {
	lsd := NewLayerDetector(getTestLayersCodesEmpty(), getTestShiftStateConfig())
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
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	// first time current layer set=> does not trigger layer change detection
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
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
	if ld.IsDetected() {
		t.Fatal("Detection not expected ")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 2, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 2 { //  layer id is changed
		t.Fatal("Layer id incorrect")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}
