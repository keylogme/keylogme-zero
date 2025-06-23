package k0

import (
	"fmt"
	"log/slog"

	"github.com/keylogme/keylogme-zero/types"
)

func Start(
	chEvt chan DeviceEvent,
	security *security,
	sd *shortcutsDetector,
	ss *shiftStateDetector,
	ld *layersDetector,
	store types.Storage,
) {
	slog.Info("Listening...")
	go func() {
		for i := range chEvt {

			sd := sd.handleKeyEvent(i)
			if sd.IsDetected() {
				slog.Info(
					fmt.Sprintf(
						"Shortcut %s detected in device %s\n",
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
						"Shift state of %d detected in device %s - auto %t",
						ssd.Code,
						ssd.DeviceId,
						ssd.Auto,
					),
				)
				slog.Debug(
					fmt.Sprintf(
						"Shift state of %d detected in device %s - auto %t | diff time : press %d μs release %d μs\n",
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
			// INFO: print info of shifted code and layer
			if ssd.IsDetected() && ssd.Auto {
				if ld.GetCurrentLayerId() == 0 {
					slog.Info(
						fmt.Sprintf(
							"Key shifted (auto): Mod: %d Code: %d not defined in any layer \n",
							ssd.Modifier,
							ssd.Code,
						),
					)
				} else {
					slog.Info(
						fmt.Sprintf(
							"Key shifted (auto): Mod: %d Code: %d in layer %d\n",
							ssd.Modifier, ssd.Code, ld.GetCurrentLayerId(),
						),
					)
				}
			}
			// INFO: block save keylog if shifted code is possible
			if ss.blockSaveKeylog() || (ssd.IsDetected() && ssd.Auto) {
				slog.Debug(
					fmt.Sprintf("Blocked keylog save | %t %t\n", ss.blockSaveKeylog(), ssd.Auto),
				)
				continue
			}
			if i.KeyRelease() {
				il := DeviceEventInLayer{DeviceEvent: i, LayerId: ld.GetCurrentLayerId()}
				isAuthorized := security.isAuthorized(&il)
				if !isAuthorized {
					continue
				}
				if ld.GetCurrentLayerId() == 0 {
					slog.Info(
						fmt.Sprintf(
							"\tDeviceId: %s \t| Key: %d (%s) \t| Layer: -\n",
							il.DeviceName,
							il.Code,
							il.KeyString(),
						),
					)
				} else {
					slog.Info(
						fmt.Sprintf(
							"\tDeviceId: %s \t| Key :%d (%s) \t| Layer: %d\n",
							il.DeviceName,
							il.Code,
							il.KeyString(),
							il.LayerId,
						),
					)
				}
				err := store.SaveKeylog(il.DeviceId, il.LayerId, il.Code)
				if err != nil {
					slog.Error(fmt.Sprintf("Error storing keylog : %s\n", err.Error()))
				}
			}
		}
	}()
}
