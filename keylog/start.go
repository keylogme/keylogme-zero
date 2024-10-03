package keylog

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog/storage"
)

func Start(ctx context.Context, store storage.Storage, config Config) ([]*Device, func()) {
	chEvt := make(chan DeviceEvent)

	sd := NewShortcutsDetector(config.Shortcuts)

	devices := []*Device{}
	for _, dev := range config.Devices {
		d := GetDevice(dev, chEvt)
		devices = append(devices, d)
	}

	modifiers := []uint16{29, 97, 42, 54, 56, 100} // ctrl, shft, alt

	slog.Info("Listening...")

	modPress := []uint16{}
	for i := range chEvt {
		if i.KeyPress() && slices.Contains(modifiers, i.Code) {
			modPress = append(modPress, i.Code)
		}
		if i.Type == evKey && i.KeyRelease() {
			start := time.Now()

			detectedShortcutID := sd.Detect(i.KeyString())
			if detectedShortcutID != 0 {
				// sendShortcut(sender, i.DeviceId, detectedShortcutID)
				store.SaveShortcut(i.DeviceId, detectedShortcutID)
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
			slog.Info(fmt.Sprintf("| %s | Key :%d %s\n", time.Since(start), i.Code, i.KeyString()))
			// Reset modPress
			modPress = []uint16{}
		}
	}
	fmt.Println("Closing...")

	return devices, func() {
		// close channels
		close(chEvt)
		for _, d := range devices {
			d.keylogger.Close()
		}
	}
}
