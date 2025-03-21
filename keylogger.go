package keylog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
)

// keyLogger wrapper around file descriptior
type keyLogger struct {
	fd *os.File
}

// newKeylogger creates a new keylogger for a device path
func newKeylogger(devPath string) (*keyLogger, error) {
	k := &keyLogger{}
	slog.Debug(fmt.Sprintf("creating keylogger with root? %t\n", k.IsRoot()))
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
func (k *keyLogger) IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
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
