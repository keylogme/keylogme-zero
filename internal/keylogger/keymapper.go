package keylogger

import "slices"

func GetAllCodes() []uint16 {
	return slices.Concat(numCodes, letterCodes, symbolCodes)
}

func GetShiftCodes() []uint16 {
	return shiftCodes
}

func GetCtrlCodes() []uint16 {
	return ctrlCodes
}

func GetAltCodes() []uint16 {
	return altCodes
}

func GetAllModifierCodes() []uint16 {
	return slices.Concat(shiftCodes, ctrlCodes, altCodes)
}
