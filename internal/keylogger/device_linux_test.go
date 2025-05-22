package keylogger

import (
	"context"
	"log/slog"
	"runtime"
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/utils"
)

func TestDisconnectionDevice(t *testing.T) {
	before := runtime.NumGoroutine()
	defer utils.CheckGoroutineLeak(t, before)

	df, err := initDeviceFile()
	if err != nil {
		t.Fatal(err)
	}
	defer df.Close()
	filepath := df.Name()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	intputDevice := DeviceInput{
		DeviceId:       "device1",
		Name:           "device1",
		KeyloggerInput: KeyloggerInput{UsbName: filepath},
	}

	chEvt := make(chan DeviceEvent, 10)
	d := GetDevice(ctx, intputDevice, chEvt)
	defer d.Close()

	time.Sleep(50 * time.Millisecond)
	// press keys
	err = writeKeyDeviceFile(df, uint16(1))
	if err != nil {
		t.Fatalf("error writing: %s\n", err.Error())
	}
	time.Sleep(50 * time.Millisecond)
	// one for keypress and other for keyrelease
	t.Log(len(chEvt))
	if len(chEvt) != 2 {
		t.Fatal("expected 2 device events")
	}
	// disconnect device
	slog.Info("Disconnecting device")
	kl, _ := d.keylogger.(*keyLogger) // cast interface to struct
	err = disconnectDeviceFile(kl.fd)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(reconnect_wait * 2)
	if d.IsConnected() {
		t.Fatal("device should not be connected")
	}
	// reconnect device
	slog.Info("Reconnecting device")
	err = reconnectDeviceFile(df)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	if !d.IsConnected() {
		t.Fatal("device should be connected")
	}
	// press keys
	err = writeKeyDeviceFile(df, uint16(1))
	if err != nil {
		t.Fatalf("error writing: %s\n", err.Error())
	}
	time.Sleep(50 * time.Millisecond)
	// one for keypress and other for keyrelease
	t.Log(len(chEvt))
	if len(chEvt) != 4 {
		t.Fatal("expected 4 device events")
	}
}
