package keylog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/keylogme/keylogme-zero/v1/keylog/storage"
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
				store.SaveShortcut(sd.DeviceId, sd.ShortcutId)
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
				store.SaveShiftState(ssd.DeviceId, ssd.Modifier, ssd.Code, ssd.Auto)

			}
			ldd := ld.isLayerChangeDetected(i)
			if ldd.IsDetected() {
				slog.Info(
					fmt.Sprintf("Layer %d detected in device %s\n", ldd.LayerId, i.DeviceId),
				)
			}
			if i.Type == evKey && i.KeyRelease() {
				start := time.Now()
				err := store.SaveKeylog(i.DeviceId, i.Code)
				if err != nil {
					fmt.Printf("error %s\n", err.Error())
				}
				slog.Info(
					fmt.Sprintf(
						"| %s | Key :%d %s\n",
						time.Since(start),
						i.Code,
						i.KeyString(),
					),
				)
			}
		}
	}()
}
