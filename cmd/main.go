package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	k0 "github.com/keylogme/keylogme-zero/v1"
	"github.com/keylogme/keylogme-zero/v1/storage"
	"github.com/keylogme/keylogme-zero/v1/utils"
)

// Use lsinput to see the usb_name to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them
// See readme

func main() {
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	// Get config
	file_config := os.Getenv("CONFIG_FILE")
	if file_config == "" {
		log.Fatal("CONFIG_FILE is not set")
	}
	var config k0.KeylogmeZeroConfig
	err := utils.ParseFromFile(file_config, &config)
	if err != nil {
		log.Fatal("Could not parse config file")
	}

	// Start logger
	ctx, cancel := context.WithCancel(context.Background())
	ffs := storage.MustGetNewFileStorage(ctx, config.Storage)

	chEvt := make(chan k0.DeviceEvent)
	devices := []k0.Device{}
	for _, dev := range config.Keylog.Devices {
		d := k0.GetDevice(ctx, dev, chEvt)
		devices = append(devices, *d)
	}
	sd := k0.MustGetNewShortcutsDetector(config.Keylog.ShortcutGroups)

	ss := k0.NewShiftStateDetector(config.Keylog.ShiftState)

	ld := k0.NewLayerDetector(config.Keylog.Layers, config.Keylog.ShiftState)

	k0.Start(chEvt, &devices, sd, ss, ld, ffs)

	// Graceful shutdown
	ctxInt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctxInt.Done()
	slog.Info("Shutting down, graceful wait...")
	cancel()
	time.Sleep(3 * time.Second) // graceful wait
	slog.Info("Logger closed.")
}
