package keylogger

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

func setDevicePathFinder(pathDevice string) {
	// This is a mock function for testing purposes
	PathFinder = func(input types.KeyloggerInput) []string {
		return []string{pathDevice}
	}
}

func TestFileDescriptor(t *testing.T) {
	k := &KeyLogger{}

	err := k.Close()
	if err != nil {
		t.Error("Closing empty file descriptor should not yield error", err)
		return
	}
}

func TestBufferParser(t *testing.T) {
	k := &KeyLogger{}

	// keyboard
	input, err := k.eventFromBuffer(
		[]byte{138, 180, 84, 92, 0, 0, 0, 0, 62, 75, 8, 0, 0, 0, 0, 0, 4, 0, 4, 0, 30, 0, 0, 0},
	)
	if err != nil {
		t.Error(err)
		return
	}
	if input != nil {
		t.Error("Event should be empty because it is not an event key")
		return
	}
}

func TestWithPermission(t *testing.T) {
	fd, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	// try to create new keylogger with file descriptor which has the permission
	setDevicePathFinder(fd.Name())
	k, err := NewKeylogger(types.KeyloggerInput{})
	if err != nil {
		t.Fatal(err)
	}
	k.Close()
	fd.Close()

	// try to create new keylogger with file descriptor which has no permission
	setDevicePathFinder("/dev/tty0")
	k, err = NewKeylogger(types.KeyloggerInput{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "permission denied. run with root permission" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeylog(t *testing.T) {
	before := runtime.NumGoroutine()
	defer CheckGoroutineLeak(t, before)

	df, err := InitDeviceFile()
	if err != nil {
		t.Fatal(err)
	}
	defer df.Close()
	deviceFile := df.Name()

	setDevicePathFinder(deviceFile)
	k, err := NewKeylogger(types.KeyloggerInput{})
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
	err = WriteKeyDeviceFile(df, uint16(1))
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
	defer CheckGoroutineLeak(t, before)

	fd, err := InitDeviceFile()
	if err != nil {
		t.Fatal(err)
	}
	deviceFile := fd.Name()

	// try to create new keylogger with file descriptor which has the permission
	setDevicePathFinder(deviceFile)
	k, err := NewKeylogger(types.KeyloggerInput{})
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
	err = DisconnectDeviceFile(k.FD[0])
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
