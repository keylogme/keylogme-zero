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

type DevicePathFinder func(input types.KeyloggerInput) []string

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

// findKeyboardDevicesById by going through each device registered on OS
func findKeyboardDevicesById(input types.KeyloggerInput) []string {
	pathProductId := "/sys/class/input/event%d/device/id/product"
	pathVendorId := "/sys/class/input/event%d/device/id/vendor"
	resolved := "/dev/input/event%d"

	listDevicesPaths := []string{}
	productIdToCompare := fmt.Sprintf("%s\n", input.ProductId)
	vendorIdToCompare := fmt.Sprintf("%s\n", input.VendorId)
	for i := 0; i < 255; i++ {
		buffProductId, err := os.ReadFile(fmt.Sprintf(pathProductId, i))
		if err != nil {
			continue
		}
		productId := string(buffProductId)

		buffVendorId, err := os.ReadFile(fmt.Sprintf(pathVendorId, i))
		if err != nil {
			continue
		}
		vendorId := string(buffVendorId)

		if productId == productIdToCompare && vendorId == vendorIdToCompare {
			listDevicesPaths = append(listDevicesPaths, fmt.Sprintf(resolved, i))
		}
	}
	return listDevicesPaths
}

var PathFinder DevicePathFinder = findKeyboardDevicesById

func getKeyLogger(input types.KeyloggerInput) (*KeyLogger, error) {
	pathsDevice := PathFinder(input)
	if len(pathsDevice) == 0 {
		return nil, fmt.Errorf(
			"Device vendor id %s and product id %s not found\n",
			input.VendorId,
			input.ProductId,
		)
	}
	k := &KeyLogger{}
	slog.Debug(fmt.Sprintf("Opening %s\n", pathsDevice))
	for _, pathDevice := range pathsDevice {
		fd, err := openDeviceFile(pathDevice)
		if err != nil {
			return nil, err
		}
		k.FD = append(k.FD, fd)
	}
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
	FD       []*os.File
	isClosed bool
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
	return getKeyLogger(kInput)
}

// Read from file descriptor
// Blocking call, returns channel
// Make sure to close channel when finish
func (k *KeyLogger) Read() chan InputEvent {
	event := make(chan InputEvent)
	for _, fdDevice := range k.FD {
		go func(event chan InputEvent, fd *os.File) {
			for {
				e, err := k.read(fd)
				if err != nil {
					slog.Debug(fmt.Sprintf("error reading from file descriptor: %s\n", err))
					if !k.isClosed {
						close(event)
						k.isClosed = true
					}
					break
				}

				if e != nil {
					event <- *e
				}
			}
		}(event, fdDevice)
	}
	return event
}

// read from file description and parse binary into go struct
func (k *KeyLogger) read(fdDevice *os.File) (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := fdDevice.Read(buffer)
	// bypass EOF, maybe keyboard is connected
	// but you don't press any key
	if err != nil && err != io.EOF {
		slog.Debug(
			fmt.Sprintf(
				"error reading from file descriptor %s: %s\n",
				fdDevice.Name(),
				err.Error(),
			),
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
		slog.Debug(fmt.Sprintf("error parsing buffer: %s\n", err.Error()))
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
	for _, fd := range k.FD {
		if fd == nil {
			continue
		}
		fd.Close()
	}
	return nil
}

// isRoot checks if the process is run with root permission
func isRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}
