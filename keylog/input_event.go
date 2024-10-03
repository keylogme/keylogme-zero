package keylog

import (
	"syscall"
	"unsafe"
)

const (
	// EvSyn is used as markers to separate events. Events may be separated in time or in space, such as with the multitouch protocol.
	EvSyn eventType = 0x00
	// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
	EvKey eventType = 0x01
	// EvRel is used to describe relative axis value changes, e.g. moving the mouse 5 units to the left.
	EvRel eventType = 0x02
	// EvAbs is used to describe absolute axis value changes, e.g. describing the coordinates of a touch on a touchscreen.
	EvAbs eventType = 0x03
	// EvMsc is used to describe miscellaneous input data that do not fit into other types.
	EvMsc eventType = 0x04
	// EvSw is used to describe binary state input switches.
	EvSw eventType = 0x05
	// EvLed is used to turn LEDs on devices on and off.
	EvLed eventType = 0x11
	// EvSnd is used to output sound to devices.
	EvSnd eventType = 0x12
	// EvRep is used for autorepeating devices.
	EvRep eventType = 0x14
	// EvFf is used to send force feedback commands to an input device.
	EvFf eventType = 0x15
	// EvPwr is a special type for power button and switch input.
	EvPwr eventType = 0x16
	// EvFfStatus is used to receive force feedback device status.
	EvFfStatus eventType = 0x17
)

// eventType are groupings of codes under a logical input construct.
// Each type has a set of applicable codes to be used in generating events.
// See the Ev section for details on valid codes for each type
type eventType uint16

// eventsize is size of structure of InputEvent
var eventsize = int(unsafe.Sizeof(inputEvent{}))

// inputEvent is the keyboard event structure itself
type inputEvent struct {
	Time  syscall.Timeval
	Type  eventType
	Code  uint16
	Value int32
}

// KeyString returns representation of pressed key as string
// eg enter, space, a, b, c...
func (i *inputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}

// KeyPress is the value when we press the key on keyboard
func (i *inputEvent) KeyPress() bool {
	return i.Value == 1
}

// KeyRelease is the value when we release the key on keyboard
func (i *inputEvent) KeyRelease() bool {
	return i.Value == 0
}

// keyEvent is the keyboard event for up/down (press/release)
type keyEvent int32

const (
	KeyPress   keyEvent = 1
	KeyRelease keyEvent = 0
)
