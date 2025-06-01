package keylogger

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/keylogme/keylogme-zero/types"
)

const (
	// evKey is used to describe state changes of keyboards, buttons, or other key-like devices.
	evKey eventType = 0x01
)

// eventType are groupings of codes under a logical input construct.
// Each type has a set of applicable codes to be used in generating events.
// See the ev section for details on valid codes for each type
type eventType uint16

// eventsize is size of structure of Inputevent
var eventsize = int(unsafe.Sizeof(InputEvent{}))

type inputEvent struct {
	Time  syscall.Timeval
	Type  eventType
	Code  uint16
	Value int32
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

func getKeyLogger(name string) (*KeyLogger, error) {
	pathDevice := getPathDevice(name)
	if pathDevice == "" {
		return nil, fmt.Errorf("Device with name %s not found\n", name)
	}
	k := &KeyLogger{}
	slog.Debug(fmt.Sprintf("Opening %s\n", pathDevice))
	fd, err := openDeviceFile(pathDevice)
	if err != nil {
		return nil, err
	}
	k.FD = fd
	return k, nil
}

func openDeviceFile(devPath string) (*os.File, error) {
	fd, err := os.OpenFile(devPath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		return nil, wrapErrorRoot(err)
	}
	return fd, nil
}

// KeyLogger wrapper around file descriptior
type KeyLogger struct {
	FD *os.File
}

func wrapErrorRoot(err error) error {
	if os.IsPermission(err) && !isRoot() {
		return errors.New(
			"permission denied. run with root permission",
		)
	}
	return err
}

// NewKeylogger creates a new keylogger for a device path
func NewKeylogger(kInput types.KeyloggerInput) (*KeyLogger, error) {
	return getKeyLogger(kInput.UsbName)
}

// Read from file descriptor
// Blocking call, returns channel
// Make sure to close channel when finish
func (k *KeyLogger) Read() chan InputEvent {
	event := make(chan InputEvent)
	go func(event chan InputEvent) {
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
func (k *KeyLogger) read() (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := k.FD.Read(buffer)
	// bypass EOF, maybe keyboard is connected
	// but you don't press any key
	if err != nil && err != io.EOF {
		slog.Debug(
			fmt.Sprintf("error reading from file descriptor %s: %s\n", k.FD.Name(), err.Error()),
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
func (k *KeyLogger) eventFromBuffer(buffer []byte) (*InputEvent, error) {
	event := &inputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	if err != nil {
		slog.Debug(fmt.Sprintf("error parsing buffer %s: %s\n", k.FD.Name(), err.Error()))
		return nil, err
	}
	if event.Type != evKey {
		slog.Debug(fmt.Sprintf("event type %d is not evKey\n", event.Type))
		return nil, nil
	}
	return &InputEvent{
		Time: time.Unix(event.Time.Sec, event.Time.Usec*1000),
		Code: event.Code,
		Type: KeyEvent(event.Value),
	}, err
}

// Close file descriptor
func (k *KeyLogger) Close() error {
	if k.FD == nil {
		return nil
	}
	return k.FD.Close()
}

// isRoot checks if the process is run with root permission
func isRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}
