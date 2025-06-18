package types

import (
	"fmt"
	"strconv"
)

// Hex allows unmarshaling hex string like "0x1A2B" into int
type Hex int

func (h *Hex) UnmarshalJSON(b []byte) error {
	// Remove quotes and parse hex string
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	val, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*h = Hex(val)
	return nil
}

func (h Hex) String() string {
	return fmt.Sprintf("%04x", int(h))
}
