package keylogme

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"syscall"
)

// KeyLogger wrapper around file descriptior
type KeyLogger struct {
	fd *os.File
}

// NewKeylogger creates a new keylogger for a device path
func NewKeylogger(devPath string) (*KeyLogger, error) {
	// TODO: input is device name so if keyboard changes  port -> device can be found by name
	k := &KeyLogger{}
	fd, err := os.OpenFile(devPath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		if os.IsPermission(err) && !k.IsRoot() {
			return nil, errors.New(
				"permission denied. run with root permission or use a user with access to " + devPath,
			)
		}
		return nil, err
	}
	k.fd = fd
	return k, nil
}

// IsRoot checks if the process is run with root permission
func (k *KeyLogger) IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
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
	n, err := k.fd.Read(buffer)
	if err != nil {
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
	event := &InputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

// Close file descriptor
func (k *KeyLogger) Close() error {
	if k.fd == nil {
		return nil
	}
	return k.fd.Close()
}
