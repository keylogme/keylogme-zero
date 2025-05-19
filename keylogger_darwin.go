package keylog

/*
// -I include current directory for headers
#cgo CFLAGS: -I.
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include "keylogger_darwin.h"
*/
import "C"

import (
	"fmt"
	"log/slog"
	"runtime"
)

var hidManager = map[int]map[int]chan inputEvent{}

type keylogger struct {
	vendorID  int
	productID int
	loop      C.CFRunLoopRef
	hid       C.IOHIDManagerRef
}

func newKeylogger(productID string) (*keylogger, error) {
	// C.ListConnectedHIDDevices()
	deviceExists := make(chan bool)
	k := &keylogger{vendorID: 0x4653, productID: 0x0001}
	go func() {
		// INFO: lock goroutine to thread so CFRunLoopRun is in
		// same goroutine's thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		exists := C.setupDevice(&k.hid, 0x4653, 0x0001)
		deviceExists <- bool(exists)
		if !exists {
			// TODO: handle error no found
			// return &keylogger{}, fmt.Errorf("Device not available\n")
			return
		}
		k.loop = C.CFRunLoopGetCurrent()
		C.Start(k.hid, k.loop)
		fmt.Println("Run loop has exited...")
	}()

	// wait for device to be found
	if !<-deviceExists {
		slog.Debug("Device not found")
		return nil, fmt.Errorf("Device not available")
	}
	slog.Debug("keylogger MacOS started")
	return k, nil
}

func (k *keylogger) Read() chan inputEvent {
	if _, ok := hidManager[k.vendorID]; !ok {
		hidManager[k.vendorID] = map[int]chan inputEvent{}
	}
	if _, ok := hidManager[k.vendorID][k.productID]; !ok {
		event := make(chan inputEvent)
		hidManager[k.vendorID][k.productID] = event
		fmt.Println("Created channel keylogger....")
	}
	return hidManager[k.vendorID][k.productID]
}

func (k *keylogger) Close() error {
	if _, ok := hidManager[k.vendorID]; !ok {
		return nil
	}
	if _, ok := hidManager[k.vendorID][k.productID]; !ok {
		return nil
	}
	delete(hidManager[k.vendorID], k.productID)
	C.stopDevice(k.loop, k.hid)
	return nil
}

//export GoHandleKeyEvent
func GoHandleKeyEvent(code, value, vendorID, productID C.int) {
	vID := int(vendorID)
	pID := int(productID)

	if _, ok := hidManager[vID]; !ok {
		slog.Debug("Vendor id not in HIDManager")
		return
	}
	if _, ok := hidManager[vID][pID]; !ok {
		slog.Debug("Vendor id and product id not in HIDManager")
		return
	}
	pressed := int32(value)
	if pressed != 0 && pressed != 1 {
		return
	}
	// fmt.Printf(
	// 	"[Event] code=%d, value=%d, VID=0x%04x, PID=0x%04x\n",
	// 	code,
	// 	value,
	// 	vendorID,
	// 	productID,
	// )
	// c := uint16(code)
	hidManager[vID][pID] <- inputEvent{Type: evKey, Code: uint16(code), Value: int32(value)}
}

//export GoHandleDeviceEvent
func GoHandleDeviceEvent(vendorID, productID, connected C.int) {
	status := "disconnected"
	if connected != 0 {
		status = "connected"
		// INFO: Read will add device to hidManager
		return
	}
	// disconnect
	fmt.Printf("[Device] %s: VID=0x%04x, PID=0x%04x\n", status, vendorID, productID)

	if _, ok := hidManager[int(vendorID)]; !ok {
		return
	}
	if _, ok := hidManager[int(vendorID)][int(productID)]; ok {
		close(hidManager[int(vendorID)][int(productID)]) // close channel
		delete(hidManager[int(vendorID)], int(productID))
	}
}
