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

func TestFileDescriptor(t *testing.T) {
	k := &keyLogger{}

	err := k.Close()
	if err != nil {
		t.Error("Closing empty file descriptor should not yield error", err)
		return
	}
}

func TestBufferParser(t *testing.T) {
	k := &keyLogger{}

	// keyboard
	input, err := k.eventFromBuffer(
		[]byte{138, 180, 84, 92, 0, 0, 0, 0, 62, 75, 8, 0, 0, 0, 0, 0, 4, 0, 4, 0, 30, 0, 0, 0},
	)
	if err != nil {
		t.Error(err)
		return
	}
	if input == nil {
		t.Error("Event is empty, expected parsed event")
		return
	}

	if input.KeyString() != "3" {
		t.Errorf("wrong input key. got %v, expected %v", input.KeyString(), "3")
		return
	}

	// if input.Type != evMsc {
	// 	t.Errorf("wrong event type. expected key press but got %v", input.Type)
	// 	return
	// }
}

func TestWithPermission(t *testing.T) {
	fd, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	// try to create new keylogger with file descriptor which has the permission
	k, err := NewKeylogger(KeyloggerInput{UsbName: fd.Name()})
	if err != nil {
		t.Fatal(err)
	}
	k.Close()
	fd.Close()

	// try to create new keylogger with file descriptor which has no permission
	_, err = NewKeylogger(KeyloggerInput{UsbName: "/dev/tty0"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "permission denied. run with root permission" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func initDeviceFile() (*os.File, error) {
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

func disconnectDeviceFile(df *os.File) error {
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

func reconnectDeviceFile(df *os.File) error {
	fd, err := os.OpenFile(df.Name(), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	*df = *fd
	return nil
}

func writeKeyDeviceFile(fd *os.File, code uint16) error {
	for _, i := range []int32{int32(KeyPress), int32(KeyRelease)} {
		slog.Info(fmt.Sprintf("writing key: %d, isRelease %d\n", code, i))
		err := binary.Write(fd, binary.LittleEndian, inputEvent{Type: evKey, Code: code, Value: i})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestKeylog(t *testing.T) {
	before := runtime.NumGoroutine()
	defer checkGoroutineLeak(t, before)

	df, err := initDeviceFile()
	if err != nil {
		t.Fatal(err)
	}
	defer df.Close()
	deviceFile := df.Name()

	k, err := NewKeylogger(KeyloggerInput{UsbName: deviceFile})
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	// run goroutine to receive keypress
	recEvt := make(chan InputEvent, 1)
	go func() {
		t.Log("Starting goroutine")
		for i := range k.Read() {
			t.Logf("Received from k.Read(): %+v\n", i)
			recEvt <- i
		}
		t.Log("Exiting goroutine")
	}()
	// test
	time.Sleep(200 * time.Millisecond)
	t.Log("writing..")
	err = writeKeyDeviceFile(df, uint16(1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("check keypress is received (keyrelease is not checked)..")
	select {
	case result := <-recEvt:
		t.Logf("Received: %+v\n", result)
		if result.Code != uint16(1) {
			t.Fatal("Wrong code")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test listener timed out")
	}
}

// when you remove a usb device from the computer, the device file is removed
func TestDisconnectionKeylogger(t *testing.T) {
	before := runtime.NumGoroutine()
	defer checkGoroutineLeak(t, before)

	fd, err := initDeviceFile()
	if err != nil {
		t.Fatal(err)
	}
	deviceFile := fd.Name()

	// try to create new keylogger with file descriptor which has the permission
	k, err := NewKeylogger(KeyloggerInput{UsbName: deviceFile})
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	// run goroutine to receive keypress
	closedSig := make(chan int)
	go func() {
		for i := range k.Read() {
			fmt.Println(i)
		}
		fmt.Println("Out of loop")
		closedSig <- 1
	}()
	time.Sleep(200 * time.Millisecond)
	// disconnect
	err = disconnectDeviceFile(k.fd)
	if err != nil {
		t.Fatal(err)
	}
	// block until decive disconnected or timeout
	select {
	case <-closedSig:
		break
	case <-time.After(3 * time.Second):
		t.Fatal("Test listener timed out")
	}
}
