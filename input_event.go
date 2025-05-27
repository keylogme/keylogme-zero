package k0

import (
	"time"
)

// InputEvent is the keyboard event structure itself
type InputEvent struct {
	Time  time.Time
	Code  uint16
	Value KeyEvent
}

// KeyString returns representation of pressed key as string
// eg enter, space, a, b, c...
func (i *InputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}

// KeyPress is the value when we press the key on keyboard
func (i *InputEvent) KeyPress() bool {
	return i.Value == 1
}

// KeyRelease is the value when we release the key on keyboard
func (i *InputEvent) KeyRelease() bool {
	return i.Value == 0
}

func (i *InputEvent) IsValid() bool {
	return i.Code != 0 && (i.KeyPress() || i.KeyRelease())
}

// KeyEvent is the keyboard event for up/down (press/release)
type KeyEvent int32

const (
	KeyPress   KeyEvent = 1
	KeyRelease KeyEvent = 0
)
