package keylogger

import "slices"

// Common key constants used in both maps
const (
	KEY_A               = "A"
	KEY_B               = "B"
	KEY_C               = "C"
	KEY_D               = "D"
	KEY_E               = "E"
	KEY_F               = "F"
	KEY_G               = "G"
	KEY_H               = "H"
	KEY_I               = "I"
	KEY_J               = "J"
	KEY_K               = "K"
	KEY_L               = "L"
	KEY_M               = "M"
	KEY_N               = "N"
	KEY_O               = "O"
	KEY_P               = "P"
	KEY_Q               = "Q"
	KEY_R               = "R"
	KEY_S               = "S"
	KEY_T               = "T"
	KEY_U               = "U"
	KEY_V               = "V"
	KEY_W               = "W"
	KEY_X               = "X"
	KEY_Y               = "Y"
	KEY_Z               = "Z"
	KEY_0               = "0"
	KEY_1               = "1"
	KEY_2               = "2"
	KEY_3               = "3"
	KEY_4               = "4"
	KEY_5               = "5"
	KEY_6               = "6"
	KEY_7               = "7"
	KEY_8               = "8"
	KEY_9               = "9"
	KEY_ENTER           = "ENTER"
	KEY_ESCAPE          = "ESCAPE"
	KEY_BACKSPACE       = "BACKSPACE"
	KEY_TAB             = "TAB"
	KEY_SPACE           = "SPACE"
	KEY_HYPHEN          = "-"
	KEY_EQUALS          = "="
	KEY_LEFT_BRACKET    = "["
	KEY_RIGHT_BRACKET   = "]"
	KEY_BACKSLASH       = "\\"
	KEY_SEMICOLON       = ";"
	KEY_SINGLE_QUOTE    = "'"
	KEY_BACKTICK        = "`"
	KEY_COMMA           = ","
	KEY_PERIOD          = "."
	KEY_SLASH           = "/"
	KEY_CAPS_LOCK       = "CAPS_LOCK"
	KEY_F1              = "F1"
	KEY_F2              = "F2"
	KEY_F3              = "F3"
	KEY_F4              = "F4"
	KEY_F5              = "F5"
	KEY_F6              = "F6"
	KEY_F7              = "F7"
	KEY_F8              = "F8"
	KEY_F9              = "F9"
	KEY_F10             = "F10"
	KEY_F11             = "F11"
	KEY_F12             = "F12"
	KEY_SCROLL_LOCK     = "SCROLL_LOCK"
	KEY_INSERT          = "INSERT"
	KEY_HOME            = "HOME"
	KEY_PAGE_UP         = "PAGE_UP"
	KEY_DELETE          = "DELETE"
	KEY_END             = "END"
	KEY_PAGE_DOWN       = "PAGE_DOWN"
	KEY_RIGHT_ARROW     = "RIGHT_ARROW"
	KEY_LEFT_ARROW      = "LEFT_ARROW"
	KEY_DOWN_ARROW      = "DOWN_ARROW"
	KEY_UP_ARROW        = "UP_ARROW"
	KEY_NUM_LOCK        = "NUM_LOCK"
	KEY_KEYPAD_SLASH    = "KEYPAD /"
	KEY_KEYPAD_ASTERISK = "KEYPAD *"
	KEY_KEYPAD_MINUS    = "KEYPAD -"
	KEY_KEYPAD_PLUS     = "KEYPAD +"
	KEY_KEYPAD_ENTER    = "KEYPAD ENTER"
	KEY_KEYPAD_0        = "KEYPAD 0"
	KEY_KEYPAD_1        = "KEYPAD 1"
	KEY_KEYPAD_2        = "KEYPAD 2"
	KEY_KEYPAD_3        = "KEYPAD 3"
	KEY_KEYPAD_4        = "KEYPAD 4"
	KEY_KEYPAD_5        = "KEYPAD 5"
	KEY_KEYPAD_6        = "KEYPAD 6"
	KEY_KEYPAD_7        = "KEYPAD 7"
	KEY_KEYPAD_8        = "KEYPAD 8"
	KEY_KEYPAD_9        = "KEYPAD 9"
	KEY_KEYPAD_PERIOD   = "KEYPAD ."
	KEY_KEYPAD_EQUALS   = "KEYPAD ="
	KEY_LEFT_CTRL       = "LEFT_CTRL"
	KEY_LEFT_SHIFT      = "LEFT_SHIFT"
	KEY_LEFT_ALT        = "LEFT_ALT"
	KEY_LEFT_GUI        = "LEFT_GUI"
	KEY_RIGHT_CTRL      = "RIGHT_CTRL"
	KEY_RIGHT_SHIFT     = "RIGHT_SHIFT"
	KEY_RIGHT_ALT       = "RIGHT_ALT"
	KEY_RIGHT_GUI       = "RIGHT_GUI"
	KEY_VOLUME_UP       = "VOLUME_UP"
	KEY_VOLUME_DOWN     = "VOLUME_DOWN"
)

var (
	// Letter keys (A-Z)
	LetterKeys = []string{
		KEY_A, KEY_B, KEY_C, KEY_D, KEY_E, KEY_F, KEY_G, KEY_H, KEY_I, KEY_J,
		KEY_K, KEY_L, KEY_M, KEY_N, KEY_O, KEY_P, KEY_Q, KEY_R, KEY_S, KEY_T,
		KEY_U, KEY_V, KEY_W, KEY_X, KEY_Y, KEY_Z,
	}

	// Number keys (0-9) from the main keyboard, not keypad
	NumberKeys = []string{
		KEY_0, KEY_1, KEY_2, KEY_3, KEY_4, KEY_5, KEY_6, KEY_7, KEY_8, KEY_9,
	}

	// Symbol keys that typically do not require the Shift key
	SymbolKeys = []string{
		KEY_HYPHEN, KEY_EQUALS, KEY_LEFT_BRACKET, KEY_RIGHT_BRACKET, KEY_BACKSLASH,
		KEY_SEMICOLON, KEY_SINGLE_QUOTE, KEY_BACKTICK, KEY_COMMA, KEY_PERIOD, KEY_SLASH,
	}

	// Modifier keys related to Shift
	ShiftKeys = []string{
		KEY_LEFT_SHIFT, KEY_RIGHT_SHIFT,
	}

	// Modifier keys related to Control (Ctrl)
	CtrlKeys = []string{
		KEY_LEFT_CTRL, KEY_RIGHT_CTRL,
	}

	// Modifier keys related to Alt
	AltKeys = []string{
		KEY_LEFT_ALT, KEY_RIGHT_ALT,
	}
)

func inverseMapCodeMap() map[string]uint16 {
	out := map[string]uint16{}
	for key, val := range keyCodeMap {
		out[val] = key
	}
	return out
}

var inversedMap = inverseMapCodeMap()

func GetAllCodes() []uint16 {
	codes := []uint16{}
	for _, val := range NumberKeys {
		codes = append(codes, inversedMap[val])
	}
	for _, val := range LetterKeys {
		codes = append(codes, inversedMap[val])
	}
	for _, val := range SymbolKeys {
		codes = append(codes, inversedMap[val])
	}
	return codes
}

func GetShiftCodes() []uint16 {
	codes := []uint16{}
	for _, val := range ShiftKeys {
		codes = append(codes, inversedMap[val])
	}
	return codes
}

func GetCtrlCodes() []uint16 {
	codes := []uint16{}
	for _, val := range CtrlKeys {
		codes = append(codes, inversedMap[val])
	}
	return codes
}

func GetAltCodes() []uint16 {
	codes := []uint16{}
	for _, val := range AltKeys {
		codes = append(codes, inversedMap[val])
	}
	return codes
}

func GetAllModifierCodes() []uint16 {
	return slices.Concat(GetShiftCodes(), GetCtrlCodes(), GetAltCodes())
}
