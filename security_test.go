package k0

import (
	"math"
	"slices"
	"testing"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
)

func TestRandSeed(t *testing.T) {
	totalRolls := 6000
	max := 6

	mapVals := map[int]int{}
	for range totalRolls {
		r := getRandInt(max)
		if _, ok := mapVals[r]; !ok {
			mapVals[r] = 0
		}
		mapVals[r] += 1
	}

	expectedFreq := totalRolls / max
	limitOffset := 10.0 // count per value can be off to expectedFreq by limitOffset %

	for _, val := range mapVals {
		absDif := math.Abs(float64(val-expectedFreq)) / float64(expectedFreq)
		// t.Logf("%d/%d %f\n", val-expectedFreq, expectedFreq, absDif)
		if absDif*100 > limitOffset {
			t.Fatal("Random value has crossed limit offset")
		}
	}
}

func TestNewBaggage(t *testing.T) {
	bagggageSize := 3
	b := newBaggage(bagggageSize)

	deviceId := "1"
	de := DeviceEventInLayer{
		DeviceEvent: DeviceEvent{
			DeviceId: deviceId,
			InputEvent: keylogger.InputEvent{
				Code: uint16(0),
				Type: keylogger.KeyPress,
			},
		},
		LayerId: 1,
	}
	for i := range bagggageSize {
		de.Code = uint16(i)
		isAuth := b.isAuthorized(&de)
		if isAuth {
			t.Fatalf("Expected baggage to not be authorized for code %d, but it was", i)
		}
	}

	resultValues := []uint16{}
	for _, b := range b.devices[deviceId] {
		resultValues = append(resultValues, b.Code)
	}

	if slices.Compare(resultValues, []uint16{0, 1, 2}) != 0 {
		t.Fatalf("Expected baggage to contain codes [0, 1, 2], but got %v", b.devices[deviceId])
	}

	authorizedCode := uint16(bagggageSize)
	de.Code = authorizedCode
	isAuth := b.isAuthorized(&de)
	if !isAuth {
		t.Fatalf("Expected authorized")
	}

	if de.Code == authorizedCode {
		t.Fatalf("Expected code being swapped")
	}
}

func TestGhostingKeys(t *testing.T) {
	g := newGhostingCodes([]uint16{1, 2})

	deviceId := "1"
	de := DeviceEventInLayer{
		DeviceEvent: DeviceEvent{
			DeviceId: deviceId,
			InputEvent: keylogger.InputEvent{
				Code: uint16(0),
				Type: keylogger.KeyPress,
			},
		},
		LayerId: 1,
	}
	isAuth := g.isAuthorized(&de)
	if !isAuth {
		t.Fatal("should be auth")
	}

	de.Code = 1
	isAuth = g.isAuthorized(&de)
	if isAuth {
		t.Fatal("should not be auth")
	}
}

func TestSecurity(t *testing.T) {
	deviceId := "1"
	si := SecurityInput{
		BaggageSize:   10,
		GhostingCodes: []uint16{25, 26, 27},
	}
	s := NewSecurity(si)

	de := DeviceEventInLayer{
		DeviceEvent: DeviceEvent{
			DeviceId: deviceId,
			InputEvent: keylogger.InputEvent{
				Code: uint16(0),
				Type: keylogger.KeyPress,
			},
		},
		LayerId: 1,
	}

	// fill baggage to size - 1
	numCodesToNotFill := si.BaggageSize - 1
	t.Log(numCodesToNotFill)
	for i := range numCodesToNotFill {
		de.Code = uint16(i)
		isAuth := s.isAuthorized(&de)
		if isAuth {
			t.Fatalf("Expected baggage to not be authorized for code %d, but it was", i)
		}
	}

	// ghost code when baggage is not full
	de.Code = si.GhostingCodes[0]
	isAuth := s.isAuthorized(&de)
	if isAuth {
		t.Fatal("should not be auth with ghost code")
	}

	// after ghost code , baggage should not be full yet
	de.Code = 0
	isAuth = s.isAuthorized(&de)
	if isAuth {
		t.Fatal("baggage should not be full yet")
	}

	// now baggage is full
	uniqueCode := uint16(100)
	de.Code = uniqueCode
	isAuth = s.isAuthorized(&de)
	if !isAuth {
		t.Fatal("should be auth")
	}

	if de.Code == uniqueCode {
		t.Fatal("code should have been swapped")
	}
}
