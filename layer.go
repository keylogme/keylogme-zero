package keylog

import (
	"fmt"
	"slices"
)

type LayerCode struct {
	Code     uint16 `json:"code"`
	Modifier uint16 `json:"modifier"`
}

type Layer struct {
	LayerId int64       `json:"id"`
	Codes   []LayerCode `json:"codes"`
}

type LayerDetected struct {
	DeviceId string
	LayerId  int64
}

func (ld *LayerDetected) IsDetected() bool {
	return ld.LayerId != 0 && ld.DeviceId != ""
}

type layerDetector struct {
	Layer         Layer
	shiftDetector *shiftStateDetector
	mapKeys       map[uint16]bool
}

func (ld *layerDetector) handleKeyEvent(ke DeviceEvent) LayerDetected {
	sd := ld.shiftDetector.handleKeyEvent(ke)
	if sd.IsDetected() && sd.Auto {
		return LayerDetected{LayerId: ld.Layer.LayerId, DeviceId: ke.DeviceId}
	}
	if ld.shiftDetector.blockSaveKeylog() {
		//  there is a potential shift state (auto) that needs to be confirmed
		return LayerDetected{}
	}
	if ke.KeyRelease() {
		if _, ok := ld.mapKeys[ke.Code]; ok {
			return LayerDetected{LayerId: ld.Layer.LayerId, DeviceId: ke.DeviceId}
		}
	}
	return LayerDetected{}
}

type layersDetector struct {
	layers               map[string][]layerDetector
	currentLayerDetected *layerDetector
}

func NewLayerDetector(devices []DeviceInput, shiftStateConfig ShiftState) *layersDetector {
	l := map[string][]layerDetector{}
	for _, dev := range devices {
		l[dev.DeviceId] = []layerDetector{}
		// each layer will have its own detector
		for _, layer := range dev.Layers {
			// for single codes we will use a map to detect
			mk := map[uint16]bool{}
			shiftedCodes := []ShortcutCodes{}
			shiftKeys := []uint16{}
			for _, lc := range layer.Codes {
				if lc.Modifier == 0 {
					mk[lc.Code] = true
					continue
				}

				fakeIDName := fmt.Sprintf("%d_%d", lc.Modifier, lc.Code)
				shiftedCodes = append(
					shiftedCodes,
					ShortcutCodes{
						Id:    fakeIDName,
						Name:  fakeIDName,
						Codes: []uint16{lc.Modifier, lc.Code},
						Type:  HoldShortcutType,
					},
				)
				if !slices.Contains(shiftKeys, lc.Modifier) {
					shiftKeys = append(shiftKeys, lc.Modifier)
				}
			}
			hsd := newHoldShortcutDetector(shiftedCodes, shiftKeys)
			// for shifted states (auto) we will use a shift state detector
			ssd := NewShiftStateDetectorWithHoldSD(hsd, shiftStateConfig)
			ld := layerDetector{
				Layer:         layer,
				shiftDetector: ssd,
				mapKeys:       mk,
			}
			l[dev.DeviceId] = append(l[dev.DeviceId], ld)
		}
	}
	return &layersDetector{layers: l}
}

func (lsd *layersDetector) isLayerChangeDetected(ke DeviceEvent) LayerDetected {
	oldLayerId := int64(0)
	if lsd.GetCurrentLayerId() != 0 {
		oldLayerId = lsd.GetCurrentLayerId()
	}
	ld := lsd.handleKeyEvent(ke)
	if oldLayerId != 0 && ld.IsDetected() {
		if oldLayerId == lsd.GetCurrentLayerId() {
			return LayerDetected{}
		}
		// if not equal then there was a change of layer
		return ld
	}
	return LayerDetected{}
}

func (lsd *layersDetector) handleKeyEvent(ke DeviceEvent) LayerDetected {
	// First check if key is in current layer (fast access to non shift keys)
	if lsd.currentLayerDetected != nil && !isShiftKey(ke.Code) {
		ld := lsd.currentLayerDetected.handleKeyEvent(ke)
		if ld.IsDetected() {
			return ld
		}
	}
	// Check in all layers
	for idx := range lsd.layers[ke.DeviceId] {
		l := &lsd.layers[ke.DeviceId][idx]
		ld := l.handleKeyEvent(ke)
		if ld.IsDetected() {
			lsd.currentLayerDetected = l
			return ld
		}
	}
	return LayerDetected{}
}

func (lsd *layersDetector) GetCurrentLayerId() int64 {
	if lsd.currentLayerDetected == nil {
		return 0
	}
	return lsd.currentLayerDetected.Layer.LayerId
}
