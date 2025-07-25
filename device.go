package k0

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/keylogme/keylogme-zero/internal/keylogger"
	"github.com/keylogme/keylogme-zero/types"
)

const (
	reconnect_wait = 300 * time.Millisecond
)

type Device struct {
	DeviceInput
	ctx       context.Context
	keylogger *keylogger.KeyLogger
	sendInput chan DeviceEvent
}

type DeviceInput struct {
	DeviceId       string               `json:"device_id"`
	Name           string               `json:"name"`
	Layers         []LayerInput         `json:"layers"`
	KeyloggerInput types.KeyloggerInput `json:"keylogger"`
}

type DeviceEvent struct {
	keylogger.InputEvent
	DeviceId   string
	DeviceName string
}

func GetDevice(
	ctx context.Context,
	input DeviceInput,
	inputChan chan DeviceEvent,
) *Device {
	device := &Device{ctx: ctx, DeviceInput: input, keylogger: nil, sendInput: inputChan}
	go device.handleReconnects()
	return device
}

func (d *Device) start() bool {
	defer d.Close()
	slog.Info(
		fmt.Sprintf(
			"🚀 Starting device %s | Vendor ID %s | Product ID %s\n",
			d.Name,
			d.KeyloggerInput.VendorId,
			d.KeyloggerInput.ProductId,
		),
	)
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
				slog.Info(fmt.Sprintf("🔚 Disconnected device %s \n", d.Name))
				return false
			}
			if !i.IsValid() {
				slog.Debug(fmt.Sprintf("Invalid input event %+v\n", i))
				continue
			}
			// Get Unix timestamp with nanoseconds and format with microseconds precision
			slog.Debug(fmt.Sprintf(
				"Current time of %d %d (microsecond precision): %s\n",
				i.Code,
				i.Type,
				i.Time.Format("2006-01-02 15:04:05.000000"),
			))

			de := DeviceEvent{InputEvent: i, DeviceId: d.DeviceId, DeviceName: d.Name}
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
		newK, err := keylogger.NewKeylogger(d.KeyloggerInput)
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
		_ = d.keylogger.Close()
	}
	d.keylogger = nil
}
