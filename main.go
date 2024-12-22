package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog"
	"github.com/keylogme/zero-trust-logger/keylog/storage"
)

type KeyLog struct {
	Code uint16 `json:"code"`
}

// Use lsinput to see which input to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them

func main() {
	// Get config
	config := keylog.Config{
		Devices: []keylog.DeviceInput{
			{DeviceId: "1", Name: "foostan Corne", UsbName: "foostan Corne"},
			{DeviceId: "2", Name: "my mouse", UsbName: "MOSART Semi. 2.4G INPUT DEVICE Mouse"},
			{DeviceId: "2", Name: "mouse at work", UsbName: "Logitech MX Master 2S"},
			// {Id: 2, Name: "Wacom Intuos BT M Pen"},
		},
		Shortcuts: []keylog.ShortcutCodes{
			{Id: 1, Codes: []uint16{36, 31}, Type: keylog.SequentialShortcutType},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	ffs := storage.NewFileStorage(ctx, "Dec21.json")

	chEvt := make(chan keylog.DeviceEvent)
	devices := []keylog.Device{}
	for _, dev := range config.Devices {
		d := keylog.GetDevice(ctx, dev, chEvt)
		devices = append(devices, *d)
	}

	sd := keylog.NewShortcutsDetector(config.Shortcuts)
	keylog.Start(chEvt, &devices, sd, ffs)

	// Graceful shutdown
	ctxInt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctxInt.Done()

	cancel()
	time.Sleep(3 * time.Second) // graceful wait
	slog.Info("Logger closed.")
}
