package k0

import (
	"fmt"
	"log/slog"
	"slices"
)

type LayerInput struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Codes []Key  `json:"codes"`
}

type LayerDetected struct {
	Id       string
	DeviceId string
	LayerId  int64
}

func (ld *LayerDetected) IsDetected() bool {
	return ld.LayerId != 0 && ld.DeviceId != ""
}

type layerDetector struct {
	Layer         LayerInput
	shiftDetector *shiftStateDetector
	mapKeys       map[uint16]bool
}

func (ld *layerDetector) handleKeyEvent(ke DeviceEvent) LayerDetected {
	sd := ld.shiftDetector.handleKeyEvent(ke)
	if sd.IsDetected() && sd.Auto {
		return LayerDetected{LayerId: ld.Layer.Id, DeviceId: ke.DeviceId, Id: sd.ShortcutId}
	}
	if ld.shiftDetector.blockSaveKeylog() {
		//  there is a potential shift state (auto) that needs to be confirmed
		return LayerDetected{}
	}
	if ke.KeyRelease() {
		if _, ok := ld.mapKeys[ke.Code]; ok {
			return LayerDetected{
				LayerId:  ld.Layer.Id,
				DeviceId: ke.DeviceId,
				Id:       fmt.Sprintf("%d", ke.Code),
			}
		}
	}
	return LayerDetected{}
}

// One layer detector per each layer your device has
type layersDetector struct {
	layers               map[string][]layerDetector
	currentLayerDetected *layerDetector
}

func NewLayerDetector(
	devices []DeviceInput,
	shiftStateConfig ShiftStateInput,
) *layersDetector {
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
			hsd := NewHoldShortcutDetector(shiftedCodes, shiftKeys)
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
	oldLayerId := lsd.GetCurrentLayerId()
	ld := lsd.handleKeyEvent(ke)
	if ld.IsDetected() {
		newLayer := lsd.GetCurrentLayerId()
		slog.Debug(fmt.Sprintf("Old layer %d - New layer %d", oldLayerId, newLayer))
		if oldLayerId == newLayer {
			return LayerDetected{}
		}
		// if not equal then there was a change of layer
		return ld
	}
	return LayerDetected{}
}

func (lsd *layersDetector) handleKeyEvent(ke DeviceEvent) LayerDetected {
	numBlockedLayers := 0
	idxPossible := 0
	possibleDetection := LayerDetected{}
	for idx := range lsd.layers[ke.DeviceId] {
		l := &lsd.layers[ke.DeviceId][idx]
		ld := l.handleKeyEvent(ke)
		if ld.IsDetected() && ld.Id != possibleDetection.Id {
			idxPossible = idx
			possibleDetection = ld
		}
		if l.shiftDetector.blockSaveKeylog() {
			numBlockedLayers++
		}
	}
	if numBlockedLayers > 0 {
		return LayerDetected{}
	}
	if possibleDetection.IsDetected() {
		lsd.currentLayerDetected = &lsd.layers[ke.DeviceId][idxPossible]
		return possibleDetection
	}
	if ke.KeyRelease() {
		lsd.currentLayerDetected = nil
	}
	return possibleDetection
}

func (lsd *layersDetector) GetCurrentLayerId() int64 {
	if lsd.currentLayerDetected == nil {
		return 0
	}
	return lsd.currentLayerDetected.Layer.Id
}
