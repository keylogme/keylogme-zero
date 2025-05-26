package shift

import (
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/shortcut"
	"github.com/keylogme/keylogme-zero/internal/types"
)

// this detects shift + keys that are triggered by
// -microcontroller of keyboard
// -human
// for example, qmk firmware allows you to define a key like '+' which is 'L_Shift' and '='
// when you press that key, the microcontroller sends then next events very fast:
// 1. key hold event for 'L_Shift'
// 2. hold/release for '='
// 3. release for 'L_Shift'
// WHY DO WE NEED THIS?
// setting up a symbol layer in qmk firmware is a common practice
// you don't need to press shift to type symbols that are common in programming
// (but you still need to press one to access the symbol layer)
type ShiftStateDetector struct {
	holdDetector           shortcut.HoldShortcutDetector
	lastModPressTime       int64 // unix micro
	lastKeyPressTime       int64 // unix micro
	lastKeyReleaseTime     int64 // unix micro
	lastModReleaseTime     int64 // unix micro
	thresholdAuto          time.Duration
	possibleAutoShiftState types.ShiftStateDetected
	mapIdToCodes           map[string]types.Key
}

func NewShiftStateDetector(config types.ShiftStateInput) *ShiftStateDetector {
	scs := getShortcutCodesForShiftState()
	mapId := getMapIdToCodes()
	return &ShiftStateDetector{
		holdDetector:  shortcut.NewHoldShortcutDetector(scs, keylogger.SHIFT_CODES),
		thresholdAuto: config.ThresholdAuto.Duration,
		mapIdToCodes:  mapId,
	}
}

// TODO: add documentation
func NewShiftStateDetectorWithHoldSD(
	hd shortcut.HoldShortcutDetector,
	config types.ShiftStateInput,
) *ShiftStateDetector {
	return &ShiftStateDetector{
		holdDetector:  hd,
		thresholdAuto: config.ThresholdAuto.Duration,
	}
}

func (skd *ShiftStateDetector) isHolded() bool {
	return skd.holdDetector.IsHolded()
}

// an auto shift state will block saving keylogs because it's not a human typing
// f.e. 42(shift) + 13(=) => +, instead of saving 42 and 13 , I will only save shifted 13
func (skd *ShiftStateDetector) BlockSaveKeylog() bool {
	return skd.possibleAutoShiftState.IsDetected()
}

func (skd *ShiftStateDetector) HandleKeyEvent(ke keylogger.DeviceEvent) types.ShiftStateDetected {
	sd := skd.holdDetector.HandleKeyEvent(ke)
	skd.setTimes(ke)
	if sd.IsDetected() && skd.isHolded() {
		k, ok := skd.mapIdToCodes[sd.ShortcutId]
		if !ok {
			skd.possibleAutoShiftState = types.ShiftStateDetected{}
			return types.ShiftStateDetected{}
		}
		auto := false
		diffTimeMicro := skd.lastKeyPressTime - skd.lastModPressTime
		sdetect := types.ShiftStateDetected{
			ShortcutId:         sd.ShortcutId,
			DeviceId:           ke.DeviceId,
			Modifier:           k.Modifier,
			Code:               k.Code,
			Auto:               auto,
			DiffTimePressMicro: diffTimeMicro,
		}
		if time.Duration(time.Microsecond*time.Duration(diffTimeMicro)) < skd.thresholdAuto {
			// auto shift needs confirmation on shift release
			skd.possibleAutoShiftState = sdetect
			return types.ShiftStateDetected{}
		}
		return sdetect
	}
	if !skd.isHolded() && skd.possibleAutoShiftState.IsDetected() {
		diffTimeMicro := skd.lastModReleaseTime - skd.lastKeyReleaseTime
		result := skd.possibleAutoShiftState
		result.Auto = false
		result.DiffTimeReleaseMicro = diffTimeMicro
		skd.possibleAutoShiftState = types.ShiftStateDetected{}
		if time.Duration(time.Microsecond*time.Duration(diffTimeMicro)) < skd.thresholdAuto {
			// confirm auto shift state
			result.Auto = true
		}
		return result
	}
	return types.ShiftStateDetected{}
}

func (skd *ShiftStateDetector) setTimes(ke keylogger.DeviceEvent) {
	t := ke.Time

	if skd.isHolded() && skd.lastModPressTime == 0 {
		skd.reset()
		skd.lastModPressTime = t.UnixMicro()
	}

	if ke.KeyRelease() && skd.isHolded() {
		skd.lastKeyReleaseTime = t.UnixMicro()
	}

	if ke.KeyPress() && skd.isHolded() {
		skd.lastKeyPressTime = t.UnixMicro()
	}
	if ke.KeyRelease() && !skd.isHolded() {
		skd.lastModPressTime = 0
		skd.lastModReleaseTime = t.UnixMicro()
	}

	if ke.KeyPress() && !skd.isHolded() {
		skd.reset()
	}
}

func (skd *ShiftStateDetector) reset() {
	skd.lastModPressTime = 0
	skd.lastKeyReleaseTime = 0
	skd.lastKeyPressTime = 0
	skd.lastModReleaseTime = 0
	skd.possibleAutoShiftState = types.ShiftStateDetected{}
}
