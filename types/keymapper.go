package types

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

var KeyCodeMapLinux = map[uint16]string{
	1:   KEY_ESCAPE,
	2:   KEY_1,
	3:   KEY_2,
	4:   KEY_3,
	5:   KEY_4,
	6:   KEY_5,
	7:   KEY_6,
	8:   KEY_7,
	9:   KEY_8,
	10:  KEY_9,
	11:  KEY_0,
	12:  KEY_HYPHEN,
	13:  KEY_EQUALS,
	14:  KEY_BACKSPACE,
	15:  KEY_TAB,
	16:  KEY_Q,
	17:  KEY_W,
	18:  KEY_E,
	19:  KEY_R,
	20:  KEY_T,
	21:  KEY_Y,
	22:  KEY_U,
	23:  KEY_I,
	24:  KEY_O,
	25:  KEY_P,
	26:  KEY_LEFT_BRACKET,
	27:  KEY_RIGHT_BRACKET,
	28:  KEY_ENTER,
	29:  KEY_LEFT_CTRL,
	30:  KEY_A,
	31:  KEY_S,
	32:  KEY_D,
	33:  KEY_F,
	34:  KEY_G,
	35:  KEY_H,
	36:  KEY_J,
	37:  KEY_K,
	38:  KEY_L,
	39:  KEY_SEMICOLON,
	40:  KEY_SINGLE_QUOTE,
	41:  KEY_BACKTICK,
	42:  KEY_LEFT_SHIFT,
	43:  KEY_BACKSLASH,
	44:  KEY_Z,
	45:  KEY_X,
	46:  KEY_C,
	47:  KEY_V,
	48:  KEY_B,
	49:  KEY_N,
	50:  KEY_M,
	51:  KEY_COMMA,
	52:  KEY_PERIOD,
	53:  KEY_SLASH,
	54:  KEY_RIGHT_SHIFT,
	55:  KEY_KEYPAD_ASTERISK,
	56:  KEY_LEFT_ALT,
	57:  KEY_SPACE,
	58:  KEY_CAPS_LOCK,
	59:  KEY_F1,
	60:  KEY_F2,
	61:  KEY_F3,
	62:  KEY_F4,
	63:  KEY_F5,
	64:  KEY_F6,
	65:  KEY_F7,
	66:  KEY_F8,
	67:  KEY_F9,
	68:  KEY_F10,
	69:  KEY_NUM_LOCK,
	70:  KEY_SCROLL_LOCK,
	71:  KEY_KEYPAD_7,
	72:  KEY_KEYPAD_8,
	73:  KEY_KEYPAD_9,
	74:  KEY_KEYPAD_MINUS,
	75:  KEY_KEYPAD_4,
	76:  KEY_KEYPAD_5,
	77:  KEY_KEYPAD_6,
	78:  KEY_KEYPAD_PLUS,
	79:  KEY_KEYPAD_1,
	80:  KEY_KEYPAD_2,
	81:  KEY_KEYPAD_3,
	82:  KEY_KEYPAD_0,
	83:  KEY_KEYPAD_PERIOD,
	87:  KEY_F11,
	88:  KEY_F12,
	96:  KEY_KEYPAD_ENTER,
	97:  KEY_RIGHT_CTRL,
	98:  KEY_KEYPAD_SLASH,
	100: KEY_RIGHT_ALT,
	102: KEY_HOME,
	103: KEY_UP_ARROW,
	104: KEY_PAGE_UP,
	105: KEY_LEFT_ARROW,
	106: KEY_RIGHT_ARROW,
	107: KEY_END,
	108: KEY_DOWN_ARROW,
	109: KEY_PAGE_DOWN,
	110: KEY_INSERT,
	111: KEY_DELETE,
	114: KEY_VOLUME_DOWN,
	115: KEY_VOLUME_UP,
	117: KEY_KEYPAD_EQUALS,
	125: KEY_LEFT_GUI,
	126: KEY_RIGHT_GUI,
}

// Source: Universal Serial Bus HID Usage Tables version 1.2 , 10 Keyboard/Keypad page (0x07)
var KeyCodeMapDarwin = map[uint16]string{
	4:   KEY_A,
	5:   KEY_B,
	6:   KEY_C,
	7:   KEY_D,
	8:   KEY_E,
	9:   KEY_F,
	10:  KEY_G,
	11:  KEY_H,
	12:  KEY_I,
	13:  KEY_J,
	14:  KEY_K,
	15:  KEY_L,
	16:  KEY_M,
	17:  KEY_N,
	18:  KEY_O,
	19:  KEY_P,
	20:  KEY_Q,
	21:  KEY_R,
	22:  KEY_S,
	23:  KEY_T,
	24:  KEY_U,
	25:  KEY_V,
	26:  KEY_W,
	27:  KEY_X,
	28:  KEY_Y,
	29:  KEY_Z,
	30:  KEY_1,
	31:  KEY_2,
	32:  KEY_3,
	33:  KEY_4,
	34:  KEY_5,
	35:  KEY_6,
	36:  KEY_7,
	37:  KEY_8,
	38:  KEY_9,
	39:  KEY_0,
	40:  KEY_ENTER,
	41:  KEY_ESCAPE,
	42:  KEY_BACKSPACE,
	43:  KEY_TAB,
	44:  KEY_SPACE,
	45:  KEY_HYPHEN,
	46:  KEY_EQUALS,
	47:  KEY_LEFT_BRACKET,
	48:  KEY_RIGHT_BRACKET,
	49:  KEY_BACKSLASH,
	51:  KEY_SEMICOLON,
	52:  KEY_SINGLE_QUOTE,
	53:  KEY_BACKTICK,
	54:  KEY_COMMA,
	55:  KEY_PERIOD,
	56:  KEY_SLASH,
	57:  KEY_CAPS_LOCK,
	58:  KEY_F1,
	59:  KEY_F2,
	60:  KEY_F3,
	61:  KEY_F4,
	62:  KEY_F5,
	63:  KEY_F6,
	64:  KEY_F7,
	65:  KEY_F8,
	66:  KEY_F9,
	67:  KEY_F10,
	68:  KEY_F11,
	69:  KEY_F12,
	71:  KEY_SCROLL_LOCK,
	73:  KEY_INSERT,
	74:  KEY_HOME,
	75:  KEY_PAGE_UP,
	76:  KEY_DELETE,
	77:  KEY_END,
	78:  KEY_PAGE_DOWN,
	79:  KEY_RIGHT_ARROW,
	80:  KEY_LEFT_ARROW,
	81:  KEY_DOWN_ARROW,
	82:  KEY_UP_ARROW,
	83:  KEY_NUM_LOCK,
	84:  KEY_KEYPAD_SLASH,
	85:  KEY_KEYPAD_ASTERISK,
	86:  KEY_KEYPAD_MINUS,
	87:  KEY_KEYPAD_PLUS,
	88:  KEY_KEYPAD_ENTER,
	89:  KEY_KEYPAD_1,
	90:  KEY_KEYPAD_2,
	91:  KEY_KEYPAD_3,
	92:  KEY_KEYPAD_4,
	93:  KEY_KEYPAD_5,
	94:  KEY_KEYPAD_6,
	95:  KEY_KEYPAD_7,
	96:  KEY_KEYPAD_8,
	97:  KEY_KEYPAD_9,
	98:  KEY_KEYPAD_0,
	99:  KEY_KEYPAD_PERIOD,
	103: KEY_KEYPAD_EQUALS,
	128: KEY_VOLUME_UP,
	129: KEY_VOLUME_DOWN,
	224: KEY_LEFT_CTRL,
	225: KEY_LEFT_SHIFT,
	226: KEY_LEFT_ALT,
	227: KEY_LEFT_GUI,
	228: KEY_RIGHT_CTRL,
	229: KEY_RIGHT_SHIFT,
	230: KEY_RIGHT_ALT,
	231: KEY_RIGHT_GUI,
}

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

func InverseCodeMap(in map[uint16]string) map[string]uint16 {
	out := map[string]uint16{}
	for key, val := range in {
		out[val] = key
	}
	return out
}

var inversedMap = InverseCodeMap(KeyCodeMap)

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
