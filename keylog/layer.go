package keylog

import (
	"time"
)

type Layer struct {
	LayerId     int64           `json:"layer_id"`
	Codes       []uint16        `json:"codes"`
	ShiftStates []ShortcutCodes `json:"shift_states"`
}

type LayerDetected struct {
	LayerId int64
}

func (ld *LayerDetected) IsDetected() bool {
	return ld.LayerId != 0
}

type layerDetector struct {
	Layer         Layer
	shiftDetector shiftStateDetector
	mapKeys       map[uint16]bool
}

func (ld *layerDetector) handleKeyEvent(ke DeviceEvent) LayerDetected {
	sd := ld.shiftDetector.handleKeyEvent(ke)
	if sd.IsDetected() {
		return LayerDetected{LayerId: ld.Layer.LayerId}
	}
	if _, ok := ld.mapKeys[ke.Code]; ok {
		return LayerDetected{LayerId: ld.Layer.LayerId}
	}
	return LayerDetected{}
}

func (ld *layerDetector) isHolded() bool {
	return ld.shiftDetector.isHolded()
}

type LayerDuration struct {
	LayerId int64
	Start   time.Time
}

type layersDetector struct {
	layers               []layerDetector
	currentLayerDetected *layerDetector
}

func NewLayerDetector(layers []Layer, thresholdShifted time.Duration) *layersDetector {
	l := []layerDetector{}
	// each layer will have its own detector
	for _, layer := range layers {
		// for single codes we will use a map to detect
		mk := map[uint16]bool{}
		for _, code := range layer.Codes {
			mk[code] = true
		}
		// for shifted states we will use a shift state detector
		hsd := holdShortcutDetector{
			shortcuts: layer.ShiftStates,
			modifiers: getShiftKeys(),
		}
		ssd := NewShiftStateDetectorWithHoldSD(hsd, thresholdShifted)
		ld := layerDetector{
			Layer:         layer,
			shiftDetector: *ssd,
			mapKeys:       mk,
		}
		l = append(l, ld)
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
	if lsd.currentLayerDetected != nil && !isShiftKey(ke.Code) {
		ld := lsd.currentLayerDetected.handleKeyEvent(ke)
		if ld.IsDetected() {
			return ld
		}
	}
	// fmt.Printf("Key %d\n", ke.Code)
	for _, l := range lsd.layers {
		ld := l.handleKeyEvent(ke)
		if ld.IsDetected() && !l.shiftDetector.isHolded() {
			// fmt.Printf("Key %d in layer %d\n", ke.Code, ld.LayerId)
			lsd.currentLayerDetected = &l
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
