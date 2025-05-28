package keylogger

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func InitDeviceFile() (*os.File, error) {
	tf, err := os.MkdirTemp("", "device_test")
	if err != nil {
		return &os.File{}, err
	}
	filename := fmt.Sprintf("device_%d", rand.Int())
	filepath := path.Join(tf, filename)
	// INFO: 0666 everyone can read and write so tests does not need to run with root privileges
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return &os.File{}, err
	}
	return fd, nil
}

func CheckGoroutineLeak(t *testing.T, before int) {
	time.Sleep(2 * time.Second)
	after := runtime.NumGoroutine()
	if after > before {
		t.Fatalf("Goroutines leak. Before: %d, After: %d", before, after)
	}
}

func DisconnectDeviceFile(df *os.File) error {
	// INFO: removing file will not close the file descriptor of keylogger
	// because if any process has the file open when this happens,
	// deletion is postponed until all processes have closed the file.
	// source: https://www.gnu.org/software/libc/manual/html_node/Deleting-Files.html
	// For simulating device disconnection, we close file descriptor
	err := df.Close()
	if err != nil {
		return err
	}
	return os.Remove(df.Name())
}

func ReconnectDeviceFile(df *os.File) error {
	fd, err := os.OpenFile(df.Name(), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	*df = *fd
	return nil
}

func WriteKeyDeviceFile(fd *os.File, code uint16) error {
	for _, i := range []int32{int32(KeyPress), int32(KeyRelease)} {
		slog.Info(fmt.Sprintf("writing key: %d, isRelease %d\n", code, i))
		err := binary.Write(
			fd,
			binary.LittleEndian,
			inputEvent{Type: evKey, Code: code, Value: i},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
