package layer

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/shift"
	"github.com/keylogme/keylogme-zero/internal/shortcut"
	"github.com/keylogme/keylogme-zero/internal/types"
)

type layerDetector struct {
	Layer         types.Layer
	shiftDetector *shift.ShiftStateDetector
	mapKeys       map[uint16]bool
}

func (ld *layerDetector) handleKeyEvent(ke keylogger.DeviceEvent) types.LayerDetected {
	sd := ld.shiftDetector.HandleKeyEvent(ke)
	if sd.IsDetected() && sd.Auto {
		return types.LayerDetected{LayerId: ld.Layer.Id, DeviceId: ke.DeviceId, Id: sd.ShortcutId}
	}
	if ld.shiftDetector.BlockSaveKeylog() {
		//  there is a potential shift state (auto) that needs to be confirmed
		return types.LayerDetected{}
	}
	if ke.KeyRelease() {
		if _, ok := ld.mapKeys[ke.Code]; ok {
			return types.LayerDetected{
				LayerId:  ld.Layer.Id,
				DeviceId: ke.DeviceId,
				Id:       fmt.Sprintf("%d", ke.Code),
			}
		}
	}
	return types.LayerDetected{}
}

// One layer detector per each layer your device has
type LayersDetector struct {
	layers               map[string][]layerDetector
	currentLayerDetected *layerDetector
}

func NewLayerDetector(
	devices []keylogger.DeviceInput,
	shiftStateConfig types.ShiftStateInput,
) *LayersDetector {
	l := map[string][]layerDetector{}
	for _, dev := range devices {
		l[dev.DeviceId] = []layerDetector{}
		// each layer will have its own detector
		for _, layer := range dev.Layers {
			// for single codes we will use a map to detect
			mk := map[uint16]bool{}
			shiftedCodes := []types.ShortcutCodes{}
			shiftKeys := []uint16{}
			for _, lc := range layer.Codes {
				if lc.Modifier == 0 {
					mk[lc.Code] = true
					continue
				}

				fakeIDName := fmt.Sprintf("%d_%d", lc.Modifier, lc.Code)
				shiftedCodes = append(
					shiftedCodes,
					types.ShortcutCodes{
						Id:    fakeIDName,
						Name:  fakeIDName,
						Codes: []uint16{lc.Modifier, lc.Code},
						Type:  types.HoldShortcutType,
					},
				)
				if !slices.Contains(shiftKeys, lc.Modifier) {
					shiftKeys = append(shiftKeys, lc.Modifier)
				}
			}
			hsd := shortcut.NewHoldShortcutDetector(shiftedCodes, shiftKeys)
			// for shifted states (auto) we will use a shift state detector
			ssd := shift.NewShiftStateDetectorWithHoldSD(hsd, shiftStateConfig)
			ld := layerDetector{
				Layer:         layer,
				shiftDetector: ssd,
				mapKeys:       mk,
			}
			l[dev.DeviceId] = append(l[dev.DeviceId], ld)
		}
	}
	return &LayersDetector{layers: l}
}

func (lsd *LayersDetector) IsLayerChangeDetected(ke keylogger.DeviceEvent) types.LayerDetected {
	oldLayerId := lsd.GetCurrentLayerId()
	ld := lsd.HandleKeyEvent(ke)
	if ld.IsDetected() {
		newLayer := lsd.GetCurrentLayerId()
		slog.Debug(fmt.Sprintf("Old layer %d - New layer %d", oldLayerId, newLayer))
		if oldLayerId == newLayer {
			return types.LayerDetected{}
		}
		// if not equal then there was a change of layer
		return ld
	}
	return types.LayerDetected{}
}

func (lsd *LayersDetector) HandleKeyEvent(ke keylogger.DeviceEvent) types.LayerDetected {
	numBlockedLayers := 0
	idxPossible := 0
	possibleDetection := types.LayerDetected{}
	for idx := range lsd.layers[ke.DeviceId] {
		l := &lsd.layers[ke.DeviceId][idx]
		ld := l.handleKeyEvent(ke)
		if ld.IsDetected() && ld.Id != possibleDetection.Id {
			idxPossible = idx
			possibleDetection = ld
		}
		if l.shiftDetector.BlockSaveKeylog() {
			numBlockedLayers++
		}
	}
	if numBlockedLayers > 0 {
		return types.LayerDetected{}
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

func (lsd *LayersDetector) GetCurrentLayerId() int64 {
	if lsd.currentLayerDetected == nil {
		return 0
	}
	return lsd.currentLayerDetected.Layer.Id
}
