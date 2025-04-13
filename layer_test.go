package keylog

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

func getTestLayers(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []Layer{
				{
					Id: 1,
					Codes: []LayerCode{
						{Code: 16},               // q
						{Code: 17},               // w
						{Code: 18},               // e
						{Code: 42},               // shift
						{Code: 13, Modifier: 42}, // +
					},
				},
				{
					Id: 2,
					Codes: []LayerCode{
						{Code: 2},               // 1
						{Code: 3},               // 2
						{Code: 4},               // 3
						{Code: 2, Modifier: 42}, //!
					},
				},
			},
		},
	}
}

func getTestLayersCodesEmpty(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []Layer{
				{
					Id: 1,
					Codes: []LayerCode{
						{Code: 16},               // q
						{Code: 17},               // w
						{Code: 18},               // e
						{Code: 42},               // shift
						{Code: 13, Modifier: 42}, // +
					},
				},
				{
					Id: 2,
					Codes: []LayerCode{
						{Code: 2},               // 1
						{Code: 2, Modifier: 42}, //!
					},
				},
			},
		},
	}
}

func getTestLayersWithCodesRepeated(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []Layer{
				{
					Id: 1,
					Codes: []LayerCode{
						{Code: 16},               // q
						{Code: 17},               // w
						{Code: 18},               // e
						{Code: 42},               // shift
						{Code: 13, Modifier: 42}, // +
					},
				},
				{
					Id: 2,
					Codes: []LayerCode{
						{Code: 16}, // q
						{Code: 17}, // w
					},
				},
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
	deviceId := "1"
	lsd := NewLayerDetector(getTestLayers(deviceId), getTestShiftStateConfig())
	// first layer - press "q" and  "w"
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyPress))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 16, KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	// change layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 3, KeyPress))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 3, KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
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
	deviceId := "1"
	lsd := NewLayerDetector(getTestLayers(deviceId), getTestShiftStateConfig())
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
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 13, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 13, KeyRelease))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
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
	deviceId := "1"
	lsd := NewLayerDetector(getTestLayersCodesEmpty(deviceId), getTestShiftStateConfig())
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
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 13, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 13, KeyRelease))
	// first time current layer set=> does not trigger layer change detection
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 42, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
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

func TestWithRepeatedCodes(t *testing.T) {
	deviceId := "1"
	lsd := NewLayerDetector(getTestLayersWithCodesRepeated(deviceId), getTestShiftStateConfig())
	// first layer - press "Q"
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, 17, KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
}
