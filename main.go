package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog"
	"github.com/keylogme/zero-trust-logger/keylog/storage"
	"github.com/keylogme/zero-trust-logger/keylog/utils"
)

// Use lsinput to see which input to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them

func main() {
	// Get config
	// config := keylog.Config{
	// 	Devices: []keylog.DeviceInput{
	// 		{DeviceId: "1", Name: "crkbd", UsbName: "foostan Corne"},
	// 		{DeviceId: "2", Name: "my mouse", UsbName: "MOSART Semi. 2.4G INPUT DEVICE Mouse"},
	// 		{DeviceId: "2", Name: "mouse at work", UsbName: "Logitech MX Master 2S"},
	// 		{DeviceId: "3", Name: "lenovo", UsbName: "LiteOn Lenovo Traditional USB Keyboard"},
	// 		// {Id: 2, Name: "Wacom Intuos BT M Pen"},
	// 	},
	// 	Shortcuts: []keylog.ShortcutCodes{
	// 		{Id: 1, Codes: []uint16{36, 31}, Type: keylog.SequentialShortcutType},
	// 	},
	// }

	// configStorage := storage.ConfigStorage{
	// 	FileOutput:        "Dec21.json",
	// 	PeriodicSaveInSec: 10,
	// }
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
	file_config := os.Getenv("CONFIG_FILE")
	if file_config == "" {
		log.Fatal("CONFIG_FILE is not set")
	}
	var config keylog.KeylogmeZeroConfig
	err := utils.ParseFromFile(file_config, &config)
	if err != nil {
		log.Fatal("Could not parse config file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ffs := storage.MustGetNewFileStorage(ctx, config.Storage)

	chEvt := make(chan keylog.DeviceEvent)
	devices := []keylog.Device{}
	for _, dev := range config.Keylog.Devices {
		d := keylog.GetDevice(ctx, dev, chEvt)
		devices = append(devices, *d)
	}

	sd := keylog.NewShortcutsDetector(config.Keylog.Shortcuts)
	keylog.Start(chEvt, &devices, sd, ffs)

	// Graceful shutdown
	ctxInt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctxInt.Done()
	slog.Info("Shutting down, graceful wait...")
	cancel()
	time.Sleep(3 * time.Second) // graceful wait
	slog.Info("Logger closed.")
}
