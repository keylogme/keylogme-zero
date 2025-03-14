package keylog

import (
	"fmt"
	"log/slog"

	"github.com/keylogme/keylogme-zero/storage"
)

func Start(
	chEvt chan DeviceEvent,
	devices *[]Device,
	sd *shortcutsDetector,
	ss *shiftStateDetector,
	ld *layersDetector,
	store storage.Storage,
) {
	slog.Info("Listening...")
	go func() {
		for i := range chEvt {
			sd := sd.handleKeyEvent(i)
			if sd.IsDetected() {
				slog.Info(
					fmt.Sprintf(
						"Shortcut %s found in device %s\n",
						sd.ShortcutId,
						sd.DeviceId,
					),
				)
				err := store.SaveShortcut(sd.DeviceId, sd.ShortcutId)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing shortcut : %s\n", err.Error()))
				}
			}
			ssd := ss.handleKeyEvent(i)
			if ssd.IsDetected() {
				slog.Info(
					fmt.Sprintf(
						"Shift state of %d found in device %s - auto %t\n",
						ssd.Code,
						ssd.DeviceId,
						ssd.Auto,
					),
				)
				err := store.SaveShiftState(ssd.DeviceId, ssd.Modifier, ssd.Code, ssd.Auto)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing shift state : %s\n", err.Error()))
				}

			}
			ldd := ld.isLayerChangeDetected(i)
			if ldd.IsDetected() {
				slog.Info(
					fmt.Sprintf("Layer %d detected in device %s\n", ldd.LayerId, i.DeviceId),
				)
				err := store.SaveLayerChange(i.DeviceId, ldd.LayerId)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing layer change : %s\n", err.Error()))
				}
			}
			if i.Type == evKey && i.KeyRelease() {
				slog.Info(
					fmt.Sprintf(
						"Key :%d %s\n",
						i.Code,
						i.KeyString(),
					),
				)
				err := store.SaveKeylog(i.DeviceId, i.Code)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing keylog : %s\n", err.Error()))
				}
			}
		}
	}()
}
