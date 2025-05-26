package types

type ShiftStateInput struct {
	ThresholdAuto Duration `json:"threshold_auto"`
}

// Auto is true when the shift state is triggered by the microcontroller
type ShiftStateDetected struct {
	ShortcutId           string
	DeviceId             string
	Modifier             uint16
	Code                 uint16
	Auto                 bool
	DiffTimePressMicro   int64
	DiffTimeReleaseMicro int64
}

func (ssd *ShiftStateDetected) IsDetected() bool {
	return ssd.DeviceId != ""
}
