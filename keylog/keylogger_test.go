package keylog

import (
	"encoding/binary"
	"fmt"
	"os"
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

	if input.Type != evMsc {
		t.Errorf("wrong event type. expected key press but got %v", input.Type)
		return
	}
}

func TestWithPermission(t *testing.T) {
	fd, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	// try to create new keylogger with file descriptor which has the permission
	k, err := newKeylogger(fd.Name())
	if err != nil {
		t.Fatal(err)
	}
	k.Close()
	fd.Close()

	// try to create new keylogger with file descriptor which has no permission
	_, err = newKeylogger("/dev/tty0")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "permission denied. run with root permission or use a user with access to /dev/tty0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeKeyOnceForTesting(filename string, code uint16) error {
	fd, err := os.OpenFile(filename, os.O_WRONLY, os.ModeCharDevice)
	if err != nil {
		return err
	}
	for _, i := range []int32{int32(KeyPress), int32(KeyRelease)} {
		err := binary.Write(fd, binary.LittleEndian, inputEvent{Type: evKey, Code: code, Value: i})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestKeylog(t *testing.T) {
	fd, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	// try to create new keylogger with file descriptor which has the permission
	k, err := newKeylogger(fd.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	// run goroutine to receive keypress
	recEvt := make(chan inputEvent)
	go func() {
		// FIXME: I added this sleep to make sure select can receive the channel
		time.Sleep(1 * time.Second)
		for i := range k.Read() {
			recEvt <- i
		}
	}()
	// test
	fmt.Println("writing..")
	err = writeKeyOnceForTesting(fd.Name(), uint16(1))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("check keypress..")
	// block until keypress received or timeout
	select {
	case result := <-recEvt:
		if result.Code != uint16(1) {
			t.Fatal("Wrong code")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test listener timed out")
	}
}

func TestDisconnection(t *testing.T) {
	fd, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	// try to create new keylogger with file descriptor which has the permission
	k, err := newKeylogger(fd.Name())
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
	// test
	err = os.Remove(fd.Name())
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
