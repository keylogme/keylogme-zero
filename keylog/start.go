package keylog

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog/storage"
)

func Start(store storage.Storage, config Config) ([]*device, func()) {
	chEvt := make(chan deviceEvent)

	sd := newShortcutsDetector(config.Shortcuts)

	devices := []*device{}
	for _, dev := range config.Devices {
		d := getDevice(dev, chEvt)
		devices = append(devices, d)
	}

	modifiers := []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt

	slog.Info("Listening...")

	modPress := []uint16{}
	go func() {
		for i := range chEvt {
			if i.KeyPress() && slices.Contains(modifiers, i.Code) {
				modPress = append(modPress, i.Code)
			}
			if i.Type == evKey && i.KeyRelease() {
				start := time.Now()

				detectedShortcut := sd.Detect(i.DeviceId, i.KeyString())
				if detectedShortcut.ShortcutId != 0 {
					// sendShortcut(sender, i.DeviceId, detectedShortcutID)
					slog.Info(
						fmt.Sprintf(
							"Shortcut %d found in device %s\n",
							detectedShortcut.ShortcutId,
							detectedShortcut.DeviceId,
						),
					)
					store.SaveShortcut(detectedShortcut.DeviceId, detectedShortcut.ShortcutId)
				}
				//
				// FIXME: mod+key is sent, but when mod is released , is sent again
				// keylogs := []uint16{i.Code}
				// keylogs = append(keylogs, modPress...)
				// err := sendKeylog(sender, i.DeviceId, i.Code)
				err := store.SaveKeylog(i.DeviceId, i.Code)
				if err != nil {
					fmt.Printf("error %s\n", err.Error())
				}
				slog.Info(
					fmt.Sprintf("| %s | Key :%d %s\n", time.Since(start), i.Code, i.KeyString()),
				)
				// Reset modPress
				modPress = []uint16{}
			}
		}
	}()

	return devices, func() {
		// close channels
		close(chEvt)
		for _, d := range devices {
			if d.keylogger != nil {
				d.keylogger.Close()
			}
			fmt.Printf("Device %s closed\n", d.Name)
		}
	}
}
