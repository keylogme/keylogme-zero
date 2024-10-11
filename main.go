package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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
			{Id: 1, Name: "foostan Corne"},
			{Id: 2, Name: "MOSART Semi. 2.4G INPUT DEVICE Mouse"},
			{Id: 2, Name: "Logitech MX Master 2S"},
			// {Id: 2, Name: "Wacom Intuos BT M Pen"},
		},
		Shortcuts: []keylog.Shortcut{
			{Id: 1, Values: []string{"J", "S"}, Type: keylog.SequentialShortcutType},
			{Id: 2, Values: []string{"J", "F"}, Type: keylog.SequentialShortcutType},
			{Id: 3, Values: []string{"J", "G"}, Type: keylog.SequentialShortcutType},
			{Id: 4, Values: []string{"J", "S", "G"}, Type: keylog.SequentialShortcutType},
		},
	}
	ctx, cancelCtx := context.WithCancel(context.Background())
	// INFO: two different types of cleanup
	// for storage, the ctx will close
	// for keylog, a cleanup function is returned
	ffs := storage.NewFileStorage(ctx, "Oct12.json")
	_, cleanup := keylog.Start(ffs, config)
	// defer cleanup()

	// fmt.Println(ds)

	// Graceful shutdown
	ctxInt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctxInt.Done()
	cancelCtx()
	cleanup()

	fmt.Println("Logger closed.")
}
