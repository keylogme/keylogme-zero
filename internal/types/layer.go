package types

type Layer struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Codes []Key  `json:"codes"`
}

type LayerDetected struct {
	Id       string
	DeviceId string
	LayerId  int64
}

func (ld *LayerDetected) IsDetected() bool {
	return ld.LayerId != 0 && ld.DeviceId != ""
}
