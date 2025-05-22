package keylog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/keylogme/keylogme-zero/utils"
)

type KeyloggerInput struct {
	UsbName string `json:"usb_name"`
}

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
	k := &keyLogger{}
	slog.Debug(fmt.Sprintf("Opening %s\n", pathDevice))
	fd, err := openDeviceFile(pathDevice)
	if err != nil {
		return nil, err
	}
	k.fd = fd
	return k, nil
}

func openDeviceFile(devPath string) (*os.File, error) {
	fd, err := os.OpenFile(devPath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		return nil, wrapErrorRoot(err)
	}
	return fd, nil
}

// keyLogger wrapper around file descriptior
type keyLogger struct {
	fd *os.File
}

func wrapErrorRoot(err error) error {
	if os.IsPermission(err) && !utils.IsRoot() {
		return errors.New(
			"permission denied. run with root permission",
		)
	}
	return err
}

// NewKeylogger creates a new keylogger for a device path
func newKeylogger(kInput KeyloggerInput) (*keyLogger, error) {
	k := &keyLogger{}
	slog.Debug(fmt.Sprintf("creating keylogger with root? %t\n", utils.IsRoot()))
	if _, err := os.Stat(kInput.UsbName); err == nil {
		fd, err := openDeviceFile(kInput.UsbName)
		if err != nil {
			return nil, err
		}
		k.fd = fd
		return k, nil
	}
	return getKeyLogger(kInput.UsbName)
}

// Read from file descriptor
// Blocking call, returns channel
// Make sure to close channel when finish
func (k *keyLogger) Read() chan inputEvent {
	event := make(chan inputEvent)
	go func(event chan inputEvent) {
		for {
			e, err := k.read()
			if err != nil {
				slog.Debug(fmt.Sprintf("error reading from file descriptor: %s\n", err))
				close(event)
				break
			}

			if e != nil {
				event <- *e
			}
		}
	}(event)
	return event
}

// read from file description and parse binary into go struct
func (k *keyLogger) read() (*inputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := k.fd.Read(buffer)
	// bypass EOF, maybe keyboard is connected
	// but you don't press any key
	if err != nil && err != io.EOF {
		slog.Debug(
			fmt.Sprintf("error reading from file descriptor %s: %s\n", k.fd.Name(), err.Error()),
		)
		return nil, err
	}
	// no input, dont send error
	if n <= 0 {
		return nil, nil
	}
	return k.eventFromBuffer(buffer)
}

// eventFromBuffer parser bytes into InputEvent struct
func (k *keyLogger) eventFromBuffer(buffer []byte) (*inputEvent, error) {
	event := &inputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

// Close file descriptor
func (k *keyLogger) Close() error {
	if k.fd == nil {
		return nil
	}
	return k.fd.Close()
}
