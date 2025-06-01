package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	k0 "github.com/keylogme/keylogme-zero"
	"github.com/keylogme/keylogme-zero/storage"
	"github.com/keylogme/keylogme-zero/utils"
)

// Linux keylogger
// Use lsinput to see the usb_name to be used
// apt install input-utils
// sudo lsinput
// If your keyboard name appeared multiple times,
// try with all of them
// See readme

// Mac keylogger
// See System Information to get Product ID and Vendor ID of keyboard and mouse devices:
// - USB
// - Bluetooth
// - SPI (Apple internal keyboard / trackpad)

func main() {
	// Get config
	file_config := os.Getenv("CONFIG_FILE")
	if file_config == "" {
		log.Fatal("CONFIG_FILE is not set")
	}
	var config k0.KeylogmeZeroConfig
	err := utils.ParseFromFile(file_config, &config)
	if err != nil {
		log.Fatal(err.Error())
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

	security := k0.NewSecurity(k0.SecurityInput{
		BaggageSize:   0,
		Baggage:       map[string][]uint16{},
		GhostingCodes: []uint16{},
	})

	sd := k0.MustGetNewShortcutsDetector(config.Keylog.ShortcutGroups)

	ss := k0.NewShiftStateDetector(config.Keylog.ShiftState)

	ld := k0.NewLayersDetector(config.Keylog.Devices, config.Keylog.ShiftState)

	k0.Start(chEvt, &devices, security, sd, ss, ld, ffs)

	// Graceful shutdown
	ctxInt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctxInt.Done()
	slog.Info("Shutting down, graceful wait...")
	cancel()
	// INFO: close channel so start loop ends
	close(chEvt)
	// FIXME: instead of graceful wait, use wg.Wait() to wait for all goroutines to finish
	time.Sleep(3 * time.Second) // graceful wait
	slog.Info("Logger closed.")
}
