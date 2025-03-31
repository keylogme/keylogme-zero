package keylog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	reconnect_wait = 300 * time.Millisecond
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

func getPathDevice(name string) string {
	// INFO: for testing purposes, we can pass the abs path of a test file.
	// So we check if the file exists, if not we try to find the device
	_, err := os.Open(name)
	if os.IsNotExist(err) {
		slog.Debug(fmt.Sprintf("file %s does not exist", name))
		return findKeyboardDevice(name)
	}
	return name
}

func getKeyLogger(name string) (*keyLogger, error) {
	pathDevice := getPathDevice(name)
	if pathDevice == "" {
		return nil, fmt.Errorf("Device with name %s not found\n", name)
	}
	slog.Debug(fmt.Sprintf("Opening %s\n", pathDevice))
	k, err := newKeylogger(pathDevice)
	if err != nil {
		return nil, fmt.Errorf("Could not set keylogger for %s. %s\n", name, err.Error())
	}
	return k, nil
}

type Device struct {
	DeviceInput
	ctx       context.Context
	keylogger *keyLogger
	sendInput chan DeviceEvent
}

type DeviceInput struct {
	DeviceId string  `json:"device_id"`
	Name     string  `json:"name"`
	UsbName  string  `json:"usb_name"`
	Layers   []Layer `json:"layers"`
}

type DeviceEvent struct {
	inputEvent
	DeviceId string
	ExecTime time.Time
}

func GetDevice(ctx context.Context, input DeviceInput, inputChan chan DeviceEvent) *Device {
	device := &Device{ctx: ctx, DeviceInput: input, keylogger: nil, sendInput: inputChan}
	go device.handleReconnects()
	return device
}

func (d *Device) start() bool {
	defer d.Close()
	slog.Info(fmt.Sprintf("🚀 Starting device %s \n", d.Name))
	if d.keylogger == nil {
		return false
	}
	keylogChan := d.keylogger.Read()
	for {
		select {
		case <-d.ctx.Done():
			return true
		case i, ok := <-keylogChan:
			if !ok {
				slog.Info("exited channel keylogger")
				return false
			}
			if !i.IsValid() {
				slog.Debug(fmt.Sprintf("Invalid input event %+v\n", i))
				continue
			}
			// Get current time with microsecond precision
			now := time.Now()

			// Get Unix timestamp with nanoseconds and format with microseconds precision
			slog.Debug(fmt.Sprintf(
				"Current time of %d %d (microsecond precision): %s\n",
				i.Code,
				i.Value,
				now.Format("2006-01-02 15:04:05.000000"),
			))

			de := DeviceEvent{inputEvent: i, DeviceId: d.DeviceId, ExecTime: now}
			d.sendInput <- de
		}
	}
}

func (d *Device) IsConnected() bool {
	return d.keylogger != nil
}

func (d *Device) handleReconnects() {
	for {
		slog.Debug(fmt.Sprintf("Reconnecting device %s\n", d.Name))
		newK, err := getKeyLogger(d.UsbName)
		if err != nil {
			slog.Debug(fmt.Sprintf("error getting keylogger : %s\n", err.Error()))
			select {
			case <-time.After(reconnect_wait):
				continue
			case <-d.ctx.Done():
				slog.Info(fmt.Sprintf("💤 Shutting down device %s\n", d.Name))
				return
			}
		}
		// newK exists
		d.keylogger = newK // connected
		shutdown := d.start()
		if shutdown {
			slog.Info(fmt.Sprintf("💤 Shutting down device %s\n", d.Name))
			return
		}
	}
}

func (d *Device) Close() {
	if d.keylogger != nil {
		d.keylogger.Close()
	}
	d.keylogger = nil
}
