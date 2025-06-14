package types

type Storage interface {
	SaveKeylog(deviceId string, layerId int64, keycode uint16) error
	SaveShortcut(deviceId string, shortcutId string) error
	SaveShiftState(deviceId string, modifier uint16, keycode uint16, auto bool) error
	SaveLayerChange(deviceId string, layerId int64) error
}
