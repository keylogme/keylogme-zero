package keylogger

import (
	"fmt"
	"log"
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

func getKeyLogger(name string) (*KeyLogger, error) {
	pathDevice := findKeyboardDevice(name)
	if pathDevice == "" {
		return nil, fmt.Errorf("Device with name %s not found\n", name)
	}
	k, err := NewKeylogger(pathDevice)
	if err != nil {
		return nil, fmt.Errorf("Could not set keylogger for %s. %s\n", name, err.Error())
	}
	return k, nil
}

type Device struct {
	DeviceInput
	Connected bool
	keylogger *KeyLogger
	sendInput chan DeviceEvent
}

type DeviceInput struct {
	Id   int64
	Name string
}

type DeviceEvent struct {
	InputEvent
	DeviceId int64
}

func MustGetDevice(input DeviceInput, inputChan chan DeviceEvent) *Device {
	k, err := getKeyLogger(input.Name)
	if err != nil {
		log.Fatal(err.Error())
	}
	device := &Device{DeviceInput: input, Connected: true, keylogger: k, sendInput: inputChan}
	go device.handleReconnects(device.start)
	return device
}

func (d *Device) start() {
	if d.keylogger == nil {
		return
	}
	for i := range d.keylogger.Read() {
		de := DeviceEvent{InputEvent: i, DeviceId: d.Id}
		d.sendInput <- de
	}
}

func (d *Device) handleReconnects(s func()) {
	if d.keylogger != nil {
		// blocking call to start reading keylogger
		d.Connected = true
		s()
		d.Connected = false
		fmt.Printf("Device %s disconnected, reconnecting...\n", d.Name)
		time.Sleep(1 * time.Second)
		d.keylogger.Close()
	}
	newK, err := getKeyLogger(d.Name)
	if err != nil {
		fmt.Printf("Device %s not connected to computer, waiting ...\n", d.Name)
		time.Sleep(5 * time.Second)
	}
	d.keylogger = newK // assign to nil if device not found
	d.handleReconnects(s)
}

// func (dm *DeviceManager)  {
//
// }
