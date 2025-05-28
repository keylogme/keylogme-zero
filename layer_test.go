package k0

import (
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/types"
)

var layer1ShiftedCode = Key{
	Code: ALL_CODES[3], Modifier: SHIFT_CODES[0],
}

var layer1 = []Key{
	{Code: ALL_CODES[0]},
	{Code: ALL_CODES[1]},
	{Code: ALL_CODES[2]},
	{Code: SHIFT_CODES[0]},
	layer1ShiftedCode,
}

var layer2ShiftedCode = Key{
	Code: ALL_CODES[8], Modifier: SHIFT_CODES[0],
}

var layer2 = []Key{
	{Code: ALL_CODES[5]},
	{Code: ALL_CODES[6]},
	{Code: ALL_CODES[7]},
	layer2ShiftedCode,
}

// layer with 2 keys: same code but with and without modifier
var layer3ShiftedCode = Key{
	Code: ALL_CODES[5], Modifier: SHIFT_CODES[0],
}

var layer3 = []Key{
	{Code: ALL_CODES[5]},
	layer3ShiftedCode,
}

func getTestLayers(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []LayerInput{
				{
					Id:    1,
					Codes: layer1,
				},
				{
					Id:    2,
					Codes: layer2,
				},
			},
		},
	}
}

// layer 3 has 2 keys: same code but with and without modifier
func getTestLayersCodesEmpty(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []LayerInput{
				{
					Id:    1,
					Codes: layer1,
				},
				{
					Id:    3,
					Codes: layer3,
				},
			},
		},
	}
}

func getTestLayersWithCodesRepeated(kdevId string) []DeviceInput {
	return []DeviceInput{
		{
			DeviceId: kdevId,
			Layers: []LayerInput{
				{
					Id:    1,
					Codes: layer1,
				},
				{
					Id:    2,
					Codes: layer1,
				},
			},
		},
	}
}

func getTestShiftStateConfig() ShiftStateInput {
	return ShiftStateInput{
		ThresholdAuto: types.Duration{Duration: 100 * time.Millisecond},
	}
}

func TestChangeLayerSingleCodes(t *testing.T) {
	deviceId := "1"
	lsd := NewLayersDetector(getTestLayers(deviceId), getTestShiftStateConfig())
	// first layer - press "q" and  "w"
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyPress))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[1].Code, keylogger.KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[1].Code, keylogger.KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	// change layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer2[0].Code, keylogger.KeyPress))
	// t.Log(lsd.GetCurrentLayerId())
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer2[0].Code, keylogger.KeyRelease))
	// t.Log(lsd.GetCurrentLayerId())
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 2 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithShiftedCodesInMultipleLayers(t *testing.T) {
	deviceId := "1"
	lsd := NewLayersDetector(getTestLayers(deviceId), getTestShiftStateConfig())
	// first layer - press "Q"
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld := lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Modifier, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Code, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Code, keylogger.KeyRelease),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Modifier, keylogger.KeyRelease),
	)
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 { // after Shift+Q press , layer id = 1 is set
		t.Fatal("Layer id incorrect")
	}
	// change layer- use shifted code in second layer
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer2ShiftedCode.Modifier, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer2ShiftedCode.Code, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer2ShiftedCode.Code, keylogger.KeyRelease),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer2ShiftedCode.Modifier, keylogger.KeyRelease),
	)
	if !ld.IsDetected() {
		t.Fatal("Detection expected ")
	}
	if lsd.GetCurrentLayerId() != 2 { //  layer id is changed
		t.Fatal("Layer id incorrect")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithShiftedCodesInMultipleLayers_CodesEmpty(t *testing.T) {
	deviceId := "1"
	lsd := NewLayersDetector(getTestLayersCodesEmpty(deviceId), getTestShiftStateConfig())
	// first layer - press "Q"
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld := lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Modifier, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Code, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Code, keylogger.KeyRelease),
	)
	// first time current layer set=> does not trigger layer change detection
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer1ShiftedCode.Modifier, keylogger.KeyRelease),
	)
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
	// change layer- use shifted code in second layer
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer3ShiftedCode.Modifier, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 1 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer3ShiftedCode.Code, keylogger.KeyPress),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected ")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer3ShiftedCode.Code, keylogger.KeyRelease),
	)
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(
		getFakeEvent(deviceId, layer3ShiftedCode.Modifier, keylogger.KeyRelease),
	)
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 3 { //  layer id is changed
		t.Fatal("Layer id incorrect")
	}
	// change back to first layer
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
}

func TestWithRepeatedCodes(t *testing.T) {
	deviceId := "1"
	lsd := NewLayersDetector(getTestLayersWithCodesRepeated(deviceId), getTestShiftStateConfig())
	// first layer - press "Q"
	if lsd.GetCurrentLayerId() != 0 {
		t.Fatal("Layer id incorrect")
	}
	ld := lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyPress))
	if ld.IsDetected() {
		t.Fatal("Detection not expected")
	}
	if lsd.GetCurrentLayerId() != 0 { // shift key are not deterministic, should not trigger a layer change
		t.Fatal("Layer id incorrect")
	}
	ld = lsd.isLayerChangeDetected(getFakeEvent(deviceId, layer1[0].Code, keylogger.KeyRelease))
	if !ld.IsDetected() {
		t.Fatal("Detection expected")
	}
	if lsd.GetCurrentLayerId() != 1 {
		t.Fatal("Layer id incorrect")
	}
}
