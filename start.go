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
						"Shift state of %d found in device %s - auto %t | diff time : press %d μs release %d μs\n",
						ssd.Code,
						ssd.DeviceId,
						ssd.Auto,
						ssd.DiffTimePressMicro,
						ssd.DiffTimeReleaseMicro,
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
					fmt.Sprintf(
						"Layer %d detected in device %s during key %d (Release? %t)\n",
						ldd.LayerId,
						ldd.DeviceId,
						i.Code,
						i.KeyRelease(),
					),
				)
				err := store.SaveLayerChange(ldd.DeviceId, ldd.LayerId)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing layer change : %s\n", err.Error()))
				}
			}
			if ss.blockSaveKeylog() || (ssd.IsDetected() && ssd.Auto) {
				slog.Debug(
					fmt.Sprintf("Blocked keylog save | %t %t\n", ss.blockSaveKeylog(), ssd.Auto),
				)
				continue
			}
			if i.Type == evKey && i.KeyRelease() {
				slog.Info(
					fmt.Sprintf(
						"Key :%d %s in layer %d\n",
						i.Code,
						i.KeyString(),
						ld.GetCurrentLayerId(),
					),
				)
				err := store.SaveKeylog(i.DeviceId, ld.GetCurrentLayerId(), i.Code)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing keylog : %s\n", err.Error()))
				}
			}
		}
	}()
}
