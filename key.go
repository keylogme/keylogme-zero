package k0

// A key pressed can be a single key or a shifted key.
// for example, "=" can be pressed as "=", or as "Shift + ="
// In QMK or ZMK firmware, you can assign a physical key to act as a single or shifted key.
type Key struct {
	Code     uint16 `json:"code"`
	Modifier uint16 `json:"modifier,omitempty"`
}
