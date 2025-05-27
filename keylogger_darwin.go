package k0

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
	"time"

	"github.com/keylogme/keylogme-zero/types"
)

type KeyloggerInput struct {
	VendorID  types.Hex `json:"vendor_id"`
	ProductID types.Hex `json:"product_id"`
}

var hidManager = map[int]map[int]chan InputEvent{}

type keylogger struct {
	vendorID  int
	productID int
}

func NewKeylogger(kInput KeyloggerInput) (*keylogger, error) {
	// C.ListConnectedHIDDevices()
	exists := C.checkDeviceIsConnected(C.int(kInput.VendorID), C.int(kInput.ProductID))
	if !exists {
		slog.Debug("Device not found")
		return nil, fmt.Errorf("Device not available")
	}
	k := &keylogger{vendorID: int(kInput.VendorID), productID: int(kInput.ProductID)}

	go func() {
		// INFO: lock goroutine to thread so CFRunLoopRun is in
		// same goroutine's thread
		// runtime.LockOSThread()
		// defer runtime.UnlockOSThread()

		// TODO: add mutex to prevent multiple calls
		C.Start()
	}()
	return k, nil
}

func (k *keylogger) Read() chan InputEvent {
	if _, ok := hidManager[k.vendorID]; !ok {
		hidManager[k.vendorID] = map[int]chan InputEvent{}
	}
	if _, ok := hidManager[k.vendorID][k.productID]; !ok {
		event := make(chan InputEvent)
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
	if len(hidManager[k.vendorID]) == 0 {
		delete(hidManager, k.vendorID)
	}
	// stop if no devices left
	if len(hidManager) == 0 {
		C.Stop()
	}
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
	hidManager[vID][pID] <- InputEvent{Time: time.Now(), Code: uint16(code), Value: KeyEvent(value)}
}

//export GoHandleDeviceEvent
func GoHandleDeviceEvent(vendorID, productID, connected C.int) {
	status := "disconnected"
	if connected != 0 {
		status = "connected"
		// INFO: Read will add device to hidManager
		fmt.Printf("[Device] %s: VID=0x%04x, PID=0x%04x\n", status, vendorID, productID)
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
		if len(hidManager[int(vendorID)]) == 0 {
			delete(hidManager, int(vendorID))
		}
	}
	// stop if no devices left
	if len(hidManager) == 0 {
		C.Stop()
	}
}
