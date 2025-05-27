package k0

import (
	"fmt"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

type ShiftStateInput struct {
	ThresholdAuto types.Duration `json:"threshold_auto"`
}

// Auto is true when the shift state is triggered by the microcontroller
type ShiftStateDetected struct {
	ShortcutId           string
	DeviceId             string
	Modifier             uint16
	Code                 uint16
	Auto                 bool
	DiffTimePressMicro   int64
	DiffTimeReleaseMicro int64
}

func (ssd *ShiftStateDetected) IsDetected() bool {
	return ssd.DeviceId != ""
}

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
type shiftStateDetector struct {
	holdDetector           holdShortcutDetector
	lastModPressTime       int64 // unix micro
	lastKeyPressTime       int64 // unix micro
	lastKeyReleaseTime     int64 // unix micro
	lastModReleaseTime     int64 // unix micro
	thresholdAuto          time.Duration
	possibleAutoShiftState ShiftStateDetected
	mapIdToCodes           map[string]Key
}

func NewShiftStateDetector(config ShiftStateInput) *shiftStateDetector {
	scs := getShortcutCodesForShiftState()
	mapId := getMapIdToCodes()
	return &shiftStateDetector{
		holdDetector:  NewHoldShortcutDetector(scs, SHIFT_CODES),
		thresholdAuto: config.ThresholdAuto.Duration,
		mapIdToCodes:  mapId,
	}
}

// TODO: add documentation
func NewShiftStateDetectorWithHoldSD(
	hd holdShortcutDetector,
	config ShiftStateInput,
) *shiftStateDetector {
	return &shiftStateDetector{
		holdDetector:  hd,
		thresholdAuto: config.ThresholdAuto.Duration,
	}
}

func (skd *shiftStateDetector) isHolded() bool {
	return skd.holdDetector.isHolded()
}

// an auto shift state will block saving keylogs because it's not a human typing
// f.e. 42(shift) + 13(=) => +, instead of saving 42 and 13 , I will only save shifted 13
func (skd *shiftStateDetector) blockSaveKeylog() bool {
	return skd.possibleAutoShiftState.IsDetected()
}

func (skd *shiftStateDetector) handleKeyEvent(ke DeviceEvent) ShiftStateDetected {
	sd := skd.holdDetector.handleKeyEvent(ke)
	skd.setTimes(ke)
	if sd.IsDetected() && skd.isHolded() {
		k, ok := skd.mapIdToCodes[sd.ShortcutId]
		if !ok {
			skd.possibleAutoShiftState = ShiftStateDetected{}
			return ShiftStateDetected{}
		}
		auto := false
		diffTimeMicro := skd.lastKeyPressTime - skd.lastModPressTime
		sdetect := ShiftStateDetected{
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
			return ShiftStateDetected{}
		}
		return sdetect
	}
	if !skd.isHolded() && skd.possibleAutoShiftState.IsDetected() {
		diffTimeMicro := skd.lastModReleaseTime - skd.lastKeyReleaseTime
		result := skd.possibleAutoShiftState
		result.Auto = false
		result.DiffTimeReleaseMicro = diffTimeMicro
		skd.possibleAutoShiftState = ShiftStateDetected{}
		if time.Duration(time.Microsecond*time.Duration(diffTimeMicro)) < skd.thresholdAuto {
			// confirm auto shift state
			result.Auto = true
		}
		return result
	}
	return ShiftStateDetected{}
}

func (skd *shiftStateDetector) setTimes(ke DeviceEvent) {
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

func (skd *shiftStateDetector) reset() {
	skd.lastModPressTime = 0
	skd.lastKeyReleaseTime = 0
	skd.lastKeyPressTime = 0
	skd.lastModReleaseTime = 0
	skd.possibleAutoShiftState = ShiftStateDetected{}
}

func getShiftCodeKey(shiftCode, code uint16) string {
	return fmt.Sprintf("%d_%d", shiftCode, code)
}

func getShortcutCodesForShiftState() []ShortcutCodes {
	listSS := []ShortcutCodes{}
	for _, sc := range SHIFT_CODES {
		for _, c := range ALL_CODES {
			scKey := getShiftCodeKey(sc, c)
			ssc := ShortcutCodes{
				Id:    scKey,
				Codes: []uint16{sc, c},
				Type:  HoldShortcutType,
			}
			listSS = append(listSS, ssc)
		}
	}
	return listSS
}

func getMapIdToCodes() map[string]Key {
	mapIdToCodes := make(map[string]Key)
	for _, sc := range SHIFT_CODES {
		for _, c := range ALL_CODES {
			scKey := getShiftCodeKey(sc, c)
			mapIdToCodes[scKey] = Key{Code: c, Modifier: sc}
		}
	}
	return mapIdToCodes
}
