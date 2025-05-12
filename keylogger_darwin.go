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
)

var hidManager = map[int]map[int]chan inputEvent{}

type keylogger struct {
	vendorID  int
	productID int
}

func newKeylogger(productID string) (*keylogger, error) {
	fmt.Println("heeheheheeee")
	C.ListConnectedHIDDevices()
	// match keyboard, mouse and trackpad devices
	C.setupDevice(0x4653, 0x0001)

	go func() {
		C.Start()
	}()
	fmt.Println("keylogger started.............")

	return &keylogger{vendorID: 0x4653, productID: 0x0001}, nil
}

func (k *keylogger) Read() chan inputEvent {
	fmt.Println(hidManager)
	if _, ok := hidManager[k.vendorID]; !ok {
		fmt.Println(1)
		hidManager[k.vendorID] = map[int]chan inputEvent{}
	}
	if _, ok := hidManager[k.vendorID][k.productID]; !ok {
		fmt.Println(2)
		event := make(chan inputEvent)
		hidManager[k.vendorID][k.productID] = event
		fmt.Println("Created channel keylogger....")
	}
	fmt.Println(3)
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
	return nil
}

//export GoHandleKeyEvent
func GoHandleKeyEvent(code, value, vendorID, productID C.int) {
	vID := int(vendorID)
	pID := int(productID)

	if _, ok := hidManager[vID]; !ok {
		return
	}
	if _, ok := hidManager[vID][pID]; !ok {
		return
	}
	pressed := int32(value)
	if pressed != 0 && pressed != 1 {
		return
	}
	fmt.Printf(
		"[Event] code=%d, value=%d, VID=0x%04x, PID=0x%04x\n",
		code,
		value,
		vendorID,
		productID,
	)
	// c := uint16(code)
	hidManager[vID][pID] <- inputEvent{Type: evKey, Code: uint16(code), Value: int32(value)}
}

//export GoHandleDeviceEvent
func GoHandleDeviceEvent(vendorID, productID, connected C.int) {
	status := "disconnected"
	if connected != 0 {
		status = "connected"
	}
	fmt.Printf("[Device] %s: VID=0x%04x, PID=0x%04x\n", status, vendorID, productID)
}
