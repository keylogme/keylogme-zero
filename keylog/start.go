package keylog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/keylogme/keylogme-zero/keylog/storage"
)

func Start(
	chEvt chan DeviceEvent,
	devices *[]Device,
	sd *shortcutsDetector,
	store storage.Storage,
) {
	slog.Info("Listening...")
	go func() {
		for i := range chEvt {
			sd := sd.handleKeyEvent(i)
			if sd.ShortcutId != "" {
				slog.Info(
					fmt.Sprintf(
						"Shortcut %d found in device %s\n",
						sd.ShortcutId,
						sd.DeviceId,
					),
				)
				store.SaveShortcut(sd.DeviceId, sd.ShortcutId)
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
