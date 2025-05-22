package keylog

import (
	"fmt"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
)

// Auto is true when the shift state is triggered by the microcontroller
type shiftStateDetected struct {
	ShortcutId           string
	DeviceId             string
	Modifier             uint16
	Code                 uint16
	Auto                 bool
	DiffTimePressMicro   int64
	DiffTimeReleaseMicro int64
}

func (ssd *shiftStateDetected) IsDetected() bool {
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
	possibleAutoShiftState shiftStateDetected
}

func getShiftCodeKey(shiftCode, code uint16) string {
	return fmt.Sprintf("%d_%d", shiftCode, code)
}

func getShortcutCodesForShiftState(shiftCodes []uint16) []ShortcutCodes {
	numCodes := []uint16{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	lettersCodes := []uint16{
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 44, 45, 46, 47, 48, 49, 50,
	}
	symbolsCodes := []uint16{12, 13, 26, 27, 28, 39, 40, 41, 43, 51, 52, 53}
	allCodes := append(numCodes, append(lettersCodes, symbolsCodes...)...)
	listSS := []ShortcutCodes{}
	for _, sc := range shiftCodes {
		for _, c := range allCodes {
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

func NewShiftStateDetector(config ShiftState) *shiftStateDetector {
	shiftMods := []uint16{42, 54}
	scs := getShortcutCodesForShiftState(shiftMods)
	return &shiftStateDetector{
		holdDetector:  newHoldShortcutDetector(scs, shiftMods),
		thresholdAuto: config.ThresholdAuto.Duration,
	}
}

func NewShiftStateDetectorWithHoldSD(
	hd holdShortcutDetector,
	config ShiftState,
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

func (skd *shiftStateDetector) handleKeyEvent(ke keylogger.DeviceEvent) shiftStateDetected {
	sd := skd.holdDetector.handleKeyEvent(ke)
	skd.setTimes(ke)
	if sd.IsDetected() && skd.isHolded() {
		mod := skd.holdDetector.modPress[0] // by default first element is the modifier
		auto := false
		diffTimeMicro := skd.lastKeyPressTime - skd.lastModPressTime
		sdetect := shiftStateDetected{
			ShortcutId:         sd.ShortcutId,
			DeviceId:           ke.DeviceId,
			Modifier:           mod,
			Code:               ke.Code,
			Auto:               auto,
			DiffTimePressMicro: diffTimeMicro,
		}
		if time.Duration(time.Microsecond*time.Duration(diffTimeMicro)) < skd.thresholdAuto {
			// auto shift needs confirmation on shift release
			skd.possibleAutoShiftState = sdetect
			return shiftStateDetected{}
		}
		return sdetect
	}
	if !skd.isHolded() && skd.possibleAutoShiftState.IsDetected() {
		diffTimeMicro := skd.lastModReleaseTime - skd.lastKeyReleaseTime
		result := skd.possibleAutoShiftState
		result.Auto = false
		result.DiffTimeReleaseMicro = diffTimeMicro
		skd.possibleAutoShiftState = shiftStateDetected{}
		if time.Duration(time.Microsecond*time.Duration(diffTimeMicro)) < skd.thresholdAuto {
			// confirm auto shift state
			result.Auto = true
		}
		return result
	}
	return shiftStateDetected{}
}

func (skd *shiftStateDetector) setTimes(ke keylogger.DeviceEvent) {
	t := ke.ExecTime

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
	skd.possibleAutoShiftState = shiftStateDetected{}
}
