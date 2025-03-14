package keylog

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"
)

// findKeyboardDevice by going through each device registered on OS
// Mostly it will contain keyword - keyboard
// Returns the file path which contains events
func findKeyboardDevice(name string) string {
	path := "/sys/class/input/event%d/device/name"
	resolved := "/dev/input/event%d"

	nameToCompare := fmt.Sprintf("%s\n", name)
	for i := 0; i < 255; i++ {
		buff, err := os.ReadFile(fmt.Sprintf(path, i))
		if err != nil {
			continue
		}

		deviceName := string(buff)
		// fmt.Printf("%#v\n", deviceName)
		if deviceName == nameToCompare {
			return fmt.Sprintf(resolved, i)
		}
	}

	return ""
}

func getKeyLogger(name string) (*keyLogger, error) {
	pathDevice := findKeyboardDevice(name)
	if pathDevice == "" {
		return nil, fmt.Errorf("Device with name %s not found\n", name)
	}
	k, err := newKeylogger(pathDevice)
	if err != nil {
		return nil, fmt.Errorf("Could not set keylogger for %s. %s\n", name, err.Error())
	}
	return k, nil
}

type Device struct {
	DeviceInput
	Connected bool
	keylogger *keyLogger
	sendInput chan DeviceEvent
}

type DeviceInput struct {
	DeviceId string `json:"device_id"`
	Name     string `json:"name"`
	UsbName  string `json:"usb_name"`
}

type DeviceEvent struct {
	inputEvent
	DeviceId string
	ExecTime time.Time
}

func GetDevice(ctx context.Context, input DeviceInput, inputChan chan DeviceEvent) *Device {
	device := &Device{DeviceInput: input, Connected: true, keylogger: nil, sendInput: inputChan}
	go device.handleReconnects(ctx, device.start)
	return device
}

func mustGetDevice(ctx context.Context, input DeviceInput, inputChan chan DeviceEvent) *Device {
	k, err := getKeyLogger(input.Name)
	if err != nil {
		log.Fatal(err.Error())
	}
	device := &Device{DeviceInput: input, Connected: true, keylogger: k, sendInput: inputChan}
	go device.handleReconnects(ctx, device.start)
	return device
}

func (d *Device) start(ctx context.Context) bool {
	if d.keylogger == nil {
		return false
	}
	keylogChan := d.keylogger.Read()
	for {
		select {
		case <-ctx.Done():
			return true
		case i, ok := <-keylogChan:
			if !ok {
				slog.Info("exited channel keylogger")
				return false
			}
			if !i.IsValid() {
				continue
			}
			// Get current time with microsecond precision
			now := time.Now()

			// Get Unix timestamp with nanoseconds and format with microseconds precision
			// fmt.Printf(
			// 	"Current time of %d %d (microsecond precision): %s\n",
			// 	i.Code,
			// 	i.Value,
			// 	now.Format("2006-01-02 15:04:05.000000"),
			// )

			de := DeviceEvent{inputEvent: i, DeviceId: d.DeviceId, ExecTime: now}
			d.sendInput <- de
		}
	}
}

func (d *Device) handleReconnects(ctx context.Context, s func(context.Context) bool) {
	if d.keylogger != nil {
		// blocking call to start reading keylogger
		d.Connected = true
		slog.Info(fmt.Sprintf("Device %s connected\n", d.Name))
		shutdown := s(ctx)
		if shutdown {
			slog.Info(fmt.Sprintf("Device %s closed\n", d.Name))
			d.keylogger.Close()
			return
		}
		d.Connected = false
		slog.Info(fmt.Sprintf("Device %s disconnected, reconnecting...\n", d.Name))
		time.Sleep(1 * time.Second)
		d.keylogger.Close()
	}
	newK, err := getKeyLogger(d.UsbName)
	if err != nil {
		slog.Info(fmt.Sprintf("Device %s not connected to computer, waiting ...\n", d.Name))
		time.Sleep(5 * time.Second)
	}
	d.keylogger = newK // assign to nil if device not found
	d.handleReconnects(ctx, s)
}
