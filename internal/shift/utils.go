package shift

import (
	"fmt"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/internal/types"
)

func getShiftCodeKey(shiftCode, code uint16) string {
	return fmt.Sprintf("%d_%d", shiftCode, code)
}

func getShortcutCodesForShiftState() []types.ShortcutCodes {
	listSS := []types.ShortcutCodes{}
	for _, sc := range keylogger.SHIFT_CODES {
		for _, c := range keylogger.ALL_CODES {
			scKey := getShiftCodeKey(sc, c)
			ssc := types.ShortcutCodes{
				Id:    scKey,
				Codes: []uint16{sc, c},
				Type:  types.HoldShortcutType,
			}
			listSS = append(listSS, ssc)
		}
	}
	return listSS
}

func getMapIdToCodes() map[string]types.Key {
	mapIdToCodes := make(map[string]types.Key)
	for _, sc := range keylogger.SHIFT_CODES {
		for _, c := range keylogger.ALL_CODES {
			scKey := getShiftCodeKey(sc, c)
			mapIdToCodes[scKey] = types.Key{Code: c, Modifier: sc}
		}
	}
	return mapIdToCodes
}
